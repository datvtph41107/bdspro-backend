package db

import (
	_db "common/db"
	qh_domain "tqd/internal/domain/qh"
)

func AutoMigrate() {
	_db.DB.AutoMigrate(
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
		&qh_domain.Parcel{},
		&qh_domain.QHParcelInfo{},
		&qh_domain.QHLayer{},
		&qh_domain.QHRegion{},
		&qh_domain.QHLabel{},
		&qh_domain.QHShape{},
	)
}
