package postgres

import (
	"context"
	"errors"
	"time"

	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/repo"

	"gorm.io/gorm"
)

type CostTypePostgreRepo struct {
	db *gorm.DB
}

func NewCostTypeRepository(db *gorm.DB) repo.TxCostTypeRepository {
	return &CostTypePostgreRepo{db: db}
}

func (r *CostTypePostgreRepo) Create(ctx context.Context, costType *domain.TxCostType) error {
	return r.db.WithContext(ctx).Create(costType).Error
}

func (r *CostTypePostgreRepo) Update(ctx context.Context, costType *domain.TxCostType) error {
	return r.db.WithContext(ctx).Save(costType).Error
}

func (r *CostTypePostgreRepo) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).
		Model(&domain.TxCostType{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", time.Now()).Error
}

func (r *CostTypePostgreRepo) List(ctx context.Context, search *dto.TxCostTypeSearchDTO) ([]*domain.TxCostType, error) {
	var costTypes []*domain.TxCostType
	query := r.db.WithContext(ctx).Where("deleted_at IS NULL")

	if search.TypeName != "" {
		query = query.Where("type_name LIKE ?", "%"+search.TypeName+"%")
	}
	if search.CostType != 0 {
		query = query.Where("cost_type = ?", search.CostType)
	}

	err := query.
		Offset(search.GetOffset()).
		Limit(int(search.GetSize())).
		Find(&costTypes).Error
	return costTypes, err
}

func (r *CostTypePostgreRepo) GetByID(ctx context.Context, id uint64) (*domain.TxCostType, error) {
	var costType domain.TxCostType
	err := GetDB(ctx, r.db).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&costType).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &costType, err
}
