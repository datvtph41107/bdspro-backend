package dto

import (
	_dto "common/dto"
	"time"
)

// CreateDealActionRequest represents the request to create a deal action
type CreateDealActionRequest struct {
	DealID          uint64    `json:"dealId" validate:"required"`
	TransactionID   uint64    `json:"transactionId" validate:"required"`
	CustomerID      uint64    `json:"customerId" validate:"required"`
	ProductIDs      uint64    `json:"productIds" validate:"required"`
	Amount          float64   `json:"amount" validate:"required"`
	TransactionStep uint32    `json:"transactionStep" validate:"required"`
	TransactionType uint32    `json:"transactionType" validate:"required"`
	Note            *string   `json:"note"`
	Timestamp       time.Time `json:"timestamp" validate:"required"`
}

// DealActionResponse represents the response for deal action
type DealActionResponse struct {
	ID              uint64  `json:"id"`
	DealID          uint64  `json:"dealId"`
	TransactionID   uint64  `json:"transactionId"`
	CustomerID      uint64  `json:"customerId"`
	ProductIDs      uint64  `json:"productIds"`
	Amount          float64 `json:"amount"`
	TransactionStep uint32  `json:"transactionStep"`
	TransactionType uint32  `json:"transactionType"`
	Note            *string `json:"note"`
	Timestamp       string  `json:"timestamp"`
	CreatedAt       string  `json:"createdAt"`
	UpdatedAt       string  `json:"updatedAt"`
	CreatedBy       uint64  `json:"createdBy"`
	OrganizationID  uint64  `json:"organizationId"`
}

// GetDealActionsRequest represents the request to get deal actions
type GetDealActionsRequest struct {
	_dto.Pagable
	DealID          *uint64 `json:"dealId"`
	CustomerID      *uint64 `json:"customerId"`
	TransactionType *uint32 `json:"transactionType"`
}

// GetDealActionsResponse represents the response for deal actions list
type GetDealActionsResponse struct {
	Data  []DealActionResponse `json:"data"`
	Total int32                `json:"total"`
}
