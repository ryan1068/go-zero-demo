-- 商城 MVP 迁移脚本（已有数据卷执行一次）
-- 用法：
--   docker compose exec -T mysql mysql -uroot -pPXDN93VRKUm8TeE7 < deploy/sql/migrations/20261018_mall.sql
--
-- 注意：
-- 1. 本脚本只执行一次；重复执行会因列/表已存在而报错。
-- 2. 旧 orders 表会改名为 orders_legacy_20261018 保留，不迁移旧演示订单。

-- ---------- 用户库 ----------
USE `looklook_usercenter`;

ALTER TABLE `user`
  ADD COLUMN `role` tinyint NOT NULL DEFAULT 0 COMMENT '角色：0=普通用户，1=管理员' AFTER `status`,
  ADD COLUMN `last_login_at` datetime DEFAULT NULL COMMENT '最后登录时间' AFTER `remark`,
  ALGORITHM=INSTANT;

UPDATE `user` SET `role` = 1 WHERE `username` = 'admin';

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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='收货地址表';

-- ---------- 商品库 ----------
CREATE DATABASE IF NOT EXISTS `looklook_product`
  DEFAULT CHARACTER SET utf8mb4
  COLLATE utf8mb4_general_ci;

USE `looklook_product`;

CREATE TABLE IF NOT EXISTS `category` (
  `id`          bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `parent_id`   bigint NOT NULL DEFAULT 0 COMMENT '父分类 ID，0=一级分类',
  `name`        varchar(64) NOT NULL COMMENT '分类名称',
  `icon`        varchar(512) NOT NULL DEFAULT '' COMMENT '分类图标',
  `sort`        int NOT NULL DEFAULT 0 COMMENT '排序，越大越靠前',
  `status`      tinyint NOT NULL DEFAULT 1 COMMENT '状态：1=启用，0=禁用',
  `created_at`  datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at`  datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at`  datetime DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_parent_status_sort` (`parent_id`, `status`, `sort`),
  KEY `idx_status_sort` (`status`, `sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品分类表';

CREATE TABLE IF NOT EXISTS `product` (
  `id`              bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `category_id`     bigint NOT NULL COMMENT '分类 ID',
  `name`            varchar(128) NOT NULL COMMENT '商品名称',
  `subtitle`        varchar(255) NOT NULL DEFAULT '' COMMENT '副标题/卖点',
  `cover_image`     varchar(512) NOT NULL DEFAULT '' COMMENT '封面图',
  `images`          json DEFAULT NULL COMMENT '图片地址数组',
  `detail`          text COMMENT '商品详情',
  `price`           bigint NOT NULL COMMENT '售价，单位分',
  `original_price`  bigint NOT NULL DEFAULT 0 COMMENT '原价，单位分',
  `stock`           int NOT NULL DEFAULT 0 COMMENT '库存',
  `sales`           int NOT NULL DEFAULT 0 COMMENT '销量',
  `unit`            varchar(16) NOT NULL DEFAULT '件' COMMENT '单位',
  `status`          tinyint NOT NULL DEFAULT 1 COMMENT '状态：1=上架，0=下架',
  `is_recommend`    tinyint NOT NULL DEFAULT 0 COMMENT '首页推荐：1=是，0=否',
  `sort`            int NOT NULL DEFAULT 0 COMMENT '排序，越大越靠前',
  `created_at`      datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at`      datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at`      datetime DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_category_status_id` (`category_id`, `status`, `id`),
  KEY `idx_status_sort` (`status`, `sort`, `id`),
  KEY `idx_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品表';

INSERT INTO `category` (`parent_id`, `name`, `icon`, `sort`, `status`) VALUES
(0, '图书', '', 30, 1),
(0, '数码', '', 20, 1),
(0, '食品', '', 10, 1);

INSERT INTO `product`
  (`category_id`, `name`, `subtitle`, `cover_image`, `images`, `detail`, `price`, `original_price`, `stock`, `sales`, `unit`, `status`, `is_recommend`, `sort`)
VALUES
(1, 'Go 语言实战', '从入门到微服务实战', '/products/book.svg', JSON_ARRAY('/products/book.svg'), 'Go 语言核心语法、并发编程与微服务实践。', 8900, 12900, 200, 0, '本', 1, 1, 30),
(2, '无线蓝牙耳机', '降噪长续航', '/products/headphones.svg', JSON_ARRAY('/products/headphones.svg'), '主动降噪，单次续航 8 小时，配合充电盒可用 32 小时。', 29900, 39900, 100, 0, '副', 1, 1, 20),
(3, '手冲咖啡豆', '中度烘焙 250g', '/products/coffee.svg', JSON_ARRAY('/products/coffee.svg'), '精选阿拉比卡豆，中度烘焙，坚果与焦糖风味。', 6800, 8800, 300, 0, '袋', 1, 1, 10);

-- ---------- 订单库 ----------
USE `looklook_order`;

-- 保留旧表；只执行一次
RENAME TABLE `orders` TO `orders_legacy_20261018`;

CREATE TABLE IF NOT EXISTS `cart_item` (
  `id`          bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id`     bigint NOT NULL COMMENT '用户 ID',
  `product_id`  bigint NOT NULL COMMENT '商品 ID',
  `quantity`    int NOT NULL DEFAULT 1 COMMENT '数量',
  `selected`    tinyint NOT NULL DEFAULT 1 COMMENT '是否选中：1=是，0=否',
  `created_at`  datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at`  datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_product` (`user_id`, `product_id`),
  KEY `idx_user_selected` (`user_id`, `selected`, `updated_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='购物车表';

CREATE TABLE IF NOT EXISTS `orders` (
  `id`                    bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `order_no`              varchar(32) NOT NULL COMMENT '订单号',
  `user_id`               bigint NOT NULL COMMENT '下单用户 ID',
  `total_amount`          bigint NOT NULL COMMENT '商品总额，单位分',
  `freight_amount`        bigint NOT NULL DEFAULT 0 COMMENT '运费，单位分',
  `discount_amount`       bigint NOT NULL DEFAULT 0 COMMENT '优惠金额，单位分',
  `pay_amount`            bigint NOT NULL COMMENT '应付金额，单位分',
  `status`                tinyint NOT NULL DEFAULT 10 COMMENT '状态：10待支付 20待发货 30待收货 40已完成 50已取消',
  `pay_type`              tinyint NOT NULL DEFAULT 0 COMMENT '支付方式：0未支付 1模拟支付 2微信 3支付宝',
  `transaction_id`        varchar(128) NOT NULL DEFAULT '' COMMENT '支付流水号',
  `receiver_name`         varchar(64) NOT NULL COMMENT '收货人快照',
  `receiver_phone`        varchar(32) NOT NULL COMMENT '收货电话快照',
  `receiver_province`     varchar(64) NOT NULL DEFAULT '' COMMENT '省快照',
  `receiver_city`         varchar(64) NOT NULL DEFAULT '' COMMENT '市快照',
  `receiver_district`     varchar(64) NOT NULL DEFAULT '' COMMENT '区县快照',
  `receiver_address`      varchar(255) NOT NULL COMMENT '详细地址快照',
  `receiver_postal_code`  varchar(16) NOT NULL DEFAULT '' COMMENT '邮编快照',
  `remark`                varchar(255) NOT NULL DEFAULT '' COMMENT '用户备注',
  `admin_remark`          varchar(255) NOT NULL DEFAULT '' COMMENT '管理员备注',
  `logistics_company`     varchar(64) NOT NULL DEFAULT '' COMMENT '物流公司',
  `tracking_no`           varchar(64) NOT NULL DEFAULT '' COMMENT '物流单号',
  `paid_at`               datetime DEFAULT NULL COMMENT '支付时间',
  `shipped_at`            datetime DEFAULT NULL COMMENT '发货时间',
  `completed_at`          datetime DEFAULT NULL COMMENT '完成时间',
  `canceled_at`           datetime DEFAULT NULL COMMENT '取消时间',
  `cancel_reason`         varchar(255) NOT NULL DEFAULT '' COMMENT '取消原因',
  `created_at`            datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at`            datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at`            datetime DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_order_no` (`order_no`),
  KEY `idx_user_status_created` (`user_id`, `status`, `created_at`),
  KEY `idx_status_created` (`status`, `created_at`),
  KEY `idx_created` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单头表';

CREATE TABLE IF NOT EXISTS `order_item` (
  `id`             bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `order_id`       bigint NOT NULL COMMENT '订单 ID',
  `product_id`     bigint NOT NULL COMMENT '商品 ID',
  `product_name`   varchar(128) NOT NULL COMMENT '商品名称快照',
  `product_cover`  varchar(512) NOT NULL DEFAULT '' COMMENT '商品封面快照',
  `price`          bigint NOT NULL COMMENT '下单单价，单位分',
  `quantity`       int NOT NULL COMMENT '数量',
  `total_amount`   bigint NOT NULL COMMENT '小计，单位分',
  `created_at`     datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_order` (`order_id`),
  KEY `idx_product` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单明细表';

CREATE TABLE IF NOT EXISTS `payment_record` (
  `id`              bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `payment_no`      varchar(32) NOT NULL COMMENT '支付单号',
  `order_id`        bigint NOT NULL COMMENT '订单 ID',
  `order_no`        varchar(32) NOT NULL COMMENT '订单号',
  `user_id`         bigint NOT NULL COMMENT '用户 ID',
  `amount`          bigint NOT NULL COMMENT '支付金额，单位分',
  `pay_type`        tinyint NOT NULL DEFAULT 1 COMMENT '支付方式：1模拟支付',
  `status`          tinyint NOT NULL DEFAULT 0 COMMENT '状态：0待支付 1成功 2失败',
  `transaction_id`  varchar(128) NOT NULL DEFAULT '' COMMENT '第三方流水号',
  `paid_at`         datetime DEFAULT NULL COMMENT '支付时间',
  `created_at`      datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at`      datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_payment_no` (`payment_no`),
  KEY `idx_order` (`order_id`),
  KEY `idx_user_created` (`user_id`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='支付记录表';

CREATE TABLE IF NOT EXISTS `order_log` (
  `id`             bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `order_id`       bigint NOT NULL COMMENT '订单 ID',
  `from_status`    tinyint NOT NULL COMMENT '原状态',
  `to_status`      tinyint NOT NULL COMMENT '新状态',
  `operator_type`  tinyint NOT NULL COMMENT '操作人类型：1=用户 2=管理员 3=系统',
  `operator_id`    bigint NOT NULL DEFAULT 0 COMMENT '操作人 ID',
  `remark`         varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
  `created_at`     datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_order_created` (`order_id`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单状态日志表';

-- 演示地址（admin 用户 ID 通常为 1）
INSERT INTO `looklook_usercenter`.`user_address`
  (`user_id`, `receiver_name`, `receiver_phone`, `province`, `city`, `district`, `detail_address`, `postal_code`, `tag`, `is_default`)
VALUES
(1, '管理员', '13800000000', '上海市', '上海市', '浦东新区', '张江高科技园区 1 号', '201203', '公司', 1);
