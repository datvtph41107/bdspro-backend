package repo

import (
	"context"
	"mdn4/modules/user/domain/user"
)

// IProfileDeletedRepo interface cho repository profile đã xóa
type IProfileDeletedRepo interface {
	// Create lưu thông tin profile đã xóa vào bảng profile_deleted
	Create(ctx context.Context, profileDeleted *user.ProfileDeleted) error

	// GetByProfileID lấy thông tin profile đã xóa theo profileID
	GetByProfileID(ctx context.Context, profileID uint64) (*user.ProfileDeleted, error)

	// GetList lấy danh sách profile đã xóa với phân trang
	GetList(ctx context.Context, offset, limit int) ([]*user.ProfileDeleted, int64, error)
}
