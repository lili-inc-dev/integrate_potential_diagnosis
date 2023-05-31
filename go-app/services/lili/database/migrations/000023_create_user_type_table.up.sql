DROP TABLE IF EXISTS `user_types`;

CREATE TABLE `user_types` (
  `id` TINYINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  `name` VARCHAR(255) UNIQUE NOT NULL COMMENT '学生 or 社会人'
);

ALTER TABLE `users` ADD FOREIGN KEY (`type_id`) REFERENCES `user_types` (`id`);
