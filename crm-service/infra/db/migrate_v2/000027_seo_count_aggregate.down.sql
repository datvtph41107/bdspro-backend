-- Keep seo_count_log; only drop aggregate table introduced by this migration.
DROP INDEX IF EXISTS idx_seo_count_updated_at;
DROP TABLE IF EXISTS seo_count;
