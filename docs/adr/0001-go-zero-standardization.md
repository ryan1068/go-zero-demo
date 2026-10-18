# ADR 0001：数据层迁移到 go-zero 标准模式

- 状态：Accepted
- 日期：2026-10-18

## 背景

项目最初只使用 go-zero 的框架层，数据访问是原生 `database/sql` + 手写 SQL，没有 goctl model，也没有 sqlx。handler 中存在业务逻辑，管理员守卫重复三份，product-api 完全手工维护。

用户要求：以后统一按 go-zero 标准实现，go-zero 没提供的能力才自行实现。

## 决策

1. **三个 RPC 服务统一使用 goctl 生成的 `internal/model` + `sqlx`**
   - 不允许在 logic 中直接使用 `database/sql`
   - 生成文件 `*_gen.go` 不手改
   - 自定义方法放同包 `*_custom.go`

2. **不启用 model cache**
   - order-rpc 会跨库直接更新 `looklook_product.product`
   - product-rpc 的缓存无法感知这些写入，启用会产生脏读
   - 等跨库写入改为 RPC + 最终一致后再评估

3. **保留 order-rpc 的跨库 SQL**
   - 库存扣减与订单创建必须在同一 MySQL 事务中
   - go-zero 不提供跨服务事务
   - 跨库语句封装在 order-rpc 的 model 自定义方法中

4. **软删除、全文检索、分页筛选、跨库查询手写**
   - goctl 生成 model 不识别 `deleted_at`，`Delete` 是物理删除
   - 这些能力在 `*_custom.go` 中用 `sqlx` 实现

5. **JWT 保留 `pkg/auth` 自定义实现**
   - go-zero 内置 JWT 不支持角色和 token 黑名单
   - 当前实现是标准形态的 go-zero middleware

6. **multipart 上传和静态文件响应保留手写 handler**
   - goctl `.api` 不支持文件类型
   - 路由仍由 goctl 生成，只有 handler 实现是手写

7. **API 层按 goctl 重新生成**
   - routes/types/handlers 由 goctl 生成
   - logic 一文件一方法
   - 管理员守卫改为 `internal/middleware`
   - 文件命名使用 goctl 默认小写风格，保证重复生成幂等

8. **指标使用 go-zero 内置 Prometheus**
   - 删除自建 `promhttp` goroutine
   - 端口：user-rpc 9100、user-api 9101、order-rpc 9102、order-api 9103、product-api 9104、product-rpc 9105

## 备选方案

| 方案 | 取舍 |
|---|---|
| 继续使用 `database/sql` | 改动小，但不符合 go-zero 标准，缺少 sqlx 的追踪/错误语义 |
| 使用 goctl model + 缓存 | 读性能更好，但跨库写入会导致缓存不一致 |
| 跨库访问改为 RPC 调用 | 服务边界更干净，但需要补偿/最终一致，改动和风险大 |

## 影响

正面：

- 数据访问统一，逻辑层不再关心连接和 SQL 细节
- 事务、错误映射、可观测性统一
- 代码生成可重复，降低手写漂移

负面/成本：

- 生成文件与自定义文件边界需要维护
- 跨库写入仍是架构债务
- 没有缓存，读性能依赖数据库

## 后续

- 补 order-rpc / product-rpc 的真实数据库集成测试
- 把跨库写入改为 product-rpc 调用 + outbox/补偿后再考虑启用缓存
- 配置中的明文密码迁移到环境变量
- 登录请求日志脱敏
