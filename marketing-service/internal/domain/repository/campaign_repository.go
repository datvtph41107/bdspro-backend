package repository

import (
	"context"
	"marketing/internal/domain/model"
	"marketing/internal/dto"
	"time"
)

type CampaignRepository interface {
	Create(ctx context.Context, campaign *model.Campaign) (*model.Campaign, error)
	Update(ctx context.Context, id uint64, campaign *model.Campaign) (*model.Campaign, error)
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*model.Campaign, error)
	Search(ctx context.Context, organizationID uint64, searchDTO dto.CampaignSearchDTO) ([]*model.Campaign, int64, error)
	GetByProductID(ctx context.Context, productID uint64) ([]*model.Campaign, error)
	UpdateStatus(ctx context.Context, id uint64, status model.CampaignStatus) error
	CheckTimeConflict(ctx context.Context, productID uint64, campaignType model.CampaignType, startDate, endDate *time.Time, excludeID *uint64) (bool, error)
} 