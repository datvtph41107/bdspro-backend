-- Discovery / related-entity / map-target đọc pois.geom (PostGIS Point 4326).
-- Bảng GORM hiện chỉ có latitude/longitude — bổ sung geom, backfill, giữ đồng bộ.

ALTER TABLE pois
    ADD COLUMN IF NOT EXISTS geom geometry(Point, 4326);

UPDATE pois
SET geom = ST_SetSRID(
        ST_MakePoint(longitude::double precision, latitude::double precision),
        4326
    )
WHERE geom IS NULL
  AND latitude IS NOT NULL
  AND longitude IS NOT NULL
  AND latitude BETWEEN -90 AND 90
  AND longitude BETWEEN -180 AND 180
  AND NOT (latitude = 0 AND longitude = 0);

CREATE INDEX IF NOT EXISTS idx_pois_geom
    ON pois USING GIST (geom)
    WHERE deleted_at IS NULL AND geom IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_discovery_pois_name_trgm
    ON pois USING GIN (LOWER(name) gin_trgm_ops)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_discovery_pois_address_trgm
    ON pois USING GIN (LOWER(COALESCE(address, '')) gin_trgm_ops)
    WHERE deleted_at IS NULL;

CREATE OR REPLACE FUNCTION pois_sync_geom_from_latlng()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
    IF NEW.latitude IS NOT NULL
       AND NEW.longitude IS NOT NULL
       AND NEW.latitude BETWEEN -90 AND 90
       AND NEW.longitude BETWEEN -180 AND 180
       AND NOT (NEW.latitude = 0 AND NEW.longitude = 0)
    THEN
        NEW.geom := ST_SetSRID(
            ST_MakePoint(NEW.longitude::double precision, NEW.latitude::double precision),
            4326
        );
    ELSIF NEW.geom IS NOT NULL THEN
        NEW.latitude := ST_Y(NEW.geom);
        NEW.longitude := ST_X(NEW.geom);
    END IF;
    RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS trg_pois_sync_geom ON pois;
CREATE TRIGGER trg_pois_sync_geom
    BEFORE INSERT OR UPDATE OF latitude, longitude, geom
    ON pois
    FOR EACH ROW
    EXECUTE FUNCTION pois_sync_geom_from_latlng();
