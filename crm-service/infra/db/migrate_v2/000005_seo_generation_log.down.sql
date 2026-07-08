-- ============================================================
-- 000005_seo_generation_log.down.sql
-- This table is introduced by this migration, so rollback drops it.
-- ============================================================

DROP INDEX IF EXISTS idx_seo_generation_log_created_at_active;
DROP INDEX IF EXISTS idx_seo_generation_log_status_active;
DROP INDEX IF EXISTS idx_seo_generation_log_job_id_active;
DROP INDEX IF EXISTS idx_seo_generation_log_domain_active;

ALTER TABLE IF EXISTS seo_generation_log DROP CONSTRAINT IF EXISTS fk_seo_generation_log_domain;
ALTER TABLE IF EXISTS seo_generation_log DROP CONSTRAINT IF EXISTS chk_seo_generation_log_duration;
ALTER TABLE IF EXISTS seo_generation_log DROP CONSTRAINT IF EXISTS chk_seo_generation_log_status;

DROP TABLE IF EXISTS seo_generation_log;