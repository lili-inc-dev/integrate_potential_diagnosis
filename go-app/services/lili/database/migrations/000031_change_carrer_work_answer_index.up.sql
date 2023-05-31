ALTER TABLE `career_work_answers` DROP INDEX `answer_group_id`;

ALTER TABLE `career_work_answers` ADD UNIQUE `answer_group` (`answer_group_id`,`question_key`, `index`);
