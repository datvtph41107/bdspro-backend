package data

import (
	_dto "common/dto"
)

type ProfileSearch struct {
	_dto.Page
	// Page int64  `form:"page"`
	// Size int64  `form:"size"`
	Text string `form:"text"`
}
