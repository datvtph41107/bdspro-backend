package repository

import (
	"context"
	"organization/internal/domain/entity"
)

type BankAccountRepository interface {
	Create(ctx context.Context, bankAccount *entity.BankAccount) (*entity.BankAccount, error)
}
