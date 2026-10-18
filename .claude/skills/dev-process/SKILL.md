---
name: dev-process
description: 开始 L2/L3 开发任务、重构、跨模块改动或需求不明确的功能开发时使用；强制执行 spec → 计划 → 测试先行 → 分步实现 → 独立 review → smoke → finish。Use for non-trivial features or refactors in this repository.
---

# 开发流程

先读 `docs/DEVELOPMENT_PROCESS.md`；本 skill 是执行清单，不是文档替代品。

## 0. 先分级

- **L1**：1-3 个文件，无接口/数据变化 → 直接做，跑相关测试和 `make check`，提交。
- **L2**：明确的中等改动 → 至少一页计划，测试，分步实现，独立 review，`make smoke`。
- **L3**：重构、跨模块、需求模糊、涉及数据迁移 → 必须走完 spec、计划、TDD、review、finish。

开始前向用户说明本次属于哪一级。L1 可以跳过第 2、3 步。

## 1. 澄清

- 明确目标、非目标、约束、验收标准。
- 一次只问一个会改变实现的问题。
- 业务规则不允许猜；猜了要显式说明是假设。

## 2. Spec（L3）

- 使用 brainstorming 流程探索方案。
- 按 `docs/superpowers/specs/TEMPLATE.md` 写到
  `docs/superpowers/specs/YYYY-MM-DD-<topic>-design.md`。
- 关键决策同步记录到 `docs/adr/`。
- 把 spec 提交到仓库。

**停点 1：等用户批准 spec，才能进入计划。**

## 3. 计划

- 拆成 3-5 个可独立验证的步骤。
- 每一步写清楚：改动文件、测试、验证命令、回滚方式。
- 每一步不超过一个提交。
- 明确「不改什么」。

L2 可直接使用 plan mode；L3 使用 writing-plans 产出计划。

**停点 2：等用户确认范围和步骤，才能写代码。**

## 4. 测试先行

- 重构：先写 characterization tests，锁住现有行为。
- 新功能：先写失败测试，并贴出「红」的结果。
- 接口：补契约测试。
- 数据层：单测 + 真实 MySQL 集成测试。

**停点 3：等用户确认测试验证的是行为，不是实现细节。**

## 5. 小步实现

每一步：

1. 只做当前步骤。
2. 运行 `make check`。
3. 按 `<type>(<scope>): <subject>` 提交。
4. 汇报改了什么、验证结果、剩余步骤。

**停点 4：每一步结束后停下，等用户确认再继续。**

## 6. 独立 review

使用新的上下文（另一个 subagent 或人），不要自己批自己。提示词：

```text
独立 review 当前分支相对 main 的 diff，不要修改代码。
重点检查：行为变化、错误处理、权限与安全、数据迁移、并发与幂等、性能与索引、测试有效性。
按严重程度输出 findings，不要参考作者已有的结论。
```

**停点 5：由用户决定哪些 findings 必须修。**

## 7. 真实验收

- 改了 `.api` 或生成代码：运行 `make gen-check`。
- 服务已启动：运行 `make smoke`。
- 涉及数据变更：先 `make db-backup`。
- 写操作流程（下单、支付、售后、上传）按 `docs/DEFINITION_OF_DONE.md` 手工验收。
- 验收后清理测试数据，恢复库存/状态。

## 8. Finish

- 按 `docs/DEFINITION_OF_DONE.md` 逐项确认。
- 更新 README / spec / ADR。
- 写变更摘要、风险、回滚方式。
- 开 PR 或合并。

## 硬性规则

- 不手改 `*_gen.go`、`routes.go`、`types.go`。
- 数据层只用 goctl model + `sqlx`，不直接使用 `database/sql`。
- 密码、密钥、token 不入仓库、不进日志。
- 生成代码和手工修改的边界要写进提交说明。
- 每个停点都停下，不要连续跑完全部步骤。
- 发现计划外的问题，先记录，不擅自扩大范围。
