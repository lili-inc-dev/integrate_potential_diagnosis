DROP TABLE IF EXISTS `personality_diagnose_choices`;

CREATE TABLE `personality_diagnose_choices` (
  `id` INTEGER UNSIGNED AUTO_INCREMENT,
  `name` VARCHAR(255) UNIQUE NOT NULL COMMENT '選択肢名（例：そう思う、どちらでもない）',
  `weight` TINYINT UNSIGNED NOT NULL COMMENT '回答の重み',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL,
  PRIMARY KEY (`id`)
);
