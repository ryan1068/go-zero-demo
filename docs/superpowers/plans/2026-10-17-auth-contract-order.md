# JWT 鉴权、契约统一与订单落库 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 让登录签发可验证 JWT，所有 HTTP 接口统一 `{code,msg,data}`，订单写入 MySQL 并只允许查询自己的订单。

**Architecture:** user-rpc 签发 JWT；`pkg/auth` 提供签名、HTTP 中间件和 gRPC metadata 工具；user-api/order-api/gateway 统一校验；order-api 把 user_id 注入 outgoing metadata，order-rpc 只从 incoming metadata 取用户身份并访问 `looklook_order.orders`。

**Tech Stack:** Go 1.23 / go-zero 1.8.5 / gRPC / MySQL / Vue 3 + TypeScript

## Global Constraints

- 所有 error 必须处理或返回，不能 `_` 忽略
- 业务逻辑只放在 `internal/logic/`，不修改生成的 handler
- 用户输入必须校验
- 密码、JWT 密钥、token 不能打印到日志
- SQL 必须参数化，禁止 `SELECT *`
- 订单金额单位为分，类型 `int64`；创建时间由服务端生成
- 响应统一为 `{code,msg,data}`
- 订单查询必须带 `id + user_id + deleted_at IS NULL`
- 提交信息遵循 `<type>(<scope>): <subject>`

---

### Task 1: 新增 `pkg/auth`

**Files:**
- Create: `pkg/auth/go.mod`, `pkg/auth/auth.go`, `pkg/auth/middleware.go`
- Test: `pkg/auth/auth_test.go`, `pkg/auth/middleware_test.go`
- Modify: `app/gateway/go.mod`, `app/user/api/go.mod`, `app/user/rpc/go.mod`, `app/order/api/go.mod`, `app/order/rpc/go.mod`（加 require + replace）

**Interfaces:**
- Produces:
  - `type Claims struct { UserID int64; Username string; jwt.RegisteredClaims }`
  - `func ValidateSecret(secret string) error`
  - `func Sign(secret string, userID int64, username string, ttl time.Duration) (string, error)`
  - `func Parse(secret, token string) (*Claims, error)`
  - `func Middleware(secret string, skipPaths ...string) func(http.HandlerFunc) http.HandlerFunc`
  - `func UserID(ctx context.Context) (int64, bool)`
  - `func Username(ctx context.Context) (string, bool)`
  - `func InjectUserID(ctx context.Context, userID int64) context.Context`
  - `func UserIDFromIncoming(ctx context.Context) (int64, bool)`
  - `const MetadataUserIDKey = "gateway-user-id"`

- [ ] **Step 1: 写失败测试**：签名/解析往返；过期 token；错误密钥；空/过短密钥；中间件放行 `/user/login`、拒绝缺失/非法 token、成功后 context 可读；`InjectUserID` 后 `UserIDFromIncoming` 可读。
- [ ] **Step 2: 运行测试确认失败**：`cd pkg/auth && go test ./...`
- [ ] **Step 3: 实现**：HS256，claims 含 `user_id/username/iat/exp`；中间件兼容 `Bearer <token>` 与裸 token；成功后覆盖 `Grpc-Metadata-User-Id` / `Grpc-Metadata-Username`，跳过路径删除这两个头；失败返回 HTTP 401 JSON。
- [ ] **Step 4: 运行测试通过**：`cd pkg/auth && go test ./... -count=1`
- [ ] **Step 5: 提交**：`feat(auth): 新增 JWT 与身份上下文工具`

---

### Task 2: 更新 proto 与 `.api` 契约

**Files:**
- Modify: `app/user/rpc/user.proto`, `app/order/rpc/order.proto`
- Modify: `app/user/api/user.api`, `app/order/api/order.api`
- Modify/regenerate: `app/user/rpc/pb/user/user.pb.go`, `app/order/rpc/pb/order/order.pb.go`
- Modify: `app/user/api/internal/types/types.go`, `app/order/api/internal/types/types.go`
- Keep: `*_grpc.pb.go`、`*client/*.go`（服务方法签名不变，只需重新编译）

**Interfaces:**
- `user.UserResp{Code int32, Msg string, Data *UserInfo}`
- `user.UserInfo{Id int64, Username string, Nickname string}`
- `user.LoginResp{Code int32, Msg string, Token string, UserId int64, Username string}`
- `order.CreateOrderReq{GoodsId int64, GoodsName string, Price int64}`
- `order.GetOrderReq{Id int64}`
- `order.OrderInfo{Id int64, GoodsId int64, GoodsName string, Price int64, Status int32, CreatedAt int64}`
- `order.CommonResp{Code int32, Msg string, Data *OrderInfo}`

- [ ] **Step 1: 修改两个 `.proto`**，按 spec 第 6 节；`CreateOrderReq` 内 `reserved 4; reserved "created_at";`
- [ ] **Step 2: 离线生成 pb.go**（已验证可用）：

```bash
cd app/gateway
GOFLAGS=-mod=mod GOPROXY=off go run ./tools/protogen \
  /tmp/protogen-build/protoc-gen-go /tmp/protogen-out ../user/rpc/user.proto
GOFLAGS=-mod=mod GOPROXY=off go run ./tools/protogen \
  /tmp/protogen-build/protoc-gen-go /tmp/protogen-out ../order/rpc/order.proto
cp /tmp/protogen-out/user.pb.go ../user/rpc/pb/user/user.pb.go
cp /tmp/protogen-out/order.pb.go ../order/rpc/pb/order/order.pb.go
```

- [ ] **Step 3: 手工同步 `.api` 与 `types.go`**（无 goctl，保持源码与生成结果一致）：用户/订单类型按 spec 第 6 节。
- [ ] **Step 4: 编译**：各模块 `go build ./...`，先修掉类型引用错误。
- [ ] **Step 5: 提交**：`feat(api): 统一用户与订单接口契约`

---

### Task 3: user-rpc 登录签发 JWT、用户查询落库

**Files:**
- Modify: `app/user/rpc/internal/config/config.go`
- Modify: `app/user/rpc/internal/svc/servicecontext.go`（如需关闭 DB）
- Modify: `app/user/rpc/internal/logic/loginlogic.go`, `userlogic.go`
- Test: `app/user/rpc/internal/logic/loginlogic_test.go`, `userlogic_test.go`
- Modify: `app/user/rpc/user.go`, `app/user/rpc/etc/user.yaml`

**Interfaces:**
- Consumes: `auth.Sign`, `auth.ValidateSecret`
- Produces: `Login` 成功时 `LoginResp{Code:200,Token,UserId,Username}`；`User` 返回 `UserResp{Code,Msg,Data}`

- [ ] **Step 1: 写失败测试**：登录成功签发可解析 token；密码错误 401；用户不存在 401；DB 错误返回 error；用户查询成功/空名 400/不存在 404。
- [ ] **Step 2: 运行失败**：`cd app/user/rpc && go test ./internal/logic -run TestLogin -count=1`
- [ ] **Step 3: 实现**：`sql.ErrNoRows` 与系统错误分流；bcrypt；ttl 取 `ExpireHours`（<=0 用 24）；查询 `SELECT id, username, password FROM user WHERE username = ?`；用户查询 `SELECT id, username, COALESCE(nickname,'') FROM user WHERE username = ?`。
- [ ] **Step 4: 配置**：`Jwt.Secret: ${JWT_SECRET}`、`Jwt.ExpireHours: 24`；`conf.UseEnv()` + `auth.ValidateSecret`。
- [ ] **Step 5: 测试通过并提交**：`feat(user): 登录签发 JWT 并实现用户查询`

---

### Task 4: user-api 鉴权中间件与统一响应

**Files:**
- Modify: `app/user/api/user.go`, `app/user/api/internal/config/config.go`, `app/user/api/etc/user-api.yaml`
- Modify: `app/user/api/internal/logic/loginlogic.go`, `userlogic.go`
- Test: `app/user/api/internal/logic/loginlogic_test.go`, `userlogic_test.go`

- [ ] **Step 1: 写失败测试**：登录透传 token/user_id/username；RPC 错误返回 error；用户查询透传 data。
- [ ] **Step 2: 实现**：`server.Use(auth.Middleware(secret, "/user/login"))`；统一错误处理器映射 gRPC 状态；`Login` 组装 `types.LoginResp{Data: types.LoginData{...}}`。
- [ ] **Step 3: 测试通过并提交**：`feat(user-api): 接入 JWT 中间件并统一响应`

---

### Task 5: order-api 鉴权与 metadata 注入

**Files:**
- Modify: `app/order/api/order.go`, `app/order/api/internal/config/config.go`, `app/order/api/etc/order.yaml`
- Modify: `app/order/api/internal/logic/createorderlogic.go`, `getorderlogic.go`
- Test: `app/order/api/internal/logic/createorderlogic_test.go`, `getorderlogic_test.go`

- [ ] **Step 1: 写失败测试**：context 有 user_id 时 metadata 注入；无 user_id 返回 Unauthenticated；RPC 结果透传 `OrderInfo`。
- [ ] **Step 2: 实现**：`server.Use(auth.Middleware(secret))`；`auth.UserID(l.ctx)`；`ctx := auth.InjectUserID(l.ctx, uid)`；映射新 request/response 类型。
- [ ] **Step 3: 测试通过并提交**：`feat(order-api): 接入 JWT 并传递用户身份`

---

### Task 6: order-rpc 接入 MySQL 并实现订单落库

**Files:**
- Modify: `app/order/rpc/internal/config/config.go`, `internal/svc/servicecontext.go`
- Modify: `app/order/rpc/internal/logic/createorderlogic.go`, `getorderlogic.go`
- Test: `app/order/rpc/internal/logic/createorderlogic_test.go`, `getorderlogic_test.go`
- Modify: `app/order/rpc/etc/order.yaml`, `app/order/rpc/go.mod`
- Create: `deploy/sql/init_order.sql`

**Interfaces:**
- Consumes: `auth.UserIDFromIncoming`
- Produces: `CreateOrder` 返回完整 `OrderInfo`；`GetOrder` 只返回当前用户订单

- [ ] **Step 1: 写失败测试**（sqlmock + `metadata.NewIncomingContext`）：缺 user_id 401；goods_id<=0 / 空名 / price<=0 400；创建成功 `INSERT` + `SELECT`；查询参数为 `id,user_id`；越权/不存在 404；DB 错误返回 error。
- [ ] **Step 2: 实现**：`INSERT INTO orders (user_id, goods_id, goods_name, price, status) VALUES (?, ?, ?, ?, 1)`；`SELECT id, goods_id, goods_name, price, status, UNIX_TIMESTAMP(created_at) FROM orders WHERE id = ? AND user_id = ? AND deleted_at IS NULL`；`sql.ErrNoRows` → 404。
- [ ] **Step 3: 建表脚本**：按 spec 第 7 节；compose 挂载为 `02-init-order.sql`。
- [ ] **Step 4: 移除 order-rpc 对 user-rpc 的依赖**（ServiceContext、config、yaml、go.mod）。
- [ ] **Step 5: 测试通过并提交**：`feat(order): 订单接入 MySQL 并做归属校验`

---

### Task 7: gateway 鉴权与错误状态

**Files:**
- Modify: `app/gateway/internal/config/config.go`, `app/gateway/gateway.go`, `app/gateway/etc/gateway.yaml`
- Test: `app/gateway/internal/middleware/`（用 `pkg/auth` 的测试覆盖即可，网关增加一个错误处理器测试可选）

- [ ] **Step 1: 配置**：自定义 `Config` 内嵌 `gateway.GatewayConf` + `Jwt.Secret`；`conf.UseEnv()`；启动校验密钥。
- [ ] **Step 2: 中间件**：`gw.Use(auth.Middleware(secret, "/user/login"))`。
- [ ] **Step 3: 错误处理**：`BizError.Code` → 对应 HTTP 状态；其他 → 500。
- [ ] **Step 4: 构建并提交**：`feat(gateway): 接入 JWT 校验与身份透传`

---

### Task 8: 前端契约对齐

**Files:**
- Modify: `app/frontend/src/types/api.ts`, `src/utils/request.ts`, `src/stores/user.ts`
- Modify: `src/views/user/UserList.vue`, `src/views/order/OrderCreate.vue`, `src/views/order/OrderList.vue`

- [ ] **Step 1: 更新类型**：`UserInfo{id,username,nickname}`；`LoginResp{code,msg,data:{token,user_id,username}}`；`OrderInfo{id,goodsId,goodsName,price,status,createdAt}`；`CreateOrderReq` 去掉 `created_at`。
- [ ] **Step 2: 更新请求**：`Authorization: Bearer <token>`。
- [ ] **Step 3: 更新 store/views**：token 从 `res.data.token` 取；用户表展示 username/nickname；订单价格按分展示；创建订单不再传 created_at。
- [ ] **Step 4: 构建并提交**：`cd app/frontend && npm run build`；`feat(frontend): 对齐统一响应与 JWT 鉴权`

---

### Task 9: 配置、Docker 与文档

**Files:**
- Modify: `docker-compose.yml`, `.gitignore`
- Create: `.env.example`
- Modify: `README.md`（JWT_SECRET、init_order.sql、RPC 端口说明）
- Modify: `app/user/rpc/etc/user.yaml`, `app/user/api/etc/user-api.yaml`, `app/order/api/etc/order.yaml`, `app/order/rpc/etc/order.yaml`, `app/gateway/etc/gateway.yaml`

- [ ] **Step 1: compose**：为 gateway/user-rpc/user-api/order-api 注入 `JWT_SECRET`；order-rpc 增加 `Mysql` 配置；user-rpc/order-rpc 的 `ports` 改为 `expose`；挂载 `init_order.sql`；补 `depends_on: mysql`。
- [ ] **Step 2: `.env.example` + `.gitignore`**：忽略 `.env`。
- [ ] **Step 3: README**：说明 `cp .env.example .env`、设置 JWT_SECRET、已有 MySQL 卷需手动执行 `init_order.sql`。
- [ ] **Step 4: 提交**：`chore(deploy): 配置 JWT 与订单数据库`

---

### Task 10: 全量校验与格式化

- [ ] **Step 1: gofmt**：对各模块执行 `gofmt -w .`（覆盖 CI 要求的 10 个文件）。
- [ ] **Step 2: 构建/测试**：

```bash
for d in app/gateway app/user/api app/user/rpc app/order/api app/order/rpc pkg/auth pkg/errcode; do
  (cd "$d" && go build ./... && go test ./... -count=1) || exit 1
done
```

- [ ] **Step 3: gofmt 检查**：所有模块 `gofmt -l .` 输出为空。
- [ ] **Step 4: 前端**：`cd app/frontend && npm run build`。
- [ ] **Step 5: 清理临时工具**：删除 `app/gateway/tools/protogen`。
- [ ] **Step 6: 提交**：`chore: 格式化并补齐构建验证`

---

## 验收

```bash
# 1. 启动依赖 + 服务
cp .env.example .env   # 设置 JWT_SECRET 和 MySQL 密码
docker compose up -d --build

# 2. 登录
curl -s -X POST localhost:8000/user/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"123456"}'
# => {"code":200,"msg":"ok","data":{"token":"...","user_id":1,"username":"admin"}}

# 3. 无 token 查询 -> 401
curl -i localhost:8000/user/admin

# 4. 带 token 创建订单
curl -s -X POST localhost:8001/order/create \
  -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"goods_id":1001,"goods_name":"book","price":9900}'
# => {"code":200,"msg":"ok","data":{"id":1,...}}

# 5. 查询自己的订单 -> 200；用另一个用户的 token 查同一订单 -> 404
```
