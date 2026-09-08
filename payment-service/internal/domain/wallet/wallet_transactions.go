package wallet

import (
	"time"
)

type TransactionType string

type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "PENDING"
	TransactionStatusCompleted TransactionStatus = "COMPLETED"
	TransactionStatusCanceled  TransactionStatus = "CANCELED"
)

const (
	// Deposit related types
	TransactionTypeDeposit TransactionType = "DEPOSIT"

	// Withdrawal related types
	TransactionTypeWithdrawal TransactionType = "WITHDRAWAL"

	// Payment related types
	TransactionTypePayment TransactionType = "PAYMENT"

	// Transfer related types
	TransactionTypeTransferOut TransactionType = "TRANSFER_OUT"
	TransactionTypeTransferIn  TransactionType = "TRANSFER_IN"
)

type WalletTransaction struct {
	Id                uint32
	WalletId          uint32
	Type              TransactionType
	Amount            int64
	Status            TransactionStatus
	RelatedService    string
	RelatedId         string
	ExternalPaymentId string
	PaymentMethod     string
	Description       string
	CreatedBy         uint32
	ApprovedBy        uint32
	TransactionCode   string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
