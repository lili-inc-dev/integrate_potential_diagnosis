ALTER TABLE `notices` DROP COLUMN `published_at`;
ALTER TABLE `notices` ADD COLUMN `title` VARCHAR(255) NOT NULL COMMENT 'お知らせタイトル' AFTER `id`;
ALTER TABLE `notices` ADD COLUMN `is_released` BOOLEAN NOT NULL DEFAULT false COMMENT '公開状態かどうか' AFTER `content`;
