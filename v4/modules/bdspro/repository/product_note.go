package repository

import (
	"context"
	"time"

	"userv4/modules/bdspro/entity"

	"gorm.io/gorm"
)

type ProductNoteRow struct {
	ID        uint64
	ProductID uint64
	Content   string
	AuthorID  uint64
	IsPinned  bool
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

type ProductNoteRepository interface {
	ListByProductID(ctx context.Context, productID uint64, page, size uint32) ([]*ProductNoteRow, int64, error)
	ListMentionsByNoteIDs(ctx context.Context, noteIDs []uint64) ([]*entity.ProductNoteMention, error)
	ListFilesByNoteIDs(ctx context.Context, noteIDs []uint64) ([]*entity.ProductNoteFile, error)
}

type productNotePostgres struct{ db *gorm.DB }

func NewProductNotePostgres(db *gorm.DB) ProductNoteRepository {
	return &productNotePostgres{db: db}
}

func (r *productNotePostgres) ListByProductID(ctx context.Context, productID uint64, page, size uint32) ([]*ProductNoteRow, int64, error) {
	if r.db == nil || productID == 0 {
		return nil, 0, nil
	}
	q := r.db.WithContext(ctx).Table("product_notes pn").
		Where("pn.product_id = ? AND pn.deleted_at IS NULL", productID)

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	limit := int(size)
	if limit <= 0 {
		limit = 20
	}
	if limit > 999 {
		limit = 999
	}
	offset := int(page) * limit

	var rows []*ProductNoteRow
	err := q.Select(`
		pn.id, pn.product_id, pn.content, pn.author_id,
		(pn.pinned_at IS NOT NULL) AS is_pinned,
		pn.created_at, pn.updated_at
	`).
		Order("pn.pinned_at DESC NULLS LAST, pn.created_at DESC").
		Limit(limit).Offset(offset).
		Scan(&rows).Error
	return rows, total, err
}

func (r *productNotePostgres) ListMentionsByNoteIDs(ctx context.Context, noteIDs []uint64) ([]*entity.ProductNoteMention, error) {
	if r.db == nil || len(noteIDs) == 0 {
		return nil, nil
	}
	var list []*entity.ProductNoteMention
	err := r.db.WithContext(ctx).
		Where("note_id IN ? AND deleted_at IS NULL", noteIDs).
		Find(&list).Error
	return list, err
}

func (r *productNotePostgres) ListFilesByNoteIDs(ctx context.Context, noteIDs []uint64) ([]*entity.ProductNoteFile, error) {
	if r.db == nil || len(noteIDs) == 0 {
		return nil, nil
	}
	var list []*entity.ProductNoteFile
	err := r.db.WithContext(ctx).
		Where("note_id IN ? AND deleted_at IS NULL", noteIDs).
		Find(&list).Error
	return list, err
}
