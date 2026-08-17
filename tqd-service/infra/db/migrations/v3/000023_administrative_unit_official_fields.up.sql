-- Thông tin hành chính chính thức + soft delete cho tỉnh/thành và phường/xã.
ALTER TABLE province_v2
    ADD COLUMN IF NOT EXISTS official_name VARCHAR(255),
    ADD COLUMN IF NOT EXISTS division_type VARCHAR(100),
    ADD COLUMN IF NOT EXISTS official_operation_date DATE,
    ADD COLUMN IF NOT EXISTS natural_area_km2 NUMERIC(12, 4),
    ADD COLUMN IF NOT EXISTS population BIGINT,
    ADD COLUMN IF NOT EXISTS population_density NUMERIC(12, 2),
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;

ALTER TABLE ward_v2
    ADD COLUMN IF NOT EXISTS official_name VARCHAR(255),
    ADD COLUMN IF NOT EXISTS division_type VARCHAR(100),
    ADD COLUMN IF NOT EXISTS official_operation_date DATE,
    ADD COLUMN IF NOT EXISTS natural_area_km2 NUMERIC(12, 4),
    ADD COLUMN IF NOT EXISTS population BIGINT,
    ADD COLUMN IF NOT EXISTS population_density NUMERIC(12, 2),
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;

UPDATE province_v2
SET official_name = full_name
WHERE (official_name IS NULL OR official_name = '')
  AND full_name IS NOT NULL
  AND full_name <> '';

UPDATE ward_v2
SET official_name = full_name
WHERE (official_name IS NULL OR official_name = '')
  AND full_name IS NOT NULL
  AND full_name <> '';

CREATE INDEX IF NOT EXISTS idx_province_v2_deleted_at
    ON province_v2(deleted_at);

CREATE INDEX IF NOT EXISTS idx_ward_v2_deleted_at
    ON ward_v2(deleted_at);

COMMENT ON COLUMN province_v2.official_name IS 'Tên gọi chính thức';
COMMENT ON COLUMN province_v2.division_type IS 'Loại đơn vị (Tỉnh, Thành phố, ...)';
COMMENT ON COLUMN province_v2.official_operation_date IS 'Ngày hoạt động chính thức';
COMMENT ON COLUMN province_v2.natural_area_km2 IS 'Diện tích tự nhiên (km²)';
COMMENT ON COLUMN province_v2.population IS 'Quy mô dân số';
COMMENT ON COLUMN province_v2.population_density IS 'Mật độ dân số (người/km²)';

COMMENT ON COLUMN ward_v2.official_name IS 'Tên gọi chính thức';
COMMENT ON COLUMN ward_v2.division_type IS 'Loại đơn vị (Phường, Xã, ...)';
COMMENT ON COLUMN ward_v2.official_operation_date IS 'Ngày hoạt động chính thức';
COMMENT ON COLUMN ward_v2.natural_area_km2 IS 'Diện tích tự nhiên (km²)';
COMMENT ON COLUMN ward_v2.population IS 'Quy mô dân số';
COMMENT ON COLUMN ward_v2.population_density IS 'Mật độ dân số (người/km²)';
