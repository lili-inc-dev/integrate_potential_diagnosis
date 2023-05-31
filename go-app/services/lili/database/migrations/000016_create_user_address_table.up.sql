DROP TABLE IF EXISTS `user_addresses`;

CREATE TABLE `user_addresses` (
  `id` CHAR(26) COMMENT 'ULID',
  `user_id` CHAR(26) UNIQUE NOT NULL,
  `postal_code` VARCHAR(255) NOT NULL COMMENT '郵便番号',
  `address` VARCHAR(255) NOT NULL COMMENT '住所',
  PRIMARY KEY (`id`)
);

ALTER TABLE `user_addresses` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);
