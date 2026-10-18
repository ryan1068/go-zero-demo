# 拾光严选 · go-zero-demo

拾光严选是基于 go-zero 的商城示例，包含用户、商品、订单三个领域服务，以及网关和 Vue 3 前端。品牌 Slogan：把好物收进时光里。

## 架构

```text
                    ┌──────────────┐
                    │   frontend   │  :80
                    │ Vue3 + Elem  │
                    └──────┬───────┘
                           │ /api/*
          ┌────────────────┼────────────────┐
          ▼                ▼                ▼
   ┌────────────┐   ┌────────────┐   ┌────────────┐
   │  user-api  │   │ order-api  │   │product-api │
   │    :8000   │   │    :8001   │   │    :8002   │
   └─────┬──────┘   └─────┬──────┘   └─────┬──────┘
         │                │                │
   ┌─────▼──────┐   ┌─────▼──────┐   ┌─────▼──────┐
   │  user-rpc  │   │ order-rpc  │   │product-rpc │
   │    :9000   │   │    :9001   │   │    :9002   │
   └─────┬──────┘   └─────┬──────┘   └─────┬──────┘
         │                │                │
         ▼                ▼                ▼
   usercenter DB      order DB         product DB
```

| 服务 | 端口 | 说明 |
|---|---|---|
| gateway | 8888 | 统一网关（JWT、日志） |
| user-api | 8000 | 用户、注册、地址 |
| user-rpc | 9000 | 用户 gRPC |
| order-api | 8001 | 购物车、订单、支付、发货 |
| order-rpc | 9001 | 订单 gRPC |
| product-api | 8002 | 分类、商品、图片上传 |
| product-rpc | 9002 | 商品 gRPC |
| frontend | 8080 / 8081 | 8080=管理后台，8081=商城 |
| mysql | 33069 | MySQL 8.0 |
| etcd | 2379 | 服务注册发现 |

## 数据库

| 数据库 | 表 |
|---|---|
| `looklook_usercenter` | `user`、`user_address`、`admin_audit_log` |
| `looklook_product` | `category`、`product`、`product_sku` |
| `looklook_order` | `cart_item`、`orders`、`order_item`、`payment_record`、`order_log`、`after_sale`、`refund_record` |

字符集统一为 `utf8mb4`，排序规则统一为 `utf8mb4_unicode_ci`。MySQL 8 建表时必须显式写 `COLLATE=utf8mb4_unicode_ci`，只写 `DEFAULT CHARSET=utf8mb4` 会使用字符集自身的默认排序规则 `utf8mb4_0900_ai_ci`，导致库与表不一致。

### 数据层

三个 RPC 服务统一使用 goctl 生成的 `internal/model` + go-zero `sqlx`：

- 不使用 `database/sql`；软删除、分页筛选、全文检索、跨库查询写在 `*Model_custom.go`
- 事务使用 `sqlx.SqlConn.TransactCtx`，业务回滚返回 `model.ErrAbortTransaction`
- 当前不启用 model cache；order-rpc 跨库更新商品库存，缓存无法失效
- 重新生成：

```bash
export MYSQL_ROOT_PASSWORD="$(grep '^MYSQL_ROOT_PASSWORD=' .env | cut -d= -f2-)"
make gen-model
```

## 快速开始

```bash
# 1. 准备环境变量（复制后设置 JWT_SECRET、MYSQL_ROOT_PASSWORD、REDIS_PASSWORD）
cp .env.example .env

# 2. 启动所有服务
docker compose up -d --build

# 3. 已有 MySQL 数据卷时，依次执行迁移（只需执行一次）
#    密码从 .env 读取；已有数据卷必须保持原来的 root 密码
export MYSQL_ROOT_PASSWORD="$(grep '^MYSQL_ROOT_PASSWORD=' .env | cut -d= -f2-)"
for f in deploy/sql/migrations/20261018_*.sql; do
  docker exec -i mysql mysql -uroot --default-character-set=utf8mb4 -p"$MYSQL_ROOT_PASSWORD" < "$f"
done
```

> 全新的 MySQL 数据卷会自动执行 `init_user.sql`、`init_order.sql`、`init_product.sql`。
> 已有数据卷不会重复执行初始化 SQL，需要手动跑迁移；`MYSQL_ROOT_PASSWORD` 只在首次初始化时生效，已有数据卷请沿用原密码。
>
> 前端 nginx 已使用 Docker 内置 DNS（`resolver 127.0.0.11`），后端容器重建后不需要再重启 frontend。

### 国内网络

Docker Hub 拉取超时时，可以先把基础镜像从国内镜像站拉下来，再打回官方标签：

```bash
docker pull docker.1panel.live/library/golang:alpine
docker tag docker.1panel.live/library/golang:alpine golang:alpine
docker pull docker.1panel.live/library/node:20-alpine
docker tag docker.1panel.live/library/node:20-alpine node:20-alpine
docker pull docker.1panel.live/library/nginx:alpine
docker tag docker.1panel.live/library/nginx:alpine nginx:alpine
# MySQL 用 daocloud 镜像（1panel 的 MySQL 层较慢）
docker pull docker.m.daocloud.io/library/mysql:8.0.36
docker tag docker.m.daocloud.io/library/mysql:8.0.36 mysql:8.0.36
```

也可以在 Docker Desktop 的 Docker Engine 配置里加 `registry-mirrors`，例如 `https://docker.1panel.live`、`https://docker.m.daocloud.io`，然后重启 Docker Desktop。

Go 服务的 Dockerfile 已把 Alpine 软件源切到清华 TUNA；如果 TUNA 也不通，可以换成阿里云或中科大源。

访问：

- 商城前台：<http://localhost:8081/mall>
- 管理后台：<http://localhost:8080/dashboard>
- 商品接口：<http://localhost:8002/products>
- 网关：<http://localhost:8888>

测试账号：

```text
管理员：admin / 123456
```

## 本地 Kubernetes 环境

用 k3d 启动 3 节点集群（1 个控制面 + 2 个工作节点），用于本地验证 k8s 部署；商城日常开发仍用 Docker Compose。两者使用相同的 host 端口，切换前必须先停止 Compose。

前置：Docker Desktop 已启动，安装 k3d：

```bash
brew install k3d
```

```bash
make k8s-up      # 创建集群并建立 mall 命名空间
make k8s-status  # 查看集群和节点
make k8s-down    # 删除集群
```

- 集群名：`mall`，配置见 `deploy/k8s/k3d-cluster.yaml`
- 镜像代理：`deploy/k8s/registries.yaml`，docker.io / ghcr.io 走国内镜像站
- 默认禁用 Traefik，避免首次启动卡在 Ingress 镜像拉取；需要 Ingress 时再单独安装
- 首次启动 CoreDNS、metrics-server 等系统镜像需要拉取，可能耗时 5–10 分钟
- host 端口：8080/8081（前端）、8000–8002（API）、8888（网关）；8090/8443 预留给后续 Ingress
- k3s 版本由 k3d 默认值决定，升级 k3d 后重建集群即可

### 部署应用到 k8s

```bash
make k8s-images            # 构建 8 个镜像并导入 3 个节点
make k8s-deploy            # 空库部署，MySQL 首次启动自动执行 init_*.sql
make k8s-deploy RESTORE=1  # 迁移 Compose 现有数据（目标库已有业务数据时需 FORCE=1）
```

也可以分开执行：

```bash
make k8s-deploy
make k8s-restore-db DB_BACKUP=backups/mysql-xxxx.sql   # 覆盖已有数据时加 FORCE=1
make k8s-backup                                         # 导出 k8s 中的数据
```

`k8s-deploy` 会重启所有 Deployment/MySQL，确保新镜像、ConfigMap 和 Secret 生效。

### 高可用拓扑

- gateway、frontend、user/order/product 的 api 和 rpc：**3 副本**，硬拓扑分散为每节点 1 个
- 滚动更新 `maxSurge: 0` / `maxUnavailable: 1`，PodDisruptionBudget `minAvailable: 2`
- product-api：**1 副本**（uploads 是 local-path RWO 卷）
- MySQL / Redis / etcd：**1 副本**，仍是单点
- 决策记录见 `docs/adr/0002-k8s-ha-topology.md`

切换前先停止 Compose 释放 host 端口（数据卷保留）：

```bash
docker compose stop
```

访问地址与 Compose 一致：`localhost:8080` 管理后台、`localhost:8081` 商城、`localhost:8000-8002` API、`localhost:8888` 网关。

回滚前先备份 k8s 侧数据；`k8s-down` 会删除 k8s 的 PVC：

```bash
make k8s-backup      # 导出 k8s MySQL 到 backups/
make k8s-down        # 删除 k8s 集群，释放端口
docker compose start # 恢复 Compose，Compose 数据卷未改动
```

注意：

- MySQL root 密码只在空数据目录初始化时生效；修改 `.env` 后需在 MySQL 内执行 `ALTER USER`，不能只重跑 `k8s-secrets`。
- MySQL 和 uploads 使用 local-path RWO 卷，Pod 绑定在单个节点；节点丢失或跨节点迁移需要手工处理 PV。

## 主要流程

```text
注册/登录 → 浏览商品 → 加入购物车 → 选择地址 → 提交订单
        → 模拟支付 → 管理员发货 → 用户确认收货
```

## 本期已补充

- 售后/退款：`POST /after-sales`、`GET /after-sales`、管理员审核/退款
- SKU：`GET /product/:id/skus`、`POST /admin/skus`、`DELETE /admin/skus/:id`
- 多图上传：`POST /admin/uploads/batch`、`DELETE /admin/uploads/:name`
- 搜索：MySQL 全文索引、`GET /products/suggest`、`GET /products/hot`
- 行政区划：`GET /user/regions?parent_id=0`
- 账号安全：`PUT /user/password`、`POST /user/logout`、登录失败限流
- 订单：`GET /order/:id/tracks`、发货 7 天后自动确认收货
- 管理端：角色 1 超级管理员 / 角色 2 运营；`admin_audit_log` 审计
- 分类：7 个一级分类 + 19 个二级分类；商城首页悬停展示子分类，后台分类树形管理
- 商品管理分类筛选：级联选择器，支持选一级或二级分类
- 后台页面：售后管理、操作审计、SKU / 多图上传、用户角色、修改密码、物流轨迹
- 显示规则：后台搜索按名称、单号、状态，商品/订单/用户/审计表格不直接显示原始 ID

## 生成 Proto

```bash
brew install protobuf
GOBIN=/tmp/protogen-build GOPROXY=https://goproxy.cn,direct \
  go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.6
GOBIN=/tmp/protogen-build GOPROXY=https://goproxy.cn,direct \
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
./scripts/protogen.sh
```

生成的 `*_grpc.pb.go` 只包含 gRPC stubs；新增 RPC 后仍需手写更新 `app/<service>/rpc/<service>client` 和 `internal/server` 的方法绑定。

## API 示例

```bash
# 注册
curl -s -X POST localhost:8000/user/register \
  -H 'Content-Type: application/json' \
  -d '{"username":"buyer","password":"123456"}'

# 登录
TOKEN=$(curl -s -X POST localhost:8000/user/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"buyer","password":"123456"}' | jq -r .data.token)

# 商品列表
curl -s 'localhost:8002/products?page=1&page_size=10'

# 加入购物车
curl -s -X POST localhost:8001/cart/items \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"product_id":1,"quantity":1}'

# 创建订单（address_id 通过 GET /addresses 获取）
curl -s -X POST localhost:8001/orders \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"address_id":1,"remark":"尽快发货","buy_now":0}'

# 模拟支付
curl -s -X POST localhost:8001/order/1/pay \
  -H "Authorization: Bearer $TOKEN"
```

管理端接口需要 `role=1` 的 token：

```bash
# 发货
curl -s -X POST localhost:8001/admin/order/1/ship \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"logistics_company":"顺丰","tracking_no":"SF123456"}'

# 看板
curl -s localhost:8001/admin/dashboard \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

## 前端页面

**商城 C 端**

```text
/mall                 首页
/mall/products        商品列表
/mall/product/:id     商品详情
/mall/cart            购物车
/mall/checkout        确认订单
/mall/orders          我的订单
/mall/order/:id       订单详情
/mall/addresses       收货地址
```

**管理后台**

```text
/dashboard            数据看板
/category             分类管理
/product              商品管理
/product/create       新增商品
/product/:id/edit     编辑商品
/order                订单管理/发货
/user                 用户管理
```

## 目录结构

```text
app/
├── gateway/          # 网关
├── user/             # 用户 api + rpc + model
├── order/            # 订单 api + rpc + model
├── product/          # 商品 api + rpc + model
└── frontend/         # Vue3 商城 + 管理后台
pkg/
├── auth/             # JWT、身份上下文
└── errcode/          # 统一错误码
deploy/
├── sql/              # 初始化 SQL + 迁移
└── ...
```

## 开发规范

- 流程：[docs/DEVELOPMENT_PROCESS.md](docs/DEVELOPMENT_PROCESS.md)
- 执行清单：`.claude/skills/dev-process/SKILL.md`
- 完成标准：[.github/pull_request_template.md](.github/pull_request_template.md)
- 决策记录：[docs/adr/](docs/adr/)

```bash
make check / gen-check / smoke / db-backup
```

模块规范见 `CLAUDE.md` 的“技术规范索引”。

## 统一错误码

| 错误码 | 含义 |
|---|---|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未授权 |
| 403 | 无权限 |
| 404 | 资源不存在 |
| 409 | 资源已存在 |
| 500 | 服务器内部错误 |
