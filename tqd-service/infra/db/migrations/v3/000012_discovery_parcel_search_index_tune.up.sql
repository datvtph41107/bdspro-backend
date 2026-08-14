-- Tune discovery parcel address search indexes.
-- SearchParcels only fuzzy-matches LOWER(adr_search); drop redundant full GIN
-- that encourages expensive scans on empty/short patterns.

CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Full-table GIN from 000002 is larger and often preferred over the partial index.
DROP INDEX IF EXISTS idx_discovery_parcel_info_adr_trgm;

CREATE INDEX IF NOT EXISTS idx_discovery_parcel_info_adr_trgm_nonempty
    ON qh_parcel_info USING GIN (LOWER(adr_search) gin_trgm_ops)
    WHERE COALESCE(adr_search, '') <> '';

CREATE INDEX IF NOT EXISTS idx_discovery_parcel_info_property_code_lower_nonempty
    ON qh_parcel_info (LOWER(property_code))
    WHERE property_code IS NOT NULL AND property_code <> '';

CREATE INDEX IF NOT EXISTS idx_discovery_parcel_info_map_land
    ON qh_parcel_info (map_number, land_number);
