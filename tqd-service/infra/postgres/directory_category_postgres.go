package postgres

import (
	"context"
	"tqd/internal/domain"
	"tqd/internal/dto"
	"tqd/internal/interface/repo"

	"gorm.io/gorm"
)

type DirectoryCategoryPostgresRepo struct {
	db *gorm.DB
}

func NewDirectoryCategoryPostgresRepo(db *gorm.DB) repo.DirectoryCategoryRepository {
	return &DirectoryCategoryPostgresRepo{db: db}
}

func (r *DirectoryCategoryPostgresRepo) Create(ctx context.Context, category *domain.DirectoryCategory) error {
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *DirectoryCategoryPostgresRepo) GetByID(ctx context.Context, id uint64) (*domain.DirectoryCategory, error) {
	var category domain.DirectoryCategory
	err := r.db.WithContext(ctx).First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *DirectoryCategoryPostgresRepo) Update(ctx context.Context, category *domain.DirectoryCategory) error {
	return r.db.WithContext(ctx).Save(category).Error
}

func (r *DirectoryCategoryPostgresRepo) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&domain.DirectoryCategory{}, id).Error
}

func (r *DirectoryCategoryPostgresRepo) List(ctx context.Context, request *dto.ListDirectoryCategoriesRequest) ([]domain.DirectoryCategory, int64, error) {
	var categories []domain.DirectoryCategory
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.DirectoryCategory{})

	// Apply filters
	if request.Search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+request.Search+"%", "%"+request.Search+"%")
	}
	if request.IsActive != nil {
		query = query.Where("is_active = ?", *request.IsActive)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := query.
		Offset(request.GetOffset()).
		Limit(request.GetLimit()).
		Order("sort_order ASC, created_at DESC").
		Find(&categories).Error; err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

func (r *DirectoryCategoryPostgresRepo) GetByParentID(ctx context.Context, parentID uint64) ([]domain.DirectoryCategory, error) {
	var categories []domain.DirectoryCategory
	err := r.db.WithContext(ctx).Where("parent_id = ?", parentID).Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}
