package repository

import (
	"context"

	"userv4/modules/notification/entity"

	"gorm.io/gorm"
)

// DealHistoryRepository persists deal timeline rows (deal_history).
type DealHistoryRepository interface {
	Create(ctx context.Context, row *entity.DealHistory) error
	ListByDealID(ctx context.Context, dealID uint64, page, size uint32) ([]*entity.DealHistory, int64, error)
}

type dealHistoryPostgres struct {
	db *gorm.DB
}

func NewDealHistoryPostgres(db *gorm.DB) DealHistoryRepository {
	return &dealHistoryPostgres{db: db}
}

func (r *dealHistoryPostgres) Create(ctx context.Context, row *entity.DealHistory) error {
	if r.db == nil || row == nil {
		return nil
	}
	return r.db.WithContext(ctx).Create(row).Error
}

func (r *dealHistoryPostgres) ListByDealID(ctx context.Context, dealID uint64, page, size uint32) ([]*entity.DealHistory, int64, error) {
	if r.db == nil || dealID == 0 {
		return nil, 0, nil
	}
	if size == 0 {
		size = 20
	}
	if size > 100 {
		size = 100
	}
	offset := int(page * size)
	limit := int(size)

	base := r.db.WithContext(ctx).Model(&entity.DealHistory{}).Where("deal_id = ?", dealID)
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []*entity.DealHistory
	err := base.Order("created_at DESC").Offset(offset).Limit(limit).Find(&rows).Error
	return rows, total, err
}
