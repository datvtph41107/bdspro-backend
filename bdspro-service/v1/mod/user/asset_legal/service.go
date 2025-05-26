package u_asset_legal

import (
	"bdspro/v1/models"
	_owner "common/case/owner"
)

type AssetLegalService struct {
	_owner.OwnerService[models.AssetLegal, SearchDTO]
}

func NewAssetLegalService(Repo *AssetLegalRepo) *AssetLegalService {
	return &AssetLegalService{
		OwnerService: _owner.OwnerService[models.AssetLegal, SearchDTO]{
			OwnerRepo: Repo,
		},
	}
}
