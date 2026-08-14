DROP INDEX IF EXISTS idx_ward_v2_seo_id;
DROP INDEX IF EXISTS idx_province_v2_seo_id;

ALTER TABLE ward_v2
    DROP COLUMN IF EXISTS seo_id;

ALTER TABLE province_v2
    DROP COLUMN IF EXISTS seo_id;
