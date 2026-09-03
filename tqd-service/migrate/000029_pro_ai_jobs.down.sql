-- 000003_pro_ai_jobs.down.sql

DROP INDEX IF EXISTS idx_pro_ai_jobs_deleted_at;
DROP INDEX IF EXISTS idx_pro_ai_jobs_planning_project_id;
DROP INDEX IF EXISTS idx_pro_ai_jobs_process_status;
DROP INDEX IF EXISTS idx_pro_ai_jobs_job_type;
DROP TABLE IF EXISTS pro_ai_jobs;
