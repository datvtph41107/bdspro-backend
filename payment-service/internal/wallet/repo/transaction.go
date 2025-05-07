package repo

import "payment/internal/wallet/domain"

type TransactionRepo interface {
	GetByID(id uint64) (*domain.Transaction, error)
	Create(account *domain.Transaction) (*domain.Transaction, error)
	Balance(accountId *uint64) (float64, error)
	ExistByID(id uint64) (bool, error)
}
