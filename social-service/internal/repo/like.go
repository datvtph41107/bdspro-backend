package repo

import (
	"context"
	"social/internal/domain"
	"social/internal/enums"
)

type LikeRepo interface {
	Create(ctx context.Context, like *domain.Like) (*domain.Like, error)
	GetByTargetIdAndTypeAndUserId(ctx context.Context, targetId uint64, targetType enums.TargetType, userId uint64) (*domain.Like, error)
	ExistByTargetIdAndTypeAndUserId(ctx context.Context, targetId uint64, targetType enums.TargetType, userId uint64) (bool, error)
	Delete(ctx context.Context, like *domain.Like) error
	Count(ctx context.Context, targetId uint64, targetType enums.TargetType, dislike bool) (int, error)
}
