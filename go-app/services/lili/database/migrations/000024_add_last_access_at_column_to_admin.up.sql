ALTER TABLE admins ADD COLUMN `last_access_at` DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL COMMENT '最終アクセス日' AFTER `is_disabled`;
