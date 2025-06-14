package transaction

import (
	"time"
)

// TransactionProduct represents the products associated with a transaction
type TransactionProduct struct {
	Id            uint32    `json:"id" gorm:"type:uint32;primary_key"`
	TransactionId uint32    `json:"transaction_id" gorm:"type:uint32"`
	ProductId     uint32    `json:"product_id" gorm:"type:uint32"`
	CreatedAt     time.Time `json:"created_at" gorm:"type:timestamp"`
} 