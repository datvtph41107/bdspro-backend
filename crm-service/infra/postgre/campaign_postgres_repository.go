package postgre

import (
	"context"
	"time"

	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/repo"

	"gorm.io/gorm"
)

type CampaignPostgresRepository struct {
	db *gorm.DB
}

func NewCampaignPostgresRepository(db *gorm.DB) repo.CampaignRepo {
	return &CampaignPostgresRepository{db: db}
}

func (r *CampaignPostgresRepository) Create(ctx context.Context, campaign *domain.Campaign) (*domain.Campaign, error) {
	if err := r.db.WithContext(ctx).Create(campaign).Error; err != nil {
		return nil, err
	}
	return campaign, nil
}

func (r *CampaignPostgresRepository) Update(ctx context.Context, id uint64, campaign *domain.Campaign) (*domain.Campaign, error) {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Updates(campaign).Error; err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *CampaignPostgresRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Campaign{}).Error
}

func (r *CampaignPostgresRepository) GetByID(ctx context.Context, id uint64) (*domain.Campaign, error) {
	var campaign domain.Campaign
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&campaign).Error; err != nil {
		return nil, err
	}
	return &campaign, nil
}

func (r *CampaignPostgresRepository) Search(ctx context.Context, organizationID uint64, searchDTO dto.CampaignSearchDTO) ([]*domain.Campaign, int64, error) {
	var campaigns []*domain.Campaign
	var total int64

	query := r.db.WithContext(ctx).Where("organization_id = ? AND deleted_at IS NULL", organizationID)

	// Apply filters
	if searchDTO.Name != "" {
		query = query.Where("name ILIKE ?", "%"+searchDTO.Name+"%")
	}
	if searchDTO.Type != "" {
		query = query.Where("type = ?", searchDTO.Type)
	}
	if searchDTO.Status != "" {
		query = query.Where("status = ?", searchDTO.Status)
	}
	if searchDTO.ProductID != nil {
		query = query.Where("product_id = ?", *searchDTO.ProductID)
	}
	if searchDTO.ProductType != "" {
		query = query.Where("product_type = ?", searchDTO.ProductType)
	}
	if searchDTO.PackageID != nil {
		query = query.Where("package_id = ?", *searchDTO.PackageID)
	}
	if searchDTO.IsActive != nil {
		query = query.Where("is_active = ?", *searchDTO.IsActive)
	}
	if searchDTO.FromDate != nil {
		query = query.Where("start_date >= ?", searchDTO.FromDate)
	}
	if searchDTO.ToDate != nil {
		query = query.Where("end_date <= ?", searchDTO.ToDate)
	}

	// Get total count
	if err := query.Model(&domain.Campaign{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := searchDTO.GetOffset()
	limit := searchDTO.GetLimit()

	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&campaigns).Error; err != nil {
		return nil, 0, err
	}

	return campaigns, total, nil
}

func (r *CampaignPostgresRepository) GetByProductID(ctx context.Context, productID uint64) ([]*domain.Campaign, error) {
	var campaigns []*domain.Campaign
	if err := r.db.WithContext(ctx).Where("product_id = ? AND deleted_at IS NULL", productID).Find(&campaigns).Error; err != nil {
		return nil, err
	}
	return campaigns, nil
}

func (r *CampaignPostgresRepository) UpdateStatus(ctx context.Context, id uint64, status domain.CampaignStatus) error {
	return r.db.WithContext(ctx).Model(&domain.Campaign{}).Where("id = ?", id).Update("status", status).Error
}

func (r *CampaignPostgresRepository) CheckTimeConflict(ctx context.Context, productID uint64, campaignType domain.CampaignType, startDate, endDate *time.Time, excludeID *uint64) (bool, error) {
	var count int64

	query := r.db.WithContext(ctx).Model(&domain.Campaign{}).
		Where("product_id = ? AND type = ? AND deleted_at IS NULL", productID, campaignType).
		Where("((start_date <= ? AND end_date >= ?) OR (start_date <= ? AND end_date >= ?) OR (start_date >= ? AND end_date <= ?))",
			startDate, startDate, endDate, endDate, startDate, endDate)

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}