package db

import (
	"bdspro/internal/domain"
	_db "common/db"
)

func MigrateDomain() {
	_db.DB.AutoMigrate(
		&domain.Asset{},
		&domain.AssetCost{},
		&domain.AssetCostType{},
		&domain.AssetIncomeType{},
		&domain.AssetExploitation{},
		&domain.AssetLegal{},
		&domain.CostDocument{},
		&domain.IncomeDocument{},
		&domain.PostMediaEntity{},
		&domain.Product{},
		&domain.AssetCostType{},
		&domain.AssetCost{},
		&domain.Post{},
		&domain.ProductUser{},
		&domain.ProductOrganization{},
		&domain.AssetOrganization{},
		&domain.AssetUser{},
		&domain.PostOrganization{},
		&domain.PostUser{},
		&domain.SharingAccess{},
		&domain.AssetCostType{},
		&domain.AssetLegal{},
		&domain.RecordHistory{},
		&domain.SharingAccess{},
		&domain.Post{},
		&domain.Product{},
		&domain.ProductHistory{},
	)
}
