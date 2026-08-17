-- Bổ sung cột hành chính chính thức + các cột domain đang dùng nhưng DB chưa có.
-- Tách từng câu ADD COLUMN IF NOT EXISTS để chạy được dù một phần schema đã tồn tại.

-- province_v2
ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS official_name VARCHAR(255);
ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS division_type VARCHAR(100);
ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS official_operation_date DATE;
ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS natural_area_km2 NUMERIC(12, 4);
ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS population BIGINT;
ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS population_density NUMERIC(12, 2);
ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;
ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS ward_count INTEGER NOT NULL DEFAULT 0;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'province_v2' AND column_name = 'description'
    ) THEN
        ALTER TABLE province_v2 ALTER COLUMN description TYPE VARCHAR(2000);
    END IF;
END $$;

-- ward_v2
ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS official_name VARCHAR(255);
ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS division_type VARCHAR(100);
ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS official_operation_date DATE;
ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS natural_area_km2 NUMERIC(12, 4);
ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS population BIGINT;
ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS population_density NUMERIC(12, 2);
ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;
ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS province_code VARCHAR(20);
ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'ward_v2' AND column_name = 'description'
    ) THEN
        ALTER TABLE ward_v2 ALTER COLUMN description TYPE VARCHAR(2000);
    END IF;
END $$;

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

UPDATE ward_v2 w
SET province_code = p.code
FROM province_v2 p
WHERE w.province_id = p.id
  AND (w.province_code IS NULL OR w.province_code = '')
  AND p.code IS NOT NULL
  AND p.code <> '';

CREATE INDEX IF NOT EXISTS idx_province_v2_deleted_at ON province_v2(deleted_at);
CREATE INDEX IF NOT EXISTS idx_ward_v2_deleted_at ON ward_v2(deleted_at);
CREATE INDEX IF NOT EXISTS idx_province_v2_division_type ON province_v2(division_type);
CREATE INDEX IF NOT EXISTS idx_ward_v2_division_type ON ward_v2(division_type);
CREATE INDEX IF NOT EXISTS idx_ward_v2_province_code ON ward_v2(province_code);
