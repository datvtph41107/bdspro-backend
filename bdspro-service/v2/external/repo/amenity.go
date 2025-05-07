package postgre_repo

import (
	"bdspro/v2/internal/domain"
	"common/case/crud3"

	"gorm.io/gorm"
)

// @bind: bdspro/v2/repo.AmenityRepo
type GormAmenityRepo struct {
	crud3.BaseRepo[domain.Amenity]
}

func NewGormAmenityRepo(db *gorm.DB) *GormAmenityRepo {
	return &GormAmenityRepo{
		BaseRepo: crud3.BaseRepo[domain.Amenity]{DB: db},
	}
}

func (r *GormAmenityRepo) InterText(text string, limit int) ([]uint64, error) {
	var ids []uint64

	query := r.DB.
		Model(&domain.Amenity{}).
		Where("LOWER(?) LIKE '%' || LOWER(name) || '%'", text)

	err := query.Limit(limit).Pluck("id", &ids).Error

	return ids, err
}

func (r *GormAmenityRepo) InterTextToEntity(text string, limit int) ([]domain.AmenityItem, error) {
	var outputs []domain.AmenityItem

	query := r.DB.
		Model(&domain.Amenity{}).
		Where("LOWER(?) LIKE '%' || LOWER(name) || '%'", text)

	err := query.Limit(limit).Scan(&outputs).Error

	return outputs, err
}
