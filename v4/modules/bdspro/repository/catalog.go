package repository

import (
	"context"

	_dto "common/domain/dto"
	"userv4/modules/bdspro/entity"

	"gorm.io/gorm"
)

// CatalogRepository loads project / project-build catalog rows from the v4 DSN.
type CatalogRepository interface {
	SearchProjects(ctx context.Context, pagable _dto.Pagable, text string) ([]entity.ProjectItem, int64, error)
	SearchProjectBuilds(ctx context.Context, pagable _dto.Pagable, text string) ([]entity.ProjectBuild, int64, error)
}

type catalogPostgres struct {
	db *gorm.DB
}

func NewCatalogPostgres(db *gorm.DB) CatalogRepository {
	return &catalogPostgres{db: db}
}

func (r *catalogPostgres) SearchProjects(ctx context.Context, pagable _dto.Pagable, text string) ([]entity.ProjectItem, int64, error) {
	if r.db == nil {
		return nil, 0, nil
	}
	q := r.db.WithContext(ctx).Model(&entity.ProjectItem{}).Where("deleted_at IS NULL")
	if text != "" {
		q = q.Where("LOWER(name) LIKE '%' || LOWER(?) || '%'", text)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []entity.ProjectItem
	err := q.Order("name ASC").
		Offset(pagable.GetOffset()).
		Limit(pagable.GetLimit()).
		Find(&items).Error
	return items, total, err
}

func (r *catalogPostgres) SearchProjectBuilds(ctx context.Context, pagable _dto.Pagable, text string) ([]entity.ProjectBuild, int64, error) {
	if r.db == nil {
		return nil, 0, nil
	}
	q := r.db.WithContext(ctx).Model(&entity.ProjectBuild{}).Where("deleted_at IS NULL")
	if text != "" {
		q = q.Where("name ILIKE ?", "%"+text+"%")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []entity.ProjectBuild
	err := q.Order("name ASC").
		Offset(pagable.GetOffset()).
		Limit(pagable.GetLimit()).
		Find(&items).Error
	return items, total, err
}
