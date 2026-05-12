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
	// CountActiveByIDs đếm label chưa xóa mềm có id thuộc ids (phục vụ validate batch).
	CountActiveByIDs(ctx context.Context, ids []uint64) (int64, error)
	GetByID(ctx context.Context, id uint64) (*qh_domain.QHLabel, error)
	GetByLayerID(ctx context.Context, layerID uint64) ([]qh_domain.QHLabel, error)
	GetByNameAndLayer(ctx context.Context, name string, layerID uint64) (*qh_domain.QHLabel, error)
	GetNamesByLayerID(ctx context.Context, layerID uint64) (map[string]struct{}, error)
	BulkCreate(ctx context.Context, labels []*qh_domain.QHLabel) error
	// SyncRegionCountByLayerID đồng bộ qh_labels.region_count theo qh_regions (logic giống infra/postgres/qh_label_region_count.sql).
	SyncRegionCountByLayerID(ctx context.Context, layerID uint64) error
	// RefreshRegionCountForLabel gán region_count = COUNT(qh_regions) cho đúng một label (vd sau merge).
	RefreshRegionCountForLabel(ctx context.Context, labelID uint64) error
}
