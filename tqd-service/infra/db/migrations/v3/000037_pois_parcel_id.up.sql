ALTER TABLE pois
    ADD COLUMN IF NOT EXISTS parcel_id BIGINT;

CREATE INDEX IF NOT EXISTS idx_pois_parcel_id
    ON pois(parcel_id)
    WHERE parcel_id IS NOT NULL AND parcel_id > 0;
