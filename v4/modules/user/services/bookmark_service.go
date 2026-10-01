package services

import (
	"context"

	_dto "common/domain/dto"
	_utils "common/utils"
	userpb "pb/types/user"
	"userv4/internal/dto"
	usecases "userv4/internal/usecase/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// BookmarkService owns bookmark-user HTTP surface migrated into modules/user.
type BookmarkService interface {
	GetBookmarkUsers(ctx context.Context, page, size uint32) (*userpb.ListAllUsersResponse, error)
	UpdateBookmarkUsers(ctx context.Context, req *userpb.UpdateBookmarkUsersRequest) error
}

// LocalBookmarkService adapts BookmarkUserUsecase into modules/user.
type LocalBookmarkService struct {
	usecase *usecases.BookmarkUserUsecase
}

func NewLocalBookmarkService(usecase *usecases.BookmarkUserUsecase) *LocalBookmarkService {
	return &LocalBookmarkService{usecase: usecase}
}

var _ BookmarkService = (*LocalBookmarkService)(nil)

func (s *LocalBookmarkService) GetBookmarkUsers(ctx context.Context, page, size uint32) (*userpb.ListAllUsersResponse, error) {
	dtoReq := &dto.BookmarkUsersRequest{
		Pagable: _dto.Pagable{Page: page, Size: size},
	}
	result, total, err := s.usecase.GetBookmarkedUsers(ctx, dtoReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get bookmarked users: %v", err)
	}

	response := &userpb.ListAllUsersResponse{
		Total: int32(total),
		Page:  dtoReq.Page,
		Size:  dtoReq.Size,
	}
	for _, item := range result {
		response.Data = append(response.Data, &userpb.AdminUserItem{
			ProfileId:      item.ProfileID,
			FullName:       item.FullName,
			Email:          item.Email,
			Phone:          item.Phone,
			Avatar:         item.Avatar,
			Address:        item.Address,
			RoleType:       item.RoleType,
			RoleRealEstate: item.RoleRealEstate,
			Position:       item.Position,
			DepartmentId:   item.DepartmentID,
			Bookmark:       item.Bookmark,
			TickVerified:   item.TickVerified,
			TotalLogin:     item.TotalLogin,
		})
	}
	return response, nil
}

func (s *LocalBookmarkService) UpdateBookmarkUsers(ctx context.Context, req *userpb.UpdateBookmarkUsersRequest) error {
	adminID := _utils.GetProfileIdWithContext(ctx)
	err := s.usecase.UpdateBookmarkUsers(ctx, adminID, &dto.BookmarkUserUpdateRequest{
		ProfileIDs: req.ProfileIds,
		Bookmark:   req.Bookmark,
	})
	if err != nil {
		return status.Errorf(codes.Internal, "Failed to update bookmark users: %v", err)
	}
	return nil
}
