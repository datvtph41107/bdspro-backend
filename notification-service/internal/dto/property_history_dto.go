package dto

import (
	_dto "common/domain/dto"
	_enum "common/domain/enum"
	"encoding/json"
	"time"
)

type CreatePropertyHistoryDTO struct {
	SubjectID   uint64
	Action      _enum.PropertyAction
	Metadata    []byte
	Description string
	ActorID     uint64
}

type PropertyHistoryDetailDTO struct {
	ID          uint64
	SubjectID   uint64
	Action      _enum.PropertyAction
	Description string
	Metadata    []byte
	ActorID     uint64
	OccurredAt  time.Time
}

type PropertyHistoryQueryDTO struct {
	_dto.Pagable
	SubjectID *uint64

	Action  *_enum.PropertyAction
	ActorID *uint64
	Keyword *string
	Limit   int
	Offset  int
}

type PropertyHistoryResponseDTO struct {
	ID          uint64           `json:"id"`
	SubjectID   uint64           `json:"subjectId"`
	Action      string           `json:"action"`
	Description string           `json:"description"`
	Metadata    json.RawMessage  `json:"metadata"`
	OccurredAt  *time.Time       `json:"occurredAt"`
	Actor       *ActorProfileDTO `json:"actor,omitempty"`
}

type ActorProfileDTO struct {
	ID       uint64 `json:"id"`
	FullName string `json:"fullName"`
	Avatar   string `json:"avatar"`
}
