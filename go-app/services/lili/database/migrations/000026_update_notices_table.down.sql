ALTER TABLE `notices` ADD COLUMN  `published_at` DATETIME NOT NULL COMMENT '公開日時' AFTER `content`;
ALTER TABLE `notices` DROP COLUMN `title`;
ALTER TABLE `notices` DROP COLUMN `is_released`;
