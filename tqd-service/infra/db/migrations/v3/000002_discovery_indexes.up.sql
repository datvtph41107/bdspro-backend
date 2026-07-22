-- Discovery V1 indexes.
-- Source domain tables remain the source of truth; these indexes accelerate spatial identify and federated search.

CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX IF NOT EXISTS idx_discovery_parcels_geometry_active
    ON parcels USING GIST (geometry)
    WHERE deleted_at IS NULL AND geometry IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_discovery_parcel_info_parcel_id
    ON qh_parcel_info (parcel_id);

CREATE INDEX IF NOT EXISTS idx_discovery_parcel_info_property_code_lower
    ON qh_parcel_info (LOWER(property_code));

CREATE INDEX IF NOT EXISTS idx_discovery_parcel_info_map_land
    ON qh_parcel_info (map_number, land_number);

CREATE INDEX IF NOT EXISTS idx_discovery_parcel_info_adr_trgm
    ON qh_parcel_info USING GIN (LOWER(adr_search) gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_discovery_regions_geometry_active
    ON qh_regions USING GIST (geometry)
    WHERE deleted_at IS NULL AND status = 10 AND is_latest = TRUE;

CREATE INDEX IF NOT EXISTS idx_discovery_regions_display_name_trgm
    ON qh_regions USING GIN (LOWER(COALESCE(display_name, '')) gin_trgm_ops)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_discovery_regions_name_trgm
    ON qh_regions USING GIN (LOWER(COALESCE(name, '')) gin_trgm_ops)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_discovery_pois_name_trgm
    ON pois USING GIN (LOWER(name) gin_trgm_ops)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_discovery_pois_address_trgm
    ON pois USING GIN (LOWER(COALESCE(address, '')) gin_trgm_ops)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_discovery_ward_name_trgm
    ON ward_v2 USING GIN (LOWER(COALESCE(full_name, '')) gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_discovery_province_name_trgm
    ON province_v2 USING GIN (LOWER(COALESCE(full_name, '')) gin_trgm_ops);
