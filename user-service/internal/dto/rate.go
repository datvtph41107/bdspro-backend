package dto

import _dto "common/dto"

type RateSearch struct {
	_dto.Pagable
	Score    *int    `form:"score"`
	ParentID *uint64 `form:"parentId"`
}
