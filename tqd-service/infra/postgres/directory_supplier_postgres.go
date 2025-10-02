package postgres

import (
	"context"
	"strings"
	"tqd/internal/domain"
	"tqd/internal/interface/repo"

	"gorm.io/gorm"
)

// directorySupplierPostgres implements DirectorySupplierRepository interface
type directorySupplierPostgres struct {
	db *gorm.DB
}

// NewDirectorySupplierPostgres creates a new directorySupplierPostgres
func NewDirectorySupplierPostgres(db *gorm.DB) repo.DirectorySupplierRepository {
	return &directorySupplierPostgres{db: db}
}

// Create creates a new directory supplier
func (r *directorySupplierPostgres) Create(ctx context.Context, supplier *domain.DirectorySupplier) error {
	return r.db.WithContext(ctx).Create(supplier).Error
}

// GetByID gets directory supplier by ID
func (r *directorySupplierPostgres) GetByID(ctx context.Context, id uint64) (*domain.DirectorySupplier, error) {
	var supplier domain.DirectorySupplier
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&supplier).Error
	if err != nil {
		return nil, err
	}
	return &supplier, nil
}

// Update updates directory supplier
func (r *directorySupplierPostgres) Update(ctx context.Context, supplier *domain.DirectorySupplier) error {
	return r.db.WithContext(ctx).Save(supplier).Error
}

// Delete deletes directory supplier by ID
func (r *directorySupplierPostgres) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&domain.DirectorySupplier{}, id).Error
}

// List gets list of directory suppliers with filters
func (r *directorySupplierPostgres) List(ctx context.Context, page, size int, search string, categories []string, isActive *bool, minRating, maxRating *float64) ([]domain.DirectorySupplier, int64, error) {
	var suppliers []domain.DirectorySupplier
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.DirectorySupplier{})

	// Apply search filter
	if search != "" {
		searchPattern := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(code) LIKE ? OR LOWER(description) LIKE ? OR LOWER(contact_person) LIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern)
	}

	// Apply categories filter
	if len(categories) > 0 {
		for _, category := range categories {
			query = query.Where("categories LIKE ?", "%\""+category+"\"%")
		}
	}

	// Apply isActive filter
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	// Apply rating filters
	if minRating != nil {
		query = query.Where("rating >= ?", *minRating)
	}
	if maxRating != nil {
		query = query.Where("rating <= ?", *maxRating)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and get results
	offset := (page - 1) * size
	err := query.Order("created_at DESC").Offset(offset).Limit(size).Find(&suppliers).Error
	if err != nil {
		return nil, 0, err
	}

	return suppliers, total, nil
}
