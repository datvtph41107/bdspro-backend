package iusecase

import (
	"context"
)

// TransactionClient defines the interface for transaction service client
type TransactionClient interface {
	// CreateTransaction creates a new transaction in transaction service
	CreateTransaction(ctx context.Context, req *CreateTransactionRequest) (*CreateTransactionResponse, error)
}

// CreateTransactionRequest represents the request to create a transaction
type CreateTransactionRequest struct {
	Amount          float64 `json:"amount"`
	CustomerID      uint64  `json:"customerId"`
	DealID          uint64  `json:"dealId"`
	Note            *string `json:"note"`
	ProductIDs      uint64  `json:"productIds"`
	Timestamp       string  `json:"timestamp"`
	TransactionStep uint32  `json:"transactionStep"`
	TransactionType uint32  `json:"transactionType"`
}

// CreateTransactionResponse represents the response from creating a transaction
type CreateTransactionResponse struct {
	TransactionID uint64 `json:"transactionId"`
}
