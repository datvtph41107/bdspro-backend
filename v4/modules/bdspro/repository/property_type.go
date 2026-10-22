package repository

import (
	"context"

	"mdn4/modules/bdspro/entity"

	"gorm.io/gorm"
)

// PropertyTypeRepository loads property type catalog rows.
type PropertyTypeRepository interface {
	List(ctx context.Context, text string) ([]entity.PropertyType, error)
}

type propertyTypePostgres struct {
	db *gorm.DB
}

func NewPropertyTypePostgres(db *gorm.DB) PropertyTypeRepository {
	return &propertyTypePostgres{db: db}
}

func (r *propertyTypePostgres) List(ctx context.Context, text string) ([]entity.PropertyType, error) {
	if r.db == nil {
		return nil, nil
	}
	q := r.db.WithContext(ctx).Model(&entity.PropertyType{}).Where("deleted_at IS NULL")
	if text != "" {
		q = q.Where("LOWER(name) LIKE '%' || LOWER(?) || '%'", text)
	}
	var items []entity.PropertyType
	err := q.Order("name ASC").Find(&items).Error
	return items, err
}
