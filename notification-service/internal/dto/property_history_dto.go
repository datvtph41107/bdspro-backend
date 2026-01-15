package dto

import (
	_dto "common/domain/dto"
	_enum "common/domain/enum"
	"time"
)

type CreatePropertyHistoryDTO struct {
	SubjectID   uint64
	SubjectType _enum.SubjectTypeProperty
	Action      _enum.PropertyAction
	Metadata    []byte
	Description string
	ActorID     uint64
}

type PropertyHistoryDetailDTO struct {
	ID          uint64
	SubjectID   uint64
	SubjectType string
	Action      _enum.PropertyAction
	Description string
	Metadata    []byte
	ActorID     uint64
	OccurredAt  time.Time
}

type PropertyHistoryQueryDTO struct {
	_dto.Pagable
	SubjectID   *uint64
	SubjectType *string
	Action      *_enum.PropertyAction
	ActorID     *uint64
	Keyword     *string
	Limit       int
	Offset      int
}
