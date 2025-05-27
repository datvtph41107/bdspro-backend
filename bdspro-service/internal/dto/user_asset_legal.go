package dto

import _dto "common/dto"

type AssetLegalGetDTO struct {
	_dto.Page
	ID      uint64 `form:"id"`
	AssetID uint64 `form:"assetId"`
}
