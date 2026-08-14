package db

import (
	"log"

	_db "common/db"
)

// CheckPostGIS kiểm tra PostGIS đã cài trong DB hiện tại (trả version, installed).
func CheckPostGIS() (version string, installed bool) {
	var extVersion string
	err := _db.DB.Raw(
		`SELECT extversion FROM pg_extension WHERE extname = 'postgis'`,
	).Scan(&extVersion).Error
	if err != nil || extVersion == "" {
		return "", false
	}
	var postgisVer string
	_ = _db.DB.Raw(`SELECT PostGIS_Version()`).Scan(&postgisVer).Error
	if postgisVer != "" {
		return postgisVer, true
	}
	return extVersion, true
}

// EnsurePostGIS bật extension PostGIS và unaccent trước khi tạo cột geometry (AutoMigrate / raw SQL).
func EnsurePostGIS() {
	if ver, ok := CheckPostGIS(); ok {
		log.Printf("✅ PostGIS đã có: %s", ver)
	} else {
		log.Println("⚠️  PostGIS chưa bật — thử CREATE EXTENSION postgis ...")
		if err := _db.DB.Exec(`CREATE EXTENSION IF NOT EXISTS postgis`).Error; err != nil {
			log.Printf("❌ không bật được PostGIS: %v", err)
		}
	}

	log.Println("⚠️  Thử CREATE EXTENSION IF NOT EXISTS unaccent ...")
	if err := _db.DB.Exec(`CREATE EXTENSION IF NOT EXISTS unaccent`).Error; err != nil {
		log.Printf("❌ không bật được unaccent: %v", err)
	}

	log.Println("⚠️  Thử CREATE EXTENSION IF NOT EXISTS pg_trgm ...")
	if err := _db.DB.Exec(`CREATE EXTENSION IF NOT EXISTS pg_trgm`).Error; err != nil {
		log.Printf("❌ không bật được pg_trgm: %v", err)
	}
}

// ensureDiscoveryParcelSearchIndexes tạo GIN/BTree phục vụ SearchParcels (idempotent).
// Không DROP index ở startup — tránh lock bảng lớn; xem migration 000012.
func ensureDiscoveryParcelSearchIndexes() {
	statements := []string{
		`CREATE EXTENSION IF NOT EXISTS pg_trgm`,
		`CREATE INDEX IF NOT EXISTS idx_discovery_parcel_info_adr_trgm_nonempty
			ON qh_parcel_info USING GIN (LOWER(adr_search) gin_trgm_ops)
			WHERE COALESCE(adr_search, '') <> ''`,
		`CREATE INDEX IF NOT EXISTS idx_discovery_parcel_info_property_code_lower_nonempty
			ON qh_parcel_info (LOWER(property_code))
			WHERE property_code IS NOT NULL AND property_code <> ''`,
		`CREATE INDEX IF NOT EXISTS idx_discovery_parcel_info_map_land
			ON qh_parcel_info (map_number, land_number)`,
	}
	for _, stmt := range statements {
		if err := _db.DB.Exec(stmt).Error; err != nil {
			log.Printf("⚠️  ensureDiscoveryParcelSearchIndexes: %v", err)
		}
	}
}

// ensureProvinceV2Columns bổ sung cột mới cho province_v2.
// Không AutoMigrate province_v2 — bảng legacy; GORM cố DROP uni_province_v2_code không tồn tại.
func ensureProvinceV2Columns() {
	statements := []string{
		`ALTER TABLE province_v2
			ADD COLUMN IF NOT EXISTS qh_order INTEGER NOT NULL DEFAULT 0,
			ADD COLUMN IF NOT EXISTS zoom INTEGER NOT NULL DEFAULT 10,
			ADD COLUMN IF NOT EXISTS description VARCHAR(255),
			ADD COLUMN IF NOT EXISTS geometry geometry(MultiPolygon,4326),
			ADD COLUMN IF NOT EXISTS seo_id BIGINT`,
		`CREATE INDEX IF NOT EXISTS idx_provinces_qh_order
			ON province_v2(qh_order DESC, full_name ASC)`,
		`CREATE INDEX IF NOT EXISTS idx_province_v2_geometry ON province_v2 USING GIST (geometry)`,
		`CREATE INDEX IF NOT EXISTS idx_province_v2_seo_id
			ON province_v2(seo_id)
			WHERE seo_id IS NOT NULL AND seo_id > 0`,
	}
	for _, stmt := range statements {
		if err := _db.DB.Exec(stmt).Error; err != nil {
			log.Printf("⚠️  ensureProvinceV2Columns: %v", err)
		}
	}
	_ = _db.DB.Exec(`
		DO $$
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_constraint
				WHERE conname = 'chk_province_v2_zoom'
				  AND conrelid = 'province_v2'::regclass
			) THEN
				ALTER TABLE province_v2
					ADD CONSTRAINT chk_province_v2_zoom
					CHECK (zoom >= 5 AND zoom <= 22);
			END IF;
		END $$`).Error
}

// ensureWardV2Columns bổ sung cột mới cho ward_v2.
// Không AutoMigrate ward_v2 — bảng legacy; GORM cố DROP uni_ward_v2_code không tồn tại.
func ensureWardV2Columns() {
	statements := []string{
		`ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS description VARCHAR(255)`,
		`ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS geometry geometry(MultiPolygon,4326)`,
		`ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS seo_id BIGINT`,
		`CREATE INDEX IF NOT EXISTS idx_ward_v2_geometry ON ward_v2 USING GIST (geometry)`,
		`CREATE INDEX IF NOT EXISTS idx_ward_v2_seo_id
			ON ward_v2(seo_id)
			WHERE seo_id IS NOT NULL AND seo_id > 0`,
	}
	for _, stmt := range statements {
		if err := _db.DB.Exec(stmt).Error; err != nil {
			log.Printf("⚠️  ensureWardV2Columns: %v", err)
		}
	}
}

func AutoMigrate() {
	EnsurePostGIS()
	ensureProvinceV2Columns()
	ensureWardV2Columns()
	ensureDiscoveryParcelSearchIndexes()
	// province_v2 / ward_v2 là bảng legacy — chỉ cập nhật schema qua ensure*() ở trên,
	// không dùng GORM AutoMigrate (GORM cố DROP uni_*_code không tồn tại).
	_db.DB.AutoMigrate(
	// &qh_domain.Parcel{},
	// &qh_domain.QHLayer{},
	// &qh_domain.Parcel{},
	)
}
