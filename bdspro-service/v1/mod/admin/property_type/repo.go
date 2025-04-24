package a_property_type

import (
	"bdspro/v1/models"
	"common/case/crud"

	"gorm.io/gorm"
)

type PropertyTypeRepo struct {
	crud.BaseRepo[models.PropertyType]
}

func NewPropertyTypeRepo(db *gorm.DB) *PropertyTypeRepo {
	return &PropertyTypeRepo{
		BaseRepo: crud.BaseRepo[models.PropertyType]{DB: db},
	}
}

func (r *PropertyTypeRepo) InterText(text string, limit int) ([]uint64, error) {
	var ids []uint64

	query := r.DB.
		Model(&models.PropertyType{}).
		Where("LOWER(?) LIKE '%' || LOWER(name) || '%'", text)

	err := query.Limit(limit).Pluck("id", &ids).Error

	return ids, err
}

func (r *PropertyTypeRepo) InterTextToItem(text string, limit int) ([]models.PropertyTypeItem, error) {
	var ids []models.PropertyTypeItem

	query := r.DB.
		Model(&models.PropertyType{}).
		Where("LOWER(?) LIKE '%' || LOWER(name) || '%'", text)

	err := query.Limit(limit).Scan(&ids).Error

	return ids, err
}
