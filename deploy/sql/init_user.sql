-- 创建用户数据库
CREATE DATABASE IF NOT EXISTS `looklook_usercenter`
  DEFAULT CHARACTER SET utf8mb4
  COLLATE utf8mb4_general_ci;

USE `looklook_usercenter`;

-- 用户表
CREATE TABLE IF NOT EXISTS `user` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `username` varchar(64) NOT NULL COMMENT '用户名',
  `password` varchar(255) NOT NULL COMMENT '密码（bcrypt）',
  `nickname` varchar(255) DEFAULT '' COMMENT '昵称',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 测试账号：admin / 123456
INSERT INTO `user` (`username`, `password`, `nickname`) VALUES
('admin', '$2a$10$7Hz4L0untJOBPCKZdzv8AeydQeUTZtk51MdX4LWiFRU6Flm7YtNo.', '管理员');
