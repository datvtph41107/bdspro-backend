DROP INDEX IF EXISTS idx_pois_adr_search_trgm;
DROP INDEX IF EXISTS idx_pois_parcel_id;

ALTER TABLE pois
    DROP COLUMN IF EXISTS adr_search,
    DROP COLUMN IF EXISTS view_in_z,
    DROP COLUMN IF EXISTS parcel_id;
