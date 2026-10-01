package repository

import (
	"context"
	domain "userv4/modules/crm/entity"
)

type RegistryRepo interface {
	Create(ctx context.Context, registry *domain.FeedbackRegistry) error
	List(ctx context.Context) ([]*domain.FeedbackRegistry, error)
	GetByRegistryName(ctx context.Context, registryName string) (*domain.FeedbackRegistry, error)
}