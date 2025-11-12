package postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"common/case/crud3"
	_dto "common/domain/dto"
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

func (r *PostgreAmenity) InterTextToItem(text string, limit int) ([]_dto.ItemDTO, error) {
	var outputs []domain.Amenity

	query := r.DB.
		Model(&domain.Amenity{}).
		Where("LOWER(?) LIKE '%' || LOWER(name) || '%'", text)

	err := query.Limit(limit).Scan(&outputs).Error
	if err != nil {
		return nil, err
	}

	items := make([]_dto.ItemDTO, 0)
	for _, output := range outputs {
		items = append(items, _dto.ItemDTO{
			ID:   output.ID,
			Name: output.Name,
		})
	}

	return items, err
}

func (r *PostgreAmenity) GetAll(c context.Context, dto *dto.AmenitySearchDTO) ([]domain.AmenityItem, error) {
	var outputs []domain.AmenityItem

	query := r.DB.
		Model(&domain.Amenity{}).
		Where("deleted_at IS NULL")

	err := query.Find(&outputs).Error

	return outputs, err
}
