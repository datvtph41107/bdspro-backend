-- ============================================================
-- 000001_seo_domain_registry.down.sql
-- Safe rollback for projects with an existing seo_domain table.
-- It removes indexes/constraints and the newly introduced SEO registry columns,
-- but does not drop seo_domain itself.
-- ============================================================

DROP INDEX IF EXISTS idx_seo_domain_updated_at_active;
DROP INDEX IF EXISTS idx_seo_domain_generated_at_active;
DROP INDEX IF EXISTS idx_seo_domain_source_updated_at_active;
DROP INDEX IF EXISTS idx_seo_domain_need_generate_active;
DROP INDEX IF EXISTS idx_seo_domain_published_sitemap_active;
DROP INDEX IF EXISTS idx_seo_domain_classify_active;
DROP INDEX IF EXISTS idx_seo_domain_scope_active;
DROP INDEX IF EXISTS idx_seo_domain_ref_active;
DROP INDEX IF EXISTS idx_seo_domain_slug_active;
DROP INDEX IF EXISTS idx_seo_domain_canonical_url_active;
DROP INDEX IF EXISTS idx_seo_domain_ref_unique_active;

ALTER TABLE IF EXISTS seo_domain DROP CONSTRAINT IF EXISTS chk_seo_domain_changefreq;
ALTER TABLE IF EXISTS seo_domain DROP CONSTRAINT IF EXISTS chk_seo_domain_sitemap_priority;
ALTER TABLE IF EXISTS seo_domain DROP CONSTRAINT IF EXISTS chk_seo_domain_classify;

-- Keep identity/source columns that may have existed before:
-- origin_url, canonical_url, ref_type, ref_id, metadata, created_at, updated_at, deleted_at.
-- This makes rollback safer for legacy systems.
ALTER TABLE IF EXISTS seo_domain
    DROP COLUMN IF EXISTS slug,
    DROP COLUMN IF EXISTS scope,
    DROP COLUMN IF EXISTS classify,
    DROP COLUMN IF EXISTS title,
    DROP COLUMN IF EXISTS description,
    DROP COLUMN IF EXISTS content,
    DROP COLUMN IF EXISTS summary,
    DROP COLUMN IF EXISTS published,
    DROP COLUMN IF EXISTS published_at,
    DROP COLUMN IF EXISTS is_site_map,
    DROP COLUMN IF EXISTS is_index,
    DROP COLUMN IF EXISTS is_robot,
    DROP COLUMN IF EXISTS site_map_lasted_at,
    DROP COLUMN IF EXISTS sitemap_priority,
    DROP COLUMN IF EXISTS sitemap_change_freq,
    DROP COLUMN IF EXISTS need_generate,
    DROP COLUMN IF EXISTS generated_at,
    DROP COLUMN IF EXISTS source_updated_at,
    DROP COLUMN IF EXISTS static_html_path,
    DROP COLUMN IF EXISTS static_html_hash,
    DROP COLUMN IF EXISTS deep_link,
    DROP COLUMN IF EXISTS note;