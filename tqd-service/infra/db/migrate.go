package db

import (
	_db "common/db"
)

func AutoMigrate() {
	_db.DB.AutoMigrate(
	// Import error log — bảng theo dõi các batch insert thất bại
	// &qh_domain.QHRegionImportErrorLog{},

	// POI related models
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
	// &qh_domain.QHRegion{},
	// &qh_domain.QHRegionLabel{},
	// &qh_domain.QHLabel{},
	// &qh_domain.QHLayer{},
	// &qh_domain.QHShape{},
	)
}
