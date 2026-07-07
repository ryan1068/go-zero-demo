# go-zero-demo

go-zero 微服务项目，包含 user（用户）、order（订单）服务和 gateway（网关），前端 Vue3 + Element Plus。
开发流程（brainstorming → planning → TDD → review → finish）由 **Superpowers** 管理，本项目 CLAUDE.md 仅提供技术规范。

## 架构

```
app/user/   → 用户服务 (API + RPC)
app/order/  → 订单服务 (API + RPC)
app/gateway → 统一网关 (鉴权/日志)
app/frontend → Vue3 前端
deploy/     → 部署脚本 + 数据库模型
```

## 常用命令

```bash
# 生成 API 代码
goctl api go -api xxx.api -dir .

# 运行服务
go run app/gateway/gateway.go
go run app/user/api/user.go
go run app/order/api/order.go

# 测试
go test ./...

# 前端
cd app/frontend && npm run dev
```

## 项目铁律（Superpowers code-review 会据此检查）

- 所有 error 必须处理，不能 `_` 忽略
- 不要在生成的 handler 中写业务逻辑，逻辑一律放 `internal/logic/`
- 用户输入必须校验
- 密码/密钥不能硬编码或打印到日志
- 数据库查询必须带索引条件，禁止 `SELECT *`
- Go 命名遵循官方规范，Vue 组件用 PascalCase

## Git Commit 规范

```
<type>(<scope>): <subject>
```

| type | 用途 |
|---|---|
| feat | 新功能 |
| fix | 修复 bug |
| refactor | 重构 |
| docs | 文档 |
| test | 测试 |
| chore | 构建/工具 |

示例：`feat(order): 新增订单创建接口` / `fix(user): 修复登录 token 过期`

分支命名：`feat/<功能>`、`fix/<问题>`、`refactor/<模块>`

## 技术规范索引

操作不同模块时，子目录 CLAUDE.md 会自动加载：

| 操作类型 | 规范文件 | 加载条件 |
|---|---|---|
| Go 后端 | `app/CLAUDE.md` | `cd app/` 时自动加载 |
| Vue 前端 | `app/frontend/CLAUDE.md` | `cd app/frontend/` 时自动加载 |
| 数据库 | `deploy/script/mysql/CLAUDE.md` | `cd deploy/script/mysql/` 时自动加载 |
| 文档 | `docs/CLAUDE.md` | `cd docs/` 时自动加载 |