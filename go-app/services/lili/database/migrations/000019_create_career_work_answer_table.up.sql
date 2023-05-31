DROP TABLE IF EXISTS `career_work_answers`;

CREATE TABLE `career_work_answers` (
  `id` INTEGER UNSIGNED AUTO_INCREMENT,
  `answer_group_id` CHAR(26) NOT NULL COMMENT 'ULID 回答グループごとに作成',
  `question_key` CHAR(255) NOT NULL COMMENT '設問のkey(フロント側で設定)',
  `answer` TEXT NOT NULL COMMENT '答え',
  `user_id` CHAR(26) NOT NULL COMMENT 'ユーザーID',
  `index` INTEGER NOT NULL COMMENT '回答順',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE (`answer_group_id`, `question_key`)
);

ALTER TABLE `career_work_answers` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);
