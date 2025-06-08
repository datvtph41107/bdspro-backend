package usecase

import (
	_routes "common/routes"
	_utils "common/utils"
	"context"
	"net/http"
	"social/internal/domain"
	"social/internal/dto"
	iusecase "social/internal/interface"
	"social/internal/repo"
)

type NewsFeedUsecase struct {
	newsFeedRepo      repo.NewsFeedRepo
	likeRepo          repo.LikeRepo
	permissionUsecase iusecase.PermissionUsecase
	friendTagRepo     repo.FriendTagRepo
	newsFeedMediaRepo repo.NewsFeedMediaRepo
	bdsproClient      iusecase.BdsproClient
	userClient        iusecase.UserClient
	transaction       iusecase.ITransaction
}

func NewNewsFeedUsecase(postRepo repo.NewsFeedRepo,
	likeRepo repo.LikeRepo,
	permissionUsecase iusecase.PermissionUsecase,
	friendTagRepo repo.FriendTagRepo,
	newsFeedMediaRepo repo.NewsFeedMediaRepo,
	bdsproClient iusecase.BdsproClient,
	userClient iusecase.UserClient,
	transaction iusecase.ITransaction,
) *NewsFeedUsecase {
	return &NewsFeedUsecase{
		newsFeedRepo:      postRepo,
		likeRepo:          likeRepo,
		permissionUsecase: permissionUsecase,
		friendTagRepo:     friendTagRepo,
		newsFeedMediaRepo: newsFeedMediaRepo,
		bdsproClient:      bdsproClient,
		userClient:        userClient,
		transaction:       transaction,
	}
}

func (u *NewsFeedUsecase) CreateNewsFeed(ctx context.Context, newsFeed *domain.NewsFeed) (*domain.NewsFeed, error) {
	profileId := _utils.GetProfileIdWithContext(ctx)

	if newsFeed.PostID != nil {
		if err := u.bdsproClient.OwnerPost(ctx, *newsFeed.PostID, profileId); err != nil {
			return nil, &_routes.Except{
				Code:    http.StatusForbidden,
				Message: "bạn không có quyền tạo/forward tin đăng này",
			}
		}
	}

	err := u.transaction.WithTransaction(ctx, func(ctx context.Context) error {
		createdFeed, err := u.newsFeedRepo.Create(ctx, newsFeed)
		if err != nil {
			return &_routes.Except{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
		}
		newsFeed = createdFeed

		if len(newsFeed.FriendTagIds) > 0 {
			isFriends, err := u.userClient.IsFriends(ctx, profileId, newsFeed.FriendTagIds)
			if err != nil {
				return &_routes.Except{
					Code:    http.StatusInternalServerError,
					Message: err.Error(),
				}
			}
			if !isFriends {
				return &_routes.Except{
					Code:    http.StatusForbidden,
					Message: "not friends",
				}
			}
			if err := u.friendTagRepo.UpdateFriendTag(ctx, newsFeed.ID, newsFeed.FriendTagIds); err != nil {
				return &_routes.Except{
					Code:    http.StatusInternalServerError,
					Message: err.Error(),
				}
			}
		}

		if len(newsFeed.NewsFeedMedias) > 0 {
			if err := u.newsFeedMediaRepo.UpdateMedia(ctx, newsFeed.ID, newsFeed.NewsFeedMedias); err != nil {
				return &_routes.Except{
					Code:    http.StatusInternalServerError,
					Message: err.Error(),
				}
			}
		}

		return nil // commit nếu mọi thứ ok
	})

	if err != nil {
		return nil, err // rollback nếu lỗi
	}
	return newsFeed, nil
}

func (u *NewsFeedUsecase) UpdatePost(ctx context.Context, post *domain.NewsFeed) (*domain.NewsFeed, error) {
	profileId := ctx.Value("profile").(*uint64)
	if err := u.bdsproClient.OwnerPost(ctx, post.ID, *profileId); err != nil {
		return nil, err
	}
	return u.newsFeedRepo.Update(ctx, post)
}

func (u *NewsFeedUsecase) DeletePost(ctx context.Context, post *domain.NewsFeed) error {
	profileId := ctx.Value("profile").(*uint64)
	if err := u.bdsproClient.OwnerPost(ctx, post.ID, *profileId); err != nil {
		return err
	}
	return u.newsFeedRepo.Delete(ctx, post)
}

func (u *NewsFeedUsecase) GetListPost(ctx context.Context, dto *dto.NewsFeedGlobalSearch) ([]*domain.NewsFeed, int32, error) {
	// tối đa 10 bản ghi
	if dto.Size > 10 || dto.Size <= 0 {
		dto.Size = 10
	}
	return u.newsFeedRepo.GetGlobal(ctx, dto)
}

func (u *NewsFeedUsecase) ViewPost(ctx context.Context, newFeedID uint64, numView int) error {
	return u.newsFeedRepo.UpdateView(ctx, newFeedID, numView)
}
