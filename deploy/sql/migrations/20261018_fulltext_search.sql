-- 2026-10-18 商品搜索改为 MySQL 全文索引（中文 ngram）
USE `looklook_product`;

ALTER TABLE `product`
  ADD FULLTEXT INDEX `ft_product_search` (`name`, `subtitle`, `detail`) WITH PARSER ngram;
