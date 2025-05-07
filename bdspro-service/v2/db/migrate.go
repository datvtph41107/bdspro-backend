package db

import (
	"bdspro/v2/internal/domain"
	_db "common/db"
)

func MigrateDomain() {
	_db.DB.AutoMigrate(
		&domain.Asset{},
		&domain.AssetCost{},
		&domain.AssetExploitation{},
		&domain.AssetLegal{},
		&domain.AssetSplitMergeHistory{},
	)
}
