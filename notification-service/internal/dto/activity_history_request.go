package dto

import (
	_dto "common/domain/dto"
	_enum "common/domain/enum"
)

type CreateActivityHistoryDTO struct {
	SubjectID   uint64
	SubjectType string
	Action      _enum.PropertyAction
	Metadata    []byte
	Description string
	ActorID     uint64
}

type ActivityHistoryQueryDTO struct {
	_dto.Pagable
	SubjectID   *uint64
	SubjectType *string
	Action      *_enum.PropertyAction
	ActorID     *uint64
	Keyword     *string
	Sort        string
}
