-- 000003_pro_ai_jobs.up.sql
-- Bảng job Pro AI: type Đồ án gắn planning_project_id; Tiến trình AI list từ đây.

-- These source-tracking fields historically came from GORM AutoMigrate. They
-- are dependencies of the backfill below, so SQL migration mode must own them
-- before reading existing planning projects.
ALTER TABLE qh_planning_projects
    ADD COLUMN IF NOT EXISTS source_folder_name VARCHAR(255),
    ADD COLUMN IF NOT EXISTS source_folder_path VARCHAR(500),
    ADD COLUMN IF NOT EXISTS process_status INT4 NOT NULL DEFAULT 10;

CREATE TABLE IF NOT EXISTS pro_ai_jobs (
    id                   BIGSERIAL PRIMARY KEY,
    job_type             INT4 NOT NULL DEFAULT 10,
    name                 VARCHAR(255) NOT NULL DEFAULT '',
    source_folder_name   VARCHAR(255),
    process_status       INT4 NOT NULL DEFAULT 10,
    planning_project_id  INT8,
    classify_error       TEXT,
    metadata             JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at           TIMESTAMPTZ,
    updated_at           TIMESTAMPTZ,
    deleted_at           TIMESTAMPTZ,
    created_by           INT8,
    updated_by           INT8
);

CREATE INDEX IF NOT EXISTS idx_pro_ai_jobs_job_type ON pro_ai_jobs(job_type);
CREATE INDEX IF NOT EXISTS idx_pro_ai_jobs_process_status ON pro_ai_jobs(process_status);
CREATE INDEX IF NOT EXISTS idx_pro_ai_jobs_planning_project_id ON pro_ai_jobs(planning_project_id);
CREATE INDEX IF NOT EXISTS idx_pro_ai_jobs_deleted_at ON pro_ai_jobs(deleted_at);

-- Backfill: mỗi đồ án đã upload folder → 1 job type Đồ án (10)
INSERT INTO pro_ai_jobs (
    job_type, name, source_folder_name, process_status, planning_project_id, metadata, created_at, updated_at
)
SELECT
    10,
    COALESCE(NULLIF(p.name, ''), COALESCE(p.source_folder_name, p.code, 'Đồ án')),
    p.source_folder_name,
    COALESCE(p.process_status, 10),
    p.id,
    '{}'::jsonb,
    COALESCE(p.created_at, NOW()),
    COALESCE(p.updated_at, NOW())
FROM qh_planning_projects p
WHERE p.deleted_at IS NULL
  AND NOT EXISTS (
      SELECT 1 FROM pro_ai_jobs j
      WHERE j.planning_project_id = p.id AND j.deleted_at IS NULL
  );
