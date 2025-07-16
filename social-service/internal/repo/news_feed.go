package repo

import (
	_dto "common/dto"
	"context"
	"social/internal/domain"
	"social/internal/dto"
	"social/internal/enums"
	"time"
)

type NewsFeedRepo interface {
	Create(ctx context.Context, newsFeed *domain.NewsFeed) (*domain.NewsFeed, error)
	Update(ctx context.Context, newsFeed *domain.NewsFeed) (*domain.NewsFeed, error)
	UpdateVisibility(ctx context.Context, newsFeedId uint64, visibility enums.Visibility) error
	Delete(ctx context.Context, newsFeed *domain.NewsFeed) error
	GetByID(ctx context.Context, id uint64) (*domain.NewsFeed, error)
	GetPublicIgnorePreloadByID(ctx context.Context, id uint64) (*domain.NewsFeed, error)
	GetPublicByUserID(ctx context.Context, userId uint64, req *dto.NewsFeedGlobalSearch) ([]*domain.NewsFeed, int32, error)
	GetGlobal(ctx context.Context, req *dto.NewsFeedGlobalSearch) ([]*dto.NewsFeedPublic, int32, error)
	UpdateComment(ctx context.Context, newFeedID uint64, numComment int) error
	UpdateShare(ctx context.Context, newFeedID uint64, numShare int) error
	UpdateView(ctx context.Context, newFeedID uint64, numView int) error
	UserAvaiableReact(ctx context.Context, profileID uint64, id uint64) (bool, error)
	ExistByID(ctx context.Context, id uint64) (bool, error)
	ExistByIDAndProfileID(ctx context.Context, id uint64, profileID uint64) (bool, error)
	GetLastUpdated(ctx context.Context, id uint64) (*time.Time, error)
	UpdateLikeNumber(ctx context.Context, newFeedID uint64, totalLike int) error
	UpdateCommentNumber(ctx context.Context, newFeedID uint64, totalComment int) error
	IncrementNumberShare(ctx context.Context, newFeedID uint64) error
	SyncData(ctx context.Context) error
	GetReel(ctx context.Context, dto _dto.Pagable) ([]*domain.ReelEntity, int64, error)
}
