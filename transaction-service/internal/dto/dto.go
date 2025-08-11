package dto

import (
	_dto "common/dto"
	"time"
	"tx/internal/enums"
)

type TransactionFilters struct {
	_dto.Pagable
	TransactionType   *string
	ApprovalStatus    *string
	TransactionStatus *uint32
	StartDate         *time.Time
	EndDate           *time.Time
	ToID              *uint64
	ToOf              enums.EOwnerType
}
