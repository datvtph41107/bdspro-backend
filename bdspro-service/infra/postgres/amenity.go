package postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"common/case/crud3"
	"context"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo.AmenityRepo
type PostgreAmenity struct {
	crud3.BaseRepo[domain.Amenity]
}

func NewPostgreAmenity(db *gorm.DB) *PostgreAmenity {
	return &PostgreAmenity{
		BaseRepo: crud3.BaseRepo[domain.Amenity]{DB: db},
	}
}

func (r *PostgreAmenity) InterText(text string, limit int) ([]uint64, error) {
	var ids []uint64

	query := r.DB.
		Model(&domain.Amenity{}).
		Where("LOWER(?) LIKE '%' || LOWER(name) || '%'", text)

	err := query.Limit(limit).Pluck("id", &ids).Error

	return ids, err
}

func (r *PostgreAmenity) InterTextToEntity(text string, limit int) ([]domain.AmenityItem, error) {
	var outputs []domain.AmenityItem

	query := r.DB.
		Model(&domain.Amenity{}).
		Where("LOWER(?) LIKE '%' || LOWER(name) || '%'", text)

	err := query.Limit(limit).Scan(&outputs).Error

	return outputs, err
}

func (r *PostgreAmenity) GetAll(c context.Context, dto *dto.AmenitySearchDTO) ([]domain.AmenityItem, error) {
	var outputs []domain.AmenityItem

	query := r.DB.
		Model(&domain.Amenity{}).
		Where("deleted_at IS NULL")

	err := query.Find(&outputs).Error

	return outputs, err
}
