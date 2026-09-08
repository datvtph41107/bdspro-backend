package postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/repo"
	"context"
	"time"

	"gorm.io/gorm"
)

type PostgreProductNoteRepo struct {
	DB *gorm.DB
}

func NewPostgreProductNoteRepo(db *gorm.DB) repo.ProductNoteRepository {
	return &PostgreProductNoteRepo{
		DB: db,
	}
}

func (r *PostgreProductNoteRepo) CreateNote(
	ctx context.Context,
	note *domain.ProductNote,
) error {
	return GetDB(ctx, r.DB).Create(note).Error
}

func (r *PostgreProductNoteRepo) UpdateContent(
	ctx context.Context,
	noteID uint64,
	content string,
) error {
	return GetDB(ctx, r.DB).
		Model(&domain.ProductNote{}).
		Where("id = ? AND deleted_at IS NULL", noteID).
		Update("content", content).Error
}

func (r *PostgreProductNoteRepo) UpdatePinnedAt(
	ctx context.Context,
	noteID uint64,
	pinnedAt *time.Time,
) error {
	return GetDB(ctx, r.DB).
		Model(&domain.ProductNote{}).
		Where("id = ? AND deleted_at IS NULL", noteID).
		Update("pinned_at", pinnedAt).Error
}

func (r *PostgreProductNoteRepo) GetByID(
	ctx context.Context,
	id uint64,
) (*domain.ProductNote, error) {
	var note domain.ProductNote

	err := GetDB(ctx, r.DB).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&note).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &note, nil
}

func (r *PostgreProductNoteRepo) ListMentionsByNoteIDs(
	ctx context.Context,
	noteIDs []uint64,
) ([]*dto.ProductNoteMentionItem, error) {
	var mentions []*dto.ProductNoteMentionItem
	if len(noteIDs) == 0 {
		return mentions, nil
	}

	err := GetDB(ctx, r.DB).
		Table("product_note_mentions").
		Where("note_id IN ? AND deleted_at IS NULL", noteIDs).
		Find(&mentions).Error

	return mentions, err
}

func (r *PostgreProductNoteRepo) GetListByProductID(
	ctx context.Context,
	filter *dto.GetProductNotesFilter,
) ([]*dto.ProductNoteRow, int64, error) {
	var rows []*dto.ProductNoteRow
	var total int64

	baseDB := GetDB(ctx, r.DB).
		Table("product_notes pn").
		Where("pn.product_id = ?", filter.ProductID).
		Where("pn.deleted_at IS NULL")

	if err := baseDB.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := baseDB.
		Select(`
			pn.id,
			pn.product_id,
			pn.content,
			pn.author_id,
			(pn.pinned_at IS NOT NULL) AS is_pinned,
			pn.created_at,
			pn.updated_at AS product_note_updated_at
		`).
		Order("pn.pinned_at DESC NULLS LAST, pn.created_at DESC").
		Limit(filter.GetLimit()).
		Offset(filter.GetOffset()).
		Scan(&rows).Error

	return rows, total, err
}

func (r *PostgreProductNoteRepo) DeleteNote(
	ctx context.Context,
	noteID uint64,
) error {
	return GetDB(ctx, r.DB).
		Where("id = ?", noteID).
		Delete(&domain.ProductNote{}).Error
}

func (r *PostgreProductNoteRepo) CreateMention(
	ctx context.Context,
	mentions []domain.ProductNoteMention,
) error {
	if len(mentions) == 0 {
		return nil
	}
	return GetDB(ctx, r.DB).Create(&mentions).Error
}

func (r *PostgreProductNoteRepo) DeleteMentionByNoteID(
	ctx context.Context,
	noteID uint64,
) error {
	return GetDB(ctx, r.DB).
		Where("note_id = ?", noteID).
		Delete(&domain.ProductNoteMention{}).Error
}

func (r *PostgreProductNoteRepo) CreateFiles(
	ctx context.Context,
	files []domain.ProductNoteFile,
) error {
	if len(files) == 0 {
		return nil
	}
	return GetDB(ctx, r.DB).Create(&files).Error
}

func (r *PostgreProductNoteRepo) DeleteFilesByNoteID(
	ctx context.Context,
	noteID uint64,
) error {
	return GetDB(ctx, r.DB).
		Where("note_id = ?", noteID).
		Delete(&domain.ProductNoteFile{}).Error
}

func (r *PostgreProductNoteRepo) ListByNoteIDs(
	ctx context.Context,
	noteIDs []uint64,
) ([]*domain.ProductNoteFile, error) {
	var files []*domain.ProductNoteFile
	if len(noteIDs) == 0 {
		return files, nil
	}

	err := GetDB(ctx, r.DB).
		Where("note_id IN ?", noteIDs).
		Where("deleted_at IS NULL").
		Find(&files).Error

	return files, err
}
