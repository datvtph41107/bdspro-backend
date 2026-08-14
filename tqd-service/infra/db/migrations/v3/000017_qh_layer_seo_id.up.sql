ALTER TABLE qh_layers
    ADD COLUMN IF NOT EXISTS seo_id BIGINT;

CREATE INDEX IF NOT EXISTS idx_qh_layers_seo_id
    ON qh_layers(seo_id)
    WHERE seo_id IS NOT NULL AND seo_id > 0;
