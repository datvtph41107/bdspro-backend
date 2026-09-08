package repo

import (
	"context"
	qh_domain "tqd/internal/domain/qh"
)

type ILayerLegalRepo interface {
	Create(ctx context.Context, legal *qh_domain.QHLayerLegal) error
	CreateBatch(ctx context.Context, legals []*qh_domain.QHLayerLegal) error
	GetByID(ctx context.Context, id uint64) (*qh_domain.QHLayerLegal, error)
	GetByLayerID(ctx context.Context, layerID uint64) ([]*qh_domain.QHLayerLegal, error)
	GetByLayerIDs(ctx context.Context, layerIDs []uint64) (map[uint64][]*qh_domain.QHLayerLegal, error)
	Delete(ctx context.Context, id uint64) error
	DeleteByLayerID(ctx context.Context, layerID uint64) error
	HardDeleteByLayerID(ctx context.Context, layerID uint64) error
}
