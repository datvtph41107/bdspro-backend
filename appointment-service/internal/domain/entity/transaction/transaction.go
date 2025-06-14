package transaction

import (
	"time"
)

// Transaction represents a financial transaction
type Transaction struct {
	Id                   uint32
	OwnerId              uint32
	OrganizationId       *uint32
	Amount               float64
	Currency             string
	TransactionName      string
	Description          string
	TransactionType      TransactionType
	CategoryId           uint32
	PaymentMethodId      uint32
	RelatedDealId        *uint32
	RelatedTransactionId *uint32
	ApprovalStatus       ApprovalStatus
	ApprovedBy           *uint32
	ApprovalDate         *time.Time
	TransactionStatus    TransactionStatus
	TransactionDate      time.Time
	IsDeleted            bool
	CreatedAt            time.Time
	UpdatedAt            time.Time
}
