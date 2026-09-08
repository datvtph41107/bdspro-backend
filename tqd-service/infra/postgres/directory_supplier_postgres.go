package postgres

import (
	"context"
	"strings"
	"tqd/internal/domain"
	"tqd/internal/dto"
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
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &supplier, nil
}

// GetByCode gets directory supplier by code
func (r *directorySupplierPostgres) GetByCode(ctx context.Context, code string) (*domain.DirectorySupplier, error) {
	var supplier domain.DirectorySupplier
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&supplier).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
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
func (r *directorySupplierPostgres) List(ctx context.Context, req *dto.DirectorySupplierFilterDTO) ([]domain.DirectorySupplier, int64, error) {
	var suppliers []domain.DirectorySupplier
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.DirectorySupplier{})

	// Apply search filter
	if req.Search != "" {
		searchPattern := "%" + strings.ToLower(req.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(code) LIKE ? OR LOWER(description) LIKE ? OR LOWER(contact_person) LIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern)
	}

	// Apply categories filter (JSON contains)
	if len(req.Categories) > 0 {
		for _, category := range req.Categories {
			query = query.Where("categories LIKE ?", "%\""+category+"\"%")
		}
	}

	// Apply isActive filter
	if req.IsActive != nil {
		query = query.Where("is_active = ?", *req.IsActive)
	}

	// Apply rating filters
	if req.MinRating != nil {
		query = query.Where("rating >= ?", *req.MinRating)
	}
	if req.MaxRating != nil {
		query = query.Where("rating <= ?", *req.MaxRating)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	switch req.SortBy {
	case "rating":
		if req.SortOrder == "desc" {
			query = query.Order("rating DESC")
		} else {
			query = query.Order("rating ASC")
		}
	case "name":
		if req.SortOrder == "desc" {
			query = query.Order("name DESC")
		} else {
			query = query.Order("name ASC")
		}
	default:
		query = query.Order("created_at DESC")
	}

	// Apply pagination
	offset := (req.Page - 1) * req.Size
	err := query.Offset(int(offset)).Limit(int(req.Size)).Find(&suppliers).Error
	if err != nil {
		return nil, 0, err
	}

	return suppliers, total, nil
}
