package repo

import (
	"context"
	"hub/internal/domain"
)

type IInteractiveEventRepo interface {
	Insert(ctx context.Context, event *domain.InteractiveEvent) error
	InsertBatch(ctx context.Context, events []*domain.InteractiveEvent) error
}
