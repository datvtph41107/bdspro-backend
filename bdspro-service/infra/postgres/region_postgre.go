package postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"common/case/crud3"
	"context"
	"strings"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo.RegionRepo
type PostgreRegion struct {
	crud3.BaseRepo[domain.Region]
}

func NewPostgreRegion(db *gorm.DB) *PostgreRegion {
	return &PostgreRegion{
		BaseRepo: crud3.BaseRepo[domain.Region]{
			DB: db,
		},
	}
}

func (r *PostgreRegion) Search(c context.Context, dto *dto.RegionRequest) ([]domain.Region, int64, error) {
	var regions []domain.Region
	var total int64

	query := r.DB.Model(&domain.Region{})
	// Preload("Parent").
	// Preload("Parent.Parent")

	if dto.Text != "" {
		text := "%" + strings.ToLower(dto.Text) + "%"
		query = query.Where("LOWER(name) LIKE ?", text)
	}

	if dto.Level != nil {
		query = query.Where("level = ?", *dto.Level)
	}

	if dto.ParentID != nil {
		query = query.Where("parent_id = ?", *dto.ParentID)
	}

	query.Count(&total)

	err := query.
		Debug().
		Offset(dto.GetOffset()).
		Limit(dto.GetLimit()).
		Find(&regions).Error

	return regions, total, err

}

func (r *PostgreRegion) InterText(text string, limit int) ([]uint64, error) {
	var ids []uint64

	query := r.DB.
		Model(&domain.Region{}).
		Where("LOWER(?) LIKE '%' || LOWER(name) || '%'", text)

	err := query.Limit(limit).Pluck("id", &ids).Error

	return ids, err
}

func (r *PostgreRegion) InterTextToItem(text string, limit int) ([]domain.Region, error) {
	var output []domain.Region

	query := r.DB.
		Model(&domain.Region{}).
		Where("LOWER(?) LIKE '%' || LOWER(name) || '%'", text)

	err := query.Limit(limit).Scan(&output).Error

	return output, err
}

func (r *PostgreRegion) GetByID(c context.Context, id uint64) (*domain.Region, error) {
	var region domain.Region

	err := r.DB.WithContext(c).
		First(&region, id).
		Error

	return &region, err
}
