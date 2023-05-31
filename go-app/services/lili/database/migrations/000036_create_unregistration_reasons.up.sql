DROP TABLE IF EXISTS `unregistration_reasons`;

CREATE TABLE `unregistration_reasons` (
  `id` TINYINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  `content` VARCHAR(255) NOT NULL COMMENT '退会理由内容',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL
);