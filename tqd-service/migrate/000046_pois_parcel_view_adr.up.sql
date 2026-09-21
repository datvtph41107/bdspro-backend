-- Port missing POI columns from legacy infra/db/migrations/v3 into canonical migrate/.
-- Domain: parcel_id, view_in_z, adr_search (see internal/domain/poi.go).

CREATE EXTENSION IF NOT EXISTS unaccent;
CREATE EXTENSION IF NOT EXISTS pg_trgm;

ALTER TABLE pois
    ADD COLUMN IF NOT EXISTS parcel_id BIGINT,
    ADD COLUMN IF NOT EXISTS view_in_z INTEGER NOT NULL DEFAULT 16,
    ADD COLUMN IF NOT EXISTS adr_search TEXT;

UPDATE pois
SET view_in_z = 16
WHERE view_in_z IS NULL OR view_in_z = 0;

UPDATE pois
SET adr_search = TRIM(
    REGEXP_REPLACE(
        LOWER(UNACCENT(COALESCE(address, ''))),
        '[^a-z0-9]+',
        ' ',
        'g'
    )
)
WHERE COALESCE(adr_search, '') = ''
  AND COALESCE(address, '') <> '';

CREATE INDEX IF NOT EXISTS idx_pois_parcel_id
    ON pois(parcel_id)
    WHERE parcel_id IS NOT NULL AND parcel_id > 0;

CREATE INDEX IF NOT EXISTS idx_pois_adr_search_trgm
    ON pois USING GIN (LOWER(adr_search) gin_trgm_ops)
    WHERE COALESCE(adr_search, '') <> '';
