package repository

import (
	"context"
	"organization/internal/domain/entity"
)

type DealMilestoneRepository interface {
	Create(ctx context.Context, milestone *entity.DealMilestone) (*entity.DealMilestone, error)
	Update(ctx context.Context, milestone *entity.DealMilestone) (*entity.DealMilestone, error)
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*entity.DealMilestone, error)
	GetByDealID(ctx context.Context, dealID uint64) ([]*entity.DealMilestone, error)
	UpdateOrder(ctx context.Context, milestoneIDs []uint64) error
} 