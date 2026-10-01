package repository

import (
	"context"
	qh "userv4/modules/tqd/entity/qh"
)

type ILayerLegalRepo interface {
	Create(ctx context.Context, legal *qh.QHLayerLegal) error
	CreateBatch(ctx context.Context, legals []*qh.QHLayerLegal) error
	GetByID(ctx context.Context, id uint64) (*qh.QHLayerLegal, error)
	GetByLayerID(ctx context.Context, layerID uint64) ([]*qh.QHLayerLegal, error)
	GetByLayerIDs(ctx context.Context, layerIDs []uint64) (map[uint64][]*qh.QHLayerLegal, error)
	Delete(ctx context.Context, id uint64) error
	DeleteByLayerID(ctx context.Context, layerID uint64) error
	HardDeleteByLayerID(ctx context.Context, layerID uint64) error
}
