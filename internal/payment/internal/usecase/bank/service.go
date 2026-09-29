package bank

import (
	"context"

	"payment/internal/domain/bank"
)

type BankUsecase struct {
	bankRepository Repository
}

func NewBankUsecase(bankRepository Repository) *BankUsecase {
	return &BankUsecase{
		bankRepository: bankRepository,
	}
}

func (b *BankUsecase) GetAll(ctx context.Context) ([]*bank.Bank, error) {
	return b.bankRepository.GetAll(ctx)
}

func (b *BankUsecase) GetByID(ctx context.Context, id uint32) (*bank.Bank, error) {
	return b.bankRepository.GetByID(ctx, id)
}

func (b *BankUsecase) Create(ctx context.Context, bank *bank.Bank) (*bank.Bank, error) {
	return b.bankRepository.Create(ctx, bank)
}

func (b *BankUsecase) Update(ctx context.Context, bank *bank.Bank) (*bank.Bank, error) {
	return b.bankRepository.Update(ctx, bank)
}

func (b *BankUsecase) Delete(ctx context.Context, id uint32) error {
	return b.bankRepository.Delete(ctx, id)
}
