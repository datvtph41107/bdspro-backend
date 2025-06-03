package postgre_repo

import (
	"bdspro/internal/domain"
	"common/case/crud"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo.PropertyTypeRepo
type GormPropertyTypeRepo struct {
	crud.BaseRepo[domain.PropertyType]
}

func NewGormPropertyTypeRepo(db *gorm.DB) *GormPropertyTypeRepo {
	return &GormPropertyTypeRepo{
		BaseRepo: crud.BaseRepo[domain.PropertyType]{DB: db},
	}
}

func (r *GormPropertyTypeRepo) InterText(text string, limit int) ([]uint64, error) {
	var ids []uint64

	query := r.DB.
		Model(&domain.PropertyType{}).
		Where("LOWER(?) LIKE '%' || LOWER(name) || '%'", text)

	err := query.Limit(limit).Pluck("id", &ids).Error

	return ids, err
}

func (r *GormPropertyTypeRepo) InterTextToItem(text string, limit int) ([]domain.PropertyTypeItem, error) {
	var ids []domain.PropertyTypeItem

	query := r.DB.
		Model(&domain.PropertyType{}).
		Where("LOWER(?) LIKE '%' || LOWER(name) || '%'", text)

	err := query.Limit(limit).Scan(&ids).Error

	return ids, err
}
