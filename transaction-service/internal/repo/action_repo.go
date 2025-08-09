package repo

import (
	"context"
	"tx/internal/domain"
)

type ActionRepository interface {
	Transfer(ctx context.Context, action *domain.TxAction) (*domain.TxAction, error)
	Approve(ctx context.Context, action *domain.TxAction) (*domain.TxAction, error)
	Reject(ctx context.Context, action *domain.TxAction) (*domain.TxAction, error)
	Sign(ctx context.Context, action *domain.TxAction) (*domain.TxAction, error)
}
