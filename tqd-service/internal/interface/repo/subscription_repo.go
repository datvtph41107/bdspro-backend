package repo

import (
	"context"
	"tqd/internal/domain"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, sub *domain.UserSubscription) error
	Update(ctx context.Context, sub *domain.UserSubscription) error
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*domain.UserSubscription, error)
	GetByUserAndTarget(ctx context.Context, userID uint64, targetType string, targetID uint64) (*domain.UserSubscription, error)
	ListByUser(ctx context.Context, userID uint64, targetType *string, status *string, page, limit int) ([]domain.UserSubscription, int64, error)
	ListByTarget(ctx context.Context, targetType string, targetID uint64) ([]domain.UserSubscription, error)

	ListByUserWithStatus(ctx context.Context, userID uint64, status *uint32, page, limit int) ([]domain.UserSubscription, int64, error)
	AdminList(ctx context.Context, userID *uint64, targetType *string, status *uint32, page, limit int) ([]domain.UserSubscription, int64, error)
	AdminDelete(ctx context.Context, id uint64) error
}
