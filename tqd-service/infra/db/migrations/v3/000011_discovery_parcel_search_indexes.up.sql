-- Accelerate discovery parcel search branches used by SearchParcels.
-- WHERE clauses in discovery/repository.go must stay aligned with these expressions.

CREATE INDEX IF NOT EXISTS idx_discovery_parcel_info_property_code_lower_nonempty
    ON qh_parcel_info (LOWER(property_code))
    WHERE property_code IS NOT NULL AND property_code <> '';

CREATE INDEX IF NOT EXISTS idx_discovery_parcel_info_adr_trgm_nonempty
    ON qh_parcel_info USING GIN (LOWER(adr_search) gin_trgm_ops)
    WHERE COALESCE(adr_search, '') <> '';

CREATE INDEX IF NOT EXISTS idx_discovery_parcel_info_address_trgm_nonempty
    ON qh_parcel_info USING GIN (LOWER(COALESCE(address_text, '')) gin_trgm_ops)
    WHERE COALESCE(address_text, '') <> '';
