-- 2026-10-18 售后与退款
USE `looklook_order`;

CREATE TABLE IF NOT EXISTS `after_sale` (
  `id`                 bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `after_sale_no`      varchar(32) NOT NULL,
  `order_id`           bigint NOT NULL,
  `order_item_id`      bigint NOT NULL DEFAULT 0,
  `user_id`            bigint NOT NULL,
  `type`               tinyint NOT NULL DEFAULT 1,
  `reason`             varchar(255) NOT NULL DEFAULT '',
  `description`        varchar(500) NOT NULL DEFAULT '',
  `images`             json DEFAULT NULL,
  `amount`             bigint NOT NULL DEFAULT 0,
  `status`             tinyint NOT NULL DEFAULT 10,
  `admin_remark`       varchar(255) NOT NULL DEFAULT '',
  `tracking_no`        varchar(64) NOT NULL DEFAULT '',
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

CREATE TABLE IF NOT EXISTS `refund_record` (
  `id`            bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `refund_no`     varchar(32) NOT NULL,
  `after_sale_id` bigint NOT NULL,
  `order_id`      bigint NOT NULL,
  `user_id`       bigint NOT NULL,
  `amount`        bigint NOT NULL,
  `channel`       varchar(32) NOT NULL DEFAULT 'mock',
  `status`        tinyint NOT NULL DEFAULT 1,
  `refunded_at`   datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_at`    datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_refund_no` (`refund_no`),
  KEY `idx_after_sale` (`after_sale_id`),
  KEY `idx_order` (`order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='退款记录表';
