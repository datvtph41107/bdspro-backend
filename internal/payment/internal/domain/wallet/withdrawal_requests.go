package wallet

import (
	"time"
)

type WithdrawalRequest struct {
	Id              uint32
	WalletId        uint32
	Amount          int64
	Status          string
	PaymentMethod   string
	BankAccountInfo string
	RequestedBy     uint32
	ApprovedBy      uint32
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
