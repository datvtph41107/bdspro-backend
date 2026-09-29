-- ============================================================
-- 000002_seo_links.down.sql
-- Safe rollback. Keeps legacy table if it existed, removes new constraints/indexes/columns.
-- ============================================================

DROP INDEX IF EXISTS idx_seo_links_pair_type_active;
DROP INDEX IF EXISTS idx_seo_links_priority_active;
DROP INDEX IF EXISTS idx_seo_links_type_active;
DROP INDEX IF EXISTS idx_seo_links_child_active;
DROP INDEX IF EXISTS idx_seo_links_parent_active;

ALTER TABLE IF EXISTS seo_links DROP CONSTRAINT IF EXISTS fk_seo_links_child_seo;
ALTER TABLE IF EXISTS seo_links DROP CONSTRAINT IF EXISTS fk_seo_links_parent_seo;
ALTER TABLE IF EXISTS seo_links DROP CONSTRAINT IF EXISTS chk_seo_links_priority;
ALTER TABLE IF EXISTS seo_links DROP CONSTRAINT IF EXISTS chk_seo_links_link_type;
ALTER TABLE IF EXISTS seo_links DROP CONSTRAINT IF EXISTS chk_seo_links_not_self;

ALTER TABLE IF EXISTS seo_links
    DROP COLUMN IF EXISTS link_type,
    DROP COLUMN IF EXISTS priority;