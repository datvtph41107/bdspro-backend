package db

import (
	"log"
	"tqd/internal/domain"
	qh_domain "tqd/internal/domain/qh"

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
}

// ensureProvinceV2PlanningColumns bổ sung cột qh_order/zoom cho province_v2.
// Không dùng AutoMigrate(&domain.Province{}) vì GORM cố DROP constraint
// uni_province_v2_code không tồn tại trên bảng legacy.
func ensureProvinceV2PlanningColumns() {
	statements := []string{
		`ALTER TABLE province_v2
			ADD COLUMN IF NOT EXISTS qh_order INTEGER NOT NULL DEFAULT 0,
			ADD COLUMN IF NOT EXISTS zoom INTEGER NOT NULL DEFAULT 10`,
		`CREATE INDEX IF NOT EXISTS idx_provinces_qh_order
			ON province_v2(qh_order DESC, full_name ASC)`,
	}
	for _, stmt := range statements {
		if err := _db.DB.Exec(stmt).Error; err != nil {
			log.Printf("⚠️  ensureProvinceV2PlanningColumns: %v", err)
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

func AutoMigrate() {
	EnsurePostGIS()
	ensureProvinceV2PlanningColumns()
	_db.DB.AutoMigrate(
		// Import error log — bảng theo dõi các batch insert thất bại
		&qh_domain.QHRegionImportErrorLog{},

		// POI related models
		&domain.POI{},
		&domain.PoiCategory{},
		&domain.OpenHour{},
		&domain.POIMedia{},
		// &domain.Feature{},

		// Amenity model
		&domain.Amenity{},

		// Contact Label model
		&domain.ContactLabel{},
		&domain.DirectoryCategory{},
		&domain.DirectorySource{},
		&domain.DirectorySupplier{},
		&qh_domain.Parcel{},
		&qh_domain.QHParcelInfo{},
		&qh_domain.QHLayer{},
		&qh_domain.QHJurisdiction{},
		&qh_domain.QHPlanningProject{},
		&qh_domain.QHRegion{},
		&qh_domain.QHRegionExtend{},
		// &qh_domain.QHRegionLabel{},
		&qh_domain.QHLabel{},
		&qh_domain.QHLabelLayer{},
		&qh_domain.QHShape{},
		&qh_domain.QHLayerFamily{},
		&qh_domain.QHLayerLegend{},
		// &qh_domain.LandUseCode{},
		// &qh_domain.LandUseGroup{},
		&domain.Province{},
		&domain.Ward{},
		&qh_domain.QHLayerLegal{},
		&qh_domain.QHLandUse{},
		&qh_domain.QHLayerResolverConfig{},
		&qh_domain.QHLayerLandUse{},
		&qh_domain.QHAuthorityIssuring{},
		&qh_domain.QHLayerLegend{},
		&qh_domain.QHAuditEntry{},
		// &qh_domain.QHLifecycleTransition{},

		&qh_domain.QHUserFollowedParcel{},
		&qh_domain.QHUserReported{},
		&qh_domain.QHUserViewHistory{},
		&qh_domain.QHUserViewEvent{},

		&qh_domain.QHPlanningProject{},
		&qh_domain.QHPlanningDocument{},
		&qh_domain.QHPlanningEvent{},
		&qh_domain.QHPlanningRelation{},
		&domain.Tag{},
		&qh_domain.QHPlanningProjectTag{},
		// province_v2 là bảng legacy — cột qh_order/zoom qua ensureProvinceV2PlanningColumns().
		&domain.Province{},
		&qh_domain.QHPlanningProjectProvince{},
		&domain.Report{},
		&domain.ReportEvent{},
		&domain.UserSubscription{},
	)
}
