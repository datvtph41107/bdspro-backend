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
		&domain.AssetHistory{},
		&domain.AssetExploitation{},
		&domain.AssetLegal{},
		&domain.CostDocument{},
		&domain.IncomeDocument{},
	)
}
