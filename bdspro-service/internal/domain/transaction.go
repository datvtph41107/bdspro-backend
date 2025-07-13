package domain

import (
	"bdspro/internal/enums"
	_models "common/models"
	"time"
)

// Transaction represents a financial transaction
type Transaction struct {
	_models.BaseEntity
	OwnerId              uint64
	OwnerType            enums.EOwnerType
	ProductID            *uint64
	DepositeAmount       *float64 // cọc
	DepositeNote         *string  // ghi chú cọc
	ContactPhone         *string  // số điện thoại liên hệ
	ContactID            *uint64
	Amount               float64
	Currency             string
	TransactionName      string
	Description          string
	TransactionType      enums.TransactionType
	CategoryId           uint32
	PaymentMethodId      uint32
	RelatedDealId        *uint32
	RelatedTransactionId *uint32
	ApprovalStatus       enums.ApprovalStatus
	ApprovedBy           *uint32
	ApprovalDate         *time.Time
	TransactionStatus    enums.TransactionStatus
	TransactionDate      time.Time
}

func (t *Transaction) TableName() string {
	return "transactions"
}
