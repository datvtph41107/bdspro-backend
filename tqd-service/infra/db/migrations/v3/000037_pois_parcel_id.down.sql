DROP INDEX IF EXISTS idx_pois_parcel_id;

ALTER TABLE pois
    DROP COLUMN IF EXISTS parcel_id;
