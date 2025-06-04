package dto

import _dto "common/dto"

type RegionRequest struct {
	_dto.Page
	Text     string `form:"text"`
	Level    *int   `form:"level"`
	ParentID *uint  `form:"parentId"`
}
