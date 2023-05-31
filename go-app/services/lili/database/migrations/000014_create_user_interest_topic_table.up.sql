DROP TABLE IF EXISTS `user_interest_topics`;

CREATE TABLE `user_interest_topics` (
  `id` BIGINT UNSIGNED AUTO_INCREMENT,
  `topic_id` INTEGER UNSIGNED NOT NULL,
  `user_id` CHAR(26) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY (`topic_id`, `user_id`)
);

ALTER TABLE `user_interest_topics` ADD FOREIGN KEY (`topic_id`) REFERENCES `interest_topics` (`id`);

ALTER TABLE `user_interest_topics` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);
