-- -- Migrate parcels → pois (thửa) — tối ưu bulk ~7M rows.
-- -- Chạy thủ công (DBeaver / psql). Không gắn auto-migrate.
-- --
-- -- Tối ưu:
-- --   1) tắt trigger geom (tránh ST_MakePoint từng row)
-- --   2) drop tạm index nặng (GIN/GIST/trgm), tạo lại sau
-- --   3) chỉ dùng parcels.lat/lon sẵn (bỏ ST_PointOnSurface khỏi đường nóng)
-- --   4) insert theo batch + COMMIT (PROCEDURE) để giảm WAL / dễ resume
-- --   5) backfill geom 1 lần sau insert
-- -- point_type = 'Thửa'
-- -- address: tạm lưu nguyên pi.address_text (không strip prefix ',')

-- CREATE EXTENSION IF NOT EXISTS unaccent;

-- ALTER TABLE pois
--     ADD COLUMN IF NOT EXISTS point_type VARCHAR(50);

-- -- Session tuning (chỉ ảnh hưởng connection hiện tại)
-- SET synchronous_commit = off;
-- SET maintenance_work_mem = '1GB';
-- SET work_mem = '256MB';

-- -- Sync sequence trước bulk insert (tránh 23505 pois_pkey)
-- SELECT setval(
--     pg_get_serial_sequence('pois', 'id'),
--     GREATEST(COALESCE((SELECT MAX(id) FROM pois), 1), 1),
--     true
-- );

-- -- Tắt trigger sync geom (bỏ qua nếu chưa có)
-- DO $$
-- BEGIN
--     ALTER TABLE pois DISABLE TRIGGER trg_pois_sync_geom;
-- EXCEPTION
--     WHEN undefined_object THEN
--         RAISE NOTICE 'trigger trg_pois_sync_geom not found, skip disable';
-- END $$;

-- -- Drop index nặng (tạo lại cuối file). PK / UNIQUE(code) giữ lại.
-- DROP INDEX IF EXISTS idx_pois_geom;
-- DROP INDEX IF EXISTS idx_pois_adr_search_trgm;
-- DROP INDEX IF EXISTS idx_discovery_pois_name_trgm;
-- DROP INDEX IF EXISTS idx_discovery_pois_address_trgm;
-- DROP INDEX IF EXISTS idx_pois_parcel_id;
-- DROP INDEX IF EXISTS idx_pois_point_type;

-- CREATE OR REPLACE PROCEDURE migrate_parcels_to_pois_thua(batch_size bigint DEFAULT 100000)
-- LANGUAGE plpgsql
-- AS $$
-- DECLARE
--     v_min bigint;
--     v_max bigint;
--     v_cur bigint;
--     v_inserted bigint;
--     v_total bigint := 0;
-- BEGIN
--     SELECT COALESCE(MIN(id), 0), COALESCE(MAX(id), -1)
--       INTO v_min, v_max
--     FROM parcels
--     WHERE deleted_at IS NULL;

--     IF v_max < v_min THEN
--         RAISE NOTICE 'No parcels to migrate';
--         RETURN;
--     END IF;

--     v_cur := v_min;
--     WHILE v_cur <= v_max LOOP
--         INSERT INTO pois (
--             code,
--             name,
--             address,
--             latitude,
--             longitude,
--             parcel_id,
--             point_type,
--             adr_search,
--             is_active,
--             is_verified,
--             is_featured,
--             rating,
--             review_count,
--             view_count,
--             like_count,
--             created_at,
--             updated_at
--         )
--         SELECT
--             'PARCEL-' || p.id::text,
--             CASE
--                 WHEN NULLIF(TRIM(COALESCE(pi.address_text, '')), '') IS NOT NULL THEN pi.address_text
--                 ELSE 'Thửa'
--             END,
--             -- Tạm lưu nguyên address_text (không strip prefix ',')
--             COALESCE(pi.address_text, ''),
--             p.lat,
--             p.lon,
--             p.id,
--             'Thửa',
--             TRIM(
--                 REGEXP_REPLACE(
--                     LOWER(UNACCENT(COALESCE(pi.address_text, ''))),
--                     '[^a-z0-9]+',
--                     ' ',
--                     'g'
--                 )
--             ),
--             TRUE,
--             FALSE,
--             FALSE,
--             0,
--             0,
--             0,
--             0,
--             NOW(),
--             NOW()
--         FROM parcels p
--         LEFT JOIN qh_parcel_info pi ON pi.parcel_id = p.id
--         -- Skip nếu đã có poi gắn parcel (kể cả soft-delete: unique code vẫn chiếm slot)
--         LEFT JOIN pois existing_by_parcel
--             ON existing_by_parcel.parcel_id = p.id
--         LEFT JOIN pois existing_by_code
--             ON existing_by_code.code = 'PARCEL-' || p.id::text
--         WHERE p.deleted_at IS NULL
--           AND p.id >= v_cur
--           AND p.id < v_cur + batch_size
--           AND p.lat IS NOT NULL
--           AND p.lon IS NOT NULL
--           AND p.lat <> 0
--           AND p.lon <> 0
--           AND p.lat BETWEEN -90 AND 90
--           AND p.lon BETWEEN -180 AND 180
--           AND existing_by_parcel.id IS NULL
--           AND existing_by_code.id IS NULL
--         ON CONFLICT (code) DO NOTHING;

--         GET DIAGNOSTICS v_inserted = ROW_COUNT;
--         v_total := v_total + v_inserted;

--         RAISE NOTICE 'batch id [% .. %) inserted=% total=%',
--             v_cur, v_cur + batch_size, v_inserted, v_total;

--         v_cur := v_cur + batch_size;
--         COMMIT;
--     END LOOP;

--     RAISE NOTICE 'migrate parcels→pois done, inserted=%', v_total;
-- END;
-- $$;

-- CALL migrate_parcels_to_pois_thua(100000);

-- DROP PROCEDURE IF EXISTS migrate_parcels_to_pois_thua(bigint);

-- -- Backfill geom 1 lần (nhanh hơn trigger từng row)
-- UPDATE pois
-- SET geom = ST_SetSRID(
--         ST_MakePoint(longitude::double precision, latitude::double precision),
--         4326
--     )
-- WHERE point_type = 'Thửa'
--   AND parcel_id IS NOT NULL
--   AND code LIKE 'PARCEL-%'
--   AND geom IS NULL
--   AND latitude IS NOT NULL
--   AND longitude IS NOT NULL
--   AND latitude BETWEEN -90 AND 90
--   AND longitude BETWEEN -180 AND 180
--   AND NOT (latitude = 0 AND longitude = 0);

-- -- Tạo lại index
-- CREATE INDEX IF NOT EXISTS idx_pois_parcel_id
--     ON pois(parcel_id)
--     WHERE parcel_id IS NOT NULL AND parcel_id > 0;

-- CREATE INDEX IF NOT EXISTS idx_pois_point_type
--     ON pois(point_type)
--     WHERE deleted_at IS NULL AND point_type IS NOT NULL;

-- CREATE INDEX IF NOT EXISTS idx_pois_geom
--     ON pois USING GIST (geom)
--     WHERE deleted_at IS NULL AND geom IS NOT NULL;

-- CREATE INDEX IF NOT EXISTS idx_pois_adr_search_trgm
--     ON pois USING GIN (LOWER(adr_search) gin_trgm_ops)
--     WHERE COALESCE(adr_search, '') <> '';

-- CREATE INDEX IF NOT EXISTS idx_discovery_pois_name_trgm
--     ON pois USING GIN (LOWER(name) gin_trgm_ops)
--     WHERE deleted_at IS NULL;

-- CREATE INDEX IF NOT EXISTS idx_discovery_pois_address_trgm
--     ON pois USING GIN (LOWER(COALESCE(address, '')) gin_trgm_ops)
--     WHERE deleted_at IS NULL;

-- -- Bật lại trigger
-- DO $$
-- BEGIN
--     ALTER TABLE pois ENABLE TRIGGER trg_pois_sync_geom;
-- EXCEPTION
--     WHEN undefined_object THEN
--         RAISE NOTICE 'trigger trg_pois_sync_geom not found, skip enable';
-- END $$;

-- SELECT setval(
--     pg_get_serial_sequence('pois', 'id'),
--     GREATEST(COALESCE((SELECT MAX(id) FROM pois), 1), 1),
--     true
-- );

-- -- ---------------------------------------------------------------------------
-- -- OPTIONAL (chậm): thửa thiếu lat/lon → lấy PointOnSurface.
-- -- Chỉ chạy nếu cần; có thể tách session riêng.
-- -- ---------------------------------------------------------------------------
-- -- INSERT INTO pois (code, name, address, latitude, longitude, parcel_id, point_type, adr_search, is_active, is_verified, is_featured, rating, review_count, view_count, like_count, created_at, updated_at)
-- -- SELECT
-- --     'PARCEL-' || p.id::text,
-- --     CASE WHEN NULLIF(TRIM(COALESCE(pi.address_text, '')), '') IS NOT NULL THEN pi.address_text ELSE 'Thửa' END,
-- --     COALESCE(pi.address_text, ''),
-- --     ST_Y(ST_PointOnSurface(p.geometry)),
-- --     ST_X(ST_PointOnSurface(p.geometry)),
-- --     p.id,
-- --     'Thửa',
-- --     TRIM(REGEXP_REPLACE(LOWER(UNACCENT(COALESCE(pi.address_text, ''))), '[^a-z0-9]+', ' ', 'g')),
-- --     TRUE, FALSE, FALSE, 0, 0, 0, 0, NOW(), NOW()
-- -- FROM parcels p
-- -- LEFT JOIN qh_parcel_info pi ON pi.parcel_id = p.id
-- -- LEFT JOIN pois existing_by_parcel ON existing_by_parcel.parcel_id = p.id
-- -- LEFT JOIN pois existing_by_code ON existing_by_code.code = 'PARCEL-' || p.id::text
-- -- WHERE p.deleted_at IS NULL
-- --   AND existing_by_parcel.id IS NULL
-- --   AND existing_by_code.id IS NULL
-- --   AND p.geometry IS NOT NULL
-- --   AND (p.lat IS NULL OR p.lat = 0 OR p.lon IS NULL OR p.lon = 0)
-- -- ON CONFLICT (code) DO NOTHING;
