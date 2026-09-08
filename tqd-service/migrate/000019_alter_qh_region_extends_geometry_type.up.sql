ALTER TABLE IF EXISTS qh_region_extends
    ALTER COLUMN geometry TYPE geometry(Geometry, 4326)
    USING ST_SetSRID(geometry, 4326)::geometry(Geometry, 4326);
