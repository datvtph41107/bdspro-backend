package iusecase

import (
	"context"
)

type ITransactionClient interface {
	CreateTransaction(ctx context.Context, transactionID uint64) error
}
