package repo

import (
	"context"
	"hub/internal/domain"
)

// IUpdateDataRepo interface cho update sync tracking
type IUpdateDataRepo interface {
	// Upsert ghi nhận update (tạo mới hoặc cập nhật nếu đã tồn tại owner+resource+resourceId)
	Upsert(c context.Context, entity *domain.UpdateData) error
	// DeleteByOwnerResourceId xóa bản ghi (khi DelUpdate)
	DeleteByOwnerResourceId(c context.Context, ownerID uint64, resourceType domain.ESyncResource, resourceID uint64) error
	// GetChangedIdsSince lấy danh sách resource IDs đã thay đổi từ lastSync
	GetChangedIdsSince(c context.Context, ownerID uint64, resourceType domain.ESyncResource, lastSync int64, limit int) ([]uint64, error)
	// FlushOld xóa các bản ghi cũ, giữ lại limit mới nhất
	FlushOld(c context.Context, ownerID uint64, resourceType domain.ESyncResource, keepLimit int64) error
}
