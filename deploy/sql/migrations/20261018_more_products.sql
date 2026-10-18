-- 2026-10-18 扩充商品分类与演示商品
-- 适用于已经初始化过的 looklook_product 库

USE `looklook_product`;

SET NAMES utf8mb4;

-- 清理早期测试数据
DELETE FROM `product` WHERE `name` = '测试商品' AND `status` = 0;
DELETE FROM `category`
WHERE `name` = '测试分类改'
  AND `id` NOT IN (SELECT DISTINCT `category_id` FROM `product`);

-- 扩展分类
INSERT INTO `category` (`id`, `parent_id`, `name`, `icon`, `sort`, `status`) VALUES
(4, 0, '家居', '', 50, 1),
(5, 0, '美妆', '', 40, 1),
(6, 0, '服饰', '', 35, 1),
(7, 0, '运动', '', 25, 1)
ON DUPLICATE KEY UPDATE
  `parent_id` = VALUES(`parent_id`),
  `name` = VALUES(`name`),
  `icon` = VALUES(`icon`),
  `sort` = VALUES(`sort`),
  `status` = VALUES(`status`);

-- 扩展商品
INSERT INTO `product`
  (`id`, `category_id`, `name`, `subtitle`, `cover_image`, `images`, `detail`, `price`, `original_price`, `stock`, `sales`, `unit`, `status`, `is_recommend`, `sort`)
VALUES
(4, 1, '代码整洁之道', '程序员必读的整洁代码实践', '/products/book2.svg', JSON_ARRAY('/products/book2.svg'), '从命名、函数到注释，系统讲解如何写出可维护的整洁代码。', 7900, 9900, 180, 3, '本', 1, 1, 29),
(5, 1, '人类简史', '从动物到上帝', '/products/book3.svg', JSON_ARRAY('/products/book3.svg'), '以宏大视角梳理人类从认知革命到科技革命的历史。', 5400, 6800, 220, 5, '本', 1, 0, 28),
(6, 1, '小王子', '写给大人的童话', '/products/book4.svg', JSON_ARRAY('/products/book4.svg'), '关于爱、责任与成长的经典童话，适合反复阅读。', 3200, 4500, 150, 6, '本', 1, 0, 27),
(7, 2, '便携蓝牙音箱', '360° 环绕声，户外随身', '/products/speaker.svg', JSON_ARRAY('/products/speaker.svg'), '小巧便携，支持蓝牙 5.3，单次续航 12 小时。', 19900, 25900, 80, 2, '台', 1, 1, 19),
(8, 2, '机械键盘', '87 键茶轴，办公游戏两用', '/products/keyboard.svg', JSON_ARRAY('/products/keyboard.svg'), '87 键紧凑布局，茶轴手感，支持有线和无线双模。', 32900, 42900, 60, 4, '把', 1, 0, 18),
(9, 2, '智能手表', '心率与睡眠监测', '/products/watch.svg', JSON_ARRAY('/products/watch.svg'), '全天候心率监测，支持睡眠分析和多种运动模式。', 59900, 79900, 45, 2, '块', 1, 1, 17),
(10, 3, '明前龙井', '春季头采 100g', '/products/tea.svg', JSON_ARRAY('/products/tea.svg'), '明前头采嫩芽，豆香清雅，汤色明亮。', 12800, 16800, 120, 4, '罐', 1, 1, 9),
(11, 3, '手工黄油曲奇', '原味 200g', '/products/cookie.svg', JSON_ARRAY('/products/cookie.svg'), '新西兰黄油手工烘焙，酥松不腻，独立小包装。', 4600, 5800, 200, 7, '盒', 1, 0, 8),
(12, 3, '每日坚果礼盒', '30 包混合装', '/products/nuts.svg', JSON_ARRAY('/products/nuts.svg'), '巴旦木、腰果、核桃仁与果干科学配比，每日一包。', 8900, 11900, 160, 5, '盒', 1, 0, 7),
(13, 4, '雪松香薰蜡烛', '木质香调 180g', '/products/candle.svg', JSON_ARRAY('/products/candle.svg'), '雪松与琥珀交织的木质香，燃烧时间约 40 小时。', 9900, 12900, 90, 3, '罐', 1, 1, 6),
(14, 4, '手作陶瓷马克杯', '月白釉 350ml', '/products/mug.svg', JSON_ARRAY('/products/mug.svg'), '手工拉坯，月白釉面，杯口圆润，适合日常咖啡与茶。', 6800, 8800, 110, 4, '只', 1, 0, 5),
(15, 4, '亚麻抱枕', '45×45cm 含芯', '/products/pillow.svg', JSON_ARRAY('/products/pillow.svg'), '天然亚麻面料，内填高弹棉，适合沙发与卧室。', 7900, 9900, 70, 2, '个', 1, 0, 4),
(16, 4, '原木收纳盒', '榫卯结构 中号', '/products/box.svg', JSON_ARRAY('/products/box.svg'), '白蜡木榫卯结构，无漆无蜡，适合桌面收纳。', 11900, 14900, 50, 1, '个', 1, 0, 3),
(17, 5, '山茶花护手霜', '滋润不黏腻 50ml', '/products/handcream.svg', JSON_ARRAY('/products/handcream.svg'), '山茶花油与乳木果油配方，吸收快，留香清淡。', 3900, 4900, 180, 6, '支', 1, 1, 6),
(18, 5, '氨基酸洁面', '温和清洁 120g', '/products/cleanser.svg', JSON_ARRAY('/products/cleanser.svg'), '氨基酸表活，泡沫细腻，洗后不紧绷。', 5900, 7900, 140, 5, '支', 1, 0, 5),
(19, 5, '丝绒口红', '豆沙色 3.5g', '/products/lipstick.svg', JSON_ARRAY('/products/lipstick.svg'), '丝绒哑光质地，显色饱满，豆沙色日常百搭。', 8900, 10900, 100, 4, '支', 1, 1, 4),
(20, 5, '白茶香氛喷雾', '清透白茶调 50ml', '/products/perfume.svg', JSON_ARRAY('/products/perfume.svg'), '白茶与鼠尾草的清透香气，适合日常与办公场景。', 12900, 15900, 80, 3, '瓶', 1, 0, 3),
(21, 6, '纯棉基础 T 恤', '260g 重磅纯棉', '/products/tshirt.svg', JSON_ARRAY('/products/tshirt.svg'), '260g 重磅纯棉，版型利落，耐洗不易变形。', 6900, 8900, 150, 5, '件', 1, 1, 4),
(22, 6, '帆布托特包', '加厚帆布 大容量', '/products/tote.svg', JSON_ARRAY('/products/tote.svg'), '加厚帆布，内置小袋，通勤与周末出行都合适。', 5900, 7900, 120, 4, '个', 1, 0, 3),
(23, 6, '羊毛围巾', '100% 羊毛 燕麦色', '/products/scarf.svg', JSON_ARRAY('/products/scarf.svg'), '100% 羊毛，柔软亲肤，燕麦色百搭。', 15900, 19900, 60, 2, '条', 1, 0, 2),
(24, 7, '天然橡胶瑜伽垫', '防滑 5mm', '/products/yoga.svg', JSON_ARRAY('/products/yoga.svg'), '天然橡胶防滑底，5mm 厚度，回弹舒适。', 12900, 16900, 70, 3, '张', 1, 1, 3),
(25, 7, '不锈钢运动水壶', '750ml 保温', '/products/bottle.svg', JSON_ARRAY('/products/bottle.svg'), '304 不锈钢内胆，保温保冷，750ml 大容量。', 7900, 9900, 100, 4, '个', 1, 0, 2),
(26, 7, '竞速跳绳', '轴承钢丝绳', '/products/rope.svg', JSON_ARRAY('/products/rope.svg'), '轴承设计，钢丝绳可调长度，适合日常燃脂训练。', 2900, 3900, 200, 6, '条', 1, 0, 1)
ON DUPLICATE KEY UPDATE
  `category_id` = VALUES(`category_id`),
  `name` = VALUES(`name`),
  `subtitle` = VALUES(`subtitle`),
  `cover_image` = VALUES(`cover_image`),
  `images` = VALUES(`images`),
  `detail` = VALUES(`detail`),
  `price` = VALUES(`price`),
  `original_price` = VALUES(`original_price`),
  `stock` = VALUES(`stock`),
  `sales` = VALUES(`sales`),
  `unit` = VALUES(`unit`),
  `status` = VALUES(`status`),
  `is_recommend` = VALUES(`is_recommend`),
  `sort` = VALUES(`sort`);
