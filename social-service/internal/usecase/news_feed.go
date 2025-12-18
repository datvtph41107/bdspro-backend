package usecase

import (
	_dto "common/domain/dto"
	_errors "common/errors"
	_routes "common/routes"
	_utils "common/utils"
	"context"
	"log"
	"net/http"
	socialpb "pb/types/social"
	"social/internal/domain"
	"social/internal/dto"
	"social/internal/enums"
	iusecase "social/internal/interface"
	"social/internal/repo"
)

type NewsFeedUsecase struct {
	newsFeedRepo        repo.NewsFeedRepo
	likeRepo            repo.LikeRepo
	permissionUsecase   iusecase.PermissionUsecase
	friendTagRepo       repo.FriendTagRepo
	newsFeedMediaRepo   repo.NewsFeedMediaRepo
	newsFeedOfUserRepo  repo.NewsFeedOfUserRepo
	newsFeedOfGroupRepo repo.NewsFeedOfGroupRepo
	bdsproClient        iusecase.BdsproClient
	userClient          iusecase.UserClient
	transaction         iusecase.ITransaction
}

func NewNewsFeedUsecase(postRepo repo.NewsFeedRepo,
	likeRepo repo.LikeRepo,
	permissionUsecase iusecase.PermissionUsecase,
	friendTagRepo repo.FriendTagRepo,
	newsFeedMediaRepo repo.NewsFeedMediaRepo,
	newsFeedOfUserRepo repo.NewsFeedOfUserRepo,
	newsFeedOfGroupRepo repo.NewsFeedOfGroupRepo,
	bdsproClient iusecase.BdsproClient,
	userClient iusecase.UserClient,
	transaction iusecase.ITransaction,
) *NewsFeedUsecase {
	return &NewsFeedUsecase{
		newsFeedRepo:        postRepo,
		likeRepo:            likeRepo,
		permissionUsecase:   permissionUsecase,
		friendTagRepo:       friendTagRepo,
		newsFeedMediaRepo:   newsFeedMediaRepo,
		newsFeedOfUserRepo:  newsFeedOfUserRepo,
		newsFeedOfGroupRepo: newsFeedOfGroupRepo,
		bdsproClient:        bdsproClient,
		userClient:          userClient,
		transaction:         transaction,
	}
}

func (u *NewsFeedUsecase) CreateGroupNewsFeed(ctx context.Context, newsFeed *domain.NewsFeed, groupId *uint64) (*domain.NewsFeed, error) {
	if groupId == nil {
		return nil, _errors.ReturnError(http.StatusBadRequest, "groupId is required")
	}
	ownerType := enums.OwnerOfGroup
	ownerID := groupId
	return u.createNewsFeed(ctx, newsFeed, ownerType, ownerID)
}

func (u *NewsFeedUsecase) CreateUserNewsFeed(ctx context.Context, newsFeed *domain.NewsFeed) (*domain.NewsFeed, error) {
	profileId := _utils.GetProfileIdWithContext(ctx)
	ownerType := enums.OwnerOfUser
	ownerID := profileId
	return u.createNewsFeed(ctx, newsFeed, ownerType, &ownerID)
}

func (u *NewsFeedUsecase) createNewsFeed(ctx context.Context, newsFeed *domain.NewsFeed, ownerType enums.OwnerOf, ownerID *uint64) (*domain.NewsFeed, error) {
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
		newsFeed.OwnerOf = ownerType
		newsFeed.OwnerID = ownerID
		createdFeed, err := u.newsFeedRepo.Create(ctx, newsFeed)
		if err != nil {
			return &_routes.Except{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
		}
		newsFeed = createdFeed
		if newsFeed.ParentID != nil {
			// This could be a group post - you might want to add more logic here
			// For example, check if ParentID refers to a group
			ownerType = enums.OwnerOfGroup
			ownerID = newsFeed.ParentID
		}

		// Create tracking record based on owner type
		if ownerType == enums.OwnerOfUser {
			newsFeedOfUser := &domain.NewsFeedOfUser{
				PostID: createdFeed.ID,
				UserID: *ownerID,
			}
			if err := u.newsFeedOfUserRepo.Create(ctx, newsFeedOfUser); err != nil {
				return &_routes.Except{
					Code:    http.StatusInternalServerError,
					Message: "không thể track tin đăng cho user",
				}
			}
		} else if ownerType == enums.OwnerOfGroup {
			newsFeedOfGroup := &domain.NewsFeedOfGroup{
				PostID:  createdFeed.ID,
				GroupID: *ownerID,
			}
			if err := u.newsFeedOfGroupRepo.Create(ctx, newsFeedOfGroup); err != nil {
				return &_routes.Except{
					Code:    http.StatusInternalServerError,
					Message: "không thể track tin đăng cho group",
				}
			}
		}

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

		// if len(newsFeed.NewsFeedMedias) > 0 {
		// 	if err := u.newsFeedMediaRepo.UpdateMedia(ctx, newsFeed.ID, newsFeed.NewsFeedMedias); err != nil {
		// 		return &_routes.Except{
		// 			Code:    http.StatusInternalServerError,
		// 			Message: err.Error(),
		// 		}
		// 	}
		// }

		return nil
	})

	if err != nil {
		return nil, err
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

func (u *NewsFeedUsecase) DeletePost(ctx context.Context, newsFeedID uint64) error {
	profileId := _utils.GetProfileIdWithContext(ctx)
	newsfeed, err := u.newsFeedRepo.GetByID(ctx, newsFeedID)
	if err != nil {
		return err
	}
	if newsfeed.CreatedBy != nil && *newsfeed.CreatedBy != profileId {
		return &_routes.Except{
			Code:    http.StatusForbidden,
			Message: "bạn không có quyền xóa bài viết này",
		}
	}
	return u.newsFeedRepo.Delete(ctx, newsfeed)
}

func (u *NewsFeedUsecase) GetListPostByUserID(ctx context.Context, userId uint64, req *dto.NewsFeedGlobalSearch) ([]*dto.NewsFeedPublic, int32, error) {
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

func (u *NewsFeedUsecase) GetListPostOrganization(ctx context.Context, organizationId uint64, req *dto.NewsFeedGlobalSearch) ([]*dto.NewsFeedPublic, int32, error) {
	if req.Size > 10 || req.Size <= 0 {
		req.Size = 10
	}
	newsFeeds, total, err := u.newsFeedRepo.GetListPostOrganization(ctx, req, organizationId)
	if err != nil {
		return nil, 0, err
	}
	return newsFeeds, int32(total), nil
}

func (u *NewsFeedUsecase) GetListPostGroup(ctx context.Context, groupId uint64, req *dto.NewsFeedGlobalSearch) ([]*dto.NewsFeedPublic, int32, error) {
	if req.Size > 10 || req.Size <= 0 {
		req.Size = 10
	}
	newsFeeds, total, err := u.newsFeedRepo.GetListPostGroup(ctx, req, groupId)
	if err != nil {
		return nil, 0, err
	}
	return newsFeeds, int32(total), nil
}

func (u *NewsFeedUsecase) GetListPostUser(ctx context.Context, req *dto.NewsFeedGlobalSearch) ([]*dto.NewsFeedPublic, int32, error) {
	if req.Size > 10 || req.Size <= 0 {
		req.Size = 10
	}
	userId := _utils.GetProfileIdWithContext(ctx)
	newsFeeds, total, err := u.newsFeedRepo.GetListPostUser(ctx, req, userId)
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
		newsFeed.LikedAt = like.CreatedAt
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

func (u *NewsFeedUsecase) SyncNewsFeed(ctx context.Context) {
	log.Printf("SyncNewsFeed start")
	err := u.newsFeedRepo.SyncData(ctx)
	if err != nil {
		log.Printf("SyncNewsFeed error: %v", err)
	}
	log.Printf("SyncNewsFeed success")
}

func (u *NewsFeedUsecase) GetReelNewsFeed(ctx context.Context, req _dto.Pagable) ([]*domain.ReelEntity, int64, error) {
	newsFeed, count, err := u.newsFeedRepo.GetReel(ctx, req)
	if err != nil {
		return nil, 0, err
	}
	return newsFeed, count, nil
}

// CountByOwner đếm số tin đăng theo ownerOf và ownerId
func (u *NewsFeedUsecase) CountByOwner(ctx context.Context, ownerOf socialpb.OwnerOf, ownerId uint64) (uint32, error) {
	count, err := u.newsFeedRepo.CountByOwner(ctx, ownerOf, ownerId)
	if err != nil {
		return 0, err
	}
	return uint32(count), nil
}

// AdminCreateNewsFeed - Admin tạo bài viết
func (u *NewsFeedUsecase) AdminCreateNewsFeed(ctx context.Context, newsFeed *domain.NewsFeed) (*domain.NewsFeed, error) {
	if !u.permissionUsecase.IsAdmin(ctx) {
		return nil, &_routes.Except{
			Code:    http.StatusForbidden,
			Message: "Bạn không có quyền admin",
		}
	}

	// Lưu createdBy từ request (nếu có) để admin có thể set user khác làm creator
	requestCreatedBy := newsFeed.CreatedBy

	// Admin có thể tạo bài viết cho bất kỳ owner nào
	ownerType := newsFeed.OwnerOf
	ownerID := newsFeed.OwnerID
	result, err := u.createNewsFeed(ctx, newsFeed, ownerType, ownerID)
	if err != nil {
		return nil, err
	}

	// Nếu admin đã set createdBy từ request, update lại sau khi create
	// (vì GORM callback sẽ tự động set createdBy từ context)
	if requestCreatedBy != nil && *requestCreatedBy > 0 {
		err = u.transaction.WithTransaction(ctx, func(ctx context.Context) error {
			return u.newsFeedRepo.UpdateCreatedBy(ctx, result.ID, *requestCreatedBy)
		})
		if err != nil {
			return nil, err
		}
		result.CreatedBy = requestCreatedBy
	}

	return result, nil
}

// AdminUpdateNewsFeed - Admin sửa bài viết
func (u *NewsFeedUsecase) AdminUpdateNewsFeed(ctx context.Context, newsFeed *domain.NewsFeed) (*domain.NewsFeed, error) {
	if !u.permissionUsecase.IsAdmin(ctx) {
		return nil, &_routes.Except{
			Code:    http.StatusForbidden,
			Message: "Bạn không có quyền admin",
		}
	}

	// Admin có thể sửa bất kỳ bài viết nào, không cần check owner
	err := u.transaction.WithTransaction(ctx, func(ctx context.Context) error {
		if len(newsFeed.FriendTagIds) > 0 {
			profileId := _utils.GetProfileIdWithContext(ctx)
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

// AdminDeleteNewsFeed - Admin xóa bài viết
func (u *NewsFeedUsecase) AdminDeleteNewsFeed(ctx context.Context, newsFeedID uint64) error {
	if !u.permissionUsecase.IsAdmin(ctx) {
		return &_routes.Except{
			Code:    http.StatusForbidden,
			Message: "Bạn không có quyền admin",
		}
	}

	newsfeed, err := u.newsFeedRepo.GetByID(ctx, newsFeedID)
	if err != nil {
		return err
	}
	// Admin có thể xóa bất kỳ bài viết nào
	return u.newsFeedRepo.Delete(ctx, newsfeed)
}

// AdminHideNewsFeed - Admin ẩn/hiện bài viết
func (u *NewsFeedUsecase) AdminHideNewsFeed(ctx context.Context, newsFeedID uint64, isHidden bool) (*domain.NewsFeed, error) {
	if !u.permissionUsecase.IsAdmin(ctx) {
		return nil, &_routes.Except{
			Code:    http.StatusForbidden,
			Message: "Bạn không có quyền admin",
		}
	}

	err := u.newsFeedRepo.HideNewsFeed(ctx, newsFeedID, isHidden)
	if err != nil {
		return nil, err
	}

	newsFeed, err := u.newsFeedRepo.GetByID(ctx, newsFeedID)
	if err != nil {
		return nil, err
	}
	return newsFeed, nil
}

// AdminGetListNewsFeed - Admin lấy danh sách bài viết
func (u *NewsFeedUsecase) AdminGetListNewsFeed(ctx context.Context, req *dto.AdminNewsFeedSearch) ([]*dto.NewsFeedPublic, int32, error) {
	if !u.permissionUsecase.IsAdmin(ctx) {
		return nil, 0, &_routes.Except{
			Code:    http.StatusForbidden,
			Message: "Bạn không có quyền admin",
		}
	}

	// Validate pagination
	if req.Size > 100 || req.Size <= 0 {
		req.Size = 100
	}
	if req.Page <= 0 {
		req.Page = 0
	}

	newsFeeds, total, err := u.newsFeedRepo.AdminGetList(ctx, req)
	if err != nil {
		return nil, 0, err
	}
	return newsFeeds, total, nil
}

// AdminGetNewsFeedByID - Admin lấy chi tiết bài viết
func (u *NewsFeedUsecase) AdminGetNewsFeedByID(ctx context.Context, newsFeedID uint64) (*domain.NewsFeed, error) {
	if !u.permissionUsecase.IsAdmin(ctx) {
		return nil, &_routes.Except{
			Code:    http.StatusForbidden,
			Message: "Bạn không có quyền admin",
		}
	}

	profileId := _utils.GetProfileIdWithContext(ctx)
	newsFeed, err := u.newsFeedRepo.AdminGetByID(ctx, newsFeedID)
	if err != nil {
		return nil, err
	}

	// Lấy thông tin like của user hiện tại
	like, err := u.likeRepo.GetByTargetIdAndTypeAndUserId(ctx, newsFeedID, enums.TargetTypeNewsFeed, profileId)
	if like != nil && err == nil {
		newsFeed.LikedAt = like.CreatedAt
		newsFeed.DisLike = &like.DisLike
	}

	return newsFeed, nil
}
