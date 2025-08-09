package repository

import (
	"context"

	"payment/internal/domain/entity"
)

type BankRepository interface {
	GetAll(ctx context.Context) ([]*entity.Bank, error)
	GetByID(ctx context.Context, id uint32) (*entity.Bank, error)
	Create(ctx context.Context, bank *entity.Bank) (*entity.Bank, error)
	Update(ctx context.Context, bank *entity.Bank) (*entity.Bank, error)
	Delete(ctx context.Context, id uint32) error
}
