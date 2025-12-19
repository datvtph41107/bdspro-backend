package repo

import (
	"bdspro/internal/domain"
	"context"
)

type BankAccountRepository interface {
	Create(ctx context.Context, bankAccount *domain.BankAccount) (*domain.BankAccount, error)
	GetByID(ctx context.Context, id *uint64) (*domain.BankAccount, error)
	DealUpdateAccount(ctx context.Context, accountId *uint64, model *domain.BankAccount) (*domain.BankAccount, error)
}
