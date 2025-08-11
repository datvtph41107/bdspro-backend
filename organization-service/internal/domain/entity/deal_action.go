package entity

import (
	_models "common/models"
	"time"
)

// DealAction represents a deal action entity
type DealAction struct {
	_models.BaseEntity
	DealID          uint64     `json:"dealId" db:"deal_id"`
	TransactionID   uint64     `json:"transactionId" db:"transaction_id"`
	CustomerID      uint64     `json:"customerId" db:"customer_id"`
	ProductID       uint64     `json:"productIds" db:"product_ids"`
	Amount          float64    `json:"amount" db:"amount"`
	TransactionStep uint32     `json:"transactionStep" db:"transaction_step"`
	TransactionType uint32     `json:"transactionType" db:"transaction_type"`
	Note            *string    `json:"note" db:"note"`
	Timestamp       *time.Time `json:"timestamp" db:"timestamp"`
	OrganizationID  uint64     `json:"organizationId" db:"organization_id"`
}
