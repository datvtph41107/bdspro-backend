DROP INDEX IF EXISTS idx_discovery_parcel_info_property_code_lower_nonempty;
DROP INDEX IF EXISTS idx_discovery_parcel_info_adr_trgm_nonempty;

-- Restore non-partial GIN dropped by up migration (best-effort).
CREATE INDEX IF NOT EXISTS idx_discovery_parcel_info_adr_trgm
    ON qh_parcel_info USING GIN (LOWER(adr_search) gin_trgm_ops);
