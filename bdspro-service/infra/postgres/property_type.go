package postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	_dto "common/domain/dto"
	"context"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo.PropertyTypeRepo
type PostgrePropertyType struct {
	// shared.PostgreCrud[domain.PropertyType, *dto.PropertyTypeSearchDTO]
	DB *gorm.DB
}

func NewPostgrePropertyType(db *gorm.DB) *PostgrePropertyType {
	repo := &PostgrePropertyType{
		DB: db,
		// PostgreCrud: shared.PostgreCrud[domain.PropertyType, *dto.PropertyTypeSearchDTO]{
		// 	DB: db,
		// },
	}
	// repo.Repo = repo
	return repo
}

func (r *PostgrePropertyType) InterText(text string, limit int) ([]uint64, error) {
	var ids []uint64

	query := r.DB.
		Model(&domain.PropertyType{}).
		Where("LOWER(?) LIKE '%' || LOWER(name) || '%'", text)

	err := query.Limit(limit).Pluck("id", &ids).Error

	return ids, err
}

func (r *PostgrePropertyType) InterTextToItem(text string, limit int) (*_dto.ItemDTO, error) {
	var item domain.PropertyTypeItem

	query := r.DB.
		Model(&domain.PropertyType{}).
		Where("LOWER(?) LIKE '%' || LOWER(name) || '%'", text)

	err := query.Limit(limit).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &_dto.ItemDTO{
		ID:   item.ID,
		Name: item.Name,
	}, nil
}

func (r *PostgrePropertyType) GetAllItem(c context.Context, dto *dto.PropertyTypeSearchDTO) ([]domain.PropertyType, error) {
	var items []domain.PropertyType

	query := r.DB.
		Model(&domain.PropertyType{}).
		Where("deleted_at IS NULL")

	err := query.Find(&items).Error

	return items, err
}
