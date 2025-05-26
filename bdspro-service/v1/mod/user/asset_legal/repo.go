package u_asset_legal

import (
	"bdspro/v1/models"
	_owner "common/case/owner"

	"gorm.io/gorm"
)

type AssetLegalRepo struct {
	_owner.OwnerRepo[models.AssetLegal]
}

func NewAssetLegalRepo(DB *gorm.DB) *AssetLegalRepo {
	return &AssetLegalRepo{
		OwnerRepo: _owner.OwnerRepo[models.AssetLegal]{
			DB: DB,
		},
	}
}
