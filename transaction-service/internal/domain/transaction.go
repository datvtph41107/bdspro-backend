package domain

import (
	_models "common/models"
	"tx/internal/enums"
)

// Transaction represents a financial transaction
type Transaction struct {
	_models.BaseEntity
	OwnerId         uint64
	OwnerType       enums.EOwnerType
	ProductID       *uint64
	DepositeAmount  *float64 // cọc
	DepositeNote    *string  // ghi chú cọc
	ContactPhone    *string  // số điện thoại liên hệ
	ContactID       *uint64
	Amount          float64
	Currency        string
	TransactionName string
	Description     string
	TransactionType enums.EMethod
	CategoryId      uint32
	PaymentMethodId uint32
}

func (t *Transaction) TableName() string {
	return "transactions"
}
