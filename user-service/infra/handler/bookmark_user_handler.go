package handler

import (
	_dto "common/domain/dto"
	_utils "common/utils"
	"context"
	userpb "pb/types/user"
	"user/internal/dto"
	"user/internal/usecases"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// BookmarkUserHandler handler cho bookmark user
type BookmarkUserHandler struct {
	userpb.UnimplementedBookmarkUserServiceServer
	BookmarkUserService *usecases.BookmarkUserUsecase
}

// NewBookmarkUserHandler tạo instance mới của BookmarkUserHandler
func NewBookmarkUserHandler(bookmarkUserService *usecases.BookmarkUserUsecase) *BookmarkUserHandler {
	return &BookmarkUserHandler{
		BookmarkUserService: bookmarkUserService,
	}
}

// UpdateBookmarkUsers implement protobuf API
func (h *BookmarkUserHandler) UpdateBookmarkUsers(ctx context.Context, req *userpb.UpdateBookmarkUsersRequest) (*emptypb.Empty, error) {
	// Lấy adminID từ context
	adminID := _utils.GetProfileIdWithContext(ctx)

	// Tạo request DTO
	dtoReq := &dto.BookmarkUserUpdateRequest{
		ProfileIDs: req.ProfileIds,
		Bookmark:   req.Bookmark,
	}

	// Gọi usecase
	err := h.BookmarkUserService.UpdateBookmarkUsers(ctx, adminID, dtoReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update bookmark users: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// BookmarkUsers implement protobuf API - lấy danh sách user có bookmark
func (h *BookmarkUserHandler) GetBookmarkUsers(ctx context.Context, req *userpb.GetBookmarkUsersRequest) (*userpb.ListAllUsersResponse, error) {
	// Tạo request DTO
	dtoReq := &dto.BookmarkUsersRequest{
		Pagable: _dto.Pagable{
			Page: req.Page,
			Size: req.Size,
		},
		// Search:    req.Search,
		// RoleType:  req.RoleType,
		// StartDate: req.StartDate,
		// EndDate:   req.EndDate,
	}

	// Gọi usecase
	result, total, err := h.BookmarkUserService.GetBookmarkedUsers(ctx, dtoReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get bookmarked users: %v", err)
	}

	// Convert DTO response to proto response
	response := &userpb.ListAllUsersResponse{
		Total: int32(total),
		Page:  dtoReq.Page,
		Size:  dtoReq.Size,
	}

	// Convert data items
	for _, item := range result {
		protoItem := &userpb.AdminUserItem{
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
		}
		response.Data = append(response.Data, protoItem)
	}

	return response, nil
}
