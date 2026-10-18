-- 创建用户数据库
CREATE DATABASE IF NOT EXISTS `looklook_usercenter`
  DEFAULT CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;

USE `looklook_usercenter`;

-- 用户表
CREATE TABLE IF NOT EXISTS `user` (
  `id`            bigint NOT NULL AUTO_INCREMENT,
  `username`      varchar(64) NOT NULL COMMENT '用户名',
  `password`      varchar(255) NOT NULL COMMENT '密码（bcrypt）',
  `nickname`      varchar(255) DEFAULT '' COMMENT '昵称',
  `email`         varchar(128) NOT NULL DEFAULT '' COMMENT '邮箱',
  `phone`         varchar(32) NOT NULL DEFAULT '' COMMENT '手机号',
  `avatar`        varchar(512) NOT NULL DEFAULT '' COMMENT '头像地址',
  `status`        tinyint NOT NULL DEFAULT 1 COMMENT '状态：1=正常，0=禁用',
  `role`          tinyint NOT NULL DEFAULT 0 COMMENT '角色：0=普通用户，1=管理员',
  `remark`        varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
  `last_login_at` datetime DEFAULT NULL COMMENT '最后登录时间',
  `created_at`    datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at`    datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`),
  KEY `idx_email` (`email`),
  KEY `idx_phone` (`phone`),
  KEY `idx_nickname` (`nickname`),
  KEY `idx_status_id` (`status`, `id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 收货地址表
CREATE TABLE IF NOT EXISTS `user_address` (
  `id`              bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id`         bigint NOT NULL COMMENT '用户 ID',
  `receiver_name`   varchar(64) NOT NULL COMMENT '收货人',
  `receiver_phone`  varchar(32) NOT NULL COMMENT '收货电话',
  `province`        varchar(64) NOT NULL DEFAULT '' COMMENT '省',
  `city`            varchar(64) NOT NULL DEFAULT '' COMMENT '市',
  `district`        varchar(64) NOT NULL DEFAULT '' COMMENT '区/县',
  `detail_address`  varchar(255) NOT NULL COMMENT '详细地址',
  `postal_code`     varchar(16) NOT NULL DEFAULT '' COMMENT '邮编',
  `tag`             varchar(16) NOT NULL DEFAULT '' COMMENT '标签：家/公司/学校',
  `is_default`      tinyint NOT NULL DEFAULT 0 COMMENT '是否默认：1=是，0=否',
  `created_at`      datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at`      datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at`      datetime DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_default` (`user_id`, `is_default`, `id`),
  KEY `idx_user_created` (`user_id`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='收货地址表';

-- 测试账号：admin / 123456，管理员
INSERT INTO `user` (`username`, `password`, `nickname`, `email`, `phone`, `status`, `role`, `remark`) VALUES
('admin', '$2a$10$7Hz4L0untJOBPCKZdzv8AeydQeUTZtk51MdX4LWiFRU6Flm7YtNo.', '管理员', 'admin@example.com', '13800000000', 1, 1, '系统管理员');

-- 管理端操作审计
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='管理端操作审计';
