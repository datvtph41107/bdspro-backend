package db

import (
	_db "common/db"
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
	// &qh_domain.QHLabel{},
	// &qh_domain.QHLandType{},
	// &qh_domain.QHLayer{}, // managed by SQL migration (000005), removing to avoid constraint name conflict
	// &qh_domain.QHParcelLabel{},
	// &qh_domain.Parcel{},
	// &qh_domain.QHRegion{},
	// &qh_domain.Region{},
	// &qh_domain.QHParcelInfo{},
	// &qh_domain.QHDirection{},

	// &qh_domain.QHPlanningUse{},
	// &domain.Report{},
	)
}
