DROP TABLE IF EXISTS `inactive_users`;

CREATE TABLE `inactive_users` (
  `id` CHAR(26) COMMENT 'ULID',
  `line_id` VARCHAR(50) NOT NULL COMMENT 'LINEユーザID',
  `type_id` TINYINT UNSIGNED NOT NULL COMMENT 'ユーザ種別',
  `firebase_uid` VARCHAR(130) UNIQUE NOT NULL COMMENT 'firebase user id',
  `name` VARCHAR(255) NOT NULL COMMENT '氏名',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL,
  PRIMARY KEY (`id`)
);

ALTER TABLE `inactive_users` ADD FOREIGN KEY (`line_id`) REFERENCES `line_accounts` (`line_id`);
ALTER TABLE `inactive_users` ADD FOREIGN KEY (`type_id`) REFERENCES `user_types` (`id`);
