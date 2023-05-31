DROP TABLE IF EXISTS `market_value_diagnose_questions`;

CREATE TABLE `market_value_diagnose_questions` (
  `id` INTEGER UNSIGNED AUTO_INCREMENT,
  `market_value_id` TINYINT UNSIGNED NOT NULL COMMENT '市場価値の項目ID',
  `index` INTEGER UNIQUE NOT NULL COMMENT '出題順',
  `content` VARCHAR(255) UNIQUE NOT NULL COMMENT '設問内容',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL,
  PRIMARY KEY (`id`)
);

ALTER TABLE `market_value_diagnose_questions` ADD FOREIGN KEY (`market_value_id`) REFERENCES `market_values` (`id`);
