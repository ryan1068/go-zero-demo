-- 媒体服务迁移（已有数据卷执行一次）
-- 用法：
--   docker compose exec -T mysql mysql -uroot -p"$MYSQL_ROOT_PASSWORD" < deploy/sql/migrations/20261019_media.sql
--   kubectl -n mall exec -i mysql-0 -- sh -c 'MYSQL_PWD="$MYSQL_ROOT_PASSWORD" mysql -uroot' < deploy/sql/migrations/20261019_media.sql

CREATE DATABASE IF NOT EXISTS `looklook_media`
  DEFAULT CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;

USE `looklook_media`;

CREATE TABLE IF NOT EXISTS `media_file` (
  `id`            bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `file_name`     varchar(128)    NOT NULL COMMENT '存储文件名',
  `original_name` varchar(255)    NOT NULL DEFAULT '' COMMENT '原始文件名',
  `mime`          varchar(100)    NOT NULL DEFAULT '' COMMENT 'MIME 类型',
  `size`          bigint          NOT NULL DEFAULT 0 COMMENT '字节数',
  `uploader_id`   bigint          NOT NULL DEFAULT 0 COMMENT '上传者用户 ID',
  `created_at`    datetime        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at`    datetime        DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_file_name` (`file_name`),
  KEY `idx_uploader` (`uploader_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='媒体文件元数据';
