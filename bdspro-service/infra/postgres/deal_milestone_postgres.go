package postgres

import (
	"bdspro/internal/domain"
	"context"
	"time"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo.DealMilestoneRepository
type DealMilestonePostgresRepository struct {
	db *gorm.DB
}

func NewDealMilestonePostgresRepository(db *gorm.DB) *DealMilestonePostgresRepository {
	return &DealMilestonePostgresRepository{db: db}
}

func (r *DealMilestonePostgresRepository) Create(ctx context.Context, milestone *domain.DealMilestone) (*domain.DealMilestone, error) {
	if err := r.db.WithContext(ctx).Create(milestone).Error; err != nil {
		return nil, err
	}
	return milestone, nil
}

func (r *DealMilestonePostgresRepository) Update(ctx context.Context, milestone *domain.DealMilestone) (*domain.DealMilestone, error) {
	if err := r.db.WithContext(ctx).Model(&domain.DealMilestone{}).Where("id = ?", milestone.ID).Updates(milestone).Error; err != nil {
		return nil, err
	}
	return milestone, nil
}

func (r *DealMilestonePostgresRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Model(&domain.DealMilestone{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			// "is_deleted": true,
			"deleted_at": time.Now(),
		}).Error
}

func (r *DealMilestonePostgresRepository) GetByID(ctx context.Context, id uint64) (*domain.DealMilestone, error) {
	var milestone domain.DealMilestone
	if err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&milestone).Error; err != nil {
		return nil, err
	}
	return &milestone, nil
}

func (r *DealMilestonePostgresRepository) GetByDealID(ctx context.Context, dealID uint64) ([]*domain.DealMilestone, error) {
	var milestones []*domain.DealMilestone
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
			if err := tx.Model(&domain.DealMilestone{}).
				Where("id = ?", id).
				Update("order_index", i).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
