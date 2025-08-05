package postgres

import (
	"context"
	"organization/internal/domain/entity"
	"organization/internal/domain/repository"
	"organization/internal/dto"
	"time"

	"gorm.io/gorm"
)

type CostTypePostgreRepo struct {
	db *gorm.DB
}

func NewCostTypeRepository(db *gorm.DB) repository.CostTypeRepository {
	return &CostTypePostgreRepo{db: db}
}

func (r *CostTypePostgreRepo) Create(ctx context.Context, costType *entity.CostType) error {
	return r.db.WithContext(ctx).Create(costType).Error
}

func (r *CostTypePostgreRepo) Update(ctx context.Context, costType *entity.CostType) error {
	return r.db.WithContext(ctx).Save(costType).Error
}

func (r *CostTypePostgreRepo) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).
		Model(&entity.CostType{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", time.Now()).Error
}

func (r *CostTypePostgreRepo) List(ctx context.Context, search *dto.CostTypeSearchDTO) ([]*entity.CostType, error) {
	var costTypes []*entity.CostType
	query := r.db.WithContext(ctx).Where("deleted_at IS NULL")

	if search.TypeName != "" {
		query = query.Where("type_name LIKE ?", "%"+search.TypeName+"%")
	}
	if search.CostType != 0 {
		query = query.Where("cost_type = ?", search.CostType)
	}

	err := query.
		Offset(search.GetOffset()).
		Limit(search.GetSize()).
		Find(&costTypes).Error
	return costTypes, err
}

func (r *CostTypePostgreRepo) GetByID(ctx context.Context, id uint64) (*entity.CostType, error) {
	var costType entity.CostType
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&costType).Error
	return &costType, err
}
