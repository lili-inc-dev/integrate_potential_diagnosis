DROP TABLE IF EXISTS `desired_annual_incomes`;

CREATE TABLE `desired_annual_incomes` (
  `id` TINYINT UNSIGNED AUTO_INCREMENT,
  `value` VARCHAR(255) UNIQUE NOT NULL COMMENT '例：〜300万円、300〜500万円',
  PRIMARY KEY (`id`)
);
