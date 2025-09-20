package postgre

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"common/case/crud"
	_db "common/db"
	"context"
)

// @bind: bdspro/internal/repo.DocTypeRepo
type PostgreDocType struct {
	*crud.CrudRepo[domain.DocType]
}

func NewPostgreDocType(db *_db.TransactionRepo) *PostgreDocType {
	x := &PostgreDocType{
		CrudRepo: &crud.CrudRepo[domain.DocType]{},
	}
	x.CrudRepo.Init(x, db)
	return x
}

func (r *PostgreDocType) GetAllItem(c context.Context) ([]domain.DocTypeItem, error) {
	var items []domain.DocTypeItem

	query := r.GetDB(c).
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

func (r *PostgreDocType) InterText(c context.Context, text string, limit int) ([]uint64, error) {
	var ids []uint64

	query := r.GetDB(c).
		Model(&domain.DocType{}).
		Where("LOWER(?) LIKE '%' || LOWER(name) || '%'", text)

	err := query.Limit(limit).Pluck("id", &ids).Error

	return ids, err
}

func (r *PostgreDocType) InterTextToItem(c context.Context, text string, limit int) ([]domain.DocTypeItem, error) {
	var ids []domain.DocTypeItem

	query := r.GetDB(c).
		Model(&domain.DocType{}).
		Where("LOWER(?) LIKE '%' || LOWER(name) || '%'", text)

	err := query.Limit(limit).Scan(&ids).Error

	return ids, err
}
