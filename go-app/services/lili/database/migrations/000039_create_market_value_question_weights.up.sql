DROP TABLE IF EXISTS `market_value_diagnose_question_weights`;

CREATE TABLE `market_value_diagnose_question_weights` (
  `id` INTEGER UNSIGNED AUTO_INCREMENT,
  `question_id` INTEGER UNSIGNED NOT NULL,
  `choice_id` INTEGER UNSIGNED NOT NULL,
  `weight` TINYINT UNSIGNED NOT NULL COMMENT '回答の重み',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE `question_choice_unique` (`question_id`, `choice_id`)
);

ALTER TABLE `market_value_diagnose_question_weights` ADD FOREIGN KEY (`question_id`) REFERENCES `market_value_diagnose_questions` (`id`);

ALTER TABLE `market_value_diagnose_question_weights` ADD FOREIGN KEY (`choice_id`) REFERENCES `market_value_diagnose_choices` (`id`);
