ALTER TABLE reports DROP COLUMN IF EXISTS closed_at;
ALTER TABLE reports DROP COLUMN IF EXISTS support_ticket_id;
ALTER TABLE reports DROP COLUMN IF EXISTS linkage;
ALTER TABLE reports DROP COLUMN IF EXISTS qa_status;
ALTER TABLE reports DROP COLUMN IF EXISTS assignee_id;
ALTER TABLE reports DROP COLUMN IF EXISTS severity;
ALTER TABLE reports DROP COLUMN IF EXISTS title;
-- keep problem_report / description (domain fields)
