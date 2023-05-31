DROP TABLE IF EXISTS `personality_diagnose_questions`;

CREATE TABLE `personality_diagnose_questions` (
  `id` INTEGER UNSIGNED AUTO_INCREMENT,
  `personality_id` TINYINT UNSIGNED NOT NULL COMMENT 'パーソナリティ',
  `index` INTEGER UNIQUE NOT NULL COMMENT '出題順',
  `content` VARCHAR(255) UNIQUE NOT NULL COMMENT '設問内容',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL,
  PRIMARY KEY (`id`)
);

ALTER TABLE `personality_diagnose_questions` ADD FOREIGN KEY (`personality_id`) REFERENCES `personalities` (`id`);
