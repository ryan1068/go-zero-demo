---
name: "go-zero-dev"
description: "go-zero microservice development expert. Invoke when creating or modifying go-zero services, APIs, RPCs, middleware, gateway config, or .api files."
---

# go-zero 微服务开发规范

你是一个 go-zero 框架专家，专注于 go-zero 微服务项目的开发。以下规范适用于本项目 `go-zero-demo`。

---

## 一、项目目录结构

```
app/
├── <service>/
│   ├── api/                        # HTTP API 服务
│   │   ├── <service>.api           # API 定义文件 (goctl 入参)
│   │   ├── <service>.go            # 入口文件 (main)
│   │   └── internal/
│   │       ├── config/config.go    # 配置结构体
│   │       ├── handler/            # HTTP handler (自动生成，勿手动改)
│   │       ├── logic/              # 业务逻辑层 (核心开发)
│   │       ├── svc/servicecontext.go  # 服务上下文 (依赖注入)
│   │       └── types/types.go      # 请求/响应类型 (自动生成)
│   └── rpc/                        # gRPC 服务
│       ├── pb/<service>/           # protobuf 生成代码
│       ├── <service>.go            # 入口文件
│       ├── <service>client/        # RPC 客户端
│       └── internal/
│           ├── config/config.go
│           ├── logic/              # RPC 业务逻辑
│           ├── server/             # gRPC server 注册
│           └── svc/servicecontext.go
└── gateway/                        # 统一网关
    ├── gateway.go
    └── internal/
        ├── config/config.go
        ├── middleware/             # 全局中间件
        ├── svc/servicecontext.go
        └── types/                  # 自定义错误类型
```

---

## 二、API 定义文件 (.api) 规范

```go
syntax = "v1"

// 请求类型：使用 json/path/form tag 标记字段来源
type CreateOrderReq {
    GoodsId   int32  `json:"goods_id"`
    GoodsName string `json:"goods_name"`
    Price     int32  `json:"price"`
}

// 统一响应结构
type CommonResp {
    Code int32  `json:"code"`
    Msg  string `json:"msg"`
    Data string `json:"data"`
}

// 服务定义
service order {
    @handler GetOrderHandler
    get /order/:id (GetOrderReq) returns (CommonResp)

    @handler CreateOrderHandler
    post /order/create (CreateOrderReq) returns (CommonResp)
}
```

**规则：**
- 使用 `syntax = "v1"`
- 每个 handler 用 `@handler` 注解，命名规范：`{Action}{Resource}Handler`
- 路径参数用 `path:"name"` tag
- JSON 请求体用 `json:"name"` tag
- 修改 `.api` 文件后，运行 `goctl api go` 重新生成代码

---

## 三、Logic 层开发规范 (核心)

Logic 层是业务逻辑的唯一位置，**不要在 handler 中写业务逻辑**。

```go
package logic

import (
    "context"
    "go-zero-demo/order-api/internal/svc"
    "go-zero-demo/order-api/internal/types"
    "github.com/zeromicro/go-zero/core/logx"
)

type CreateOrderLogic struct {
    logx.Logger
    ctx    context.Context
    svcCtx *svc.ServiceContext
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
    return &CreateOrderLogic{
        Logger: logx.WithContext(ctx),
        ctx:    ctx,
        svcCtx: svcCtx,
    }
}

func (l *CreateOrderLogic) CreateOrder(req *types.CreateOrderReq) (resp *types.CommonResp, err error) {
    // 1. 参数校验
    // 2. 调用 RPC / 数据库操作
    result, err := l.svcCtx.OrderRpc.CreateOrder(l.ctx, &order.CreateOrderReq{...})
    if err != nil {
        return nil, err
    }
    // 3. 组装响应
    return &types.CommonResp{Code: 0, Msg: "success"}, nil
}
```

**规则：**
- 结构体必须包含 `logx.Logger`、`ctx context.Context`、`svcCtx *svc.ServiceContext`
- 使用构造函数模式 `NewXxxLogic`
- **错误必须处理**，不能忽略 `err`
- 通过 `l.svcCtx` 访问所有依赖（RPC 客户端、DB、Redis 等）
- 使用 `logx.WithContext(ctx)` 获取带 trace 的 logger

---

## 四、ServiceContext 依赖注入

```go
package svc

import (
    "github.com/zeromicro/go-zero/zrpc"
    "go-zero-demo/order-api/internal/config"
    "go-zero-demo/order-rpc/orderclient"
)

type ServiceContext struct {
    Config   config.Config
    OrderRpc orderclient.Order    // RPC 客户端
    // DB    *sql.DB               // 数据库连接
    // Redis *redis.Redis          // 缓存
}

func NewServiceContext(c config.Config) *ServiceContext {
    return &ServiceContext{
        Config:   c,
        OrderRpc: orderclient.NewOrder(zrpc.MustNewClient(c.OrderRpcConf)),
    }
}
```

**规则：**
- 所有依赖（RPC 客户端、DB、Redis）在 `ServiceContext` 中统一管理
- 在 `NewServiceContext` 中完成初始化
- Config 配置结构体通过 `github.com/zeromicro/go-zero/core/conf` 加载

---

## 五、中间件规范

```go
package middleware

import (
    "net/http"
    "github.com/zeromicro/go-zero/rest/httpx"
)

type AuthMiddleware struct{}

func NewAuthMiddleware() *AuthMiddleware {
    return &AuthMiddleware{}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // 1. 前置处理（鉴权、参数校验）
        token := r.Header.Get("Authorization")
        if token == "" {
            httpx.Error(w, errors.New("unauthorized"))
            return
        }

        // 2. 通过 context 传递数据
        ctx := context.WithValue(r.Context(), "user_id", userId)
        r = r.WithContext(ctx)

        // 3. 调用下一个处理器
        next(w, r)
    }
}
```

**规则：**
- 中间件签名：`func (m *M) Handle(next http.HandlerFunc) http.HandlerFunc`
- 中断请求用 `httpx.Error(w, err)` 返回错误
- 传递数据通过 `context.WithValue`，不要用全局变量
- 在网关入口 `gw.Use()` 注册全局中间件

---

## 六、错误处理规范

```go
// 自定义业务错误类型
type BizError struct {
    Code int    `json:"code"`
    Msg  string `json:"msg"`
}

func (e *BizError) Error() string {
    return e.Msg
}

func NewBizError(code int, msg string) *BizError {
    return &BizError{Code: code, Msg: msg}
}
```

**规则：**
- 定义 `BizError` 类型统一业务错误
- 在网关注册自定义错误处理器 `httpx.SetErrorHandler`
- 业务错误返回非零 Code，系统错误返回 500
- RPC 调用失败必须返回 error，不要吞掉

---

## 七、RPC 调用规范

```go
// 在 logic 中调用 RPC
result, err := l.svcCtx.OrderRpc.CreateOrder(l.ctx, &order.CreateOrderReq{
    GoodsId:   req.GoodsId,
    GoodsName: req.GoodsName,
    Price:     req.Price,
})
if err != nil {
    return nil, err
}
```

**规则：**
- RPC 客户端通过 `ServiceContext` 注入
- 使用 `zrpc.MustNewClient` 创建客户端连接
- 传递 `l.ctx` 保持链路追踪
- 不要在 logic 中直接 new RPC 客户端

---

## 八、配置规范

```go
package config

import (
    "github.com/zeromicro/go-zero/rest"
    "github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
    rest.RestConf                    // HTTP 服务配置
    OrderRpcConf zrpc.RpcClientConf  // RPC 客户端配置
    // DB  DataSource                  // 数据库配置
}
```

**规则：**
- 嵌入 go-zero 内置配置结构体（`rest.RestConf`、`zrpc.RpcClientConf`）
- 配置通过 `conf.MustLoad(*configFile, &c)` 加载
- 敏感信息通过环境变量或配置中心管理

---

## 九、命名规范

| 元素 | 规范 | 示例 |
|------|------|------|
| 服务名 | 小写短横线 | `order-api`、`user-rpc` |
| Handler | `{Action}{Resource}Handler` | `CreateOrderHandler` |
| Logic | `{Action}{Resource}Logic` | `CreateOrderLogic` |
| API 文件 | 服务名.api | `order.api` |
| 目录 | 小写单数 | `handler/`、`logic/`、`types/` |

---

## 十、开发流程

1. **定义 API**：在 `.api` 文件中定义接口和类型
2. **生成代码**：运行 `goctl api go -api xxx.api -dir .`
3. **编写 Logic**：在 `internal/logic/` 中实现业务逻辑
4. **注册依赖**：在 `ServiceContext` 中添加需要的依赖
5. **配置中间件**：在 `gateway/internal/middleware/` 中编写
6. **测试**：编写单元测试，确保 error 路径覆盖

---

## 注意事项

- **不要在生成的 handler 中写业务逻辑**，handler 只负责参数绑定和调用 logic
- **所有 error 必须处理**，不允许 `_` 忽略
- **使用 `logx.WithContext(ctx)` 获取 logger**，自动携带 trace_id
- **RPC 调用必须传递 ctx**，保持全链路追踪
- **修改 .api 后重新生成代码**，不要手动修改 types.go 和 handler