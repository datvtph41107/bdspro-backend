-- ST_DWithin(geom::geography, ...) needs a geography GIST, not geometry-only.
CREATE INDEX IF NOT EXISTS idx_pois_geom_geography
    ON pois USING GIST ((geom::geography));
