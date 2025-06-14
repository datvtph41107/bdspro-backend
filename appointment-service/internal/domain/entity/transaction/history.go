package transaction

import (
	"time"
)

// TransactionHistory represents the history of changes to a transaction
type TransactionHistory struct {
	Id            uint32           `json:"id" gorm:"type:uint32;primary_key"`
	TransactionId uint32           `json:"transaction_id" gorm:"type:uint32"`
	Action        TransactionAction `json:"action" gorm:"type:varchar(20)"`
	UserId        uint32           `json:"user_id" gorm:"type:uint32"`
	UpdatedAt     time.Time        `json:"updated_at" gorm:"type:datetime"`
} 