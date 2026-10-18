# ADR 0002：Kubernetes 应用层多副本拓扑

- 状态：Accepted
- 日期：2026-10-18

## 背景

迁移到 k8s 后所有服务都是单副本，任意节点故障都会导致对应服务中断。用户要求无状态服务具备多节点冗余，同时接受 MySQL/Redis/etcd 和上传存储在本地环境保持单实例。

## 决策

1. **无状态服务 3 副本**
   - gateway、frontend、user/order/product 的 api 和 rpc
   - 3 个节点上每个服务正好 1 个 Pod

2. **硬拓扑分散**
   - `topologySpreadConstraints` 使用 `whenUnsatisfiable: DoNotSchedule`、`maxSkew: 1`、`topologyKey: kubernetes.io/hostname`
   - 保证同一服务不会把两个 Pod 放在同一节点

3. **滚动更新优先保可用**
   - `maxSurge: 0`、`maxUnavailable: 1`
   - 更新时先停 1 个、再补 1 个，任一时点至少 2 个可用

4. **PodDisruptionBudget**
   - 上述 7 个 Deployment 设置 `minAvailable: 2`
   - 节点维护/驱逐时至少保留 2 个实例

5. **product-api 保持 1 副本**
   - 上传目录是 local-path RWO 卷，跨节点多副本无法挂载
   - 使用 `strategy: Recreate` 避免滚动更新卡在卷亲和性

6. **MySQL / Redis / etcd 保持单实例**
   - 本地开发可接受；生产需要主从、Sentinel 或托管服务
   - 服务发现仍是单 etcd，etcd 故障会影响新注册和客户端发现

## 后果

- 单个无状态 Pod 或节点故障时，服务仍有 2/3 实例可用
- 硬拓扑分散意味着节点故障后第 3 个副本会 Pending，直到节点恢复或节点数增加；服务本身不受影响
- 24 个 Pod（21 个无状态 + product-api + MySQL/Redis/etcd）对 8GB Docker Desktop 仍有压力，必要时降低副本数
- 数据层仍是单点；不要把本拓扑当作生产高可用方案

## 备选

- 软拓扑（`ScheduleAnyway`）：允许故障后 2+1 分布，但初始分布不保证每节点 1 个
- 全部服务 3 副本：需要 RWX 共享存储或对象存储，并修改 product 服务上传逻辑
- MySQL 主从 / Redis Sentinel / etcd 集群：改动和运维成本高，另立 spec
