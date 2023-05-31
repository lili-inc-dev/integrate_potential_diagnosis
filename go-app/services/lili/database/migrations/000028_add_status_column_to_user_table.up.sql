ALTER TABLE `users` ADD COLUMN `status` ENUM('registered', 'banned', 'unregistered') DEFAULT 'registered' COMMENT 'ステータス（例：registered、banned、unregistered）';
