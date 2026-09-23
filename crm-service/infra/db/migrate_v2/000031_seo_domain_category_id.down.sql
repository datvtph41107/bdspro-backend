BEGIN;

DROP INDEX IF EXISTS idx_seo_domain_category_id;

ALTER TABLE seo_domain DROP CONSTRAINT IF EXISTS fk_seo_domain_category;
ALTER TABLE seo_domain DROP COLUMN IF EXISTS category_id;

COMMIT;
