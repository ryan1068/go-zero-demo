# 媒体服务独立（media）

日期：2026-10-19
状态：已实现

## 背景与目标

图片上传/读取目前耦合在 product-api。其他模块（头像、售后凭证、评价图）也需要上传，不应各自重复实现或依赖商品服务。

目标：新增独立 `app/media/api`，统一负责文件上传、读取、删除和元数据；调用方只保存返回的 URL。本期存储使用本机磁盘，保留旧 URL 兼容。

## 非目标

- 不接对象存储/CDN、不做预签名直传（后续演进）
- 不迁移前端静态商品图 `/products/*.svg`
- 不做图片裁剪、压缩、水印、审核
- 不改变现有商品/用户/订单业务接口语义

## 方案与取舍

| 项 | 方案 | 取舍 |
|---|---|---|
| 服务 | `app/media/api`，HTTP 8003，独立 go.mod | 不建 RPC；调用方只需要 URL |
| 存储 | 本机磁盘 `/app/uploads`，`Storage` 接口封装 | 单副本 + RWO PVC；后续可换 S3 |
| 元数据 | `looklook_media.media_file` | 支持归属、审计、软删除 |
| 鉴权 | JWT；上传/删除需登录，读取公开 | 与现有 `pkg/auth` 一致 |
| 删除权限 | 上传者本人或管理员 | 防止越权删除 |
| 新 URL | `/api/media/files/<name>` | 前端 nginx 代理 |
| 旧 URL | `/api/product/uploads/<name>` 由 product-api 反向代理到 media | 历史数据不失效 |
| 入口 | 前端 nginx 增加 `/api/media/` 路由 | 不经过 gateway |
| k8s | media 1 副本 + `media-files` PVC；product-api 移除上传卷 | 本机磁盘限制多副本 |

## 接口与数据变化

### media-api

```text
POST   /upload                multipart: file        -> {code,msg,url}
POST   /uploads/batch         multipart: files[]     -> {code,msg,urls}
DELETE /admin/files/:name                            -> {code,msg}
GET    /files/:name                                 -> 文件流（公开）
```

- 文件限制：jpg/jpeg/png/gif/webp，默认 10MB
- 文件名：随机 hex + 原扩展名；拒绝路径穿越
- 元数据：file_name、original_name、mime、size、uploader_id、created_at、deleted_at

### 数据库

新建 `looklook_media`，表 `media_file`；迁移脚本 `deploy/sql/migrations/20261019_media.sql`，初始化脚本 `deploy/sql/init_media.sql`。

### product-api

- 删除 `POST /admin/upload`、`POST /admin/uploads/batch`、`DELETE /admin/uploads/:name`
- 保留 `GET /uploads/:name`，改为反向代理到 `Media.Url + /files/:name`
- 配置新增 `Media.Url`

### 前端

- 新增 `mediaRequest = createRequest('/api/media')` 和 `src/api/media.ts`
- `ProductEdit.vue` 改用 media 上传接口
- nginx 增加 `/api/media/` 代理

## 错误处理与安全

- 类型、大小、空文件校验；上传失败清理半成品文件
- 文件名随机化，读取时强制 `filepath.Base` 并拒绝 `..`
- JWT 校验；上传者从 token 取，不信任请求参数
- 删除校验归属；软删除元数据，物理文件同步删除
- 旧 URL 代理在 media 不可用时返回 502，不静默 404

## 测试策略

- 失败测试先行：类型/大小校验、保存成功、元数据写入、删除归属、越权删除、文件读取
- 集成：上传 → 新 URL 200；旧 URL 200；删除后 404；未登录上传 401
- 前端：`npm run build`；商品编辑上传流程手工验证
- 门禁：`make check`、`make gen-check`、`make smoke`
- 数据：迁移前备份；复制现有上传文件到 media 卷

## 风险与回滚

- **本机磁盘仍是单副本**：media 只能 1 副本；换对象存储后才能 3 副本
- **旧 URL 依赖 media**：media 故障时历史图片 502
- **数据迁移**：先复制文件、再切前端；product-api 上传接口在验证通过后再删除
- **回滚**：git revert；Compose 下 `./data/uploads` 文件仍在；k8s 下保留 `media-files` PVC，product-api 恢复挂载即可

## 验收标准

1. media Pod 1/1 Ready
2. `POST /api/media/upload` 返回 URL
3. 新 URL 200；旧 `/api/product/uploads/<name>` 200
4. 商品编辑页上传、删除正常
5. 未登录上传 401；越权删除被拒绝
6. `make check`、`make gen-check`、`make smoke` 通过
7. users/products/orders 计数不变
