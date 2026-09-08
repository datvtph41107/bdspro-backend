package repo

import (
	"context"
	"social/internal/domain"
)

type NewsFeedOfUserRepo interface {
	Create(ctx context.Context, newsFeedOfUser *domain.NewsFeedOfUser) error
	GetByUserID(ctx context.Context, userID uint64, page, size int) ([]domain.NewsFeedOfUser, int64, error)
	GetByPostID(ctx context.Context, postID uint64) (*domain.NewsFeedOfUser, error)
	DeleteByPostID(ctx context.Context, postID uint64) error
}
