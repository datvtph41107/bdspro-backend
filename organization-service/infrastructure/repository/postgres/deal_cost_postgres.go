package postgres

import (
	"context"
	"time"

	"organization/internal/domain/entity"
	"organization/internal/domain/repository"
	"organization/internal/dto"

	"gorm.io/gorm"
)

type DealCostPostgreRepo struct {
	db *gorm.DB
}

func NewDealCostRepository(db *gorm.DB) repository.DealCostRepository {
	return &DealCostPostgreRepo{db: db}
}

func (r *DealCostPostgreRepo) Create(ctx context.Context, dealCost *entity.DealCost) error {
	return r.db.WithContext(ctx).Create(dealCost).Error
}

func (r *DealCostPostgreRepo) Update(ctx context.Context, dealCost *entity.DealCost) error {
	return r.db.WithContext(ctx).Save(dealCost).Error
}

func (r *DealCostPostgreRepo) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).
		Model(&entity.DealCost{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", time.Now()).Error
}

func (r *DealCostPostgreRepo) GetByID(ctx context.Context, id uint64) (*entity.DealCost, error) {
	var cost entity.DealCost
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&cost).Error
	return &cost, err
}

func (r *DealCostPostgreRepo) ListByDealID(ctx context.Context, search *dto.DealCostSearchDTO) ([]*entity.DealCost, error) {
	var costs []*entity.DealCost
	err := r.db.WithContext(ctx).
		Preload("CostType").
		Where("deal_id = ? AND deleted_at IS NULL", search.DealId).
		Offset(search.GetOffset()).
		Limit(int(search.GetSize())).
		Find(&costs).Error
	return costs, err
}
