DROP TABLE IF EXISTS `interest_topics`;

CREATE TABLE `interest_topics` (
  `id` INTEGER UNSIGNED AUTO_INCREMENT,
  `name` VARCHAR(255) UNIQUE NOT NULL COMMENT '項目名（例：ビジネス、ファッション）',
  `display_order` INTEGER UNSIGNED UNIQUE NOT NULL COMMENT '表示順番',
  PRIMARY KEY (`id`)
);
