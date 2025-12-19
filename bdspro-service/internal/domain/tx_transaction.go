package domain

import (
	"bdspro/internal/enums"
	_models "common/models"
)

// Transaction represents a financial transaction
type TxTransaction struct {
	_models.BaseEntity
	OwnerId         uint64
	OwnerType       enums.TxOwnerType
	ProductID       *uint64
	DepositeAmount  *float64 // cọc
	DepositeNote    *string  // ghi chú cọc
	ContactPhone    *string  // số điện thoại liên hệ
	ContactID       *uint64
	Amount          float64
	Currency        string
	TransactionName string
	Description     string
	TransactionType enums.TxMethod
	CategoryId      uint32
	PaymentMethodId uint32
}

func (t *TxTransaction) TableName() string {
	return "transactions"
}
