package usecase

import (
	"context"

	"payment/internal/domain/entity"
	"payment/internal/domain/repository"
)

type BankUsecase struct {
	bankRepository repository.BankRepository
}

func NewBankUsecase(bankRepository repository.BankRepository) *BankUsecase {
	return &BankUsecase{
		bankRepository: bankRepository,
	}
}

func (b *BankUsecase) GetAll(ctx context.Context) ([]*entity.Bank, error) {
	return b.bankRepository.GetAll(ctx)
}

func (b *BankUsecase) GetByID(ctx context.Context, id uint32) (*entity.Bank, error) {
	return b.bankRepository.GetByID(ctx, id)
}

func (b *BankUsecase) Create(ctx context.Context, bank *entity.Bank) (*entity.Bank, error) {
	return b.bankRepository.Create(ctx, bank)
}

func (b *BankUsecase) Update(ctx context.Context, bank *entity.Bank) (*entity.Bank, error) {
	return b.bankRepository.Update(ctx, bank)
}

func (b *BankUsecase) Delete(ctx context.Context, id uint32) error {
	return b.bankRepository.Delete(ctx, id)
}
