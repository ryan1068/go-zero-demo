---
name: dev-process
description: 开始 L2/L3 开发任务、重构、跨模块改动或需求不明确的功能开发时使用；强制走 spec/计划、测试先行、小步实现、独立 review、smoke。Use for non-trivial features or refactors in this repository.
---

# 开发流程

先读 `docs/DEVELOPMENT_PROCESS.md`；这里只是清单。

1. **分级**：L1 直接做；L2 一页计划；L3 spec + 3-5 步计划。开始前说明级别。
2. **测试先行**：重构写 characterization tests；新功能先红后绿。
3. **小步实现**：每步 `make check` → 提交 → 汇报；高风险步骤停下。
4. **Review + 验收**：独立上下文 review；`make gen-check` / `make smoke`；数据变更前备份，验收后清理。
5. **收尾**：按 PR 模板检查，更新文档，写回滚。

停点：spec/计划批准、测试先红、最终验收。L1 可跳过第 1 步的计划部分。

硬性规则：Go 代码遵循 `app/CLAUDE.md`；不改生成文件；密码/密钥不入仓库和日志。
