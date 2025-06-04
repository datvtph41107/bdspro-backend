package dto

import _dto "common/dto"

type AssetCostTypeSearchDTO struct {
	_dto.Page
	AssetID *uint64   `form:"assetId"`
	TypeIDs *[]uint64 `form:"typeIds" parser:"uint64s"`
	Types   *[]uint64 `form:"types" parser:"uint64s"`
	Text    string    `form:"text"`
}
