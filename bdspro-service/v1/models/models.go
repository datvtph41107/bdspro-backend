package models

import (
	_db "common/db"
)

func AutoMigrate() {
	_db.DB.AutoMigrate(
		// &Deposite{},
		// &ProductChildHistory{},
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
		// &ProductPrice{},
		// &Post{},
		// &Action{},
		// &_owner.OwnerModel[AssetCost]{},
		// &_category.CategoryModel[AssetCostType]{},
		// &_owner.OwnerModel[AssetExploitation]{},
		// &_owner.OwnerModel[AssetLegal]{},
		// &_category.CategoryModel[AssetIncomeType]{},
		// &_owner.OwnerModel[AssetShare]{},
		// &AssetLegal{},
		// &AssetHistory{},
		// &IncomeDocument{},
		// &CostDocument{},
		// &Asset{},
		// &ProjectBuild{},
		// &AssetExploitation{},
		// &AssetCost{},
		&AssetCostType{},
	)
}
