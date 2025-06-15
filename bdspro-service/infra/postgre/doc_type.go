package postgre

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"common/case/crud"
	"context"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo.DocTypeRepo
type PostgreDocType struct {
	crud.BaseRepo[domain.DocType]
}

func NewPostgreDocType(db *gorm.DB) *PostgreDocType {
	return &PostgreDocType{
		BaseRepo: crud.BaseRepo[domain.DocType]{DB: db},
	}
}

func (r *PostgreDocType) GetAllItem(c context.Context) ([]domain.DocTypeItem, error) {
	var items []domain.DocTypeItem

	query := r.DB.
		Model(&domain.DocType{}).
		Where("deleted_at IS NULL")

	err := query.Find(&items).Error

	return items, err
}

func (r *PostgreDocType) SearchItem(c context.Context, dto *dto.DocTypeSearchDTO) ([]domain.DocTypeItem, int64, error) {
	// var items []domain.DocTypeItem
	// var total int64

	// query := r.DB.Model(&domain.DocType{})

	// if dto.Text != "" {
	// 	text := "%" + strings.ToLower(dto.Text) + "%"
	// 	query = query.Where("LOWER(name) LIKE ?", text)
	// }

	// query.Count(&total)
	return nil, 0, nil
}

func (r *PostgreDocType) InterText(text string, limit int) ([]uint64, error) {
	var ids []uint64

	query := r.DB.
		Model(&domain.DocType{}).
		Where("LOWER(?) LIKE '%' || LOWER(name) || '%'", text)

	err := query.Limit(limit).Pluck("id", &ids).Error

	return ids, err
}

func (r *PostgreDocType) InterTextToItem(text string, limit int) ([]domain.DocTypeItem, error) {
	var ids []domain.DocTypeItem

	query := r.DB.
		Model(&domain.DocType{}).
		Where("LOWER(?) LIKE '%' || LOWER(name) || '%'", text)

	err := query.Limit(limit).Scan(&ids).Error

	return ids, err
}
