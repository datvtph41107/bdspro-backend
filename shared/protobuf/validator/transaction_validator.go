package validator

import (
	"errors"
	transactionpb "pb/types/transaction"
)

type TransactionValidator struct {
}

func NewTransactionValidator() *TransactionValidator {
	return &TransactionValidator{}
}

func (v *TransactionValidator) ValidateProductTransactionRequest(req *transactionpb.CreateProductTransactionRequest) error {
	if req.FromId == 0 {
		return errors.New("fromId is required")
	}
	if req.ToId == 0 {
		return errors.New("toId is required")
	}
	if req.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	if req.Currency == "" {
		return errors.New("currency is required")
	}
	if req.Note == "" {
		return errors.New("note is required")
	}
	if req.TransactionType == "" {
		return errors.New("transactionType is required")
	}
	if req.DealId == 0 {
		return errors.New("dealId is required")
	}
	if req.CustomerId == 0 {
		return errors.New("customerId is required")
	}
	if req.ProductId == 0 {
		return errors.New("productId is required")
	}
	return nil
}
