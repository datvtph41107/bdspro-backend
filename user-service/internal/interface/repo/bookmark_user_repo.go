package repo

import (
	"context"
	"user/internal/dto"
	models "user/internal/models"
)

// IBookmarkUserRepo interface cho bookmark user repository
type IBookmarkUserRepo interface {
	// GetBookmarkUsers lấy danh sách bookmark user với phân trang
	GetBookmarkUsers(ctx context.Context, adminID uint64, page, size int) ([]models.BookmarkUserEntity, int64, error)

	// CreateBookmarkUser tạo bookmark user mới
	CreateBookmarkUser(ctx context.Context, bookmark *models.BookmarkUserEntity) error

	// DeleteBookmarkUser xóa bookmark user
	DeleteBookmarkUser(ctx context.Context, adminID, userID uint64) error

	// DeleteAllBookmarkUsers xóa tất cả bookmark user của admin
	DeleteAllBookmarkUsers(ctx context.Context, adminID uint64) error

	// CheckBookmarkUserExists kiểm tra bookmark user đã tồn tại chưa
	CheckBookmarkUserExists(ctx context.Context, adminID, userID uint64) (bool, error)

	// GetBookmarkedUsers lấy danh sách user có bookmark với phân trang và tìm kiếm
	GetBookmarkedUsers(ctx context.Context, adminID uint64, req *dto.BookmarkUsersRequest) ([]dto.AdminUserItem, int64, error)
}
