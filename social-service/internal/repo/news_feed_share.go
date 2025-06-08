package repo

import (
	"context"
	"social/internal/domain"
)

type NewsFeedShareRepo interface {
	Create(ctx context.Context, newsFeedShare *domain.NewsFeedShare) (*domain.NewsFeedShare, error)
}
