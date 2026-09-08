package postgres

import (
	"context"
	"organization/internal/domain/entity"
	"time"

	"gorm.io/gorm"
)

// @bind: organization/internal/domain/repository.DealMilestoneRepository
type DealMilestonePostgresRepository struct {
	db *gorm.DB
}

func NewDealMilestonePostgresRepository(db *gorm.DB) *DealMilestonePostgresRepository {
	return &DealMilestonePostgresRepository{db: db}
}

func (r *DealMilestonePostgresRepository) Create(ctx context.Context, milestone *entity.DealMilestone) (*entity.DealMilestone, error) {
	if err := r.db.WithContext(ctx).Create(milestone).Error; err != nil {
		return nil, err
	}
	return milestone, nil
}

func (r *DealMilestonePostgresRepository) Update(ctx context.Context, milestone *entity.DealMilestone) (*entity.DealMilestone, error) {
	if err := r.db.WithContext(ctx).Model(&entity.DealMilestone{}).Where("id = ?", milestone.ID).Updates(milestone).Error; err != nil {
		return nil, err
	}
	return milestone, nil
}

func (r *DealMilestonePostgresRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Model(&entity.DealMilestone{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_at": time.Now(),
		}).Error
}

func (r *DealMilestonePostgresRepository) GetByID(ctx context.Context, id uint64) (*entity.DealMilestone, error) {
	var milestone entity.DealMilestone
	if err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&milestone).Error; err != nil {
		return nil, err
	}
	return &milestone, nil
}

func (r *DealMilestonePostgresRepository) GetByDealID(ctx context.Context, dealID uint64) ([]*entity.DealMilestone, error) {
	var milestones []*entity.DealMilestone
	if err := r.db.WithContext(ctx).
		Where("deal_id = ? AND deleted_at IS NULL", dealID).
		Order("order_index asc").
		Find(&milestones).Error; err != nil {
		return nil, err
	}
	return milestones, nil
}

func (r *DealMilestonePostgresRepository) UpdateOrder(ctx context.Context, milestoneIDs []uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i, id := range milestoneIDs {
			if err := tx.Model(&entity.DealMilestone{}).
				Where("id = ?", id).
				Update("order_index", i).Error; err != nil {
				return err
			}
		}
		return nil
	})
} 