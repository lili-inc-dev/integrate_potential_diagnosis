DROP TABLE IF EXISTS `admin_roles`;

CREATE TABLE `admin_roles` (
  `id` INTEGER UNSIGNED AUTO_INCREMENT,
  `name` VARCHAR(255) UNIQUE NOT NULL,
  `admin_browsable` BOOLEAN NOT NULL DEFAULT false,
  `admin_editable` BOOLEAN NOT NULL DEFAULT false,
  `user_browsable` BOOLEAN NOT NULL DEFAULT false,
  `user_editable` BOOLEAN NOT NULL DEFAULT false,
  `company_browsable` BOOLEAN NOT NULL DEFAULT false,
  `company_editable` BOOLEAN NOT NULL DEFAULT false,
  `project_browsable` BOOLEAN NOT NULL DEFAULT false,
  `project_editable` BOOLEAN NOT NULL DEFAULT false,
  `project_disclosable` BOOLEAN NOT NULL DEFAULT false COMMENT 'プロジェクト公開権限',
  `project_comment_browsable` BOOLEAN NOT NULL DEFAULT false,
  `project_comment_editable` BOOLEAN NOT NULL DEFAULT false,
  `project_comment_postable` BOOLEAN NOT NULL DEFAULT false COMMENT 'プロジェクトコメント投稿権限',
  `notice_browsable` BOOLEAN NOT NULL DEFAULT false,
  `notice_editable` BOOLEAN NOT NULL DEFAULT false,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL,
  PRIMARY KEY (`id`)
);

