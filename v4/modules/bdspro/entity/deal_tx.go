package entity

import (
	"time"

	"gorm.io/gorm"
)

// TxTransaction maps to tx_transaction.
type TxTransaction struct {
	ID              uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedAt       *time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt       *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at;index"`
	CreatedBy       *uint64        `gorm:"column:created_by"`
	UpdatedBy       *uint64        `gorm:"column:updated_by"`
	TransactionName string         `gorm:"column:transaction_name"`
	FromID          uint64         `gorm:"column:from_id;default:0"`
	FromOf          uint32         `gorm:"column:from_of;default:10"`
	ToID            uint64         `gorm:"column:to_id;default:0"`
	ToOf            uint32         `gorm:"column:to_of;default:10"`
	Method          uint32         `gorm:"column:method;default:10"`
	Status          uint32         `gorm:"column:status;default:10"`
	Timestamp       *time.Time     `gorm:"column:timestamp"`
	LastActionID    *uint64        `gorm:"column:last_action_id"`
}

func (TxTransaction) TableName() string { return "tx_transaction" }

// TxTransactionAction maps to tx_transaction_action.
type TxTransactionAction struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedAt *time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
	CreatedBy *uint64        `gorm:"column:created_by"`
	UpdatedBy *uint64        `gorm:"column:updated_by"`
	TxID      uint64         `gorm:"column:tx_id;index;not null"`
	FromID    uint64         `gorm:"column:from_id;default:0"`
	FromOf    uint32         `gorm:"column:from_of;default:10"`
	ToID      uint64         `gorm:"column:to_id;default:0"`
	ToOf      uint32         `gorm:"column:to_of;default:10"`
	Action    uint32         `gorm:"column:action;default:10"`
	Value     float64        `gorm:"column:value;default:0"`
	Note      string         `gorm:"column:note;type:text"`
	Message   string         `gorm:"column:message;type:text"`
	Timestamp *time.Time     `gorm:"column:timestamp"`
}

func (TxTransactionAction) TableName() string { return "tx_transaction_action" }

// TxContractDeal maps to tx_contract_deals.
type TxContractDeal struct {
	ID            uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedAt     *time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;index"`
	CreatedBy     *uint64        `gorm:"column:created_by"`
	UpdatedBy     *uint64        `gorm:"column:updated_by"`
	DealID        uint64         `gorm:"column:deal_id;index;not null"`
	Type          uint32         `gorm:"column:type;default:10"`
	Note          string         `gorm:"column:note;type:text"`
	ProductID     *uint64        `gorm:"column:product_id"`
	TransactionID *uint64        `gorm:"column:transaction_id;index"`
	CustomerID    uint64         `gorm:"column:customer_id;default:0"`
	ApprovedBy    *uint64        `gorm:"column:approved_by"`
	ApprovedAt    *time.Time     `gorm:"column:approved_at"`
	RejectReason  *string        `gorm:"column:reject_reason;type:text"`
}

func (TxContractDeal) TableName() string { return "tx_contract_deals" }
