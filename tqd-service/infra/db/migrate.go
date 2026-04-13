package db

import (
	_db "common/db"
	"tqd/internal/domain"
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
		&domain.Parcel{},
	)
}
