package repository

import (
	"context"
	"organization/internal/domain/entity"
)

type BankAccountRepository interface {
	Create(ctx context.Context, bankAccount *entity.BankAccount) (*entity.BankAccount, error)
	GetByID(ctx context.Context, id *uint64) (*entity.BankAccount, error)
	DealUpdateAccount(ctx context.Context, accountId *uint64, model *entity.BankAccount) (*entity.BankAccount, error)
}
