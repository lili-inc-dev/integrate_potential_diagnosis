DROP TABLE IF EXISTS `market_value_diagnose_answers`;

CREATE TABLE `market_value_diagnose_answers` (
  `id` CHAR(26) COMMENT 'ULID',
  `answer_group_id` CHAR(26) NOT NULL COMMENT 'ULID 回答グループごとに作成',
  `question_id` INTEGER UNSIGNED NOT NULL COMMENT '設問ID',
  `choice_id` INTEGER UNSIGNED NOT NULL COMMENT '選択した選択肢のID',
  `user_id` CHAR(26) NOT NULL,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE `answer_unique` (`answer_group_id`, `question_id`)
);

ALTER TABLE `market_value_diagnose_answers` ADD FOREIGN KEY (`question_id`) REFERENCES `market_value_diagnose_questions` (`id`);

ALTER TABLE `market_value_diagnose_answers` ADD FOREIGN KEY (`choice_id`) REFERENCES `market_value_diagnose_choices` (`id`);

ALTER TABLE `market_value_diagnose_answers` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);
