-- db_v4 picked up CRM-shaped columns (seo_id/clicked/timestamp, seo_id/count)
-- on top of the 000037 schema. They are NOT NULL without defaults, so every
-- insert written against seo_domain_id fails. modules/seo owns these tables.
ALTER TABLE seo_count_log
    DROP COLUMN IF EXISTS seo_id,
    DROP COLUMN IF EXISTS clicked,
    DROP COLUMN IF EXISTS "timestamp";

ALTER TABLE seo_count
    DROP COLUMN IF EXISTS seo_id,
    DROP COLUMN IF EXISTS count;
