-- Optional ops script (not run by golang-migrate). Batch until "0".
-- Usage: psql "$DATABASE_URL" -f migrate/ops/backfill_parcel_spatial_cache.sql
-- Or: while psql ... -c "SELECT parcels_backfill_spatial_read_cache(500);" | grep -qv '^0$'; do sleep 1; done

-- Fast path: no preview_geojson, ST_Centroid (not PointOnSurface), 500 rows/round.
SELECT parcels_backfill_spatial_read_cache(500) AS rows_updated;
