package repo

import (
	_dto "common/domain/dto"
	"context"
	socialpb "pb/types/social"
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
	GetPublicByUserID(ctx context.Context, userId uint64, req *dto.NewsFeedGlobalSearch) ([]*dto.NewsFeedPublic, int32, error)
	GetGlobal(ctx context.Context, req *dto.NewsFeedGlobalSearch) ([]*dto.NewsFeedPublic, int32, error)
	GetListPostGroup(ctx context.Context, req *dto.NewsFeedGlobalSearch, groupId uint64) ([]*dto.NewsFeedPublic, int32, error)
	GetListPostOrganization(ctx context.Context, req *dto.NewsFeedGlobalSearch, organizationId uint64) ([]*dto.NewsFeedPublic, int32, error)
	GetListPostUser(ctx context.Context, req *dto.NewsFeedGlobalSearch, userId uint64) ([]*dto.NewsFeedPublic, int32, error)
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

	// Đếm số tin đăng theo ownerOf và ownerId
	CountByOwner(ctx context.Context, ownerOf socialpb.OwnerOf, ownerId uint64) (int64, error)

	// Admin methods
	HideNewsFeed(ctx context.Context, newsFeedId uint64, isHidden bool) error
	AdminGetList(ctx context.Context, req *dto.AdminNewsFeedSearch) ([]*dto.NewsFeedPublic, int32, error)
	AdminGetByID(ctx context.Context, id uint64) (*domain.NewsFeed, error)
	UpdateCreatedBy(ctx context.Context, newsFeedId uint64, createdBy uint64) error
}
