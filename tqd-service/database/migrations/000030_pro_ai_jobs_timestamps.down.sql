-- 000004_pro_ai_jobs_timestamps.down.sql

DROP INDEX IF EXISTS idx_pro_ai_jobs_approved_at;
DROP INDEX IF EXISTS idx_pro_ai_jobs_completed_at;

ALTER TABLE pro_ai_jobs
    DROP COLUMN IF EXISTS approved_at,
    DROP COLUMN IF EXISTS completed_at;
