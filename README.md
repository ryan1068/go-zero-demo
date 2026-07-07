# go-zero-demo

go-zero 微服务示例项目，包含用户服务、订单服务、网关和前端。

## 架构

```
┌──────────┐   ┌──────────┐   ┌──────────┐
│ 用户 API  │   │ 订单 API  │   │  Gateway  │
│  :8080   │   │  :8081   │   │  :8888   │
└────┬─────┘   └────┬─────┘   └──────────┘
     │              │
┌────▼─────┐   ┌────▼─────┐
│ 用户 RPC  │   │ 订单 RPC  │
│  :9090   │   │  :9091   │
└────┬─────┘   └──────────┘
     │
┌────▼─────┐
│  MySQL   │
└──────────┘
```

| 服务 | 端口 | 说明 |
|---|---|---|
| gateway | 8888 | 统一网关（鉴权、日志） |
| user-api | 8080 | 用户 HTTP 服务 |
| user-rpc | 9090 | 用户 gRPC 服务 |
| order-api | 8081 | 订单 HTTP 服务 |
| order-rpc | 9091 | 订单 gRPC 服务 |
| frontend | 3000 | Vue3 + Element Plus 前端 |

## 快速开始

```bash
# 1. 启动基础设施（MySQL、Redis、Jaeger、Prometheus）
docker-compose up -d

# 2. 初始化数据库
mysql -u root < deploy/script/mysql/user.sql

# 3. 启动服务
make run-user-rpc &
make run-user-api &
make run-order-rpc &
make run-order-api &
make run-gateway &

# 4. 启动前端
make frontend-dev
```

## 常用命令

```bash
make build          # 编译所有服务
make test           # 运行所有测试
make lint           # 代码检查
make fmt            # 格式化代码

make run-gateway    # 启动网关
make run-user-api   # 启动用户 API
make run-user-rpc   # 启动用户 RPC
make run-order-api  # 启动订单 API
make run-order-rpc  # 启动订单 RPC

make frontend-dev   # 启动前端开发服务器
make frontend-build # 构建前端
```

## 目录结构

```
app/
├── gateway/          # 网关（鉴权、日志中间件）
├── user/
│   ├── api/          # 用户 HTTP 服务
│   │   ├── user.api  # API 定义
│   │   └── internal/
│   │       ├── handler/  # 自动生成，勿改
│   │       ├── logic/    # 业务逻辑
│   │       ├── svc/      # 依赖注入
│   │       └── config/   # 配置
│   └── rpc/          # 用户 gRPC 服务
│       ├── user.proto
│       └── internal/logic/
├── order/            # 订单服务（结构同上）
└── frontend/         # Vue3 前端
pkg/
└── errcode/          # 统一错误码
deploy/
└── script/mysql/     # 数据库脚本
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

所有错误码统一定义在 `pkg/errcode/`，禁止魔法数字。

| 错误码 | 含义 |
|---|---|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未授权 |
| 403 | 无权限 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |