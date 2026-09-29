package domain

import (
	"time"
)

// TransactionProduct represents the products associated with a transaction
type TxTransactionProduct struct {
	Id            uint32
	TransactionId uint32
	ProductId     uint32
	CreatedAt     time.Time
}
