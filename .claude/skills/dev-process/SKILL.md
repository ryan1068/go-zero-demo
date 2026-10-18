---
name: dev-process
description: 开始 L2/L3 开发任务、重构、跨模块改动或需求不明确的功能开发时使用；强制执行 spec → 计划 → 测试先行 → 分步实现 → 独立 review → smoke → finish。Use for non-trivial features or refactors in this repository.
---

# 开发流程

先读 `docs/DEVELOPMENT_PROCESS.md`；这里只是执行清单。

1. **分级**：L1 直接做；L2 计划 + 测试；L3 先写 spec。开始前说明级别。
2. **澄清**：目标、非目标、约束、验收标准。
3. **Spec（L3）**：写到 `docs/superpowers/specs/`，等用户批准。
4. **计划**：3-5 步，每步写文件、测试、验证命令、回滚，等用户确认。
5. **测试先行**：重构写 characterization tests；新功能先红后绿。
6. **小步实现**：每步 `make check` → 提交 → 停下。
7. **独立 review**：用新上下文 review diff，重点行为、安全、迁移、并发、测试有效性。
8. **验收**：`make gen-check`、`make smoke`、数据变更前 `make db-backup`、清理测试数据。
9. **收尾**：按 PR 模板检查，更新文档，写回滚方式。

停点：spec、计划、测试、每步、review、验收。L1 可跳过第 3、4 步。

硬性规则：Go 代码按 go-zero 标准，go-zero 未提供才自定义（见 `app/CLAUDE.md`）；不改生成文件；数据层只用 model + `sqlx`；密码/密钥不入仓库和日志。
