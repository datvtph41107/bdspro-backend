package usecases

import (
	"errors"
	"payment/internal/wallet/domain"
	"payment/internal/wallet/repo"
	"time"
)

type TransactionUsecase struct {
	repo repo.TransactionRepo
}

func NewTransactionUsecase(repo repo.TransactionRepo) *TransactionUsecase {
	return &TransactionUsecase{repo: repo}
}

func (uc *TransactionUsecase) MakeTransaction(tx *domain.Transaction) (*domain.Transaction, error) {
	amount, err := uc.repo.Balance(&tx.AccountId)
	if err != nil {
		return nil, err
	}

	if amount+tx.Amount < 0 {
		return nil, errors.New("số dư không đủ")
	}

	now := time.Now()
	domain, err := uc.repo.Create(&domain.Transaction{
		Amount:    tx.Amount,
		AccountId: tx.AccountId,
		Timestamp: &now,
	})

	return domain, err
}

func (uc *TransactionUsecase) Balance(userId *uint64) (float64, error) {
	amount, err := uc.repo.Balance(userId)
	if err != nil {
		return 0, err
	}

	return amount, err
}
