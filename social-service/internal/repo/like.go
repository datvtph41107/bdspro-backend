package repo

import (
	"context"
	"social/internal/domain"
)

type LikeRepo interface {
	Create(ctx context.Context, like *domain.Like) (*domain.Like, error)
	ExistByTargetIdAndTypeAndUserId(ctx context.Context, targetId uint64, likeType uint32, userId uint64) (bool, error)
}
