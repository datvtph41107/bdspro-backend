package postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"context"
	"strings"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo.BdsDomainRepo
type BdsDomainPostgres struct {
	db *gorm.DB
}

func NewBdsDomainPostgres(db *gorm.DB) *BdsDomainPostgres {
	return &BdsDomainPostgres{db: db}
}

// Create tạo mới bất động sản
func (r *BdsDomainPostgres) Create(ctx context.Context, bds *domain.BDSDomain) error {
	return r.db.WithContext(ctx).Create(bds).Error
}

// GetList lấy danh sách bất động sản với phân trang và bộ lọc
func (r *BdsDomainPostgres) GetList(ctx context.Context, request *dto.BdsDomainListRequest) ([]domain.BDSDomain, int64, error) {
	var results []domain.BDSDomain
	var total int64

	query := r.db.WithContext(ctx).
		Model(&domain.BDSDomain{}).
		Where("deleted_at IS NULL")

	// Áp dụng các filter
	if request.Text != nil && *request.Text != "" {
		text := "%" + strings.ToLower(*request.Text) + "%"
		query = query.Where(`LOWER(address) LIKE ? OR LOWER(street) LIKE ? OR LOWER(note) LIKE ?`, text, text, text)
	}

	if request.ProductID != nil {
		query = query.Where("product_id = ?", *request.ProductID)
	}

	if request.AssetID != nil {
		query = query.Where("asset_id = ?", *request.AssetID)
	}

	if request.ProvinceID != nil {
		query = query.Where("province_id = ?", *request.ProvinceID)
	}

	if request.WardID != nil {
		query = query.Where("ward_id = ?", *request.WardID)
	}

	if request.DistrictID != nil {
		query = query.Where("district_id = ?", *request.DistrictID)
	}

	// Preload relationships
	query = query.
		Preload("Province", "deleted_at IS NULL").
		Preload("Ward", "deleted_at IS NULL")

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if request.Size > 0 {
		query = query.Limit(int(request.GetLimit())).Offset(int(request.GetOffset()))
	}

	// Apply sorting
	if request.Sort != "" {
		query = query.Order(request.Sort)
	} else {
		query = query.Order("updated_at DESC")
	}

	// Execute query
	if err := query.Find(&results).Error; err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

// GetByID lấy chi tiết bất động sản theo ID
func (r *BdsDomainPostgres) GetByID(ctx context.Context, id uint64) (*domain.BDSDomain, error) {
	var result domain.BDSDomain
	err := r.db.WithContext(ctx).
		Preload("Province", "deleted_at IS NULL").
		Preload("Ward", "deleted_at IS NULL").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetByProductID lấy bất động sản theo Product ID
func (r *BdsDomainPostgres) GetByProductID(ctx context.Context, productID uint64) (*domain.BDSDomain, error) {
	var result domain.BDSDomain
	err := r.db.WithContext(ctx).
		Preload("Province", "deleted_at IS NULL").
		Preload("Ward", "deleted_at IS NULL").
		Where("product_id = ? AND deleted_at IS NULL", productID).
		First(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

// GetByAssetID lấy bất động sản theo Asset ID
func (r *BdsDomainPostgres) GetByAssetID(ctx context.Context, assetID uint64) (*domain.BDSDomain, error) {
	var result domain.BDSDomain
	err := r.db.WithContext(ctx).
		Preload("Province", "deleted_at IS NULL").
		Preload("Ward", "deleted_at IS NULL").
		Where("asset_id = ? AND deleted_at IS NULL", assetID).
		First(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}
