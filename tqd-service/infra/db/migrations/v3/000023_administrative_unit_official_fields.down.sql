DROP INDEX IF EXISTS idx_ward_v2_deleted_at;
DROP INDEX IF EXISTS idx_province_v2_deleted_at;

ALTER TABLE ward_v2
    DROP COLUMN IF EXISTS deleted_at,
    DROP COLUMN IF EXISTS population_density,
    DROP COLUMN IF EXISTS population,
    DROP COLUMN IF EXISTS natural_area_km2,
    DROP COLUMN IF EXISTS official_operation_date,
    DROP COLUMN IF EXISTS division_type,
    DROP COLUMN IF EXISTS official_name;

ALTER TABLE province_v2
    DROP COLUMN IF EXISTS deleted_at,
    DROP COLUMN IF EXISTS population_density,
    DROP COLUMN IF EXISTS population,
    DROP COLUMN IF EXISTS natural_area_km2,
    DROP COLUMN IF EXISTS official_operation_date,
    DROP COLUMN IF EXISTS division_type,
    DROP COLUMN IF EXISTS official_name;
