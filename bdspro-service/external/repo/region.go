package postgre_repo

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"common/case/crud3"
	"strings"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo.RegionRepo
type GormRegionRepo struct {
	crud3.BaseRepo[domain.Region]
}

func NewGormRegionRepo(db *gorm.DB) *GormRegionRepo {
	return &GormRegionRepo{
		BaseRepo: crud3.BaseRepo[domain.Region]{
			DB: db,
		},
	}
}

func (r *GormRegionRepo) Search(dto *dto.RegionRequest) ([]domain.Region, int64, error) {
	var regions []domain.Region
	var total int64
	if dto.PageSize == 0 {
		dto.PageSize = 20
	}

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

	offset := (dto.Page - 1) * dto.PageSize
	err := query.Offset(offset).Limit(dto.PageSize).Find(&regions).Error

	return regions, total, err

}

func (r *GormRegionRepo) InterText(text string, limit int) ([]uint64, error) {
	var ids []uint64

	query := r.DB.
		Model(&domain.Region{}).
		Where("LOWER(?) LIKE '%' || LOWER(name) || '%'", text)

	err := query.Limit(limit).Pluck("id", &ids).Error

	return ids, err
}

func (r *GormRegionRepo) InterTextToItem(text string, limit int) ([]domain.Region, error) {
	var output []domain.Region

	query := r.DB.
		Model(&domain.Region{}).
		Where("LOWER(?) LIKE '%' || LOWER(name) || '%'", text)

	err := query.Limit(limit).Scan(&output).Error

	return output, err
}
