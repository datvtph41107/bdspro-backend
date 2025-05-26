package u_asset_legal

import (
	"bdspro/v1/models"
	_owner "common/case/owner"
)

type AssetLegalRouter struct {
	_owner.OwnerRouter[models.AssetLegal, SearchDTO]
}

func NewAssetLegalRouter(Service *AssetLegalService) *AssetLegalRouter {
	return &AssetLegalRouter{
		OwnerRouter: _owner.OwnerRouter[models.AssetLegal, SearchDTO]{
			Service: Service,
		},
	}
}
