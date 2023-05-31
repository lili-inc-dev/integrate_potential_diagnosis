-- 先に外キー制約を外す
ALTER TABLE `inactive_users` DROP FOREIGN KEY `inactive_users_ibfk_1`;
-- ユニークキー制約を外す
ALTER TABLE `inactive_users` DROP INDEX `line_id_unique`;
-- 外キー制約を付け直す
ALTER TABLE `inactive_users` ADD CONSTRAINT `inactive_users_ibfk_1` FOREIGN KEY (`line_id`) REFERENCES `line_accounts` (`line_id`);
