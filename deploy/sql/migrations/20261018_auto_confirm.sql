-- 2026-10-18 发货订单自动确认收货
USE `looklook_order`;

ALTER TABLE `orders`
  ADD KEY `idx_status_shipped` (`status`, `shipped_at`);
