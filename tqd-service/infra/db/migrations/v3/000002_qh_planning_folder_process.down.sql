-- 000002_qh_planning_folder_process.down.sql

DROP INDEX IF EXISTS idx_qh_planning_documents_process_status;
DROP INDEX IF EXISTS idx_qh_planning_projects_process_status;

ALTER TABLE qh_planning_documents
    DROP COLUMN IF EXISTS classify_error,
    DROP COLUMN IF EXISTS process_status,
    DROP COLUMN IF EXISTS relative_path;

ALTER TABLE qh_planning_projects
    DROP COLUMN IF EXISTS process_status,
    DROP COLUMN IF EXISTS source_folder_path,
    DROP COLUMN IF EXISTS source_folder_name;
