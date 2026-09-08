-- ============================================================
-- 000003_seo_relative.down.sql
-- Safe rollback. Keeps legacy table if it existed, removes new constraints/indexes/columns.
-- ============================================================

DROP INDEX IF EXISTS idx_seo_relative_pair_type_active;
DROP INDEX IF EXISTS idx_seo_relative_priority_active;
DROP INDEX IF EXISTS idx_seo_relative_type_active;
DROP INDEX IF EXISTS idx_seo_relative_child_active;
DROP INDEX IF EXISTS idx_seo_relative_parent_active;

ALTER TABLE IF EXISTS seo_relative DROP CONSTRAINT IF EXISTS fk_seo_relative_child_seo;
ALTER TABLE IF EXISTS seo_relative DROP CONSTRAINT IF EXISTS fk_seo_relative_parent_seo;
ALTER TABLE IF EXISTS seo_relative DROP CONSTRAINT IF EXISTS chk_seo_relative_priority;
ALTER TABLE IF EXISTS seo_relative DROP CONSTRAINT IF EXISTS chk_seo_relative_type;
ALTER TABLE IF EXISTS seo_relative DROP CONSTRAINT IF EXISTS chk_seo_relative_not_self;

ALTER TABLE IF EXISTS seo_relative
    DROP COLUMN IF EXISTS relation_type,
    DROP COLUMN IF EXISTS priority;