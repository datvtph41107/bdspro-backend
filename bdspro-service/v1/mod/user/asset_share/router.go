package u_asset_share

import (
	"bdspro/v1/models"
	_owner "common/case/owner"
)

type AssetShareRouter struct {
	_owner.OwnerRouter[models.AssetShare, SearchDTO]
}

func NewAssetShareRouter(Service *AssetShareService) *AssetShareRouter {
	return &AssetShareRouter{
		OwnerRouter: _owner.OwnerRouter[models.AssetShare, SearchDTO]{
			Service: Service,
		},
	}
}
