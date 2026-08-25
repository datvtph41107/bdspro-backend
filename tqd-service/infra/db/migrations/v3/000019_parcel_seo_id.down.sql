DROP INDEX IF EXISTS idx_parcels_seo_id;

ALTER TABLE parcels
    DROP COLUMN IF EXISTS seo_id;
