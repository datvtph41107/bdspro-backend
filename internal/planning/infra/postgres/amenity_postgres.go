package postgres

import (
	_db "common/db"
	_provider "common/provider"
	"context"
	"tqd/internal/domain"
	"tqd/internal/dto"
	"tqd/internal/interface/repo"
)

type AmenityRepo struct {
	_provider.CrudRepo[domain.Amenity]
}

// NewAmenityRepo creates a new amenity repository
func NewAmenityRepo(db *_db.TransactionRepo) repo.IAmenityRepo {
	repo := &AmenityRepo{}
	repo.Init(repo, db)
	return repo
}

// GetByCategory gets amenities by category
func (r *AmenityRepo) GetByCategory(ctx context.Context, category string) ([]domain.Amenity, error) {
	var amenities []domain.Amenity
	query := r.GetDB(ctx).
		Where("category = ?", category).
		Order("sort_order ASC, created_at DESC")

	if err := query.Find(&amenities).Error; err != nil {
		return nil, err
	}
	return amenities, nil
}

// ListAmenities lists amenities with filtering
func (r *AmenityRepo) ListAmenities(ctx context.Context, filter *dto.AmenityFilterDTO) ([]domain.Amenity, int64, error) {
	var amenities []domain.Amenity
	var total int64

	query := r.GetDB(ctx).Model(&domain.Amenity{})

	// Search by name or description
	if filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", searchPattern, searchPattern)
	}

	// Filter by category
	if filter.Category != "" {
		query = query.Where("category = ?", filter.Category)
	}

	// Filter by active status
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Paginate
	offset := (filter.Page - 1) * filter.Size
	query = query.Offset(int(offset)).Limit(int(filter.Size))

	// Sort
	query = query.Order("sort_order ASC, created_at DESC")

	// Fetch results
	if err := query.Find(&amenities).Error; err != nil {
		return nil, 0, err
	}

	return amenities, total, nil
}

// GetActive gets all active amenities
func (r *AmenityRepo) GetActive(ctx context.Context) ([]domain.Amenity, error) {
	var amenities []domain.Amenity
	query := r.GetDB(ctx).
		Where("is_active = ?", true).
		Order("sort_order ASC, created_at DESC")

	if err := query.Find(&amenities).Error; err != nil {
		return nil, err
	}
	return amenities, nil
}

// GetByName gets amenities by name (partial match)
func (r *AmenityRepo) GetByName(ctx context.Context, name string) ([]domain.Amenity, error) {
	var amenities []domain.Amenity
	query := r.GetDB(ctx).
		Where("name ILIKE ?", "%"+name+"%").
		Order("sort_order ASC, created_at DESC")

	if err := query.Find(&amenities).Error; err != nil {
		return nil, err
	}
	return amenities, nil
}
