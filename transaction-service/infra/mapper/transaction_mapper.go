package mapper

import (
	transactionpb "pb/types/transaction"
	domain "tx/internal/domain"
)

type TransactionMapper struct {
}

func NewTransactionMapper() *TransactionMapper {
	return &TransactionMapper{}
}

func (m *TransactionMapper) ToTransaction(req *transactionpb.CreateTransactionRequest) *domain.Transaction {
	return &domain.Transaction{}
}
