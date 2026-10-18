-- 媒体服务初始化（空数据卷或首次部署）
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
