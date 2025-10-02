package dto

import _dto "common/domain/dto"

type AssetLegalGetDTO struct {
	_dto.Pagable
	ID      uint64 `form:"id"`
	AssetID uint64 `form:"assetId"`
}
