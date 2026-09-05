-- Rollback migrate parcels → pois (thửa).

DELETE FROM pois
WHERE point_type = 'Thửa'
  AND parcel_id IS NOT NULL
  AND code LIKE 'PARCEL-%';

DROP INDEX IF EXISTS idx_pois_point_type;

-- Không drop point_type nếu còn data khác dùng cột này; chỉ drop khi chắc chắn.
-- ALTER TABLE pois DROP COLUMN IF EXISTS point_type;

SELECT setval(
    pg_get_serial_sequence('pois', 'id'),
    GREATEST(COALESCE((SELECT MAX(id) FROM pois), 1), 1),
    true
);
