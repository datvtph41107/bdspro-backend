package usecases

import (
	"fmt"
	"payment/internal/wallet/domain"
	"payment/internal/wallet/repo"
	"time"

	"errors"
)

type AcountUsecase struct {
	repo repo.AccountRepo
}

func NewAccountUsecase(repo repo.AccountRepo) *AcountUsecase {
	return &AcountUsecase{repo: repo}
}

func (uc *AcountUsecase) Register(account *domain.Account) (*domain.Account, error) {
	now := time.Now()
	existed, err := uc.repo.ExistByID(account.UserId)
	if err != nil {
		return nil, err
	}
	if existed {
		return nil, errors.New("tài khoản đã tồn tại")
	}
	return uc.repo.Create(&domain.Account{
		Fullname:      account.Fullname,
		Phone:         account.Phone,
		UserId:        account.UserId,
		AccountNumber: fmt.Sprintf("500%d", account.UserId),
		CreatedAt:     &now,
		UpdatedAt:     &now,
	})
}
