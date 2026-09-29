package dto

import _dto "common/domain/dto"

type RegionRequest struct {
	_dto.Pagable
	Text     string  `form:"text"`
	Level    *int    `form:"level"`
	ParentID *uint64 `form:"parentId"`
}
