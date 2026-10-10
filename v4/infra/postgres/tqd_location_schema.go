package postgres

import (
	"fmt"
	"log/slog"

	"gorm.io/gorm"
)

// ensureTqdLocationCatalogSchema aligns province_v2 / ward_v2 with modules/tqd location
// queries on shared-DB develop, where v4/migrate (hub) may have created the tables first.
func ensureTqdLocationCatalogSchema(database *gorm.DB) error {
	if database == nil {
		return nil
	}
	var hasProvinceV2 bool
	if err := database.Raw(`
SELECT EXISTS (
  SELECT 1 FROM information_schema.tables
  WHERE table_schema = 'public' AND table_name = 'province_v2'
)`).Scan(&hasProvinceV2).Error; err != nil {
		return fmt.Errorf("detect province_v2: %w", err)
	}
	if !hasProvinceV2 {
		return nil
	}

	for _, stmt := range tqdProvinceV2CatalogStatements() {
		if err := database.Exec(stmt).Error; err != nil {
			slog.Warn(
				"TQD location schema ensure statement failed",
				"category", "lifecycle",
				"event_name", "database.schema.tqd_location.ensure",
				"error", err.Error(),
			)
		}
	}
	for _, stmt := range tqdWardV2CatalogStatements() {
		if err := database.Exec(stmt).Error; err != nil {
			slog.Warn(
				"TQD location schema ensure statement failed",
				"category", "lifecycle",
				"event_name", "database.schema.tqd_location.ensure",
				"error", err.Error(),
			)
		}
	}
	return nil
}

func tqdProvinceV2CatalogStatements() []string {
	return []string{
		`ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS full_name VARCHAR(255)`,
		`ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS short_name VARCHAR(100)`,
		`ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS lat NUMERIC(10,8)`,
		`ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS lng NUMERIC(11,8)`,
		`ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS ward_count BIGINT DEFAULT 0`,
		`ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS qh_order INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS zoom INTEGER NOT NULL DEFAULT 10`,
		`ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS description VARCHAR(2000)`,
		`ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS seo_id BIGINT`,
		`ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS unit_type INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS data_source INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS plan_doc_id BIGINT`,
		`ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS official_name VARCHAR(255)`,
		`ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS division_type VARCHAR(100)`,
		`ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS official_operation_date DATE`,
		`ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS natural_area_km2 NUMERIC(12, 4)`,
		`ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS population BIGINT`,
		`ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS population_density NUMERIC(12, 2)`,
		`ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS official_role VARCHAR(500) NOT NULL DEFAULT ''`,
		`ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS official_address VARCHAR(500) NOT NULL DEFAULT ''`,
		`DO $$
BEGIN
  IF EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'province_v2' AND column_name = 'name'
  ) THEN
    UPDATE province_v2
    SET full_name = COALESCE(NULLIF(btrim(full_name), ''), NULLIF(btrim(name), ''))
    WHERE COALESCE(NULLIF(btrim(full_name), ''), '') = ''
      AND name IS NOT NULL AND btrim(name) <> '';
  END IF;
  IF EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'province_v2' AND column_name = 'codename'
  ) THEN
    UPDATE province_v2
    SET short_name = COALESCE(NULLIF(btrim(short_name), ''), NULLIF(btrim(codename), ''))
    WHERE COALESCE(NULLIF(btrim(short_name), ''), '') = ''
      AND codename IS NOT NULL AND btrim(codename) <> '';
  END IF;
  IF EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'province_v2' AND column_name = 'code'
  ) THEN
    UPDATE province_v2
    SET full_name = COALESCE(NULLIF(btrim(full_name), ''), code::text)
    WHERE COALESCE(NULLIF(btrim(full_name), ''), '') = ''
      AND code IS NOT NULL;
  END IF;
END $$`,
		`UPDATE province_v2
SET official_name = full_name
WHERE (official_name IS NULL OR official_name = '')
  AND full_name IS NOT NULL AND full_name <> ''`,
	}
}

func tqdWardV2CatalogStatements() []string {
	return []string{
		`ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS full_name VARCHAR(255)`,
		`ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS short_name VARCHAR(100)`,
		`ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS lat NUMERIC(10,8)`,
		`ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS lng NUMERIC(11,8)`,
		`ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS province_id BIGINT`,
		`ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS province_code VARCHAR(20)`,
		`ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS description VARCHAR(2000)`,
		`ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS summary VARCHAR(200)`,
		`ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS seo_id BIGINT`,
		`ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS official_name VARCHAR(255)`,
		`ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS division_type VARCHAR(100)`,
		`ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS official_operation_date DATE`,
		`ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS natural_area_km2 NUMERIC(12, 4)`,
		`ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS population BIGINT`,
		`ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS population_density NUMERIC(12, 2)`,
		`ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS unit_type INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS data_source INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS official_role VARCHAR(500) NOT NULL DEFAULT ''`,
		`ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS official_address VARCHAR(500) NOT NULL DEFAULT ''`,
		`DO $$
BEGIN
  IF EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'ward_v2' AND column_name = 'name'
  ) THEN
    UPDATE ward_v2
    SET full_name = COALESCE(NULLIF(btrim(full_name), ''), NULLIF(btrim(name), ''))
    WHERE COALESCE(NULLIF(btrim(full_name), ''), '') = ''
      AND name IS NOT NULL AND btrim(name) <> '';
  END IF;
  IF EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'ward_v2' AND column_name = 'codename'
  ) THEN
    UPDATE ward_v2
    SET short_name = COALESCE(NULLIF(btrim(short_name), ''), NULLIF(btrim(codename), ''))
    WHERE COALESCE(NULLIF(btrim(short_name), ''), '') = ''
      AND codename IS NOT NULL AND btrim(codename) <> '';
  END IF;
END $$`,
	}
}
