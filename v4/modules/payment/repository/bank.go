package repository

import (
	"context"

	"userv4/modules/payment/entity"
)

// BankRepository reads the bank catalog owned by payment-service schema.
type BankRepository interface {
	ListActive(ctx context.Context) ([]*entity.Bank, error)
}
