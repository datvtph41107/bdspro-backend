package dto

import _dto "common/domain/dto"

type AssetCostSearchDTO struct {
	_dto.Pagable
	// Page    int     `form:"page"`
	// Size    int     `form:"size"`
	AssetID *uint64   `form:"assetId"`
	TypeIDs *[]uint64 `form:"typeIds" parser:"uint64s"`
	Types   *[]uint64 `form:"types" parser:"uint64s"`
	// Text    string    `form:"text"`
}
