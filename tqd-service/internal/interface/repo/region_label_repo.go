package repo

import (
	"context"
	qh_domain "tqd/internal/domain/qh"
)

// RegionLabelRepository - Quản lý quan hệ nhiều-nhiều giữa region và label
type RegionLabelRepository interface {
	// AssignLabels - Gán labels cho region (thay thế toàn bộ labels cũ)
	AssignLabels(ctx context.Context, regionID uint64, labelIDs []uint64, assignedBy uint64, sourceFile, notes string) error

	// AddLabels - Thêm labels cho region (giữ nguyên labels cũ)
	AddLabels(ctx context.Context, regionID uint64, labelIDs []uint64, assignedBy uint64, sourceFile, notes string) error

	// RemoveLabel - Gỡ một label khỏi region
	RemoveLabel(ctx context.Context, regionID, labelID uint64) error

	// BatchAssignLabels - Gán labels hàng loạt cho nhiều regions
	BatchAssignLabels(ctx context.Context, regionIDs, labelIDs []uint64, assignedBy uint64, sourceFile, notes string) (int, error)

	// GetLabelsByRegion - Lấy danh sách labels của một region
	GetLabelsByRegion(ctx context.Context, regionID uint64) ([]qh_domain.QHLabel, error)

	// GetRegionIDsByLabel - Lấy danh sách region IDs theo label
	GetRegionIDsByLabel(ctx context.Context, labelID uint64) ([]uint64, error)

	// DeleteByRegion - Xóa tất cả labels của một region
	DeleteByRegion(ctx context.Context, regionID uint64) error

	// DeleteByLabel - Xóa tất cả quan hệ của một label
	DeleteByLabel(ctx context.Context, labelID uint64) error
}
