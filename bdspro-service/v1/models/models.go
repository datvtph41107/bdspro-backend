package models

import (
	_db "common/db"
)

func AutoMigrate() {
	_db.DB.AutoMigrate(
		// &Product{},
		// &Province{},
		// &District{},
		// &Ward{},
		// &ProductAccess{},
		// &DocType{},
		// &PropertyType{},
		// &Amenity{},
		// &ProductAmenity{},
		// &Project{},
		// &ProjectInventory{},
		// &Product{},
		// &ProductMedia{},
		// &HouseInfo{},
		// &ProjectBuild{},
		// &BuildRoom{},
		// &ApartmentAttribute{},
		// &Apartment{},
		// &Developer{},
		// &ProductPrivate{},
		// &Region{},
		// &Product{},
		&ProductPrice{},
	)
}
