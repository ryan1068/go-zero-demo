-- 2026-10-18 商品 SKU
USE `looklook_product`;

CREATE TABLE IF NOT EXISTS `product_sku` (
  `id`             bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `product_id`     bigint NOT NULL,
  `sku_code`       varchar(64) NOT NULL,
  `specs`          varchar(500) NOT NULL DEFAULT '{}',
  `price`          bigint NOT NULL,
  `original_price` bigint NOT NULL DEFAULT 0,
  `stock`          int NOT NULL DEFAULT 0,
  `sales`          int NOT NULL DEFAULT 0,
  `status`         tinyint NOT NULL DEFAULT 1,
  `created_at`     datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at`     datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at`     datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_sku_code` (`sku_code`),
  KEY `idx_product_status` (`product_id`, `status`, `id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品 SKU 表';
