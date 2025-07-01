package dto

import (
	_dto "common/dto"
	"time"
)

type RateSearchDTO struct {
	_dto.Pagable
	TargetId   uint64    `json:"targetId"`
	TargetType uint32    `json:"targetType"`
	ParentId   *uint64   `json:"parentId"`
	CreatedBy  uint64    `json:"createdBy"`
	UpdatedBy  uint64    `json:"updatedBy"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
