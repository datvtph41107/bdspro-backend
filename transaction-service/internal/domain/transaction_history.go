package domain

import (
	"time"
	"tx/internal/enums"
)

// TransactionHistory represents the history of changes to a transaction
type TransactionHistory struct {
	Id            uint32
	TransactionId uint32
	Action        enums.EMethod
	UserId        uint32
	UpdatedAt     time.Time
}
