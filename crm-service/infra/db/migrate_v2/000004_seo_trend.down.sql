-- ============================================================
-- 000004_seo_trend.down.sql
-- This table is introduced by this migration, so rollback drops it.
-- ============================================================

DROP INDEX IF EXISTS idx_seo_trend_domain_date_source_active;
DROP INDEX IF EXISTS idx_seo_trend_date_active;
DROP INDEX IF EXISTS idx_seo_trend_domain_active;

ALTER TABLE IF EXISTS seo_trend DROP CONSTRAINT IF EXISTS fk_seo_trend_domain;
ALTER TABLE IF EXISTS seo_trend DROP CONSTRAINT IF EXISTS chk_seo_trend_non_negative;

DROP TABLE IF EXISTS seo_trend;