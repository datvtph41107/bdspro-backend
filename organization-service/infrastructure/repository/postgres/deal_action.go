package postgres

import (
	"context"
	"organization/internal/domain/entity"
	"organization/internal/repo"

	"gorm.io/gorm"
)

type DealActionRepository struct {
	db *gorm.DB
}

func NewDealActionRepository(db *gorm.DB) repo.DealActionRepository {
	return &DealActionRepository{db: db}
}

// Create creates a new deal action
func (r *DealActionRepository) Create(ctx context.Context, dealAction *entity.DealAction) error {
	return r.db.WithContext(ctx).Create(dealAction).Error
}

// GetByID gets a deal action by ID
func (r *DealActionRepository) GetByID(ctx context.Context, id uint64) (*entity.DealAction, error) {
	var dealAction entity.DealAction
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&dealAction).Error
	if err != nil {
		return nil, err
	}
	return &dealAction, nil
}

// GetByDealID gets deal actions by deal ID with pagination
func (r *DealActionRepository) GetByDealID(ctx context.Context, dealID uint64, page, size int32) ([]*entity.DealAction, int32, error) {
	var dealActions []*entity.DealAction
	var total int64

	// Get total count
	err := r.db.WithContext(ctx).Model(&entity.DealAction{}).Where("deal_id = ?", dealID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * size
	err = r.db.WithContext(ctx).Where("deal_id = ?", dealID).
		Order("created_at DESC").
		Offset(int(offset)).
		Limit(int(size)).
		Find(&dealActions).Error
	if err != nil {
		return nil, 0, err
	}

	return dealActions, int32(total), nil
}

// GetByTransactionID gets a deal action by transaction ID
func (r *DealActionRepository) GetByTransactionID(ctx context.Context, transactionID uint64) (*entity.DealAction, error) {
	var dealAction entity.DealAction
	err := r.db.WithContext(ctx).Where("transaction_id = ?", transactionID).First(&dealAction).Error
	if err != nil {
		return nil, err
	}
	return &dealAction, nil
}

// Update updates a deal action
func (r *DealActionRepository) Update(ctx context.Context, dealAction *entity.DealAction) error {
	return r.db.WithContext(ctx).Save(dealAction).Error
}

// Delete deletes a deal action
func (r *DealActionRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&entity.DealAction{}, id).Error
}
