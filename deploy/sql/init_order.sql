-- 创建订单数据库
CREATE DATABASE IF NOT EXISTS `looklook_order`
  DEFAULT CHARACTER SET utf8mb4
  COLLATE utf8mb4_general_ci;

USE `looklook_order`;

-- 购物车
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

-- 订单头
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
  KEY `idx_status_shipped` (`status`, `shipped_at`),
  KEY `idx_created` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单头表';

-- 订单明细
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

-- 支付记录
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

-- 订单状态日志
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

-- 售后申请
CREATE TABLE IF NOT EXISTS `after_sale` (
  `id`                 bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `after_sale_no`      varchar(32) NOT NULL COMMENT '售后单号',
  `order_id`           bigint NOT NULL COMMENT '订单 ID',
  `order_item_id`      bigint NOT NULL DEFAULT 0 COMMENT '订单明细 ID',
  `user_id`            bigint NOT NULL COMMENT '申请用户 ID',
  `type`               tinyint NOT NULL DEFAULT 1 COMMENT '1=仅退款 2=退货退款',
  `reason`             varchar(255) NOT NULL DEFAULT '' COMMENT '售后原因',
  `description`        varchar(500) NOT NULL DEFAULT '' COMMENT '问题描述',
  `images`             json DEFAULT NULL COMMENT '凭证图片',
  `amount`             bigint NOT NULL DEFAULT 0 COMMENT '退款金额，单位分',
  `status`             tinyint NOT NULL DEFAULT 10 COMMENT '10待审核 20已通过 30已拒绝 40买家已寄回 50商家已收货 60已退款 70已关闭',
  `admin_remark`       varchar(255) NOT NULL DEFAULT '' COMMENT '商家备注',
  `tracking_no`        varchar(64) NOT NULL DEFAULT '' COMMENT '退回物流单号',
  `apply_at`           datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `review_at`          datetime DEFAULT NULL,
  `buyer_shipped_at`   datetime DEFAULT NULL,
  `seller_received_at` datetime DEFAULT NULL,
  `refunded_at`        datetime DEFAULT NULL,
  `created_at`         datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at`         datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at`         datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_after_sale_no` (`after_sale_no`),
  KEY `idx_user_status_created` (`user_id`, `status`, `created_at`),
  KEY `idx_order` (`order_id`),
  KEY `idx_status_created` (`status`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='售后申请表';

-- 退款记录
CREATE TABLE IF NOT EXISTS `refund_record` (
  `id`            bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `refund_no`     varchar(32) NOT NULL COMMENT '退款单号',
  `after_sale_id` bigint NOT NULL COMMENT '售后 ID',
  `order_id`      bigint NOT NULL COMMENT '订单 ID',
  `user_id`       bigint NOT NULL COMMENT '用户 ID',
  `amount`        bigint NOT NULL COMMENT '退款金额，单位分',
  `channel`       varchar(32) NOT NULL DEFAULT 'mock' COMMENT '退款渠道',
  `status`        tinyint NOT NULL DEFAULT 1 COMMENT '1=成功',
  `refunded_at`   datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_at`    datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_refund_no` (`refund_no`),
  KEY `idx_after_sale` (`after_sale_id`),
  KEY `idx_order` (`order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='退款记录表';
