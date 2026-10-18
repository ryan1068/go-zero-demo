# 开发流程

一句话：分级 → 澄清 → spec（L3）→ 计划 → 测试先行 → 小步实现 → review + 验收 → 收尾。

执行清单见 `.claude/skills/dev-process/SKILL.md`；决策记录见 `docs/adr/`。

## 分级

| 级别 | 判断 | 走法 |
|---|---|---|
| L1 | 1-3 个文件，无接口/数据变化 | 直接改 + 相关测试 + `make check` |
| L2 | 多文件、目标明确的中等改动 | 计划 → 测试 → 分步实现 → review + `make smoke` |
| L3 | 重构、跨模块、需求模糊、数据迁移 | spec → 计划 → TDD → 分阶段 → review + finish |

开始前先说明本次属于哪一级。

## 流程

1. **澄清**：目标、非目标、约束、验收标准。只问会改变实现的问题。
2. **Spec（L3）**：写到 `docs/superpowers/specs/YYYY-MM-DD-<topic>-design.md`，用下面的骨架；关键决策写 `docs/adr/`。
3. **计划**：3-5 步，每步写文件、测试、验证命令、回滚；一步一个提交。
4. **测试先行**：重构先写 characterization tests；新功能先写失败测试；数据层补集成测试。
5. **小步实现**：每步 `make check` → 提交 → 停下。
6. **Review**：独立上下文 review diff，重点看行为、安全、迁移、并发、测试有效性。
7. **验收**：`make gen-check`、`make smoke`；写操作手工验证；数据变更前 `make db-backup`；清理测试数据。
8. **收尾**：按 PR 模板检查，更新文档，写风险和回滚，开 PR。

停点：spec 批准、计划确认、测试确认、每步结束、review 结论、验收结果。

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

- Go 代码一律按 go-zero 标准：优先使用 go-zero/goctl 提供的能力（model/sqlx、TransactCtx、middleware、Prometheus 等）；go-zero 未提供时才自定义，详见 `app/CLAUDE.md`
- 不手改 `*_gen.go`、`routes.go`、`types.go`
- 数据层只用 goctl model + `sqlx`，不直接使用 `database/sql`
- 密码、密钥、token 不入仓库、不进日志
- 发现计划外问题先记录，不擅自扩大范围
