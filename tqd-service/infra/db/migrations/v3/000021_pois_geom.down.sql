DROP TRIGGER IF EXISTS trg_pois_sync_geom ON pois;
DROP FUNCTION IF EXISTS pois_sync_geom_from_latlng();

DROP INDEX IF EXISTS idx_discovery_pois_address_trgm;
DROP INDEX IF EXISTS idx_discovery_pois_name_trgm;
DROP INDEX IF EXISTS idx_pois_geom;

ALTER TABLE pois
    DROP COLUMN IF EXISTS geom;
