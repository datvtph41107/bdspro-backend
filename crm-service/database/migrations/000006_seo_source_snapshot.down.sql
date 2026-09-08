-- ============================================================
-- 000006_seo_source_snapshot.down.sql
-- Roll back source-aware SEO binding fields.
-- ============================================================

DROP INDEX IF EXISTS idx_seo_domain_quality_score_active;
DROP INDEX IF EXISTS idx_seo_domain_ref_last_synced_active;
DROP INDEX IF EXISTS idx_seo_domain_ref_hash_active;
DROP INDEX IF EXISTS idx_seo_domain_source_status_active;
DROP INDEX IF EXISTS idx_seo_domain_ref_source_active;
DROP INDEX IF EXISTS uniq_seo_domain_active_ref_source;

ALTER TABLE IF EXISTS seo_domain DROP CONSTRAINT IF EXISTS chk_seo_domain_quality_score;
ALTER TABLE IF EXISTS seo_domain DROP CONSTRAINT IF EXISTS chk_seo_domain_source_status;

ALTER TABLE IF EXISTS seo_domain
    DROP COLUMN IF EXISTS quality_score,
    DROP COLUMN IF EXISTS ref_last_synced_at,
    DROP COLUMN IF EXISTS ref_hash,
    DROP COLUMN IF EXISTS ref_snapshot_json,
    DROP COLUMN IF EXISTS ref_missing,
    DROP COLUMN IF EXISTS source_status,
    DROP COLUMN IF EXISTS ref_url,
    DROP COLUMN IF EXISTS ref_label,
    DROP COLUMN IF EXISTS ref_source;