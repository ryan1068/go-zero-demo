# go-zero 后端开发规范

## 核心原则

**优先使用 go-zero/goctl 提供的能力；go-zero 未提供时才自定义，且优先复用 `sqlx`。**

- 数据访问：goctl 生成的 `internal/model` + `core/stores/sqlx`，禁止 `database/sql`
- 配置：`sqlx.SqlConf`、`rest.RestConf`、`zrpc.RpcServerConf`
- 事务：`sqlx.SqlConn.TransactCtx`；业务回滚返回 `model.ErrAbortTransaction`
- 指标：内置 `Prometheus`，不手写 promhttp
- 错误：用 `sqlx.ErrNotFound`，不用 `sql.ErrNoRows`
- go-zero 未提供的能力（软删除、全文检索、跨库查询、multipart 上传）：在 `internal/model` 或 `internal/logic` 用 `sqlx`/标准库实现

## 目录与生成文件

- 目录：`api/internal/{config,handler,logic,middleware,svc,types}`、`rpc/internal/{config,logic,model,server,svc}`
- handler、types、`*_gen.go`、`vars.go`、`*client` 由 goctl 生成，禁止手改
- 业务逻辑只写 `internal/logic/`，一文件一方法
- model 自定义查询/软删除/事务方法写同包 `*_custom.go`
- 重新生成用 `make gen-model` / `make gen-*-api`；当前不启用 model cache（order-rpc 跨库更新库存，缓存无法失效）

## API 与 Logic

- `.api` 是接口事实来源；改完必须 `goctl api go` 重新生成
- 用户输入必须带 `validate` 校验
- Logic 结构体包含 `logx.Logger`、`ctx`、`svcCtx`，依赖一律从 `l.svcCtx` 取
- 错误必须返回，不能吞掉；业务错误用 `BizError{Code,Msg}`，网关用 `httpx.SetErrorHandler` 统一映射，系统错误转 500
- 中间件中断请求用 `httpx.Error`，传数据用 `context.WithValue`
- 调用 RPC 必须传 `l.ctx` 保持链路追踪

## 测试

- 命名 `Test<Function>_<Scenario>`，优先表驱动
- 必须覆盖：正常路径、边界值、错误路径、空值/nil
- 核心 logic 覆盖率 > 80%
