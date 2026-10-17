# 鉴权、接口契约与订单落库设计

- 日期：2026-10-17
- 状态：已确认，待实现
- 范围：第一批（JWT 鉴权闭环、前后端响应契约统一、订单服务接入 MySQL）

## 1. 背景

当前项目能编译、单测能通过，但主链路没有闭环：

- 网关只认硬编码的 `valid_token`，登录接口本身又被全局鉴权挡住
- user-api / order-api 没有鉴权，前端 token 只用于路由守卫
- 登录返回的 `token_<username>` 可预测，且 user-rpc 不使用 proto 里已有的 `token` 字段
- user-api 返回 `{message}`，前端按 `{code,msg,data}` 解析；order-api 的 `data` 是字符串，前端却当对象用
- order-rpc 不落库：创建只调用一次 user-rpc 并吞掉错误，查询固定返回 `"1111"`
- 订单金额是 `int32`，创建时间由客户端传入，无任何参数校验

本设计只解决上述主链路问题，不扩展到注册、分页、链路追踪等后续工作。

## 2. 目标与非目标

### 目标

1. 登录后签发可验证的 JWT；user-api、order-api、gateway 都校验 JWT
2. 所有 HTTP 接口统一为 `{code, msg, data}` 响应结构
3. 订单写入 MySQL，订单归属当前登录用户；越权查询返回 404
4. 服务端完成参数校验、金额单位统一为“分”、创建时间由服务端生成
5. 保证现有 Go 模块构建、测试通过，`gofmt` 干净，前端可构建

### 非目标（本批不做）

- 幂等键、refresh token、服务端登出/黑名单
- 数据库密码/DSN 环境变量化、密钥轮换（第二批）
- 用户注册、用户列表分页、用户软删除
- 订单列表、状态流转、支付
- 链路追踪、Prometheus 配置修复
- 前端 UI 重设计

## 3. 关键决策

### 3.1 JWT 由 user-rpc 签发

user-rpc 的 `LoginResp` 已有 `token` 字段。登录成功后由 user-rpc 签发 JWT：

- 算法：HS256
- Claims：`user_id`、`username`、`iat`、`exp`
- 有效期：默认 24 小时，配置项 `Jwt.ExpireHours`
- 密钥：环境变量 `JWT_SECRET`，通过 go-zero `conf.UseEnv()` 注入
- 启动时校验密钥长度至少 32 字符，不满足直接启动失败

user-api 只透传 token；gateway 的 `/user/login` 映射到 user-rpc 后也能拿到同一个 token，避免两条登录链路不一致。

### 3.2 用户身份通过 gRPC metadata 传给 order-rpc

- HTTP 中间件校验 JWT 后，把 `user_id` 放入 `context`
- order-api 调用 order-rpc 前通过 `metadata.AppendToOutgoingContext` 写入 `gateway-user-id`
- gateway 中间件写入 HTTP 头 `Grpc-Metadata-User-Id`，go-zero gateway 会转换成同名 gRPC metadata `gateway-user-id`
- order-rpc 只从 incoming metadata 读取 `user_id`，不接受请求体里的用户 ID，避免客户端伪造
- order-rpc 的端口只在内网暴露，不作为公网入口

### 3.3 订单归属校验

- `orders.user_id` 来自 metadata
- 查询固定条件：`id = ? AND user_id = ? AND deleted_at IS NULL`
- 记录不存在或不属于当前用户，统一返回业务码 404，不暴露记录是否存在
- 本批不做幂等键

## 4. 总体架构

```
前端 ──> user-api  ──> user-rpc ──> MySQL(user)
  │        │ 校验JWT      │ 登录校验密码并签发JWT
  │        │
  └────> order-api ──> order-rpc ──> MySQL(order)
           │ 校验JWT      │ 从 metadata 取 user_id
           │ 注入 user_id

外部 ──> gateway ──> user-rpc / order-rpc
            │ 校验JWT，写 Grpc-Metadata-User-Id
```

可独立测试的单元：

| 单元 | 职责 | 依赖 |
|---|---|---|
| `pkg/auth` | 签发/解析 JWT、HTTP 中间件、context 与 metadata 工具 | 标准库、jwt/v4、gRPC |
| user-rpc logic | 登录校验、签发 token、用户查询 | MySQL、`pkg/auth` |
| user-api logic | 透传 RPC 结果、组装统一响应 | user-rpc |
| order-api logic | 注入 user_id、透传 RPC 结果 | order-rpc |
| order-rpc logic | 参数校验、订单写入/查询、归属校验 | MySQL、metadata |

## 5. 组件设计

### 5.1 新增 `pkg/auth`

独立 Go module：`go-zero-demo/pkg/auth`，各服务通过 `replace` 引用，方式与 `pkg/errcode` 一致。

对外接口：

```go
const MetadataUserIDKey = "gateway-user-id"

type Claims struct {
    UserID   int64  `json:"user_id"`
    Username string `json:"username"`
    jwt.RegisteredClaims
}

func ValidateSecret(secret string) error
func Sign(secret string, userID int64, username string, ttl time.Duration) (string, error)
func Parse(secret, token string) (*Claims, error)

// HTTP 中间件，签名与 go-zero 的 rest.Middleware / gateway 中间件一致
func Middleware(secret string, skipPaths ...string) func(http.HandlerFunc) http.HandlerFunc

func UserID(ctx context.Context) (int64, bool)
func Username(ctx context.Context) (string, bool)

// 写入 outgoing gRPC metadata
func InjectUserID(ctx context.Context, userID int64) context.Context
// 从 incoming gRPC metadata 读取
func UserIDFromIncoming(ctx context.Context) (int64, bool)
```

中间件行为：

1. `r.URL.Path` 命中 `skipPaths` 时直接放行
2. 从 `Authorization` 读取 token，兼容 `Bearer <token>` 与裸 token
3. 解析失败、过期、`user_id <= 0` 时返回 HTTP 401：`{"code":401,"msg":"未授权"}`
4. 成功后把 `user_id`、`username` 写入 context
5. 成功后**覆盖**请求头 `Grpc-Metadata-User-Id` 和 `Grpc-Metadata-Username`；跳过鉴权的路径则删除这两个头，防止客户端伪造后透传给 gRPC
6. JSON 编码失败通过 `logx.Errorf` 记录，不静默忽略

### 5.2 配置

所有服务将 `conf.MustLoad` 改为带 `conf.UseEnv()` 的形式：

```yaml
Jwt:
  Secret: ${JWT_SECRET}
```

- user-rpc 额外配置 `Jwt.ExpireHours`，缺省 24
- gateway 使用自定义 `config.Config`，内嵌 `gateway.GatewayConf` 并增加 `Jwt` 字段
- 每个服务启动时调用 `auth.ValidateSecret`，失败则退出

本地开发需要 `export JWT_SECRET=...`（至少 32 字符）；Docker 环境在 compose 中为相关服务设置同一密钥。

### 5.3 user-rpc

`Login`：

1. 查询 `SELECT id, username, password FROM user WHERE username = ?`
2. `sql.ErrNoRows` → 业务码 401
3. 其他数据库错误 → 返回 error（不伪装成 401）
4. bcrypt 校验失败 → 业务码 401
5. 签发 JWT，返回 `code=200, msg="ok", token, user_id, username`

`User`：

1. 校验 `name` 非空
2. 查询 `SELECT id, username, nickname FROM user WHERE username = ?`
3. 不存在 → 业务码 404；数据库错误 → 返回 error
4. 成功返回 `UserInfo`

不记录密码、token。

### 5.4 user-api

- 全局启用 `auth.Middleware(secret, "/user/login")`
- `Login` 透传 RPC 的 `code/msg/token/user_id/username`
- `User` 透传 RPC 的 `code/msg/data`
- 统一错误处理器：
  - gRPC `Unauthenticated` → 401，`PermissionDenied` → 403，`NotFound` → 404，`InvalidArgument` → 400
  - gRPC `Internal` / `Unknown` 和连接错误 → 500，只返回通用提示
  - 非 gRPC 错误（参数解析、校验）→ 400
  - 系统错误不返回内部细节
- 修正 `.api` 与生成代码漂移：`LoginResp` 补上 `Token`

### 5.5 order-api

- 全局启用 JWT 中间件
- 从 context 取 `user_id`，取不到则返回业务码 401；调用 RPC 前 `auth.InjectUserID`
- 创建请求只接收 `goods_id`、`goods_name`、`price`
- 响应统一为 `{code,msg,data}`，`data` 为 `OrderInfo`

### 5.6 order-rpc

- 新增 MySQL 配置与连接，连接参数与 user-rpc 保持一致（最大连接 20、空闲 10）
- `CreateOrder`：
  1. 从 incoming metadata 取 `user_id`，缺失或非法 → 业务码 401
  2. 校验 `goods_id > 0`、`goods_name` 去空格后非空且不超过 128 字符、`price > 0`
  3. 校验失败 → 业务码 400
  4. `INSERT INTO orders (user_id, goods_id, goods_name, price, status) VALUES (?, ?, ?, ?, 1)`
  5. 使用 `LastInsertId` 查询并返回完整 `OrderInfo`；创建时间由数据库生成
  6. 数据库错误 → 返回 error
- `GetOrder`：
  1. 从 metadata 取 `user_id`
  2. `SELECT id, goods_id, goods_name, price, status, UNIX_TIMESTAMP(created_at) FROM orders WHERE id = ? AND user_id = ? AND deleted_at IS NULL`
  3. 无记录 → 业务码 404
  4. 数据库错误 → 返回 error

### 5.7 gateway

- 自定义配置加载 `Jwt.Secret`
- 全局启用 `auth.Middleware(secret, "/user/login")`
- 中间件在鉴权成功后写入 `Grpc-Metadata-User-Id`，供 order-rpc 读取
- 错误处理器按 `BizError.Code` 返回对应的 HTTP 状态（401/403/404/400），不再一律 400
- 保留现有 RPC 映射

## 6. 接口契约

### 6.1 HTTP

所有响应统一为：

```json
{ "code": 200, "msg": "ok", "data": {} }
```

| 方法 | 路径 | 鉴权 | 请求 | 成功 data |
|---|---|---|---|---|
| POST | `/user/login` | 否 | `{username,password}` | `{token,user_id,username}` |
| GET | `/user/:name` | 是 | path `name` | `{id,username,nickname}` |
| POST | `/order/create` | 是 | `{goods_id,goods_name,price}` | `OrderInfo` |
| GET | `/order/:id` | 是 | path `id` | `OrderInfo` |

`OrderInfo`：

```json
{
  "id": 1,
  "goods_id": 1001,
  "goods_name": "book",
  "price": 9900,
  "status": 1,
  "created_at": 1760688000
}
```

- `price` 单位为分，`int64`
- `created_at` 为 Unix 秒，由服务端生成

### 6.2 user.proto

```proto
message UserReq { string name = 1; }

message UserInfo {
  int64 id = 1;
  string username = 2;
  string nickname = 3;
}

message UserResp {
  int32 code = 1;
  string msg = 2;
  UserInfo data = 3;
}

message LoginReq {
  string username = 1;
  string password = 2;
}

message LoginResp {
  int32 code = 1;
  string msg = 2;
  string token = 3;
  int64 user_id = 4;
  string username = 5;
}
```

### 6.3 order.proto

```proto
message GetOrderReq { int64 id = 1; }

message CreateOrderReq {
  int64 goods_id = 1;
  string goods_name = 2;
  int64 price = 3;
  reserved 4;
  reserved "created_at";
}

message OrderInfo {
  int64 id = 1;
  int64 goods_id = 2;
  string goods_name = 3;
  int64 price = 4;
  int32 status = 5;
  int64 created_at = 6;
}

message CommonResp {
  int32 code = 1;
  string msg = 2;
  OrderInfo data = 3;
}
```

注意：`GetOrderReq.id`、`CreateOrderReq.price` 从 `int32` 改为 `int64`，需要重新生成 pb 与 client。

## 7. 数据模型

新增 `deploy/sql/init_order.sql`：

```sql
CREATE DATABASE IF NOT EXISTS `looklook_order`
  DEFAULT CHARACTER SET utf8mb4
  COLLATE utf8mb4_general_ci;

USE `looklook_order`;

CREATE TABLE IF NOT EXISTS `orders` (
  `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id`    BIGINT NOT NULL COMMENT '下单用户 ID',
  `goods_id`   BIGINT NOT NULL COMMENT '商品 ID',
  `goods_name` VARCHAR(128) NOT NULL COMMENT '商品名称',
  `price`      BIGINT NOT NULL COMMENT '金额，单位分',
  `status`     TINYINT NOT NULL DEFAULT 1 COMMENT '订单状态：1=已创建',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_created` (`user_id`, `created_at`),
  KEY `idx_goods` (`goods_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单表';
```

- docker-compose 将文件挂载为 `/docker-entrypoint-initdb.d/02-init-order.sql`
- 已有 MySQL 数据卷不会自动执行初始化脚本，需要在 README 中说明手动执行
- order-rpc 使用独立数据库 `looklook_order`

## 8. 错误处理

| 场景 | HTTP / RPC 表现 |
|---|---|
| 缺少 token、token 非法或过期 | HTTP 401，业务码 401 |
| 请求参数解析或校验失败 | 业务码 400 |
| 用户不存在 | 业务码 404 |
| 订单不存在或不属于当前用户 | 业务码 404 |
| metadata 缺少 user_id | order-rpc 业务码 401 |
| 数据库错误、RPC 连接错误 | 返回 error；API 统一 500，客户端只看到通用提示 |

- user-api / order-api 的业务错误（400/401/404）放在响应体 `code/msg` 中，HTTP 状态为 200；鉴权中间件直接返回 HTTP 401
- gateway 的错误处理器同时设置响应体 code 和对应的 HTTP 状态
- 系统错误不返回 SQL、连接串等内部信息
- 所有 error 都必须处理或返回，不允许 `_` 忽略

## 9. 前端变更

- `src/types/api.ts`：
  - `ApiResponse<T>` 保持 `{code,msg,data}`
  - `UserInfo` 改为 `{id,username,nickname}`
  - `LoginResp` 改为 `{code,msg,data:{token,user_id,username}}`
  - `OrderInfo` 改为 `{id,goodsId,goodsName,price,status,createdAt}`
  - `CreateOrderReq` 去掉 `created_at`
- `src/utils/request.ts`：发送 `Authorization: Bearer <token>`
- `src/stores/user.ts`：从 `res.data` 读取 token
- `src/views/user/UserList.vue`：展示 `username`、`nickname`
- `src/views/order/OrderCreate.vue`：价格按“分”输入，不再传 `created_at`
- `src/views/order/OrderList.vue`：按新 `OrderInfo` 字段展示
- `src/views/Dashboard.vue`：不受影响，仍通过 `getUser('admin')` 验证连通性

## 10. 测试策略

按 TDD 执行，先写失败测试再实现。

| 模块 | 测试 |
|---|---|
| `pkg/auth` | 签名/解析往返、过期 token、错误密钥、中间件放行登录、拦截缺失/非法 token、context 读写、metadata 注入与读取 |
| user-rpc | 登录成功签发 token、密码错误 401、用户不存在 401、数据库错误返回 error、用户查询成功/404 |
| user-api | 登录透传 token、RPC 错误透传、用户查询透传 |
| order-rpc | 参数校验（user_id/goods_id/名称/价格）、创建成功、查询成功、越权 404、数据库错误 |
| order-api | metadata 注入、创建/查询透传、未登录 context 返回错误 |
| 构建 | 所有模块 `go build ./...`、`go test ./...`、`gofmt -l` 为空 |
| 前端 | `npm run build`（含 vue-tsc） |

## 11. 验收标准

1. 所有 Go 模块构建、测试通过，`gofmt -l` 无输出
2. 前端构建通过
3. 通过 user-api：
   - 不带 token 访问 `/user/admin` → 401
   - `admin / 123456` 登录 → 获得 JWT
   - 带 token 创建订单 → 返回订单 ID 与完整信息
   - 查询自己的订单 → 200
   - 查询他人的订单 → 404
4. 通过 gateway 重复登录、创建、查询链路，行为一致
5. 旧 token（伪造或过期）访问受保护接口 → 401

## 12. 实施顺序

1. 安装 `goctl`、`protoc`、`protoc-gen-go`、`protoc-gen-go-grpc`
2. 新建 `pkg/auth` 与测试
3. 更新 proto / `.api` 定义并重新生成代码
4. user-rpc 登录签发 JWT、用户查询
5. user-api / order-api 中间件与契约
6. order-rpc 数据库与订单逻辑
7. gateway 配置与中间件
8. 数据库脚本与 compose 挂载
9. 前端契约对齐
10. 全量验证与 `gofmt`

## 13. 风险与缓解

| 风险 | 缓解 |
|---|---|
| 重新生成 proto 会覆盖业务代码 | 只重新生成 pb/client/server，logic 手写内容先提交/比对；生成后跑测试 |
| 已有 MySQL 数据卷不执行新 SQL | README 说明手动执行 `init_order.sql`；验收前先执行 |
| RPC 端口公网可访问时 metadata 可被伪造 | 本批将 RPC 端口改为仅内网暴露；生产环境不发布 RPC 端口 |
| 前后端同时改契约导致中间态不可用 | 同一分支内一次完成，验收脚本覆盖登录→下单→查询 |
| JWT_SECRET 未设置导致启动失败 | README、compose 和启动检查给出明确错误信息 |
