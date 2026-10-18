-- 2026-10-18 把演示商品归入二级分类
USE `looklook_product`;

UPDATE `product` SET `category_id` = CASE `name`
  WHEN 'Go 语言实战' THEN 9
  WHEN '无线蓝牙耳机' THEN 11
  WHEN '手冲咖啡豆' THEN 15
  WHEN '代码整洁之道' THEN 9
  WHEN '人类简史' THEN 10
  WHEN '小王子' THEN 8
  WHEN '便携蓝牙音箱' THEN 11
  WHEN '机械键盘' THEN 12
  WHEN '智能手表' THEN 13
  WHEN '明前龙井' THEN 14
  WHEN '手工黄油曲奇' THEN 16
  WHEN '每日坚果礼盒' THEN 16
  WHEN '雪松香薰蜡烛' THEN 17
  WHEN '手作陶瓷马克杯' THEN 18
  WHEN '亚麻抱枕' THEN 19
  WHEN '原木收纳盒' THEN 19
  WHEN '山茶花护手霜' THEN 20
  WHEN '氨基酸洁面' THEN 20
  WHEN '丝绒口红' THEN 21
  WHEN '白茶香氛喷雾' THEN 21
  WHEN '纯棉基础 T 恤' THEN 22
  WHEN '帆布托特包' THEN 23
  WHEN '羊毛围巾' THEN 23
  WHEN '天然橡胶瑜伽垫' THEN 24
  WHEN '不锈钢运动水壶' THEN 25
  WHEN '竞速跳绳' THEN 24
  ELSE `category_id`
END
WHERE `deleted_at` IS NULL;
