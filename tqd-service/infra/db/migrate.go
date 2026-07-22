package db

import (
	"log"
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

func AutoMigrate() {
	EnsurePostGIS()
	_db.DB.AutoMigrate(
		// Import error log — bảng theo dõi các batch insert thất bại
		// &qh_domain.QHRegionImportErrorLog{},

		// // POI related models
		// &domain.POI{},
		// &domain.PoiCategory{},
		// &domain.OpenHour{},
		// &domain.POIMedia{},
		// // &domain.Feature{},

		// // Amenity model
		// &domain.Amenity{},

		// // Contact Label model
		// &domain.ContactLabel{},
		// &domain.DirectoryCategory{},
		// &domain.DirectorySource{},
		// &domain.DirectorySupplier{},
		// &qh_domain.Parcel{},
		// &qh_domain.QHParcelInfo{},
		// &qh_domain.QHLayer{},
		// &qh_domain.QHJurisdiction{},
		// &qh_domain.QHPlanningProject{},
		// &qh_domain.QHRegion{},
		// &qh_domain.QHRegionExtend{},
		// // &qh_domain.QHRegionLabel{},
		// &qh_domain.QHLabel{},
		// &qh_domain.QHLabelLayer{},
		// &qh_domain.QHShape{},
		// &qh_domain.QHLayerFamily{},
		// &qh_domain.QHLayerLegend{},
		// // &qh_domain.LandUseCode{},
		// // &qh_domain.LandUseGroup{},
		// &domain.Province{},
		// &domain.Ward{},
		// &qh_domain.QHLayerLegal{},
		// &qh_domain.QHLandUse{},
		// &qh_domain.QHLayerResolverConfig{},
		// &qh_domain.QHLayerLandUse{},
		// &qh_domain.QHAuthorityIssuring{},
		// &qh_domain.QHLayerLegend{},
		// &qh_domain.QHAuditEntry{},
		// &qh_domain.QHLifecycleTransition{},

		// &qh_domain.QHUserFollowedParcel{},
		// &qh_domain.QHUserReported{},
		// &qh_domain.QHUserViewHistory{},
		// &qh_domain.QHUserViewEvent{},

		// &qh_domain.QHPlanningProject{},
		&qh_domain.QHPlanningDocument{},
	// &domain.UserSubscription{},
	)
}
