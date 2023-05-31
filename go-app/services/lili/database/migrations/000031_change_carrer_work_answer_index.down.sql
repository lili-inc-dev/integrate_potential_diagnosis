ALTER TABLE `career_work_answers` DROP INDEX `answer_group`;

ALTER TABLE `career_work_answers` ADD UNIQUE `answer_group_id` (`answer_group_id`,`question_key`);
