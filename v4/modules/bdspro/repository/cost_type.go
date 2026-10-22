package repository

import (
	"context"
	"errors"
	"time"

	"mdn4/modules/bdspro/dto"
	"mdn4/modules/bdspro/entity"

	"gorm.io/gorm"
)

// CostTypeRepository is the tx_cost_types store on the v4 DSN.
type CostTypeRepository interface {
	Create(ctx context.Context, costType *entity.TxCostType) error
	Update(ctx context.Context, costType *entity.TxCostType) error
	Delete(ctx context.Context, id uint64) error
	List(ctx context.Context, search dto.CostTypeSearchRequest) ([]*entity.TxCostType, error)
	GetByID(ctx context.Context, id uint64) (*entity.TxCostType, error)
}

type costTypePostgres struct {
	db *gorm.DB
}

func NewCostTypePostgres(db *gorm.DB) CostTypeRepository {
	return &costTypePostgres{db: db}
}

func (r *costTypePostgres) Create(ctx context.Context, costType *entity.TxCostType) error {
	if r.db == nil {
		return nil
	}
	return r.db.WithContext(ctx).Create(costType).Error
}

func (r *costTypePostgres) Update(ctx context.Context, costType *entity.TxCostType) error {
	if r.db == nil {
		return nil
	}
	return r.db.WithContext(ctx).Save(costType).Error
}

func (r *costTypePostgres) Delete(ctx context.Context, id uint64) error {
	if r.db == nil || id == 0 {
		return nil
	}
	return r.db.WithContext(ctx).
		Model(&entity.TxCostType{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", time.Now()).Error
}

func (r *costTypePostgres) List(ctx context.Context, search dto.CostTypeSearchRequest) ([]*entity.TxCostType, error) {
	if r.db == nil {
		return nil, nil
	}
	query := r.db.WithContext(ctx).Where("deleted_at IS NULL")
	if search.TypeName != "" {
		query = query.Where("type_name LIKE ?", "%"+search.TypeName+"%")
	}
	if search.CostType != 0 {
		query = query.Where("cost_type = ?", search.CostType)
	}
	var rows []*entity.TxCostType
	err := query.
		Order("id ASC").
		Offset(search.GetOffset()).
		Limit(search.GetLimit()).
		Find(&rows).Error
	return rows, err
}

func (r *costTypePostgres) GetByID(ctx context.Context, id uint64) (*entity.TxCostType, error) {
	if r.db == nil || id == 0 {
		return nil, nil
	}
	var row entity.TxCostType
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &row, nil
}
