package postgres

import (
	_db "common/db"
	_dto "common/domain/dto"
	_provider "common/provider"
	"context"
	"strings"
	"tqd/internal/domain"
	"tqd/internal/dto"
	"tqd/internal/interface/repo"
)

type poiCategoryPostgres struct {
	_provider.CrudRepo[domain.PoiCategory]
}

func NewPoiCategoryPostgres(db *_db.TransactionRepo) repo.IPoiCategoryRepo {
	r := &poiCategoryPostgres{}
	r.Init(r, db)
	return r
}

// GetByID override để thêm Preload relations
func (r *poiCategoryPostgres) GetByID(ctx context.Context, id uint64) (*domain.PoiCategory, error) {
	var category domain.PoiCategory
	err := r.GetDB(ctx).
		Where("id = ? and deleted_at is null", id).
		Preload("Parent").
		Preload("Children").
		First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// GetDetail override để thêm Preload relations
func (r *poiCategoryPostgres) GetDetail(ctx context.Context, id uint64) (*domain.PoiCategory, error) {
	return r.GetByID(ctx, id)
}

// GetAll override để thêm Preload relations và sorting
func (r *poiCategoryPostgres) GetAll(ctx context.Context) ([]domain.PoiCategory, error) {
	var categories []domain.PoiCategory
	err := r.GetDB(ctx).
		Where("deleted_at is null").
		Preload("Parent").
		Preload("Children").
		Order("sort_order ASC, created_at DESC").
		Find(&categories).Error
	return categories, err
}

// GetList override để thêm custom filters
func (r *poiCategoryPostgres) GetList(ctx context.Context, pagable _dto.IPagable) ([]domain.PoiCategory, int64, error) {
	var categories []domain.PoiCategory
	var total int64

	query := r.GetDB(ctx).Model(&domain.PoiCategory{}).Where("deleted_at is null")

	// Type assertion để lấy custom filters
	if req, ok := pagable.(*dto.PoiCategoryListRequest); ok {
		// Apply filters
		if req.Code != "" {
			query = query.Where("LOWER(code) LIKE ?", "%"+strings.ToLower(req.Code)+"%")
		}

		if req.Name != "" {
			searchTerm := "%" + strings.ToLower(req.Name) + "%"
			query = query.Where("LOWER(name) LIKE ?", searchTerm)
		}

		if req.ParentID != nil {
			query = query.Where("parent_id = ?", *req.ParentID)
		}

		if req.IsActive != nil {
			query = query.Where("is_active = ?", *req.IsActive)
		}
	}

	// Count total
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Apply pagination and load relations
	err = query.
		Preload("Parent").
		Order("sort_order ASC, created_at DESC").
		Offset(pagable.GetOffset()).
		Limit(pagable.GetLimit()).
		Find(&categories).Error
	if err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

// GetTree lấy tất cả category dưới dạng cấu trúc cha-con
func (r *poiCategoryPostgres) GetTree(ctx context.Context) ([]domain.PoiCategory, error) {
	var categories []domain.PoiCategory
	err := r.GetDB(ctx).
		Where("deleted_at is null AND parent_id IS NULL").
		Preload("Children", "deleted_at is null").
		Order("sort_order ASC, created_at DESC").
		Find(&categories).Error
	return categories, err
}
