DROP TABLE IF EXISTS `market_values`;

CREATE TABLE `market_values` (
  `id` TINYINT UNSIGNED AUTO_INCREMENT,
  `display_order` TINYINT UNSIGNED UNIQUE NOT NULL COMMENT '何ページ目に表示されるか',
  `name` VARCHAR(255) UNIQUE NOT NULL COMMENT '市場価値診断の項目名',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL,
  PRIMARY KEY (`id`)
);
