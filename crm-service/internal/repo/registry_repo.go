package repo

import (
	"context"
	"crm/internal/domain"
)

type RegistryRepo interface {
	Create(ctx context.Context, registry *domain.FeedbackRegistry) error
	List(ctx context.Context) ([]*domain.FeedbackRegistry, error)
	GetByRegistryName(ctx context.Context, registryName string) (*domain.FeedbackRegistry, error)
}