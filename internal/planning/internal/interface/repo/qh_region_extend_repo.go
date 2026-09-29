package repo

import (
	"context"

	_dto "common/domain/dto"
	qh_domain "tqd/internal/domain/qh"
)

type RegionExtendRepository interface {
	Create(ctx context.Context, record *qh_domain.QHRegionExtend) error
	// CreateMergedReplaceSources gộp geometry các sourceIDs (+ geometry request nếu có) thành MultiLineString, tạo record mới, xóa mềm source cũ.
	CreateMergedReplaceSources(ctx context.Context, record *qh_domain.QHRegionExtend, sourceIDs []uint64, extraGeometry []byte) error
	CreateBatch(ctx context.Context, records []*qh_domain.QHRegionExtend) error
	Update(ctx context.Context, record *qh_domain.QHRegionExtend) error
	Delete(ctx context.Context, id uint64) error
	List(ctx context.Context, filter *RegionExtendFilter) ([]qh_domain.QHRegionExtend, int64, error)
	GetByID(ctx context.Context, id uint64) (*qh_domain.QHRegionExtend, error)
	HardDeleteAllByLayerID(ctx context.Context, layerID uint64) error
}

type RegionExtendFilter struct {
	LayerID  *uint64
	LabelID  *uint64
	GeomType *string
	Pagable  *_dto.Pagable
}
