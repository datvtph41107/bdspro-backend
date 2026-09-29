package repo

import (
	"bdspro/internal/domain"
	"context"
)

type DealMilestoneRepository interface {
	Create(ctx context.Context, milestone *domain.DealMilestone) (*domain.DealMilestone, error)
	Update(ctx context.Context, milestone *domain.DealMilestone) (*domain.DealMilestone, error)
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*domain.DealMilestone, error)
	GetByDealID(ctx context.Context, dealID uint64) ([]*domain.DealMilestone, error)
	UpdateOrder(ctx context.Context, milestoneIDs []uint64) error
}
