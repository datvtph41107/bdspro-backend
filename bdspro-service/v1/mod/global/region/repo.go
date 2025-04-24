package g_region

import (
	"bdspro/v1/data"
	"bdspro/v1/models"
	"common/case/crud2"
	"strings"

	"gorm.io/gorm"
)

type RegionRepo struct {
	crud2.BaseRepo[models.Region]
}

func NewRegionRepo(db *gorm.DB) *RegionRepo {
	return &RegionRepo{
		BaseRepo: crud2.BaseRepo[models.Region]{
			DB: db,
		},
	}
}

func (r *RegionRepo) Search(dto *data.RegionRequest) ([]models.Region, int64, error) {
	var regions []models.Region
	var total int64
	if dto.PageSize == 0 {
		dto.PageSize = 20
	}

	query := r.DB.Model(&models.Region{})
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

func (r *RegionRepo) InterText(text string, limit int) ([]uint64, error) {
	var ids []uint64

	query := r.DB.
		Model(&models.Region{}).
		Where("LOWER(?) LIKE '%' || LOWER(name) || '%'", text)

	err := query.Limit(limit).Pluck("id", &ids).Error

	return ids, err
}

func (r *RegionRepo) InterTextToItem(text string, limit int) ([]models.Region, error) {
	var output []models.Region

	query := r.DB.
		Model(&models.Region{}).
		Where("LOWER(?) LIKE '%' || LOWER(name) || '%'", text)

	err := query.Limit(limit).Scan(&output).Error

	return output, err
}
