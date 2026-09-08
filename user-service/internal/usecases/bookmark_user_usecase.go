package usecases

import (
	_dto "common/domain/dto"
	_utils "common/utils"
	"context"
	"errors"
	"user/internal/dto"
	"user/internal/interface/repo"
	"user/internal/models"
)

// BookmarkUserUsecase usecase cho bookmark user
type BookmarkUserUsecase struct {
	BookmarkUserRepo repo.IBookmarkUserRepo
}

// NewBookmarkUserUsecase tạo instance mới của BookmarkUserUsecase
func NewBookmarkUserUsecase(bookmarkUserRepo repo.IBookmarkUserRepo) *BookmarkUserUsecase {
	return &BookmarkUserUsecase{
		BookmarkUserRepo: bookmarkUserRepo,
	}
}

// GetBookmarkUsers lấy danh sách bookmark user với phân trang
func (u *BookmarkUserUsecase) GetBookmarkUsers(ctx context.Context, req *dto.BookmarkUserListRequest) (*dto.BookmarkUserListResponse, error) {
	bookmarks, total, err := u.BookmarkUserRepo.GetBookmarkUsers(ctx, req.AdminID, int(req.Page), int(req.Size))
	if err != nil {
		return nil, err
	}

	// Convert to response
	responses := make([]dto.BookmarkUserResponse, len(bookmarks))
	for i, bookmark := range bookmarks {
		responses[i] = dto.BookmarkUserResponse{
			AdminID: bookmark.AdminID,
			UserID:  bookmark.UserID,
		}

		// Map user info if available
		if bookmark.User != nil {
			responses[i].User = &models.UserItem{
				ProfileID:    bookmark.User.ProfileID,
				FullName:     bookmark.User.FullName,
				Avatar:       bookmark.User.Avatar,
				TickVerified: bookmark.User.TickVerified,
			}
		}
	}

	return &dto.BookmarkUserListResponse{
		Pagable: _dto.Pagable{
			Page: req.Page,
			Size: req.Size,
		},
		Data:  responses,
		Total: total,
	}, nil
}

// UpdateBookmarkUsers cập nhật bookmark users (add/remove bookmark)
func (u *BookmarkUserUsecase) UpdateBookmarkUsers(ctx context.Context, adminID uint64, req *dto.BookmarkUserUpdateRequest) error {
	if req.Bookmark {
		// Add bookmark cho các user - chỉ tạo mới nếu chưa có
		for _, profileID := range req.ProfileIDs {
			// Kiểm tra xem đã bookmark chưa
			exists, err := u.BookmarkUserRepo.CheckBookmarkUserExists(ctx, adminID, profileID)
			if err != nil {
				return err
			}

			// Nếu chưa bookmark thì tạo mới
			if !exists {
				bookmark := &models.BookmarkUserEntity{
					AdminID: adminID,
					UserID:  profileID,
				}

				if err := u.BookmarkUserRepo.CreateBookmarkUser(ctx, bookmark); err != nil {
					return err
				}
			}
		}
	} else {
		// Remove bookmark cho các user - xóa trực tiếp luôn
		for _, profileID := range req.ProfileIDs {
			if err := u.BookmarkUserRepo.DeleteBookmarkUser(ctx, adminID, profileID); err != nil {
				return err
			}
		}
	}

	return nil
}

// GetBookmarkedUsers lấy danh sách user có bookmark với phân trang và tìm kiếm
func (u *BookmarkUserUsecase) GetBookmarkedUsers(ctx context.Context, req *dto.BookmarkUsersRequest) ([]dto.AdminUserItem, int64, error) {
	// Lấy adminID từ context
	adminID := _utils.GetProfileIdWithContext(ctx)
	if adminID == 0 {
		return nil, 0, errors.New("unauthorized")
	}

	// Lấy danh sách user có bookmark
	users, total, err := u.BookmarkUserRepo.GetBookmarkedUsers(ctx, adminID, req)
	if err != nil {
		return nil, 0, err
	}

	// Convert sang AdminUserItem format
	var adminUserItems []dto.AdminUserItem
	for _, user := range users {
		adminUserItem := dto.AdminUserItem{
			ProfileID:      user.ProfileID,
			FullName:       user.FullName,
			Email:          user.Email,
			Phone:          user.Phone,
			Avatar:         user.Avatar,
			Address:        user.Address,
			RoleType:       user.RoleType,
			RoleRealEstate: user.RoleRealEstate,
			Position:       user.Position,
			Workplace:      user.Workplace,
			DepartmentID:   user.DepartmentID,
			Bookmark:       true, // Tất cả user trong danh sách này đều có bookmark
			TickVerified:   user.TickVerified,
		}
		adminUserItems = append(adminUserItems, adminUserItem)
	}

	return adminUserItems, total, nil
}
