package a_doc_type

import (
	"bdspro/v1/models"
	"common/case/crud"

	"gorm.io/gorm"
)

type DocTypeRepo struct {
	crud.BaseRepo[models.DocType]
}

func NewDocTypeRepo(db *gorm.DB) *DocTypeRepo {
	return &DocTypeRepo{
		BaseRepo: crud.BaseRepo[models.DocType]{DB: db},
	}
}

func (r *DocTypeRepo) InterText(text string, limit int) ([]uint64, error) {
	var ids []uint64

	query := r.DB.
		Model(&models.DocType{}).
		Where("LOWER(?) LIKE '%' || LOWER(name) || '%'", text)

	err := query.Limit(limit).Pluck("id", &ids).Error

	return ids, err
}

func (r *DocTypeRepo) InterTextToItem(text string, limit int) ([]models.DocTypeItem, error) {
	var ids []models.DocTypeItem

	query := r.DB.
		Model(&models.DocType{}).
		Where("LOWER(?) LIKE '%' || LOWER(name) || '%'", text)

	err := query.Limit(limit).Scan(&ids).Error

	return ids, err
}
