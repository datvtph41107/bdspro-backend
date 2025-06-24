package transaction

import (
	"time"
)

// TransactionHistory represents the history of changes to a transaction
type TransactionHistory struct {
	Id            uint32           
	TransactionId uint32           
	Action        TransactionAction
	UserId        uint32           
	UpdatedAt     time.Time        
} 