-- ============================================================
-- 000007_seo_lifecycle_render_fields.down.sql
-- Purpose:
--   Roll back explicit lifecycle/render fields.
--   This keeps seo_generation_log table itself because it may contain
--   operational history. Only newly added columns/indexes are removed.
-- ============================================================

BEGIN;

DROP INDEX IF EXISTS idx_seo_generation_log_static_hash;
DROP INDEX IF EXISTS idx_seo_generation_log_trigger_type;
DROP INDEX IF EXISTS idx_seo_generation_log_status;
DROP INDEX IF EXISTS idx_seo_generation_log_domain_id;

ALTER TABLE seo_generation_log
    DROP COLUMN IF EXISTS metadata,
    DROP COLUMN IF EXISTS static_html_hash,
    DROP COLUMN IF EXISTS static_html_path,
    DROP COLUMN IF EXISTS trigger_type;

DROP INDEX IF EXISTS idx_seo_domain_render_queue;
DROP INDEX IF EXISTS idx_seo_domain_page_render_status;
DROP INDEX IF EXISTS idx_seo_domain_template_key;
DROP INDEX IF EXISTS idx_seo_domain_render_status;
DROP INDEX IF EXISTS idx_seo_domain_page_status;

ALTER TABLE seo_domain
    DROP CONSTRAINT IF EXISTS chk_seo_domain_render_status,
    DROP CONSTRAINT IF EXISTS chk_seo_domain_page_status;

ALTER TABLE seo_domain
    DROP COLUMN IF EXISTS template_version,
    DROP COLUMN IF EXISTS template_key,
    DROP COLUMN IF EXISTS last_render_error,
    DROP COLUMN IF EXISTS rendered_html,
    DROP COLUMN IF EXISTS render_status,
    DROP COLUMN IF EXISTS page_status;

COMMIT;