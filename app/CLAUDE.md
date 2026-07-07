# go-zero 后端开发规范

## 目录结构

```
<service>/
├── api/                     # HTTP 服务
│   ├── <service>.api        # API 定义 (goctl 入参)
│   ├── <service>.go         # 入口 main
│   └── internal/
│       ├── config/config.go
│       ├── handler/          # 自动生成，勿手动改
│       ├── logic/            # 业务逻辑 (核心)
│       ├── svc/servicecontext.go
│       └── types/types.go    # 自动生成
└── rpc/                     # gRPC 服务
    └── internal/
        ├── logic/
        ├── server/
        └── svc/
```

## API 定义 (.api)

```go
syntax = "v1"

type CreateReq {
    Name  string `json:"name"`
    Price int32  `json:"price"`
}

type CommonResp {
    Code int32  `json:"code"`
    Msg  string `json:"msg"`
    Data string `json:"data"`
}

service xxx {
    @handler CreateHandler
    post /xxx/create (CreateReq) returns (CommonResp)
}
```

修改 .api 后必须运行 `goctl api go` 重新生成。

## Logic 层

```go
type CreateLogic struct {
    logx.Logger
    ctx    context.Context
    svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
    return &CreateLogic{
        Logger: logx.WithContext(ctx),
        ctx:    ctx,
        svcCtx: svcCtx,
    }
}

func (l *CreateLogic) Create(req *types.CreateReq) (*types.CommonResp, error) {
    // 1. 参数校验
    // 2. 调用 RPC / DB
    // 3. 组装响应
}
```

- 结构体必须包含 `logx.Logger`、`ctx`、`svcCtx`
- 用 `l.svcCtx` 访问所有依赖
- 错误必须返回，不能吞掉

## ServiceContext

```go
type ServiceContext struct {
    Config    config.Config
    OrderRpc  orderclient.Order
}

func NewServiceContext(c config.Config) *ServiceContext {
    return &ServiceContext{
        Config:   c,
        OrderRpc: orderclient.NewOrder(zrpc.MustNewClient(c.OrderRpcConf)),
    }
}
```

## 中间件

```go
func (m *M) Handle(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // 前置处理
        next(w, r)
    }
}
```

中断请求用 `httpx.Error(w, err)`，传数据用 `context.WithValue`。

## 错误处理

- 自定义 `BizError{Code, Msg}` 类型
- 网关注册 `httpx.SetErrorHandler` 统一处理
- 业务错误返回非零 Code，系统错误 500

## RPC 调用

```go
result, err := l.svcCtx.OrderRpc.CreateOrder(l.ctx, &order.CreateOrderReq{...})
if err != nil {
    return nil, err
}
```

必须传 `l.ctx` 保持链路追踪。

## 测试规范

```go
func TestCreateOrder_Success(t *testing.T) {
    // 正常路径
}

func TestCreateOrder_InvalidPrice(t *testing.T) {
    // 错误路径
}
```

- 命名：`Test<Function>_<Scenario>`
- 必须覆盖：正常路径、边界值、错误路径、空值/nil
- 核心 logic 覆盖率 > 80%
- 使用表驱动测试：

```go
func TestLogin(t *testing.T) {
    tests := []struct {
        name     string
        username string
        password string
        wantCode int32
    }{
        {"success", "admin", "123456", 0},
        {"wrong password", "admin", "wrong", 401},
        {"empty username", "", "123456", 401},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // ...
        })
    }
}
```