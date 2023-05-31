CREATE TABLE `users` (
  `id` CHAR(26) COMMENT 'ULID',
  `line_id` VARCHAR(50) UNIQUE NOT NULL COMMENT 'LINEユーザID',
  `type_id` TINYINT UNSIGNED NOT NULL COMMENT 'ユーザ種別',
  `firebase_uid` VARCHAR(130) UNIQUE NOT NULL COMMENT 'firebase user id',
  `email` VARCHAR(255) UNIQUE NOT NULL,
  `name` VARCHAR(255) NOT NULL COMMENT '氏名',
  `name_kana` VARCHAR(255) DEFAULT null COMMENT '氏名フリガナ',
  `gender_id` TINYINT UNSIGNED NOT NULL COMMENT '性別',
  `age` TINYINT NOT NULL COMMENT '年齢',
  `birthday` DATE NOT NULL COMMENT '生年月日',
  `phone_number` VARCHAR(255) NOT NULL COMMENT '電話番号',
  `memo` TEXT DEFAULT null COMMENT 'メモ（管理画面で利用）',
  `is_disabled` BOOLEAN NOT NULL DEFAULT false COMMENT '無効状態かどうか',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL,
  PRIMARY KEY (`id`)
);

ALTER TABLE `users` ADD FOREIGN KEY (`line_id`) REFERENCES `line_accounts` (`line_id`);
