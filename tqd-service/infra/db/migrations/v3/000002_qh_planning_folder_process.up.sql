-- 000002_qh_planning_folder_process.up.sql
-- Thêm cột hỗ trợ upload folder + job phân loại AI cho qh_planning_projects / qh_planning_documents.

ALTER TABLE qh_planning_projects
    ADD COLUMN IF NOT EXISTS source_folder_name VARCHAR(255),
    ADD COLUMN IF NOT EXISTS source_folder_path VARCHAR(500),
    ADD COLUMN IF NOT EXISTS process_status INT4 NOT NULL DEFAULT 10;

ALTER TABLE qh_planning_documents
    ADD COLUMN IF NOT EXISTS relative_path VARCHAR(500),
    ADD COLUMN IF NOT EXISTS process_status INT4 NOT NULL DEFAULT 10,
    ADD COLUMN IF NOT EXISTS classify_error TEXT;

CREATE INDEX IF NOT EXISTS idx_qh_planning_projects_process_status ON qh_planning_projects(process_status);
CREATE INDEX IF NOT EXISTS idx_qh_planning_documents_process_status ON qh_planning_documents(process_status);
