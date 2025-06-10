package usecase

import (
	_routes "common/routes"
	_utils "common/utils"
	"context"
	"net/http"
	"social/internal/domain"
	"social/internal/dto"
	"social/internal/enums"
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

func (u *NewsFeedUsecase) CheckOwner(ctx context.Context, newsFeed *domain.NewsFeed) error {
	profileId := _utils.GetProfileIdWithContext(ctx)
	exist, err := u.newsFeedRepo.ExistByIDAndProfileID(ctx, newsFeed.ID, profileId)
	if err != nil {
		return err
	}
	if !exist {
		return &_routes.Except{
			Code:    http.StatusForbidden,
			Message: "bài viết không tồn tại hoặc bạn không có quyền cập nhật",
		}
	}
	return nil
}

func (u *NewsFeedUsecase) UpdatePost(ctx context.Context, newsFeed *domain.NewsFeed) (*domain.NewsFeed, error) {
	profileId := _utils.GetProfileIdWithContext(ctx)
	if newsFeed.PostID != nil {
		if err := u.bdsproClient.OwnerPost(ctx, *newsFeed.PostID, profileId); err != nil {
			return nil, err
		}
	}

	if err := u.CheckOwner(ctx, newsFeed); err != nil {
		return nil, err
	}

	err := u.transaction.WithTransaction(ctx, func(ctx context.Context) error {
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
		_, err := u.newsFeedRepo.Update(ctx, newsFeed)
		if err != nil {
			return &_routes.Except{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return newsFeed, nil
}

func (u *NewsFeedUsecase) DeletePost(ctx context.Context, post *domain.NewsFeed) error {
	profileId := ctx.Value("profile").(*uint64)
	if err := u.bdsproClient.OwnerPost(ctx, post.ID, *profileId); err != nil {
		return err
	}
	return u.newsFeedRepo.Delete(ctx, post)
}

func (u *NewsFeedUsecase) GetListPostByUserID(ctx context.Context, userId uint64, req *dto.NewsFeedGlobalSearch) ([]*domain.NewsFeed, int32, error) {
	newsFeeds, total, err := u.newsFeedRepo.GetPublicByUserID(ctx, userId, req)
	if err != nil {
		return nil, 0, err
	}

	return newsFeeds, int32(total), nil
}

func (u *NewsFeedUsecase) GetListPost(ctx context.Context, req *dto.NewsFeedGlobalSearch) ([]*dto.NewsFeedPublic, int32, error) {
	if req.Size > 10 || req.Size <= 0 {
		req.Size = 10
	}

	newsFeeds, total, err := u.newsFeedRepo.GetGlobal(ctx, req)
	if err != nil {
		return nil, 0, err
	}
	return newsFeeds, int32(total), nil
}

func (u *NewsFeedUsecase) ViewPost(ctx context.Context, newFeedID uint64, numView int) error {
	return u.newsFeedRepo.UpdateView(ctx, newFeedID, numView)
}

func (u *NewsFeedUsecase) GetNewsFeedByID(ctx context.Context, newsFeedID uint64) (*domain.NewsFeed, error) {
	profileId := _utils.GetProfileIdWithContext(ctx)
	newsFeed, err := u.newsFeedRepo.GetByID(ctx, newsFeedID)
	if err != nil {
		return nil, err
	}
	like, err := u.likeRepo.GetByTargetIdAndTypeAndUserId(ctx, newsFeedID, enums.TargetTypeNewsFeed, profileId)
	if like != nil && err == nil {
		newsFeed.LikeAt = like.CreatedAt
		newsFeed.DisLike = &like.DisLike
	}
	return newsFeed, nil
}

func (u *NewsFeedUsecase) UpdateVisibility(ctx context.Context, newsFeedID uint64, visibility enums.Visibility) (*domain.NewsFeed, error) {
	newsFeed, err := u.newsFeedRepo.GetByID(ctx, newsFeedID)
	if err != nil {
		return nil, err
	}
	if err := u.CheckOwner(ctx, newsFeed); err != nil {
		return nil, err
	}
	err = u.newsFeedRepo.UpdateVisibility(ctx, newsFeedID, visibility)
	if err != nil {
		return nil, err
	}
	return newsFeed, nil
}
