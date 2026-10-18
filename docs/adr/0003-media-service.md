# ADR 0003：媒体服务独立

- 状态：Accepted
- 日期：2026-10-19

## 背景

图片上传/读取原本在 product-api 内实现，文件写在节点本地磁盘。这带来两个问题：

- 任何需要上传的模块都要依赖 product-api，职责耦合
- product-api 因为有 RWO 上传卷，只能单副本

## 决策

1. **新增独立 `app/media/api`（8003）**
   - `POST /upload`、`POST /uploads/batch`、`DELETE /admin/files/:name`、`GET /files/:name`
   - 上传需要 JWT；读取公开
   - 删除仅管理员（不记录上传者，无法做归属校验）

2. **本机磁盘 + Storage 接口**
   - 当前用 local-path PVC（Compose 用 `./data/uploads`）
   - 代码通过 `Storage` 接口访问，后续可替换 S3/OSS
   - 本机磁盘限制下 media-api 单副本

3. **不落库**
   - 文件系统是唯一状态，上传只负责保存文件并返回 URL
   - media-api 不依赖 MySQL，减少启动依赖和故障面
   - 不提供文件列表、审计、软删除；删除是物理删除

4. **旧 URL 兼容**
   - `/api/product/uploads/<name>` 由 product-api 反向代理到 media-api
   - product-api 的上传/批量/删除接口移除

5. **职责调整**
   - product-api 变为无状态，恢复 3 副本 + PDB
   - media-api 单副本（RWO 卷），后续换对象存储后再扩副本

## 后果

- 所有模块可以统一使用 media-api，不再依赖商品服务
- 旧图片 URL 不失效
- media-api 不依赖数据库，但也没有文件归属信息
- 本机磁盘仍是单点；对象存储迁移是后续演进
- media-api 故障时，旧 URL 代理会返回 502

## 备选

- 预签名直传对象存储：需要前端/ CORS 改造，本期不做
- 继续放在 product-api 并加 NFS RWX：耦合仍在，不解决职责问题
