DROP TABLE IF EXISTS `service_triggers`;

CREATE TABLE `service_triggers` (
  `id` TINYINT UNSIGNED AUTO_INCREMENT,
  `name` VARCHAR(255) UNIQUE NOT NULL COMMENT '例：インスタグラム、その他',
  PRIMARY KEY (`id`)
);
