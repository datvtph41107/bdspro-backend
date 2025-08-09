package db

import (
	"bdspro/internal/domain"
	tx_domain "bdspro/internal/domain/transaction"
	_db "common/db"
)

func MigrateDomain() {
	_db.DB.AutoMigrate(
		// &domain.Asset{},
		// &domain.AssetCost{},
		// &domain.AssetCostType{},
		// &domain.AssetIncomeType{},
		// &domain.AssetHistory{},
		// &domain.AssetExploitation{},
		// &domain.AssetLegal{},
		// &domain.CostDocument{},
		// &domain.IncomeDocument{},
		// &domain.PostMediaEntity{},
		// &domain.Product{},
		// &domain.AssetCostType{},
		// &domain.AssetCost{},
		// &domain.Post{},
		// &domain.SharingAccess{},
		// &domain.AssetCostType{},
		// &domain.AssetLegal{},
		// &domain.RecordHistory{},
		// &domain.SharingAccess{},
		// &domain.Post{},
		&tx_domain.Tx{},
		&domain.Product{},
	)
}
