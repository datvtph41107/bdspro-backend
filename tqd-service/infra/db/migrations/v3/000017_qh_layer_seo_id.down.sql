DROP INDEX IF EXISTS idx_qh_layers_seo_id;

ALTER TABLE qh_layers
    DROP COLUMN IF EXISTS seo_id;
