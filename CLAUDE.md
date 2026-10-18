# go-zero-demo

go-zero 微服务商城：user（用户）、product（商品）、order（订单）、gateway（网关），前端 Vue3 + Element Plus。
开发流程见 `docs/DEVELOPMENT_PROCESS.md`，执行清单见 `.claude/skills/dev-process/SKILL.md`。本文件只写技术约束。

## 开发流程（摘要）

- L1 直接改；L2 走计划 + 测试 + 自查；L3 先 spec。
- 每阶段 `make check`；完成前 L3/高风险做独立 review + `make smoke`。
- 关键决策写入 `docs/adr/`。

## 架构

```
app/user/     → 用户服务 (API + RPC)
app/product/  → 商品服务 (API + RPC)
app/order/    → 订单服务 (API + RPC)
app/gateway/  → 统一网关 (鉴权/日志)
app/frontend/ → Vue3 商城 + 管理后台
deploy/       → 部署脚本 + 数据库模型
```

## 常用命令

- 服务：`make run-*`；测试：`make test`；全量校验：`make check`
- 生成：`make gen-*-api` / `make gen-model`；前端：`cd app/frontend && npm run dev`
- 其余见 `Makefile` 与 `README.md`

## 项目铁律（Superpowers code-review 会据此检查）

- Go 代码一律按 go-zero 标准：优先使用 go-zero/goctl 提供的能力；go-zero 未提供时才自定义（详见 `app/CLAUDE.md`）
- 所有 error 必须处理，不能 `_` 忽略
- 不要在生成的 handler 中写业务逻辑，逻辑一律放 `internal/logic/`
- 用户输入必须校验
- 密码/密钥不能硬编码或打印到日志
- 数据库查询必须带索引条件，禁止 `SELECT *`
- Go 命名遵循官方规范，Vue 组件用 PascalCase

## Git Commit 规范

`<type>(<scope>): <subject>`，type 用 `feat|fix|refactor|docs|test|chore`。

示例：`feat(order): 新增订单创建接口` / `fix(user): 修复登录 token 过期`

分支命名：`feat/<功能>`、`fix/<问题>`、`refactor/<模块>`

## 技术规范

子目录 `CLAUDE.md` 按工作目录自动加载：Go 后端 `app/`、Vue 前端 `app/frontend/`、数据库 `deploy/script/mysql/`、文档 `docs/`。