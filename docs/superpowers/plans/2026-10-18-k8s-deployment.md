# Kubernetes 运行环境迁移 Implementation Plan

> 状态：已实现。规格见 `docs/superpowers/specs/2026-10-18-k8s-deployment.md`。

**Goal:** 把 Compose 中的 8 个应用服务 + etcd/MySQL/Redis 迁到 k3d `mall` 集群，host 端口不变，数据保留，Compose 可回滚。

**Architecture:** 单命名空间 `mall`；Service 名沿用 Compose 的 DNS 名；go-zero 继续通过 etcd 做服务发现；配置用 ConfigMap，密钥用 Secret；frontend/API 用 NodePort + k3d loadbalancer 暴露。

**Tech Stack:** k3d v5.9 / k3s v1.36.4 / kubectl / MySQL 8.0.36 / Redis 6.2.5 / etcd 3.5

## Global Constraints

- 不改业务代码、`.api`、proto、DDL
- 密钥不进仓库、不打印
- 镜像 arm64；不依赖 Docker Hub 直连
- MySQL 迁移前必须备份；验收前不删除 Compose 卷
- 每步验证后再进入下一步；部署失败保留现场

---

### Task 1: 准备镜像、密钥和集群端口

**Files:**
- Modify: `deploy/k8s/k3d-cluster.yaml`（增加端口映射）
- Create: `scripts/k8s-images.sh`、`scripts/k8s-secrets.sh`
- Modify: `Makefile`（`k8s-images`、`k8s-secret`、`k8s-config`）

**Steps:**
- [x] 备份数据：`make db-backup`，记录 users/products/orders 计数
- [x] 从当前 3 个节点导出系统镜像（coredns、metrics-server、local-path、pause）到 `/tmp/k8s-system-images.tar`
- [x] 更新 k3d 端口映射：8080→30080、8081→30081、8000→30000、8001→30001、8002→30002、8888→30088
- [x] 重建集群，导入系统镜像，确认 3 节点 Ready
- [x] `scripts/k8s-images.sh`：构建 8 个镜像，`docker save --platform linux/arm64`，`ctr -n k8s.io images import` 到 3 个节点
- [x] `scripts/k8s-secrets.sh`：从 `.env` 创建 `mall-secrets`（JWT/MYSQL/REDIS）
- [x] `make k8s-config`：把 `app/*/etc/*.yaml` 创建为各服务 ConfigMap（不复制配置内容到仓库）
- [x] 验证：`kubectl -n mall get secret,configmap`；节点上 `ctr images ls` 能看到 8 个镜像

---

### Task 2: 基础设施（MySQL / Redis / etcd）

**Files:**
- Create: `deploy/k8s/base/infra.yaml`（MySQL StatefulSet + PVC + Service、Redis、etcd）
- Create: `deploy/k8s/base/namespace.yaml`

**Steps:**
- [x] apply namespace + infra；MySQL 用 local-path PVC 10Gi，Redis/etcd 用 emptyDir
- [x] 等待 `mysql-0` Ready；配置 TCP probe
- [x] 恢复数据：`kubectl -n mall exec -i mysql-0 -- mysql -uroot -p"$MYSQL_ROOT_PASSWORD" < backups/...sql`
- [x] 验证：三库存在；users/products/orders 计数与备份前一致

---

### Task 3: 应用服务与入口

**Files:**
- Create: `deploy/k8s/base/apps.yaml`
- Create: `deploy/k8s/base/frontend.yaml`
- Create: `deploy/k8s/base/kustomization.yaml`

**Steps:**
- [x] 7 个 Deployment + Service，镜像 `:local`，`imagePullPolicy: IfNotPresent`，ConfigMap 挂 `/app/etc`
- [x] product-api 挂 uploads PVC；所有 Pod 加 TCP readiness/liveness
- [x] frontend Deployment + 两个 NodePort Service（30080/30081）；API/gateway NodePort 30000/30001/30002/30088
- [x] 等待全部 Pod Ready；直接访问 8000/8001/8002 验证 200
- [x] 验证 etcd 注册：`kubectl -n mall exec` 或日志确认 user.rpc/order.rpc/product.rpc 已注册

---

### Task 4: 切换与验收

**Steps:**
- [x] `docker compose stop` 释放 8080/8081/8000-8002/8888（保留卷）
- [x] 验证 host 入口：8080 后台、8081 商城、8888 网关
- [x] `USER_API=http://localhost:8000 ORDER_API=http://localhost:8001 PRODUCT_API=http://localhost:8002 make smoke`
- [x] 手工流程：登录 → 加购 → 下单 → 支付 → 发货 → 确认收货
- [x] 对比 users/products/orders 计数
- [x] 失败时保留现场并回滚（`docker compose up -d`）
- [x] 通过后提交：`feat(k8s): 迁移运行环境到 Kubernetes`

---

### Task 5: 文档、回滚与 review

**Files:**
- Modify: `README.md`（k8s 部署/切换/回滚 runbook）
- Modify: `Makefile`（`k8s-deploy` / `k8s-delete` / `k8s-images`）

**Steps:**
- [x] 文档写清：前置、部署顺序、数据迁移、切换、回滚、常见问题
- [x] 独立 review manifest 和脚本
- [x] 验证回滚命令可用（不实际删除数据）
- [x] `git diff --check`；提交文档

---

## 验收

```bash
kubectl get nodes
kubectl -n mall get pods
curl -s localhost:8002/products?page=1
curl -s -o /dev/null -w '%{http_code}\n' localhost:8081/mall
make smoke
# 计数对比
docker compose up -d   # 回滚验证
```
