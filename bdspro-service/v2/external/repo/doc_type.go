package postgre_repo

import (
	"bdspro/v2/internal/domain"
	"common/case/crud"

	"gorm.io/gorm"
)

// @bind: bdspro/v2/repo.DocTypeRepo
type GormDocTypeRepo struct {
	crud.BaseRepo[domain.DocType]
}

func NewGormDocTypeRepo(db *gorm.DB) *GormDocTypeRepo {
	return &GormDocTypeRepo{
		BaseRepo: crud.BaseRepo[domain.DocType]{DB: db},
	}
}

func (r *GormDocTypeRepo) InterText(text string, limit int) ([]uint64, error) {
	var ids []uint64

	query := r.DB.
		Model(&domain.DocType{}).
		Where("LOWER(?) LIKE '%' || LOWER(name) || '%'", text)

	err := query.Limit(limit).Pluck("id", &ids).Error

	return ids, err
}

func (r *GormDocTypeRepo) InterTextToItem(text string, limit int) ([]domain.DocTypeItem, error) {
	var ids []domain.DocTypeItem

	query := r.DB.
		Model(&domain.DocType{}).
		Where("LOWER(?) LIKE '%' || LOWER(name) || '%'", text)

	err := query.Limit(limit).Scan(&ids).Error

	return ids, err
}
