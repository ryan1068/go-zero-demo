# 数据库设计规范

## 表设计

```sql
CREATE TABLE `user` (
    `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
    `username`   VARCHAR(64) NOT NULL DEFAULT '' COMMENT '用户名',
    `password`   VARCHAR(256) NOT NULL DEFAULT '' COMMENT '密码(bcrypt)',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '软删除',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';
```

- 主键：`id BIGINT UNSIGNED AUTO_INCREMENT`
- 时间：`created_at` + `updated_at` 必须
- 软删除：`deleted_at DATETIME DEFAULT NULL`
- 引擎：InnoDB，字符集：utf8mb4
- 字段必须有 COMMENT

## 索引规范

- 高频查询字段建索引
- 联合索引遵循最左前缀
- 唯一约束用 UNIQUE KEY
- 命名：`uk_xxx`（唯一）、`idx_xxx`（普通）

## SQL 编写

```sql
-- ✅ 正确
SELECT id, username FROM user WHERE id = ? AND deleted_at IS NULL

-- ❌ 禁止
SELECT * FROM user
```

- 禁止 `SELECT *`
- 必须带 `deleted_at IS NULL`（软删除表）
- 必须使用参数化查询，禁止拼接 SQL
- 大批量操作分批处理，单次不超过 1000 条

## Go 模型生成

```bash
# 使用 goctl 生成 model
goctl model mysql ddl -src user.sql -dir .
```

生成的文件：
- `xxxModel_gen.go` — 自动生成，勿手动改
- `xxxModel.go` — 自定义方法扩展

## 迁移规范

- 所有 DDL 变更写 SQL 脚本，放在 `deploy/script/mysql/`
- 脚本命名：`V{版本号}__{描述}.sql`
- 变更必须可回滚