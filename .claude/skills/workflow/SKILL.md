---
name: workflow
description: 开始 L2/L3 开发任务、重构、跨模块改动或需求不明确的功能开发时使用；强制走 spec/计划、测试先行、小步实现、按级别 review；环境可用时默认 smoke。Use for non-trivial features or refactors in this repository.
---

# 开发流程

先读 `docs/DEVELOPMENT_PROCESS.md`；这里只是清单。

1. **分级**：L1 直接做；L2 一页计划；L3 spec + 3-5 步计划。开始前说明级别。
2. **测试先行**：重构写 characterization tests；新功能先红后绿。
3. **小步实现**：迭代中跑改动模块的相关测试，提交前 `make check`；高风险步骤停下。
4. **Review + 验收**：L2 自查 + `make check`；环境可用时加 `make smoke`；L3 或高风险改动才做独立上下文 review，范围限定本次 diff；生成代码改动加 `make gen-check`；部署/配置/路由改动或发布前必须跑 `make smoke`；数据变更前备份，验收后清理。
5. **收尾**：按 PR 模板检查，更新文档，写回滚。

停点：spec/计划批准、测试先红、最终验收。L1 可跳过第 1 步的计划部分。

验证分层：迭代中只跑相关测试；提交前 `make check`；环境已启动且部署当前代码时默认 smoke；环境未启动或未部署时不要用 smoke 代替 check。

硬性规则：Go 代码遵循 `app/CLAUDE.md`；不改生成文件；密码/密钥不入仓库和日志。
