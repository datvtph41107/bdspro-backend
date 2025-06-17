package dto

import (
	_dto "common/dto"
)

type ProfileSearch struct {
	_dto.Pagable
	// Page int64  `form:"page"`
	// Size int64  `form:"size"`
	Text string `form:"text"`
}

type AccountItemDTO struct {
	ID       uint64 `json:"id"`
	FullName string `json:"fullName"`
	Avatar   string `json:"avatar"`
	Phone    string `json:"phone"`
}
