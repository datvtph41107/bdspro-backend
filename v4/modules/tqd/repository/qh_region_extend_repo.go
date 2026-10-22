package repository

import (
	"context"

	_dto "common/domain/dto"
	qh "mdn4/modules/tqd/entity/qh"
)

type RegionExtendRepository interface {
	Create(ctx context.Context, record *qh.QHRegionExtend) error
	// CreateMergedReplaceSources gộp geometry các sourceIDs (+ geometry request nếu có) thành MultiLineString, tạo record mới, xóa mềm source cũ.
	CreateMergedReplaceSources(ctx context.Context, record *qh.QHRegionExtend, sourceIDs []uint64, extraGeometry []byte) error
	CreateBatch(ctx context.Context, records []*qh.QHRegionExtend) error
	Update(ctx context.Context, record *qh.QHRegionExtend) error
	Delete(ctx context.Context, id uint64) error
	List(ctx context.Context, filter *RegionExtendFilter) ([]qh.QHRegionExtend, int64, error)
	GetByID(ctx context.Context, id uint64) (*qh.QHRegionExtend, error)
	HardDeleteAllByLayerID(ctx context.Context, layerID uint64) error
}

type RegionExtendFilter struct {
	LayerID  *uint64
	LabelID  *uint64
	GeomType *string
	Pagable  *_dto.Pagable
}
