-- 统一字符集与排序规则：utf8mb4 / utf8mb4_unicode_ci
--
-- 背景：MySQL 8 中建表只写 DEFAULT CHARSET=utf8mb4 时，会使用字符集自身的
-- 默认排序规则 utf8mb4_0900_ai_ci，而不是数据库的默认排序规则，导致
-- 库默认（utf8mb4_general_ci）与表实际（utf8mb4_0900_ai_ci）不一致。
-- 本脚本把三个库及其全部表的排序规则统一为 utf8mb4_unicode_ci。

ALTER DATABASE `looklook_usercenter` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
ALTER DATABASE `looklook_order` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
ALTER DATABASE `looklook_product` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

ALTER TABLE `looklook_usercenter`.`user` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
ALTER TABLE `looklook_usercenter`.`user_address` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
ALTER TABLE `looklook_usercenter`.`admin_audit_log` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

ALTER TABLE `looklook_product`.`category` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
ALTER TABLE `looklook_product`.`product` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
ALTER TABLE `looklook_product`.`product_sku` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

ALTER TABLE `looklook_order`.`cart_item` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
ALTER TABLE `looklook_order`.`orders` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
ALTER TABLE `looklook_order`.`order_item` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
ALTER TABLE `looklook_order`.`payment_record` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
ALTER TABLE `looklook_order`.`order_log` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
ALTER TABLE `looklook_order`.`after_sale` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
ALTER TABLE `looklook_order`.`refund_record` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
