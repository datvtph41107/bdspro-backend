package repository

import (
	"context"

	"github.com/tuandq2112/payment-service/internal/domain/entity"
)

type WalletsRepository interface {
	GetWalletById(ctx context.Context, id uint32) (*entity.Wallet, error)
	GetWalletByIdForUpdate(ctx context.Context, id uint32) (*entity.Wallet, error)
	UpdateWallet(ctx context.Context, wallet *entity.Wallet) error
	BeginTx(ctx context.Context) Transaction
	CreateWallet(ctx context.Context, wallet *entity.Wallet) (*entity.Wallet, error)
	GetWalletByUserId(ctx context.Context, userId uint32) (*entity.Wallet, error)
}

type Transaction interface {
	Commit() error
	Rollback() error
}