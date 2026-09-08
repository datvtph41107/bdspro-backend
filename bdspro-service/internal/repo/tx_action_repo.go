package repo

import (
	"bdspro/internal/domain"
	"bdspro/internal/enums"
	"context"
)

type TxActionRepository interface {
	Transfer(ctx context.Context, action *domain.TxAction) (*domain.TxAction, error)
	Approve(ctx context.Context, action *domain.TxAction) (*domain.TxAction, error)
	Reject(ctx context.Context, action *domain.TxAction) (*domain.TxAction, error)
	Sign(ctx context.Context, action *domain.TxAction) (*domain.TxAction, error)
	CountByTransactionIdAndActions(ctx context.Context, transactionId uint64, actions []enums.TxAction) (int64, error)
}
