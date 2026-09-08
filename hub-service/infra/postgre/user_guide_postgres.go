package postgres

import (
	"context"
	"strings"

	_db "common/db"
	_provider "common/provider"
	"hub/internal/domain"
	_repo "hub/internal/repo"

	_dto "common/domain/dto"
)

type UserGuideRepo struct {
	_provider.CrudRepo[domain.UserGuideEntity]
}

func NewUserGuideRepo(db *_db.TransactionRepo) _repo.IUserGuideRepo {
	repo := &UserGuideRepo{}
	repo.Init(repo, db)
	return repo
}

// GetListWithFilter retrieves user guides with filters and pagination
func (r *UserGuideRepo) GetListWithFilter(ctx context.Context, title string, groupKey string, pagable _dto.IPagable) ([]*domain.UserGuideEntity, int64, error) {
	var data []*domain.UserGuideEntity
	var total int64

	db := r.GetDB(ctx).Model(&domain.UserGuideEntity{}).Where("deleted_at IS NULL")

	// Apply filters
	if title != "" {
		db = db.Where("title ILIKE ?", "%"+title+"%")
	}

	if groupKey != "" {
		db = db.Where("group_key = ?", groupKey)
	}

	// Count total
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := pagable.GetOffset()
	limit := pagable.GetLimit()

	db = db.Offset(offset).Limit(limit)

	// Get data
	if err := db.Order("created_at DESC").Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

// GetGroups retrieves all unique groups with count
func (r *UserGuideRepo) GetGroups(ctx context.Context) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	db := r.GetDB(ctx).Model(&domain.UserGuideEntity{})

	err := db.
		Select("group_key, COUNT(*) as count").
		Where("deleted_at IS NULL").
		Group("group_key").
		Order("group_key ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}

// GetSimpleListWithText retrieves user guides with text search (title) and pagination
func (r *UserGuideRepo) GetSimpleListWithText(ctx context.Context, text string, groupKey string, pagable _dto.IPagable) ([]*domain.UserGuideEntity, int64, error) {
	var data []*domain.UserGuideEntity
	var total int64

	db := r.GetDB(ctx).Model(&domain.UserGuideEntity{}).Where("deleted_at IS NULL")

	// Apply text search: lowercase và thay khoảng trống bằng %
	if text != "" {
		// Lowercase
		searchText := strings.ToLower(text)
		// Thay khoảng trống bằng %
		searchText = strings.ReplaceAll(searchText, " ", "%")
		// Tạo pattern LIKE
		likePattern := "%" + searchText + "%"

		// Search trong title (lowercase)
		db = db.Where("LOWER(title) LIKE ?", likePattern)
	}

	// Apply groupKey filter
	if groupKey != "" {
		db = db.Where("group_key = ?", groupKey)
	}

	// Count total
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := pagable.GetOffset()
	limit := pagable.GetLimit()
	db = db.Offset(offset).Limit(limit)

	// Order by created_at desc
	db = db.Order("created_at DESC")

	// Execute query
	if err := db.Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

// GetByKey retrieves user guide by key
func (r *UserGuideRepo) GetByKey(ctx context.Context, key string) (*domain.UserGuideEntity, error) {
	var entity domain.UserGuideEntity
	err := r.GetDB(ctx).
		Where("key = ? AND deleted_at IS NULL", key).
		First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// GetAllWithKey lấy toàn bộ user guide có key không rỗng (để sync theo timestamp).
func (r *UserGuideRepo) GetAllWithKey(ctx context.Context) ([]*domain.UserGuideEntity, error) {
	var data []*domain.UserGuideEntity
	err := r.GetDB(ctx).
		Model(&domain.UserGuideEntity{}).
		Where("deleted_at IS NULL AND key IS NOT NULL AND key != ''").
		Order("updated_at DESC").
		Find(&data).Error
	return data, err
}
