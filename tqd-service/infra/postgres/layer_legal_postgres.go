package postgres

import (
	"context"
	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/interface/repo"

	"gorm.io/gorm"
)

type LayerLegalPostgres struct {
	db *gorm.DB
}

func NewLayerLegalPostgres(db *gorm.DB) repo.ILayerLegalRepo {
	return &LayerLegalPostgres{db: db}
}

func (r *LayerLegalPostgres) Create(ctx context.Context, legal *qh_domain.QHLayerLegal) error {
	return r.db.WithContext(ctx).Create(legal).Error
}

func (r *LayerLegalPostgres) CreateBatch(ctx context.Context, legals []*qh_domain.QHLayerLegal) error {
	if len(legals) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).CreateInBatches(legals, 100).Error
}

func (r *LayerLegalPostgres) GetByID(ctx context.Context, id uint64) (*qh_domain.QHLayerLegal, error) {
	var legal qh_domain.QHLayerLegal
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&legal).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &legal, nil
}

func (r *LayerLegalPostgres) GetByLayerID(ctx context.Context, layerID uint64) ([]*qh_domain.QHLayerLegal, error) {
	var legals []*qh_domain.QHLayerLegal
	err := r.db.WithContext(ctx).
		Where("layer_id = ? AND deleted_at IS NULL", layerID).
		Order("created_at DESC").
		Find(&legals).Error
	return legals, err
}

func (r *LayerLegalPostgres) GetByLayerIDs(ctx context.Context, layerIDs []uint64) (map[uint64][]*qh_domain.QHLayerLegal, error) {
	if len(layerIDs) == 0 {
		return make(map[uint64][]*qh_domain.QHLayerLegal), nil
	}

	var legals []*qh_domain.QHLayerLegal
	err := r.db.WithContext(ctx).
		Where("layer_id IN ? AND deleted_at IS NULL", layerIDs).
		Order("created_at DESC").
		Find(&legals).Error
	if err != nil {
		return nil, err
	}

	result := make(map[uint64][]*qh_domain.QHLayerLegal)
	for _, legal := range legals {
		result[legal.LayerID] = append(result[legal.LayerID], legal)
	}
	return result, nil
}

func (r *LayerLegalPostgres) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).
		Model(&qh_domain.QHLayerLegal{}).
		Where("id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (r *LayerLegalPostgres) DeleteByLayerID(ctx context.Context, layerID uint64) error {
	return r.db.WithContext(ctx).
		Model(&qh_domain.QHLayerLegal{}).
		Where("layer_id = ?", layerID).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (r *LayerLegalPostgres) HardDeleteByLayerID(ctx context.Context, layerID uint64) error {
	if layerID == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Unscoped().Where("layer_id = ?", layerID).Delete(&qh_domain.QHLayerLegal{}).Error
}
