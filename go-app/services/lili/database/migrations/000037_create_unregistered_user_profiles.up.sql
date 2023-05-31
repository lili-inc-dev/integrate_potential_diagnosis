DROP TABLE IF EXISTS `unregistered_user_profiles`;

CREATE TABLE `unregistered_user_profiles` (
  `id` CHAR(26) PRIMARY KEY COMMENT 'ULID',
  `user_id` CHAR(26) UNIQUE NOT NULL,
  `reason_id` TINYINT UNSIGNED NOT NULL,
  `line_id` VARCHAR(50) NOT NULL,
  `email` VARCHAR(255) NOT NULL,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL
);

ALTER TABLE `unregistered_user_profiles` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `unregistered_user_profiles` ADD FOREIGN KEY (`reason_id`) REFERENCES `unregistration_reasons` (`id`);

ALTER TABLE `unregistered_user_profiles` ADD FOREIGN KEY (`line_id`) REFERENCES `line_accounts` (`line_id`);