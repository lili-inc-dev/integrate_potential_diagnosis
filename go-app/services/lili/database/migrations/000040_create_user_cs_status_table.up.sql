DROP TABLE IF EXISTS `user_cs_statuses`;

CREATE TABLE `user_cs_statuses` (
  `id` CHAR(26),
  `user_id` CHAR(26) UNIQUE NOT NULL,
  `status` ENUM ('register_only', 'diagnosing', 'normal') NOT NULL DEFAULT "register_only" COMMENT 'CSステータス',
  PRIMARY KEY (`id`)
);

ALTER TABLE `user_cs_statuses` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);
