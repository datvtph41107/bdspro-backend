ALTER TABLE province_v2
    ADD COLUMN IF NOT EXISTS seo_id BIGINT;

ALTER TABLE ward_v2
    ADD COLUMN IF NOT EXISTS seo_id BIGINT;

CREATE INDEX IF NOT EXISTS idx_province_v2_seo_id
    ON province_v2(seo_id)
    WHERE seo_id IS NOT NULL AND seo_id > 0;

CREATE INDEX IF NOT EXISTS idx_ward_v2_seo_id
    ON ward_v2(seo_id)
    WHERE seo_id IS NOT NULL AND seo_id > 0;
