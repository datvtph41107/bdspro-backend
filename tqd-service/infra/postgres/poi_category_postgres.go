package postgres

import (
	"context"
	"fmt"
	"strings"

	_db "common/db"
	_provider "common/provider"
	"tqd/internal/domain"
	"tqd/internal/dto"
	"tqd/internal/interface/repo"

	"gorm.io/gorm"
)

// poiCategoryPostgres implements IPoiCategoryRepo
type poiCategoryPostgres struct {
	_provider.CrudRepo[domain.PoiCategory]
}

// NewPoiCategoryPostgres creates new POI category repository
func NewPoiCategoryPostgres(db *_db.TransactionRepo) repo.IPoiCategoryRepo {
	r := &poiCategoryPostgres{}
	r.Init(r, db)
	return r
}

// BeforeSave hook - called before create/update
func (r *poiCategoryPostgres) BeforeSave(ctx context.Context, id *uint64, entity *domain.PoiCategory) error {
	// Validate entity
	if err := entity.Validate(); err != nil {
		return err
	}

	// Check if code already exists
	existing, err := r.FindByCode(ctx, entity.Code)
	if err == nil && existing != nil {
		if id == nil || existing.ID != *id {
			return fmt.Errorf("category with code %s already exists", entity.Code)
		}
	}

	return nil
}

// AfterSave hook - called after create/update
func (r *poiCategoryPostgres) AfterSave(ctx context.Context, id *uint64, entity *domain.PoiCategory) error {
	// Update path and level for this category and its children
	if err := r.updateCategoryPath(ctx, entity); err != nil {
		return err
	}
	return nil
}

// Create overrides base Create
func (r *poiCategoryPostgres) Create(ctx context.Context, entity *domain.PoiCategory) error {
	if err := r.BeforeSave(ctx, nil, entity); err != nil {
		return err
	}

	if err := r.GetDB(ctx).Create(entity).Error; err != nil {
		return err
	}

	return r.AfterSave(ctx, &entity.ID, entity)
}

// Update overrides base Update
func (r *poiCategoryPostgres) Update(ctx context.Context, id uint64, entity *domain.PoiCategory) error {
	// Get existing
	existing, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return gorm.ErrRecordNotFound
	}

	// Preserve immutable fields
	entity.CreatedAt = existing.CreatedAt
	entity.CreatedBy = existing.CreatedBy
	entity.POICount = existing.POICount

	if err := r.BeforeSave(ctx, &id, entity); err != nil {
		return err
	}

	if err := r.GetDB(ctx).Save(entity).Error; err != nil {
		return err
	}

	// If parent changed, update paths for this category and its children.
	if parentIDChanged(existing.ParentID, entity.ParentID) {
		if err := r.updateCategoryPath(ctx, entity); err != nil {
			return err
		}
	}

	return r.AfterSave(ctx, &id, entity)
}

func parentIDChanged(before, after *uint64) bool {
	if before == nil || after == nil {
		return before != nil || after != nil
	}
	return *before != *after
}

// ==================== Custom Methods ====================

// FindByCode finds category by code
func (r *poiCategoryPostgres) FindByCode(ctx context.Context, code string) (*domain.PoiCategory, error) {
	var category domain.PoiCategory
	err := r.GetDB(ctx).
		Where("code = ? AND deleted_at IS NULL", code).
		First(&category).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

// FindByParent finds categories by parent ID
func (r *poiCategoryPostgres) FindByParent(ctx context.Context, parentID *uint64) ([]domain.PoiCategory, error) {
	var categories []domain.PoiCategory
	query := r.GetDB(ctx).Where("deleted_at IS NULL")

	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", *parentID)
	}

	err := query.Order("sort_order ASC, created_at DESC").Find(&categories).Error
	return categories, err
}

// GetTree gets category tree
func (r *poiCategoryPostgres) GetTree(ctx context.Context) ([]domain.PoiCategory, error) {
	var categories []domain.PoiCategory

	// Get all categories
	err := r.GetDB(ctx).
		Where("deleted_at IS NULL").
		Order("sort_order ASC, created_at DESC").
		Find(&categories).Error
	if err != nil {
		return nil, err
	}

	// Build tree in memory
	categoryMap := make(map[uint64]*domain.PoiCategory)
	for i := range categories {
		categoryMap[categories[i].ID] = &categories[i]
	}

	// Set children
	var roots []domain.PoiCategory
	for i := range categories {
		if categories[i].ParentID != nil {
			if parent, ok := categoryMap[*categories[i].ParentID]; ok {
				parent.Children = append(parent.Children, categories[i])
			}
		} else {
			roots = append(roots, categories[i])
		}
	}

	return roots, nil
}

// UpdatePath updates category path
func (r *poiCategoryPostgres) UpdatePath(ctx context.Context, id uint64, path string) error {
	return r.GetDB(ctx).Model(&domain.PoiCategory{}).
		Where("id = ?", id).
		Update("path", path).Error
}

// UpdateLevel updates category level
func (r *poiCategoryPostgres) UpdateLevel(ctx context.Context, id uint64, level int) error {
	return r.GetDB(ctx).Model(&domain.PoiCategory{}).
		Where("id = ?", id).
		Update("level", level).Error
}

// IncrementPOICount increments POI count
func (r *poiCategoryPostgres) IncrementPOICount(ctx context.Context, id uint64) error {
	return r.GetDB(ctx).Model(&domain.PoiCategory{}).
		Where("id = ?", id).
		UpdateColumn("poi_count", gorm.Expr("poi_count + 1")).Error
}

// DecrementPOICount decrements POI count
func (r *poiCategoryPostgres) DecrementPOICount(ctx context.Context, id uint64) error {
	return r.GetDB(ctx).Model(&domain.PoiCategory{}).
		Where("id = ? AND poi_count > 0", id).
		UpdateColumn("poi_count", gorm.Expr("poi_count - 1")).Error
}

// ListWithFilter lists categories with filters
func (r *poiCategoryPostgres) ListWithFilter(ctx context.Context, filter *dto.PoiCategoryFilter) ([]domain.PoiCategory, int64, error) {
	var categories []domain.PoiCategory
	var total int64

	query := r.GetDB(ctx).Model(&domain.PoiCategory{}).Where("deleted_at IS NULL")

	// Apply filters
	if filter.Search != "" {
		searchPattern := "%" + strings.ToLower(filter.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(code) LIKE ?", searchPattern, searchPattern)
	}
	if filter.Code != "" {
		query = query.Where("LOWER(code) LIKE ?", "%"+strings.ToLower(filter.Code)+"%")
	}
	if filter.Name != "" {
		query = query.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(filter.Name)+"%")
	}
	if filter.ParentID != nil {
		query = query.Where("parent_id = ?", *filter.ParentID)
	}
	if filter.Level != nil {
		query = query.Where("level = ?", *filter.Level)
	}
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and sorting
	err := query.
		Order("sort_order ASC, created_at DESC").
		Offset(filter.GetOffset()).
		Limit(filter.GetLimit()).
		Find(&categories).Error

	return categories, total, err
}

// ==================== Helper Methods ====================

// updateCategoryPath updates path and level for category and its children
func (r *poiCategoryPostgres) updateCategoryPath(ctx context.Context, category *domain.PoiCategory) error {
	// Calculate path and level
	path := fmt.Sprintf("/%d", category.ID)
	level := 1

	if category.ParentID != nil {
		var parent domain.PoiCategory
		if err := r.GetDB(ctx).First(&parent, *category.ParentID).Error; err != nil {
			return err
		}
		path = parent.Path + path
		level = parent.Level + 1
	}

	// Update current category
	if err := r.GetDB(ctx).Model(category).Updates(map[string]interface{}{
		"path":  path,
		"level": level,
	}).Error; err != nil {
		return err
	}

	// Update children recursively
	var children []domain.PoiCategory
	if err := r.GetDB(ctx).Where("parent_id = ?", category.ID).Find(&children).Error; err != nil {
		return err
	}

	for i := range children {
		if err := r.updateCategoryPath(ctx, &children[i]); err != nil {
			return err
		}
	}

	return nil
}
