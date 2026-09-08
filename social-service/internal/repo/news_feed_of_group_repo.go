package repo

import (
	"context"
	"social/internal/domain"
)

type NewsFeedOfGroupRepo interface {
	Create(ctx context.Context, newsFeedOfGroup *domain.NewsFeedOfGroup) error
	GetByGroupID(ctx context.Context, groupID uint64, page, size int) ([]domain.NewsFeedOfGroup, int64, error)
	GetByPostID(ctx context.Context, postID uint64) (*domain.NewsFeedOfGroup, error)
	DeleteByPostID(ctx context.Context, postID uint64) error
}
