# Kubernetes 运行环境迁移

日期：2026-10-18
状态：已实现

## 背景与目标

当前运行环境是 Docker Compose：gateway、user/product/order 各 api+rpc、frontend、etcd、MySQL、Redis。本地已有 k3d `mall` 集群（1 控制面 + 2 工作节点）。

目标：把上述服务整体迁到 `mall` 命名空间运行，host 端口保持 8080 / 8081 / 8000 / 8001 / 8002 / 8888，前端 `portal.ts` 依赖 `window.location.port` 判断管理后台/商城，不能改端口。Compose 保留，作为回滚方案。

## 非目标

- 不接 CI/CD、Ingress、Prometheus/Grafana
- 不改业务代码、`.api`、proto、数据库结构
- 不迁移 Kafka/ES/Kibana/Filebeat/go-stash/Asynqmon 等未使用服务
- 不做多副本、高可用、HPA（MySQL/Redis/etcd 单副本）
- 不处理 HTTPS、对象存储

## 方案与取舍

| 项 | 方案 | 取舍 |
|---|---|---|
| 命名空间 | `mall` | 与现有 `mall` 命名空间一致 |
| 服务发现 | 继续用 etcd，RPC 注册 key 不变 | 配置无需改 |
| Service 名 | `etcd` / `mysql` / `redis` / `user-api` / `order-api` / `product-api` / `gateway` / `frontend` | 现有 `etc/*.yaml` 里的 DNS 名直接可用 |
| 配置 | 每个服务的 yaml 做 ConfigMap，挂到 `/app/etc` | 不重建镜像即可改配置 |
| 密钥 | `.env` 生成 Secret，注入 `JWT_SECRET` / `MYSQL_ROOT_PASSWORD` / `REDIS_PASSWORD` | 不入仓库、不打印 |
| 镜像 | 本地构建 7 个 Go 镜像 + frontend，arm64 单平台导出后 `ctr` 导入 3 个节点 | `k3d image import` 对多架构镜像报 digest 缺失，走已验证的 `docker save --platform` + `ctr images import` |
| MySQL | StatefulSet 1 副本 + local-path PVC 10Gi | 单点，本地开发可接受 |
| Redis | Deployment 1 副本 + emptyDir | 缓存/限流/黑名单数据可重建 |
| etcd | Deployment 1 副本 + emptyDir | 服务注册可重建 |
| uploads | product-api 挂 local-path PVC 2Gi | 单副本 RWO，避免 Pod 重建丢图 |
| 入口 | 不用 Ingress；frontend 两个 NodePort（30080/30081），API/gateway 用 NodePort | host 端口通过 k3d loadbalancer 映射到 NodePort，绕过 Traefik 慢镜像 |
| k3d 端口映射 | 增加 8080→30080、8081→30081、8000→30000、8001→30001、8002→30002、8888→30088 | 需要重建一次集群；重建前导出已缓存系统镜像 |
| 数据 | 空库由 MySQL 执行 `init_*.sql`；迁移时先 `make db-backup`，恢复脚本会预备份并要求 `FORCE=1` | 保留现有 10 万用户、26 商品、订单 |

## 接口与数据变化

- 无 API / proto / 表结构变化
- 数据从 Compose 卷一次性迁到 k8s PVC
- host 端口不变；切换时先 `docker compose stop` 释放端口，卷保留

## 错误处理与安全

- Secret 从 `.env` 生成；不在 manifest 中硬编码密码
- 所有 Pod 配 readiness/liveness TCP probe；失败时不接流量
- 空库由 `init_*.sql` 建表；恢复脚本覆盖已有数据前会预备份，恢复后校验计数
- 镜像 `FROM scratch`，probe 只能用 TCP，不能用 exec

## 测试策略

- 代码测试不变：`make check`
- 部署验收：
  - `kubectl -n mall get pods` 全部 Ready
  - `scripts/smoke.sh` 指向 NodePort（`USER_API=http://localhost:8000` 等）
  - 登录、创建订单、后台看板、商城首页
  - 对比 Compose 与 k8s 的 users / products / orders 计数
- 失败时保留现场，不删除 PVC

## 风险与回滚

- **系统镜像重复拉取**：重建集群会清空 containerd 缓存；先从当前节点导出 coredns/metrics-server/local-path/pause，重建后导入
- **端口冲突**：切换前停止 Compose；若需并行，改用 9080/9081 等备用 host 端口
- **数据丢失**：迁移前备份；验收前不删除 Compose 数据卷
- **回滚**：`kubectl delete namespace mall` 或 `make k8s-down`，然后 `docker compose up -d` 恢复原环境

## 验收标准

1. `kubectl get nodes` 3 个 Ready
2. `kubectl -n mall get pods` 全部 Running/Ready
3. `localhost:8000/users?page=1`、`localhost:8001/admin/orders`、`localhost:8002/products` 返回 200
4. 登录拿 token，创建订单成功
5. `localhost:8080/dashboard`、`localhost:8081/mall` 返回 200
6. `make smoke` 通过
7. users / products / orders 计数与 Compose 一致
8. Compose 仍可一键恢复

> 2026-10-18 后续：无状态服务已扩到 3 副本，见 `docs/adr/0002-k8s-ha-topology.md`。
