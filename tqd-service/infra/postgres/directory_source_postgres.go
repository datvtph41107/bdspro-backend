package postgres

import (
	"context"
	"fmt"
	"tqd/internal/domain"
	"tqd/internal/dto"
	"tqd/internal/interface/repo"

	"gorm.io/gorm"
)

// directorySourcePostgres implements DirectorySourceRepository interface
type directorySourcePostgres struct {
	db *gorm.DB
}

// NewDirectorySourcePostgres creates a new directorySourcePostgres
func NewDirectorySourcePostgres(db *gorm.DB) repo.DirectorySourceRepository {
	return &directorySourcePostgres{db: db}
}

// Create creates a new directory source
func (r *directorySourcePostgres) Create(ctx context.Context, directorySource *domain.DirectorySource) error {
	return r.db.WithContext(ctx).Create(directorySource).Error
}

// GetByID gets directory source by ID
func (r *directorySourcePostgres) GetByID(ctx context.Context, id uint64) (*domain.DirectorySource, error) {
	var directorySource domain.DirectorySource
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&directorySource).Error
	if err != nil {
		return nil, err
	}
	return &directorySource, nil
}

// GetByCode gets directory source by code
func (r *directorySourcePostgres) GetByCode(ctx context.Context, code string) (*domain.DirectorySource, error) {
	var directorySource domain.DirectorySource
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&directorySource).Error
	if err != nil {
		return nil, err
	}
	return &directorySource, nil
}

// Update updates directory source
func (r *directorySourcePostgres) Update(ctx context.Context, directorySource *domain.DirectorySource) error {
	return r.db.WithContext(ctx).Save(directorySource).Error
}

// Delete soft deletes directory source
func (r *directorySourcePostgres) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&domain.DirectorySource{}, id).Error
}

// List lists directory sources with pagination and filters
func (r *directorySourcePostgres) List(ctx context.Context, req *dto.ListDirectorySourcesRequestDTO) ([]dto.DirectorySourceDTO, int64, error) {
	var directorySources []dto.DirectorySourceDTO
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.DirectorySource{})

	// Apply filters
	if req.Search != "" {
		searchPattern := "%" + req.Search + "%"
		query = query.Where("name ILIKE ? OR code ILIKE ? OR description ILIKE ?", searchPattern, searchPattern, searchPattern)
	}

	if req.Category != nil && *req.Category != "" {
		query = query.Where("category_id = ?", *req.Category)
	}

	if req.Type != nil && *req.Type != "" {
		query = query.Where("type = ?", *req.Type)
	}

	if req.IsActive != nil {
		query = query.Where("is_active = ?", *req.IsActive)
	}

	if req.IsRecurring != nil {
		query = query.Where("is_recurring = ?", *req.IsRecurring)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and ordering
	err := query.
		Preload("Category").
		Order("sort_order ASC, created_at DESC").
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Find(&directorySources).Error

	if err != nil {
		return nil, 0, err
	}

	return directorySources, total, nil
}

// GetActiveSources gets all active directory sources
func (r *directorySourcePostgres) GetActiveSources(ctx context.Context) ([]domain.DirectorySource, error) {
	var directorySources []domain.DirectorySource
	err := r.db.WithContext(ctx).Where("is_active = ?", true).Order("sort_order ASC, created_at DESC").Find(&directorySources).Error
	if err != nil {
		return nil, err
	}
	return directorySources, nil
}

// GetByCategory gets directory sources by category
func (r *directorySourcePostgres) GetByCategory(ctx context.Context, category string) ([]domain.DirectorySource, error) {
	var directorySources []domain.DirectorySource
	err := r.db.WithContext(ctx).Where("category = ? AND is_active = ?", category, true).Order("sort_order ASC, created_at DESC").Find(&directorySources).Error
	if err != nil {
		return nil, err
	}
	return directorySources, nil
}

// GetByType gets directory sources by type
func (r *directorySourcePostgres) GetByType(ctx context.Context, typeStr string) ([]domain.DirectorySource, error) {
	var directorySources []domain.DirectorySource
	err := r.db.WithContext(ctx).Where("type = ? AND is_active = ?", typeStr, true).Order("sort_order ASC, created_at DESC").Find(&directorySources).Error
	if err != nil {
		return nil, err
	}
	return directorySources, nil
}

// GetRecurringSources gets all recurring directory sources
func (r *directorySourcePostgres) GetRecurringSources(ctx context.Context) ([]domain.DirectorySource, error) {
	var directorySources []domain.DirectorySource
	err := r.db.WithContext(ctx).Where("is_recurring = ? AND is_active = ?", true, true).Order("sort_order ASC, created_at DESC").Find(&directorySources).Error
	if err != nil {
		return nil, err
	}
	return directorySources, nil
}

// UpdateAmounts updates expected and actual amounts
func (r *directorySourcePostgres) UpdateAmounts(ctx context.Context, id uint64, expectedAmount, actualAmount int64) error {
	updates := map[string]interface{}{}
	if expectedAmount >= 0 {
		updates["expected_amount"] = expectedAmount
	}
	if actualAmount >= 0 {
		updates["actual_amount"] = actualAmount
	}

	if len(updates) == 0 {
		return fmt.Errorf("no amounts to update")
	}

	return r.db.WithContext(ctx).Model(&domain.DirectorySource{}).Where("id = ?", id).Updates(updates).Error
}
