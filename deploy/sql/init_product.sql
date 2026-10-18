-- 创建商品数据库
CREATE DATABASE IF NOT EXISTS `looklook_product`
  DEFAULT CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;

USE `looklook_product`;

SET NAMES utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 商品分类
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商品分类表';

-- 商品
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
  KEY `idx_name` (`name`),
  FULLTEXT KEY `ft_product_search` (`name`, `subtitle`, `detail`) WITH PARSER ngram
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商品表';

-- 演示分类
INSERT INTO `category` (`parent_id`, `name`, `icon`, `sort`, `status`) VALUES
(0, '图书', '', 30, 1),
(0, '数码', '', 20, 1),
(0, '食品', '', 10, 1);

-- 扩展分类
INSERT INTO `category` (`id`, `parent_id`, `name`, `icon`, `sort`, `status`) VALUES
(4, 0, '家居', '', 50, 1),
(5, 0, '美妆', '', 40, 1),
(6, 0, '服饰', '', 35, 1),
(7, 0, '运动', '', 25, 1);

-- 演示商品
INSERT INTO `product`
  (`category_id`, `name`, `subtitle`, `cover_image`, `images`, `detail`, `price`, `original_price`, `stock`, `sales`, `unit`, `status`, `is_recommend`, `sort`)
VALUES
(1, 'Go 语言实战', '从入门到微服务实战', '/products/book.svg', JSON_ARRAY('/products/book.svg'), 'Go 语言核心语法、并发编程与微服务实践。', 8900, 12900, 200, 0, '本', 1, 1, 30),
(2, '无线蓝牙耳机', '降噪长续航', '/products/headphones.svg', JSON_ARRAY('/products/headphones.svg'), '主动降噪，单次续航 8 小时，配合充电盒可用 32 小时。', 29900, 39900, 100, 0, '副', 1, 1, 20),
(3, '手冲咖啡豆', '中度烘焙 250g', '/products/coffee.svg', JSON_ARRAY('/products/coffee.svg'), '精选阿拉比卡豆，中度烘焙，坚果与焦糖风味。', 6800, 8800, 300, 0, '袋', 1, 1, 10);

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
(26, 7, '竞速跳绳', '轴承钢丝绳', '/products/rope.svg', JSON_ARRAY('/products/rope.svg'), '轴承设计，钢丝绳可调长度，适合日常燃脂训练。', 2900, 3900, 200, 6, '条', 1, 0, 1);

-- 商品 SKU / 多规格
CREATE TABLE IF NOT EXISTS `product_sku` (
  `id`             bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `product_id`     bigint NOT NULL COMMENT '商品 ID',
  `sku_code`       varchar(64) NOT NULL COMMENT 'SKU 编码',
  `specs`          varchar(500) NOT NULL DEFAULT '{}' COMMENT '规格 JSON，如 {"颜色":"胭脂红","尺寸":"M"}',
  `price`          bigint NOT NULL COMMENT '售价，单位分',
  `original_price` bigint NOT NULL DEFAULT 0 COMMENT '原价，单位分',
  `stock`          int NOT NULL DEFAULT 0 COMMENT '库存',
  `sales`          int NOT NULL DEFAULT 0 COMMENT '销量',
  `status`         tinyint NOT NULL DEFAULT 1 COMMENT '1=启用 0=停用',
  `created_at`     datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at`     datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at`     datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_sku_code` (`sku_code`),
  KEY `idx_product_status` (`product_id`, `status`, `id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商品 SKU 表';

-- 二级分类
INSERT INTO `category` (`id`, `parent_id`, `name`, `icon`, `sort`, `status`) VALUES
(8, 1, '文学小说', '', 30, 1),
(9, 1, '技术编程', '', 20, 1),
(10, 1, '人文社科', '', 10, 1),
(26, 1, '教育考试', '', 5, 1),
(11, 2, '音频设备', '', 30, 1),
(12, 2, '电脑外设', '', 20, 1),
(13, 2, '智能穿戴', '', 10, 1),
(14, 3, '茶饮冲调', '', 30, 1),
(15, 3, '咖啡烘焙', '', 20, 1),
(16, 3, '休闲零食', '', 10, 1),
(17, 4, '香薰蜡烛', '', 30, 1),
(18, 4, '餐厨器物', '', 20, 1),
(19, 4, '收纳整理', '', 10, 1),
(20, 5, '护肤护理', '', 30, 1),
(21, 5, '彩妆香氛', '', 20, 1),
(22, 6, '上装', '', 30, 1),
(23, 6, '配饰', '', 20, 1),
(24, 7, '瑜伽健身', '', 30, 1),
(25, 7, '户外运动', '', 20, 1);
