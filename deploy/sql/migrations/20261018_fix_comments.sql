-- 修复历史遗留的乱码表注释/列注释
--
-- 原因：早期通过默认字符集为 latin1 的 mysql 客户端执行含中文的 DDL，
-- MySQL 的 latin1 实为 cp1252，UTF-8 字节被按 cp1252 存成了乱码。
-- 本脚本把注释还原为正确的中文；执行时客户端必须使用 utf8mb4。

SET NAMES utf8mb4 COLLATE utf8mb4_unicode_ci;

ALTER TABLE `looklook_product`.`category` COMMENT='商品分类表';
ALTER TABLE `looklook_product`.`product` COMMENT='商品表';
ALTER TABLE `looklook_usercenter`.`user` MODIFY COLUMN `username` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户名';
ALTER TABLE `looklook_usercenter`.`user` MODIFY COLUMN `password` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '密码（bcrypt）';
ALTER TABLE `looklook_usercenter`.`user` MODIFY COLUMN `nickname` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '昵称';
ALTER TABLE `looklook_usercenter`.`user` MODIFY COLUMN `email` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '邮箱';
ALTER TABLE `looklook_usercenter`.`user` MODIFY COLUMN `phone` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '手机号';
ALTER TABLE `looklook_usercenter`.`user` MODIFY COLUMN `avatar` varchar(512) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '头像地址';
ALTER TABLE `looklook_usercenter`.`user` MODIFY COLUMN `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态：1=正常，0=禁用';
ALTER TABLE `looklook_usercenter`.`user` MODIFY COLUMN `remark` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '备注';
ALTER TABLE `looklook_usercenter`.`user` MODIFY COLUMN `role` tinyint NOT NULL DEFAULT '0' COMMENT '角色：0=普通用户，1=管理员';
ALTER TABLE `looklook_usercenter`.`user` MODIFY COLUMN `last_login_at` datetime DEFAULT NULL COMMENT '最后登录时间';
ALTER TABLE `looklook_usercenter`.`user` COMMENT='用户表';
ALTER TABLE `looklook_usercenter`.`user_address` MODIFY COLUMN `user_id` bigint NOT NULL COMMENT '用户 ID';
ALTER TABLE `looklook_usercenter`.`user_address` MODIFY COLUMN `receiver_name` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '收货人';
ALTER TABLE `looklook_usercenter`.`user_address` MODIFY COLUMN `receiver_phone` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '收货电话';
ALTER TABLE `looklook_usercenter`.`user_address` MODIFY COLUMN `province` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '省';
ALTER TABLE `looklook_usercenter`.`user_address` MODIFY COLUMN `city` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '市';
ALTER TABLE `looklook_usercenter`.`user_address` MODIFY COLUMN `district` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '区/县';
ALTER TABLE `looklook_usercenter`.`user_address` MODIFY COLUMN `detail_address` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '详细地址';
ALTER TABLE `looklook_usercenter`.`user_address` MODIFY COLUMN `postal_code` varchar(16) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '邮编';
ALTER TABLE `looklook_usercenter`.`user_address` MODIFY COLUMN `tag` varchar(16) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '标签';
ALTER TABLE `looklook_usercenter`.`user_address` MODIFY COLUMN `is_default` tinyint NOT NULL DEFAULT '0' COMMENT '是否默认：1=是，0=否';
ALTER TABLE `looklook_usercenter`.`user_address` COMMENT='收货地址表';
