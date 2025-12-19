package domain

import (
	"bdspro/internal/enums"
	_models "common/models"
	"time"
)

type Tx struct {
	_models.BaseEntity
	TransactionName string
	FromID          uint64
	FromOf          enums.TxOwnerType
	ToID            uint64
	ToOf            enums.TxOwnerType
	Method          enums.TxMethod
	Status          enums.TxStatus
	Timestamp       time.Time
	LastActionID    *uint64
	LastAction      *TxAction  `gorm:"foreignKey:LastActionID"`
	Actions         []TxAction `gorm:"foreignKey:TxId"`
	Note            string     `gorm:"-"`
	Value           float64    `gorm:"-"`
}

func (Tx) TableName() string {
	return "tx_transaction"
}

// TransferTransactionRequest represents the request for creating a transfer transaction
type TransferTransactionRequest struct {
	FromID          uint64
	ToID            uint64
	Amount          float64
	Currency        string
	Note            string
	TransactionType string
	DealID          uint64
	CustomerID      uint64
	ProductIDs      uint64
}

// SignContractRequest represents the request for signing a contract
type SignContractRequest struct {
	DealID        uint64
	CustomerID    uint64
	ProductID     uint64
	ContractValue float64
	Currency      string
	Note          string
	ContractType  string
	SignDate      time.Time
}

// PaymentRequest represents the request for making a payment
type PaymentRequest struct {
	DealID        uint64
	CustomerID    uint64
	ProductID     uint64
	Amount        float64
	Currency      string
	PaymentMethod string
	Note          string
	PaymentDate   time.Time
}

// HandOverRequest represents the request for handover process
type HandOverRequest struct {
	DealID       uint64
	CustomerID   uint64
	ProductID    uint64
	HandOverDate time.Time
	Note         string
	HandOverType string
	Location     string
}
