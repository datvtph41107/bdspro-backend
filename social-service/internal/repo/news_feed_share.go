package repo

import (
	"context"
	"social/internal/domain"
	"social/internal/enums"
)

type NewsFeedShareRepo interface {
	Create(ctx context.Context, newsFeedShare *domain.NewsFeedShare) (*domain.NewsFeedShare, error)
	GetNearestUserIds(ctx context.Context, newsFeedId uint64, targetType enums.TargetType) ([]uint64, error)
	SumUserComment(ctx context.Context, newsFeedId uint64, targetType enums.TargetType) (int, error)
}
