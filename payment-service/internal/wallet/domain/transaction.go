package domain

import "time"

type Transaction struct {
	// Deposite bool    `json:"deposite"`
	TransactionID   uint64     `gorm:"primaryKey;autoIncrement;not null" json:"transactionId"`
	Amount          float64    `json:"amount" binding:"required"`
	AccountId       uint64     `json:"accountId" binding:"required"`
	TransactionType uint       `json:"transactionType" binding:"required"`
	TargetId        uint64     `json:"targetId" binding:"required"`
	Timestamp       *time.Time `json:"timestamp"`
	DeletedAt       *time.Time `json:"deletedAt"`
	// Time
}

func (Transaction) TableName() string {
	return "wallet_transaction"
}
