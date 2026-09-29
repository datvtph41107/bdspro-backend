package repo

import (
	"context"
	qh_domain "tqd/internal/domain/qh"
)

type QHLabelRepository interface {
	Create(ctx context.Context, label *qh_domain.QHLabel) error
	Update(ctx context.Context, label *qh_domain.QHLabel) error
	Delete(ctx context.Context, id uint64) error
	// DeleteBatch soft-delete nhiều label còn active trong một câu lệnh (id IN (...)).
	DeleteBatch(ctx context.Context, ids []uint64) error
	// DeleteLabelLayerLinksByLayerID chỉ xóa bản ghi trong qh_label_layers (không xóa qh_labels).
	DeleteLabelLayerLinksByLayerID(ctx context.Context, layerID uint64) error
	// HardDeleteAllByLayerID xóa thật mọi label thuộc layer (layer_id) mà không được dùng bởi layer khác.
	HardDeleteAllByLayerID(ctx context.Context, layerID uint64) error
	// CountActiveByIDs đếm label chưa xóa mềm có id thuộc ids (phục vụ validate batch).
	CountActiveByIDs(ctx context.Context, ids []uint64) (int64, error)
	GetByID(ctx context.Context, id uint64) (*qh_domain.QHLabel, error)
	// GetByName tìm nhãn theo name (chưa xóa mềm). Nhiều bản ghi cùng tên: lấy id nhỏ nhất.
	GetByName(ctx context.Context, name string) (*qh_domain.QHLabel, error)
	GetByLayerID(ctx context.Context, layerID uint64) ([]qh_domain.QHLabel, error)
	// ListAdmin — danh sách label cho admin; layerID nil = mọi layer. Phân trang tại DB (offset/limit).
	ListAdmin(ctx context.Context, layerID *uint64, offset, limit int, includeInactive bool) ([]qh_domain.QHLabel, int64, error)
	GetByNameAndLayer(ctx context.Context, name string, layerID uint64) (*qh_domain.QHLabel, error)
	// EnsureLayerLink thêm (label_id, layer_id) vào qh_label_layers nếu chưa có.
	EnsureLayerLink(ctx context.Context, labelID, layerID uint64) error
	UpdateLayerLinkMetadata(ctx context.Context, labelID, layerID uint64, landUseID, legendID *uint64) error
	// EnsureLayerLinksBatch đảm bảo mọi label_id đều có bản ghi nối tới layer_id (idempotent, một transaction).
	EnsureLayerLinksBatch(ctx context.Context, layerID uint64, labelIDs []uint64) error
	GetNamesByLayerID(ctx context.Context, layerID uint64) (map[string]struct{}, error)
	BulkCreate(ctx context.Context, labels []*qh_domain.QHLabel) error
	// SyncRegionCountByLayerID đồng bộ qh_labels.region_count theo qh_regions (logic giống infra/postgres/qh_label_region_count.sql).
	SyncRegionCountByLayerID(ctx context.Context, layerID uint64) error
	// RefreshRegionCountForLabel gán region_count = COUNT(qh_regions) cho đúng một label (vd sau merge).
	RefreshRegionCountForLabel(ctx context.Context, labelID uint64) error
}
