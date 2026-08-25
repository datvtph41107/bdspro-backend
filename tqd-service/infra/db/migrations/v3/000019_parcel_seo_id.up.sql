ALTER TABLE parcels
    ADD COLUMN IF NOT EXISTS seo_id BIGINT;

CREATE INDEX IF NOT EXISTS idx_parcels_seo_id
    ON parcels(seo_id)
    WHERE seo_id IS NOT NULL AND seo_id > 0;
