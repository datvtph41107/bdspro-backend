package repo

import (
	"context"
	"user/internal/models"
)

// IProfileDeletedRepo interface cho repository profile đã xóa
type IProfileDeletedRepo interface {
	// Create lưu thông tin profile đã xóa vào bảng profile_deleted
	Create(ctx context.Context, profileDeleted *models.ProfileDeleted) error

	// GetByProfileID lấy thông tin profile đã xóa theo profileID
	GetByProfileID(ctx context.Context, profileID uint64) (*models.ProfileDeleted, error)

	// GetList lấy danh sách profile đã xóa với phân trang
	GetList(ctx context.Context, offset, limit int) ([]*models.ProfileDeleted, int64, error)
}
