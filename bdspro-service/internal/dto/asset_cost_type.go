package dto

import _dto "common/domain/dto"

type AssetCostTypeSearchDTO struct {
	_dto.Pagable
	AssetID *uint64   `form:"assetId"`
	TypeIDs *[]uint64 `form:"typeIds" parser:"uint64s"`
	Type    *uint32   `form:"type" parser:"uint32"`
	Text    string    `form:"text"`
}
