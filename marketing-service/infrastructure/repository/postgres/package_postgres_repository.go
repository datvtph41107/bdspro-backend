package postgres

import (
	"context"
	"fmt"
	"marketing/internal/domain/model"
	"marketing/internal/domain/repository"
	"marketing/internal/dto"

	"gorm.io/gorm"
)

type PackagePostgresRepository struct {
	db *gorm.DB
}

func NewPackagePostgresRepository(db *gorm.DB) repository.PackageRepository {
	return &PackagePostgresRepository{
		db: db,
	}
}

func (r *PackagePostgresRepository) Create(ctx context.Context, pkg *model.Package) (*model.Package, error) {
	if err := r.db.WithContext(ctx).Create(pkg).Error; err != nil {
		return nil, fmt.Errorf("failed to create package: %w", err)
	}
	return pkg, nil
}

func (r *PackagePostgresRepository) Update(ctx context.Context, id uint64, pkg *model.Package) (*model.Package, error) {
	if err := r.db.WithContext(ctx).Model(&model.Package{}).Where("id = ?", id).Updates(pkg).Error; err != nil {
		return nil, fmt.Errorf("failed to update package: %w", err)
	}
	return r.GetByID(ctx, id)
}

func (r *PackagePostgresRepository) Delete(ctx context.Context, id uint64) error {
	if err := r.db.WithContext(ctx).Delete(&model.Package{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete package: %w", err)
	}
	return nil
}

func (r *PackagePostgresRepository) GetByID(ctx context.Context, id uint64) (*model.Package, error) {
	var pkg model.Package
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&pkg).Error; err != nil {
		return nil, fmt.Errorf("package not found: %w", err)
	}
	return &pkg, nil
}

func (r *PackagePostgresRepository) Search(ctx context.Context, organizationID uint64, searchDTO dto.PackageSearchDTO) ([]*model.Package, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.Package{})

	// Apply filters
	if searchDTO.Type != "" {
		query = query.Where("type = ?", searchDTO.Type)
	}
	if searchDTO.Status != "" {
		query = query.Where("status = ?", searchDTO.Status)
	}
	if searchDTO.IsPopular != nil && *searchDTO.IsPopular {
		query = query.Where("is_popular = ?", true)
	}
	if searchDTO.IsRecommended != nil && *searchDTO.IsRecommended {
		query = query.Where("is_recommended = ?", true)
	}
	if searchDTO.MinPrice != nil {
		query = query.Where("price >= ?", *searchDTO.MinPrice)
	}
	if searchDTO.MaxPrice != nil {
		query = query.Where("price <= ?", *searchDTO.MaxPrice)
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count packages: %w", err)
	}

	// Apply pagination
	offset := (searchDTO.Page - 1) * searchDTO.Limit
	query = query.Offset(offset).Limit(searchDTO.Limit)

	// Apply ordering
	query = query.Order("created_at DESC")

	var packages []*model.Package
	if err := query.Find(&packages).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to search packages: %w", err)
	}

	return packages, total, nil
}

func (r *PackagePostgresRepository) GetByType(ctx context.Context, packageType model.PackageType) ([]*model.Package, error) {
	var packages []*model.Package
	if err := r.db.WithContext(ctx).Where("type = ? AND status = ?", packageType, model.PackageStatusActive).Find(&packages).Error; err != nil {
		return nil, fmt.Errorf("failed to get packages by type: %w", err)
	}
	return packages, nil
}

func (r *PackagePostgresRepository) GetPopularPackages(ctx context.Context, limit int) ([]*model.Package, error) {
	var packages []*model.Package
	query := r.db.WithContext(ctx).
		Where("is_popular = ? AND status = ?", true, model.PackageStatusActive).
		Order("created_at DESC").
		Limit(limit)

	if err := query.Find(&packages).Error; err != nil {
		return nil, fmt.Errorf("failed to get popular packages: %w", err)
	}
	return packages, nil
}

func (r *PackagePostgresRepository) GetRecommendedPackages(ctx context.Context, limit int) ([]*model.Package, error) {
	var packages []*model.Package
	query := r.db.WithContext(ctx).
		Where("is_recommended = ? AND status = ?", true, model.PackageStatusActive).
		Order("created_at DESC").
		Limit(limit)

	if err := query.Find(&packages).Error; err != nil {
		return nil, fmt.Errorf("failed to get recommended packages: %w", err)
	}
	return packages, nil
}

func (r *PackagePostgresRepository) GetBestValuePackages(ctx context.Context, limit int) ([]*model.Package, error) {
	var packages []*model.Package
	query := r.db.WithContext(ctx).
		Where("status = ?", model.PackageStatusActive).
		Order("price ASC").
		Limit(limit)

	if err := query.Find(&packages).Error; err != nil {
		return nil, fmt.Errorf("failed to get best value packages: %w", err)
	}
	return packages, nil
}

func (r *PackagePostgresRepository) GetPackagesByProductType(ctx context.Context, productType string) ([]*model.Package, error) {
	// TODO: Implement logic to get packages based on product type
	// For now, return all active packages
	var packages []*model.Package
	if err := r.db.WithContext(ctx).Where("status = ?", model.PackageStatusActive).Find(&packages).Error; err != nil {
		return nil, fmt.Errorf("failed to get packages by product type: %w", err)
	}
	return packages, nil
}

func (r *PackagePostgresRepository) UpdateStatus(ctx context.Context, id uint64, status model.PackageStatus) error {
	if err := r.db.WithContext(ctx).Model(&model.Package{}).Where("id = ?", id).Update("status", status).Error; err != nil {
		return fmt.Errorf("failed to update package status: %w", err)
	}
	return nil
} 