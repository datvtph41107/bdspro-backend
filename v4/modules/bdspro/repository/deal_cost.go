package repository

import (
	"context"
	"errors"
	"strings"

	"userv4/modules/bdspro/dto"
	"userv4/modules/bdspro/entity"

	"gorm.io/gorm"
)

// DealCostRepository is tx_deal_costs CRUD on the v4 DSN.
type DealCostRepository interface {
	Create(ctx context.Context, row *entity.TxDealCost) error
	Update(ctx context.Context, row *entity.TxDealCost) error
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*entity.TxDealCost, error)
	ListByDealID(ctx context.Context, search dto.DealCostSearchRequest) ([]*entity.TxDealCost, error)
}

type dealCostPostgres struct {
	db *gorm.DB
}

func NewDealCostPostgres(db *gorm.DB) DealCostRepository {
	return &dealCostPostgres{db: db}
}

func (r *dealCostPostgres) Create(ctx context.Context, row *entity.TxDealCost) error {
	if r.db == nil || row == nil {
		return nil
	}
	row.ID = 0
	err := r.db.WithContext(ctx).Omit("ID").Create(row).Error
	if err == nil {
		return nil
	}
	// Dev/import: explicit ids can leave tx_deal_costs_id_seq behind MAX(id).
	if strings.Contains(err.Error(), "tx_deal_costs_pkey") {
		_ = r.db.WithContext(ctx).Exec(`
			SELECT setval(
				pg_get_serial_sequence('tx_deal_costs', 'id'),
				COALESCE((SELECT MAX(id) FROM tx_deal_costs), 1)
			)`).Error
		row.ID = 0
		return r.db.WithContext(ctx).Omit("ID").Create(row).Error
	}
	return err
}

func (r *dealCostPostgres) Update(ctx context.Context, row *entity.TxDealCost) error {
	if r.db == nil {
		return nil
	}
	return r.db.WithContext(ctx).Save(row).Error
}

func (r *dealCostPostgres) Delete(ctx context.Context, id uint64) error {
	if r.db == nil || id == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.TxDealCost{}).Error
}

func (r *dealCostPostgres) GetByID(ctx context.Context, id uint64) (*entity.TxDealCost, error) {
	if r.db == nil || id == 0 {
		return nil, nil
	}
	var row entity.TxDealCost
	err := r.db.WithContext(ctx).
		Preload("CostType").
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

func (r *dealCostPostgres) ListByDealID(ctx context.Context, search dto.DealCostSearchRequest) ([]*entity.TxDealCost, error) {
	if r.db == nil || search.DealID == 0 {
		return nil, nil
	}
	var rows []*entity.TxDealCost
	err := r.db.WithContext(ctx).
		Preload("CostType").
		Where("deal_id = ? AND deleted_at IS NULL", search.DealID).
		Order("id ASC").
		Offset(search.GetOffset()).
		Limit(search.GetLimit()).
		Find(&rows).Error
	return rows, err
}
