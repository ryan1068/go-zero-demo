-- 2026-10-18 管理端操作审计
USE `looklook_usercenter`;

CREATE TABLE IF NOT EXISTS `admin_audit_log` (
  `id`          bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `admin_id`    bigint NOT NULL COMMENT '操作人 ID',
  `action`      varchar(64) NOT NULL COMMENT '动作',
  `target_type` varchar(32) NOT NULL COMMENT '目标类型',
  `target_id`   bigint NOT NULL DEFAULT 0 COMMENT '目标 ID',
  `detail`      varchar(512) NOT NULL DEFAULT '' COMMENT '详情',
  `created_at`  datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_admin_created` (`admin_id`, `created_at`),
  KEY `idx_target` (`target_type`, `target_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='管理端操作审计';
