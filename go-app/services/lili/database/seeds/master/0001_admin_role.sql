INSERT INTO `admin_roles` (
  `name`,
  `description`,
  `admin_browsable`,
  `admin_editable`,
  `user_browsable`,
  `user_editable`,
  `company_browsable`,
  `company_editable`,
  `project_browsable`,
  `project_editable`,
  `project_disclosable`,
  `project_comment_browsable`,
  `project_comment_editable`,
  `project_comment_postable`,
  `notice_browsable`,
  `notice_editable`
) VALUES
 ('システム管理者', null, true, true, true, true, true, true, true, true, true, true, true, true, true, true),
 ('運営', '運営担当者用のロール', true, false, true, true, true, true, true, true, true, true, true, false, true, true);
--  ('講師', false, false, false, false, false, false, true, false, false, true, false, true, false, false);
