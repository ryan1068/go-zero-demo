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

## 快速开始

```bash
# 1. 准备环境变量（JWT_SECRET 至少 32 个字符）
cp .env.example .env

# 2. 启动所有服务
docker compose up -d --build

# 3. 已有 MySQL 数据卷时，依次执行迁移（只需执行一次）
#    先把 docker-compose.yml 中 mysql 服务的 root 密码导出为环境变量
export MYSQL_ROOT_PASSWORD='<docker-compose.yml 中 mysql 的 root 密码>'
for f in deploy/sql/migrations/20261018_*.sql; do
  docker exec -i mysql mysql -uroot -p"$MYSQL_ROOT_PASSWORD" < "$f"
done
```

> 全新的 MySQL 数据卷会自动执行 `init_user.sql`、`init_order.sql`、`init_product.sql`。
> 已有数据卷不会重复执行初始化 SQL，需要手动跑迁移。
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
├── user/             # 用户 api + rpc
├── order/            # 订单 api + rpc
├── product/          # 商品 api + rpc
└── frontend/         # Vue3 商城 + 管理后台
pkg/
├── auth/             # JWT、身份上下文
└── errcode/          # 统一错误码
deploy/
├── sql/              # 初始化 SQL + 迁移
└── ...
```

## 开发规范

本项目使用 **Superpowers** 管理开发流程，**CLAUDE.md** 管理技术规范。

| 规范 | 文件 |
|---|---|
| 项目铁律 + Git | `CLAUDE.md` |
| Go 后端 | `app/CLAUDE.md` |
| Vue 前端 | `app/frontend/CLAUDE.md` |
| 数据库 | `deploy/script/mysql/CLAUDE.md` |
| 文档 | `docs/CLAUDE.md` |

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
