CREATE TABLE `email_authentication_codes` (
  `id` CHAR(26) COMMENT 'ULID',
  `inactive_user_id` CHAR(26) NOT NULL,
  `code_hash` BINARY(60) NOT NULL COMMENT '認証コードのハッシュ値',
  `attempt_count` TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '試行回数',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL,
   PRIMARY KEY (`id`)
);

ALTER TABLE `email_authentication_codes` ADD FOREIGN KEY (`inactive_user_id`) REFERENCES `inactive_users` (`id`);
