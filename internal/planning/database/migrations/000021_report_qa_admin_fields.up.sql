CREATE TABLE IF NOT EXISTS reports (
    id BIGSERIAL PRIMARY KEY
);

ALTER TABLE reports ADD COLUMN IF NOT EXISTS problem_report INT NOT NULL DEFAULT 10;
ALTER TABLE reports ADD COLUMN IF NOT EXISTS description TEXT;
ALTER TABLE reports ADD COLUMN IF NOT EXISTS title VARCHAR(255) DEFAULT '';
ALTER TABLE reports ADD COLUMN IF NOT EXISTS severity INT NOT NULL DEFAULT 20;
ALTER TABLE reports ADD COLUMN IF NOT EXISTS assignee_id BIGINT;
ALTER TABLE reports ADD COLUMN IF NOT EXISTS qa_status INT NOT NULL DEFAULT 10;
ALTER TABLE reports ADD COLUMN IF NOT EXISTS linkage JSONB DEFAULT '{}'::jsonb;
ALTER TABLE reports ADD COLUMN IF NOT EXISTS support_ticket_id BIGINT;
ALTER TABLE reports ADD COLUMN IF NOT EXISTS closed_at TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS idx_reports_qa_status ON reports(qa_status) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_reports_assignee ON reports(assignee_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_reports_problem ON reports(problem_report) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_reports_support_ticket ON reports(support_ticket_id) WHERE deleted_at IS NULL;
