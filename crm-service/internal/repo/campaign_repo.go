package repo

import (
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
	"time"
)

type CampaignRepo interface {
	Create(ctx context.Context, campaign *domain.Campaign) (*domain.Campaign, error)
	Update(ctx context.Context, id uint64, campaign *domain.Campaign) (*domain.Campaign, error)
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*domain.Campaign, error)
	Search(ctx context.Context, organizationID uint64, searchDTO dto.CampaignSearchDTO) ([]*domain.Campaign, int64, error)
	GetByProductID(ctx context.Context, productID uint64) ([]*domain.Campaign, error)
	UpdateStatus(ctx context.Context, id uint64, status domain.CampaignStatus) error
	CheckTimeConflict(ctx context.Context, productID uint64, campaignType domain.CampaignType, startDate, endDate *time.Time, excludeID *uint64) (bool, error)
}