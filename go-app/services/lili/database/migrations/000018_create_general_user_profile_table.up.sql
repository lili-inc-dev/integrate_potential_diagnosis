DROP TABLE IF EXISTS `general_user_profiles`;

CREATE TABLE `general_user_profiles` (
  `id` CHAR(26) COMMENT 'ULID',
  `user_id` CHAR(26) UNIQUE NOT NULL,
  `nickname` VARCHAR(255) NOT NULL COMMENT 'ニックネーム',
  `university` VARCHAR(255) NOT NULL COMMENT '大学',
  `faculty` VARCHAR(255) NOT NULL COMMENT '学部',
  `department` VARCHAR(255) NOT NULL COMMENT '学科',
  `graduation_year` INTEGER UNSIGNED NOT NULL COMMENT '大学卒業年度',
  `service_trigger_id` TINYINT UNSIGNED NOT NULL COMMENT 'サービスを知ったきっかけ',
  `service_trigger_detail` VARCHAR(255) DEFAULT null COMMENT 'サービスを知ったきっかけ（自由記述）',
  `introducer` VARCHAR(255) DEFAULT NULL COMMENT '紹介者',
  `desired_annual_income_id` TINYINT UNSIGNED NOT NULL COMMENT '希望の年収',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL,
  PRIMARY KEY (`id`)
);

ALTER TABLE `general_user_profiles` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `general_user_profiles` ADD FOREIGN KEY (`service_trigger_id`) REFERENCES `service_triggers` (`id`);

ALTER TABLE `general_user_profiles` ADD FOREIGN KEY (`desired_annual_income_id`) REFERENCES `desired_annual_incomes` (`id`);
