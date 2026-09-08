package bank

import (
	"context"
	bankdomain "payment/internal/domain/bank"
)

type Repository interface {
	GetAll(ctx context.Context) ([]*bankdomain.Bank, error)
	GetByID(ctx context.Context, id uint32) (*bankdomain.Bank, error)
	Create(ctx context.Context, bank *bankdomain.Bank) (*bankdomain.Bank, error)
	Update(ctx context.Context, bank *bankdomain.Bank) (*bankdomain.Bank, error)
	Delete(ctx context.Context, id uint32) error
}
