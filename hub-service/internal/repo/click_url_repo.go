package repo

import (
	"context"

	"hub/internal/domain"
)

type IClickURLRepo interface {
	Insert(ctx context.Context, event *domain.ClickURL) error
	InsertBatch(ctx context.Context, events []*domain.ClickURL) error
}
