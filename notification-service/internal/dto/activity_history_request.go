package dto

import (
	_dto "common/domain/dto"
	"notification/internal/enums"
)

type CreateActivityHistoryDTO struct {
	SubjectID   uint64
	SubjectType string
	Action      enums.PropertyAction
	Metadata    []byte
	Description string
	ActorID     uint64
	ActorType   enums.ActorType
}

type ActivityHistoryQueryDTO struct {
	_dto.Pagable
	SubjectID   *uint64
	SubjectType *string
	Action      *enums.PropertyAction
	ActorID     *uint64
	ActorType   *enums.ActorType
	Keyword     *string
	Sort        string
}
