# 媒体服务独立 Implementation Plan

> 状态：已实现。规格见 `docs/superpowers/specs/2026-10-19-media-service.md`。

**Goal:** 新增 `app/media/api`，把上传/读取/删除从 product-api 解耦；本机磁盘存储，旧 URL 兼容。

**Architecture:** media-api（HTTP 8003）负责文件与元数据；product-api 只保留商品接口和旧 URL 反代；前端改走 `/api/media/`。

**Tech Stack:** Go 1.23 / go-zero 1.8.5 / MySQL / Vue3

## Global Constraints

- 业务逻辑只放 `internal/logic/`；生成文件不手改
- 用户输入校验；错误必须处理；密码/密钥不入日志
- 本机磁盘存储；单副本；不引入对象存储
- 旧 URL `/api/product/uploads/<name>` 必须继续可用
- 数据迁移前备份；验收后清理测试文件

---

### Task 1: 媒体服务骨架与失败测试

**Files:**
- Create: `app/media/api/go.mod`, `media.api`, `etc/media-api.yaml`, `Dockerfile`
- Create: `app/media/api/internal/{config,svc,types,handler,logic,model}`
- Test: `app/media/api/internal/logic/uploadlogic_test.go`, `deletelogic_test.go`

**Steps:**
- [x] 写失败测试：类型/大小校验、保存成功、元数据写入、删除归属、越权删除
- [x] 运行测试确认失败
- [x] 建 go.mod（依赖 go-zero、pkg/auth、pkg/errcode + replace）
- [x] 定义 `media.api`：upload / batch / delete / get file
- [x] 生成 handler/types（goctl 或按现有服务手工对齐），multipart 与文件流 handler 保持手写并从 gen-check 排除
- [x] 配置 `Upload.Dir`、`Upload.MaxSizeMB`、`Jwt.Secret`

---

### Task 2: 实现上传、元数据与删除

**Files:**
- Create: `app/media/api/internal/logic/{uploadlogic,deletefilelogic}.go`
- Create: `app/media/api/internal/handler/servefilehandler.go`

**Steps:**
- [x] `Storage` 接口 + `LocalStorage`：Save/Open/Delete
- [x] 上传校验（扩展名、Content-Type、大小、空文件）并清理半成品
- [x] 不写数据库；上传只保存文件并返回 URL
- [x] 删除：仅管理员，物理删除；文件不存在返回 404
- [x] 读取：`filepath.Base` 防穿越，设置 Cache-Control
- [x] 测试通过；`go build ./...`

---

### Task 3: product-api 兼容与前端切换

**Files:**
- Modify: `app/product/api/etc/product.yaml`（`Media.Url`）
- Modify: `app/product/api/product.api`（移除上传接口）
- Modify: `app/product/api/internal/handler/serveuploadhandler.go`（反代到 media）
- Modify: `app/frontend/src/api/product.ts`, `src/api/media.ts`, `src/views/admin/ProductEdit.vue`
- Modify: `app/frontend/nginx.conf`, `deploy/k8s/frontend-nginx.conf`

**Steps:**
- [x] product-api 保留 `GET /uploads/:name`，改为反代 `Media.Url/files/:name`
- [x] 删除 product-api 上传/批量/删除接口及前端调用
- [x] 前端新增 media 上传 API 并替换 ProductEdit
- [x] nginx 增加 `/api/media/` 代理（Compose 与 k8s 两份配置）
- [x] 更新 gen-check 排除项；`make gen-check` 通过

---

### Task 4: Compose / k8s 部署与数据迁移

**Files:**
- Modify: `docker-compose.yml`, `scripts/k8s-config.sh`, `scripts/k8s-images.sh`
- Modify: `deploy/k8s/base/apps.yaml`, `infra.yaml`（mysql init configmap）
- Create: `deploy/k8s/base/media-files` PVC（放在 apps.yaml）

**Steps:**
- [x] Compose 增加 media 服务，复用 `./data/uploads`；product-api 移除上传卷
- [x] k8s 增加 media Deployment/Service（1 副本、Recreate、`media-files` PVC）
- [x] product-api 移除 uploads 卷，配置 `Media.Url: http://media:8003`
- [x] media 配置加入 ConfigMap 脚本；镜像加入 `k8s-images.sh`
- [x] 复制现有上传文件到 media 卷
- [x] 部署并等待 Ready

---

### Task 5: 验收、文档与回滚

**Steps:**
- [x] 集成验证：上传、新 URL、旧 URL、删除、越权、未登录
- [x] 前端商品编辑上传流程手工验证
- [x] `make check`、`make gen-check`、`make smoke`
- [x] 清理测试文件；计数核对
- [x] ADR 0003 + README 更新；独立 review；提交

---

## 验收

```bash
kubectl -n mall get pods | grep media
curl -F file=@test.png localhost:8003/upload -H "Authorization: Bearer $TOKEN"
curl -o /dev/null -w '%{http_code}' localhost:8081/api/media/files/<name>
curl -o /dev/null -w '%{http_code}' localhost:8081/api/product/uploads/<old-name>
make check && make gen-check && make smoke
```
