package data

import _dto "common/dto"

type RateSearch struct {
	_dto.Page
	Score    *int    `form:"score"`
	ParentID *uint64 `form:"parentId"`
}
