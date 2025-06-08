package repo

import (
	"context"
	"social/internal/domain"
	"social/internal/dto"
)

type NewsFeedRepo interface {
	Create(ctx context.Context, newsFeed *domain.NewsFeed) (*domain.NewsFeed, error)
	Update(ctx context.Context, newsFeed *domain.NewsFeed) (*domain.NewsFeed, error)
	Delete(ctx context.Context, newsFeed *domain.NewsFeed) error
	GetByID(ctx context.Context, id uint64) (*domain.NewsFeed, error)
	GetGlobal(ctx context.Context, filter *dto.NewsFeedGlobalSearch) ([]*domain.NewsFeed, int32, error)
	UpdateLike(ctx context.Context, newFeedID uint64, numLike int) error
	UpdateComment(ctx context.Context, newFeedID uint64, numComment int) error
	UpdateShare(ctx context.Context, newFeedID uint64, numShare int) error
	UpdateView(ctx context.Context, newFeedID uint64, numView int) error
	UserAvaiableReact(ctx context.Context, profileID uint64, id uint64) (bool, error)
}
