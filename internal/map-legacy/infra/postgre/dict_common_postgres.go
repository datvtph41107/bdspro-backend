package postgre

import (
	"context"
	"fmt"
	"strings"

	"map/internal/domain"
	"map/internal/interface/repo"

	"gorm.io/gorm"
)

// DictCommonPostgres implements DictCommonRepository interface
type DictCommonPostgres struct {
	db *gorm.DB
}

// NewDictCommonPostgres creates a new dictionary common postgres repository
// @bind: map/internal/interface/repo.DictCommonRepository
func NewDictCommonPostgres(db *gorm.DB) repo.DictCommonRepository {
	return &DictCommonPostgres{
		db: db,
	}
}

// GetByCode gets a dictionary item by code
func (r *DictCommonPostgres) GetByCode(ctx context.Context, code string) (*domain.DictCommon, error) {
	var item domain.DictCommon
	result := r.db.WithContext(ctx).Where("code = ?", code).First(&item)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get dictionary item by code: %w", result.Error)
	}
	return &item, nil
}

// GetByCategory gets all dictionary items by category
func (r *DictCommonPostgres) GetByCategory(ctx context.Context, category string) ([]domain.DictCommon, error) {
	var items []domain.DictCommon
	result := r.db.WithContext(ctx).Where("category = ?", category).Order("sort_order ASC, name ASC").Find(&items)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get dictionary items by category: %w", result.Error)
	}
	return items, nil
}

// GetByCategoryAndParent gets dictionary items by category and parent
func (r *DictCommonPostgres) GetByCategoryAndParent(ctx context.Context, category string, parentID *uint) ([]domain.DictCommon, error) {
	var items []domain.DictCommon
	query := r.db.WithContext(ctx).Where("category = ?", category)

	if parentID != nil {
		query = query.Where("parent_id = ?", *parentID)
	} else {
		query = query.Where("parent_id IS NULL")
	}

	result := query.Order("sort_order ASC, name ASC").Find(&items)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get dictionary items by category and parent: %w", result.Error)
	}
	return items, nil
}

// GetHierarchical gets dictionary items in hierarchical structure
func (r *DictCommonPostgres) GetHierarchical(ctx context.Context, category string) ([]domain.DictCommonWithChildren, error) {
	// Get root items (parent_id IS NULL)
	rootItems, err := r.GetByCategoryAndParent(ctx, category, nil)
	if err != nil {
		return nil, err
	}

	var result []domain.DictCommonWithChildren
	for _, item := range rootItems {
		children, err := r.getChildrenRecursive(ctx, item.ID)
		if err != nil {
			return nil, err
		}

		result = append(result, domain.DictCommonWithChildren{
			DictCommon: item,
			Children:   children,
		})
	}

	return result, nil
}

// getChildrenRecursive recursively gets children of a dictionary item
func (r *DictCommonPostgres) getChildrenRecursive(ctx context.Context, parentID uint) ([]domain.DictCommonWithChildren, error) {
	children, err := r.GetByCategoryAndParent(ctx, "", &parentID)
	if err != nil {
		return nil, err
	}

	var result []domain.DictCommonWithChildren
	for _, child := range children {
		grandChildren, err := r.getChildrenRecursive(ctx, child.ID)
		if err != nil {
			return nil, err
		}

		result = append(result, domain.DictCommonWithChildren{
			DictCommon: child,
			Children:   grandChildren,
		})
	}

	return result, nil
}

// Create creates a new dictionary item
func (r *DictCommonPostgres) Create(ctx context.Context, item *domain.DictCommon) (*domain.DictCommon, error) {
	result := r.db.WithContext(ctx).Create(item)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to create dictionary item: %w", result.Error)
	}
	return item, nil
}

// Update updates an existing dictionary item
func (r *DictCommonPostgres) Update(ctx context.Context, id uint, item *domain.DictCommon) (*domain.DictCommon, error) {
	result := r.db.WithContext(ctx).Model(&domain.DictCommon{}).Where("id = ?", id).Updates(item)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to update dictionary item: %w", result.Error)
	}

	// Get updated item
	updatedItem, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return updatedItem, nil
}

// Delete deletes a dictionary item by ID
func (r *DictCommonPostgres) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&domain.DictCommon{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete dictionary item: %w", result.Error)
	}
	return nil
}

// GetByID gets a dictionary item by ID
func (r *DictCommonPostgres) GetByID(ctx context.Context, id uint) (*domain.DictCommon, error) {
	var item domain.DictCommon
	result := r.db.WithContext(ctx).First(&item, id)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get dictionary item by ID: %w", result.Error)
	}
	return &item, nil
}

// List gets dictionary items with filtering and pagination
func (r *DictCommonPostgres) List(ctx context.Context, filter *domain.DictCommonFilter, page, size int) ([]domain.DictCommon, int64, error) {
	var items []domain.DictCommon
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.DictCommon{})

	// Apply filters
	if filter != nil {
		if filter.Category != nil {
			query = query.Where("category = ?", *filter.Category)
		}
		if filter.IsActive != nil {
			query = query.Where("is_active = ?", *filter.IsActive)
		}
		if filter.ParentID != nil {
			query = query.Where("parent_id = ?", *filter.ParentID)
		}
		if filter.Search != "" {
			searchTerm := "%" + strings.ToLower(filter.Search) + "%"
			query = query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", searchTerm, searchTerm)
		}
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count dictionary items: %w", err)
	}

	// Apply pagination and ordering
	offset := (page - 1) * size
	result := query.Order("category ASC, sort_order ASC, name ASC").Offset(offset).Limit(size).Find(&items)
	if result.Error != nil {
		return nil, 0, fmt.Errorf("failed to list dictionary items: %w", result.Error)
	}

	return items, total, nil
}

// GetSummary gets summary statistics of dictionary data
func (r *DictCommonPostgres) GetSummary(ctx context.Context) (*domain.DictCommonSummary, error) {
	var summary domain.DictCommonSummary
	var totalItems int64
	var activeItems int64

	// Get total items
	if err := r.db.WithContext(ctx).Model(&domain.DictCommon{}).Count(&totalItems).Error; err != nil {
		return nil, fmt.Errorf("failed to count total items: %w", err)
	}

	// Get active items
	if err := r.db.WithContext(ctx).Model(&domain.DictCommon{}).Where("is_active = ?", true).Count(&activeItems).Error; err != nil {
		return nil, fmt.Errorf("failed to count active items: %w", err)
	}
	summary.TotalItems = int(totalItems)
	summary.ActiveItems = int(activeItems)

	// Get category counts
	var categoryCounts []struct {
		Category string `json:"category"`
		Count    int    `json:"count"`
	}

	if err := r.db.WithContext(ctx).Model(&domain.DictCommon{}).
		Select("category, COUNT(*) as count").
		Group("category").
		Order("count DESC").
		Scan(&categoryCounts).Error; err != nil {
		return nil, fmt.Errorf("failed to get category counts: %w", err)
	}

	summary.Categories = make(map[string]int)
	summary.TopCategories = make([]domain.CategoryCount, len(categoryCounts))

	for i, cc := range categoryCounts {
		summary.Categories[cc.Category] = cc.Count
		summary.TopCategories[i] = domain.CategoryCount{
			Category: cc.Category,
			Count:    cc.Count,
		}
	}

	return &summary, nil
}

// BulkCreate creates multiple dictionary items
func (r *DictCommonPostgres) BulkCreate(ctx context.Context, items []domain.DictCommon) error {
	if len(items) == 0 {
		return nil
	}

	result := r.db.WithContext(ctx).CreateInBatches(items, 100)
	if result.Error != nil {
		return fmt.Errorf("failed to bulk create dictionary items: %w", result.Error)
	}
	return nil
}

// BulkUpdate updates multiple dictionary items
func (r *DictCommonPostgres) BulkUpdate(ctx context.Context, items []domain.DictCommon) error {
	if len(items) == 0 {
		return nil
	}

	for _, item := range items {
		if err := r.db.WithContext(ctx).Model(&domain.DictCommon{}).Where("id = ?", item.ID).Updates(item).Error; err != nil {
			return fmt.Errorf("failed to bulk update dictionary item %d: %w", item.ID, err)
		}
	}
	return nil
}

// BulkDelete deletes multiple dictionary items
func (r *DictCommonPostgres) BulkDelete(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}

	result := r.db.WithContext(ctx).Delete(&domain.DictCommon{}, ids)
	if result.Error != nil {
		return fmt.Errorf("failed to bulk delete dictionary items: %w", result.Error)
	}
	return nil
}

// GetActiveByCategory gets active dictionary items by category
func (r *DictCommonPostgres) GetActiveByCategory(ctx context.Context, category string) ([]domain.DictCommon, error) {
	var items []domain.DictCommon
	result := r.db.WithContext(ctx).Where("category = ? AND is_active = ?", category, true).Order("sort_order ASC, name ASC").Find(&items)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get active dictionary items by category: %w", result.Error)
	}
	return items, nil
}

// Search searches dictionary items by name or description
func (r *DictCommonPostgres) Search(ctx context.Context, query string, category *string) ([]domain.DictCommon, error) {
	var items []domain.DictCommon
	dbQuery := r.db.WithContext(ctx)

	searchTerm := "%" + strings.ToLower(query) + "%"
	dbQuery = dbQuery.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", searchTerm, searchTerm)

	if category != nil {
		dbQuery = dbQuery.Where("category = ?", *category)
	}

	result := dbQuery.Order("category ASC, sort_order ASC, name ASC").Find(&items)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to search dictionary items: %w", result.Error)
	}
	return items, nil
}
