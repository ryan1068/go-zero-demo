# 开发流程

分级 → 澄清/计划 → 测试先行 → 小步实现 → review + 验收 + 收尾。执行清单见 `.claude/skills/dev-process/SKILL.md`；决策记录见 `docs/adr/`。

## 分级

| 级别 | 判断 | 走法 |
|---|---|---|
| L1 | 1-3 个文件，无接口/数据变化 | 直接改 + 相关测试 + `make check` |
| L2 | 多文件、目标明确 | 一页计划 → 测试 → 分步实现 → review + `make smoke` |
| L3 | 重构、跨模块、需求模糊、数据迁移 | spec → 计划 → TDD → 分阶段 → review + finish |

## 流程

1. **分级与计划**：明确目标、非目标、验收标准。L2 写一页计划；L3 先写 spec 到 `docs/superpowers/specs/YYYY-MM-DD-*.md`（骨架见下），再拆 3-5 步计划。
2. **测试先行**：重构先写 characterization tests；新功能先写失败测试；数据层补集成测试。
3. **小步实现**：每步 `make check` → 提交 → 汇报；高风险步骤停下。
4. **Review 与验收**：独立上下文 review diff；`make gen-check`、`make smoke`；写操作手工验证；数据变更前 `make db-backup`；清理测试数据。
5. **收尾**：按 PR 模板检查，更新文档，写回滚方式，开 PR。

## 停点

1. spec/计划批准
2. 测试先红
3. 最终验收

## Spec 骨架

- 背景与目标 / 非目标
- 方案与取舍
- 接口与数据变化
- 错误处理与安全
- 测试策略
- 风险与回滚
- 验收标准

## 命令

```bash
make check       # fmt-check + build + test + lint
make gen-check   # 生成代码漂移检查
make smoke       # 只读冒烟
make db-backup   # 备份数据库（需要 MYSQL_ROOT_PASSWORD）
```

## 硬性规则

- Go 代码遵循 `app/CLAUDE.md`；不改生成文件；密码/密钥不入仓库和日志
- 发现计划外问题先记录，不擅自扩大范围
