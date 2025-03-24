package u_asset

import (
	"bdspro/models"
	"common/case/crud"
)

type AssetService struct {
	crud.BaseService[models.Asset]
}

func NewAssetService(repo *AssetRepo) *AssetService {
	return &AssetService{
		BaseService: crud.BaseService[models.Asset]{Repo: repo.BaseRepo},
	}
}
