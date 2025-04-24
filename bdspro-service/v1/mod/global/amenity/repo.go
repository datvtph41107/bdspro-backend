package amenity

import (
	"bdspro/v1/models"
	"common/case/crud2"

	"gorm.io/gorm"
)

type AmenityRepo struct {
	crud2.BaseRepo[models.Amenity]
}

func NewAmenityRepo(db *gorm.DB) *AmenityRepo {
	return &AmenityRepo{
		BaseRepo: crud2.BaseRepo[models.Amenity]{DB: db},
	}
}

func (r *AmenityRepo) InterText(text string, limit int) ([]uint64, error) {
	var ids []uint64

	query := r.DB.
		Model(&models.Amenity{}).
		Where("LOWER(?) LIKE '%' || LOWER(name) || '%'", text)

	err := query.Limit(limit).Pluck("id", &ids).Error

	return ids, err
}

func (r *AmenityRepo) InterTextToEntity(text string, limit int) ([]models.AmenityItem, error) {
	var outputs []models.AmenityItem

	query := r.DB.
		Model(&models.Amenity{}).
		Where("LOWER(?) LIKE '%' || LOWER(name) || '%'", text)

	err := query.Limit(limit).Scan(&outputs).Error

	return outputs, err
}
