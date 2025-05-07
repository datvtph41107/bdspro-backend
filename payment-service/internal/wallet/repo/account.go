package repo

import "payment/internal/wallet/domain"

type AccountRepo interface {
	Create(account *domain.Account) (*domain.Account, error)
	Update(account *domain.Account) (*domain.Account, error)
	ExistByID(id uint64) (bool, error)
	Lock(accountId uint64) error
}
