package repo

import (
	"context"
	"hub/internal/domain"
	"hub/internal/dto"
)

type IInteractiveEventRepo interface {
	Insert(ctx context.Context, event *domain.InteractiveEvent) error
	InsertBatch(ctx context.Context, events []*domain.InteractiveEvent) error

	GetProductStatsView(ctx context.Context, productID uint64, lastEventTime int64) (*dto.ViewStatsProduct, error)
}
