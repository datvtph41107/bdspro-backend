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
	return r.DB.WithContext(ctx).Create(note).Error
}

func (r *PostgreProductNoteRepo) GetByID(ctx context.Context, id uint64) (*domain.ProductNote, error) {
	var note domain.ProductNote
	err := r.DB.WithContext(ctx).
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

func (r *PostgreProductNoteRepo) UpdateContent(
	ctx context.Context,
	noteID uint64,
	content string,
) error {
	return r.DB.WithContext(ctx).
		Model(&domain.ProductNote{}).
		Where("id = ?", noteID).
		Update("content", content).Error
}

func (r *PostgreProductNoteRepo) UpdatePinnedAt(
	ctx context.Context,
	noteID uint64,
	pinnedAt *time.Time,
) error {
	return r.DB.WithContext(ctx).
		Model(&domain.ProductNote{}).
		Where("id = ?", noteID).
		Update("pinned_at", pinnedAt).Error
}

func (r *PostgreProductNoteRepo) DeleteNote(
	ctx context.Context,
	noteID uint64,
) error {
	return r.DB.WithContext(ctx).
		Delete(&domain.ProductNote{}, noteID).Error
}

func (r *PostgreProductNoteRepo) CreateMention(
	ctx context.Context,
	mentions []domain.ProductNoteMention,
) error {
	if len(mentions) == 0 {
		return nil
	}
	return r.DB.WithContext(ctx).Create(&mentions).Error
}

func (r *PostgreProductNoteRepo) DeleteMentionByNoteID(
	ctx context.Context,
	noteID uint64,
) error {
	return r.DB.WithContext(ctx).
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
	return r.DB.WithContext(ctx).Create(&files).Error
}

func (r *PostgreProductNoteRepo) DeleteFilesByNoteID(
	ctx context.Context,
	noteID uint64,
) error {
	return r.DB.WithContext(ctx).
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

	err := r.DB.WithContext(ctx).
		Where("note_id IN ?", noteIDs).
		Find(&files).Error

	return files, err
}

func (r *PostgreProductNoteRepo) GetListByProductID(
	ctx context.Context,
	filter *dto.GetProductNotesFilter,
) ([]*dto.ProductNoteItemResponse, int64, error) {

	var (
		result []*dto.ProductNoteItemResponse
		total  int64
	)

	baseQuery := GetDB(ctx, r.DB).
		Table("product_notes pn").
		Joins("LEFT JOIN product_note_mentions pm ON pm.note_id = pn.id").
		Joins("LEFT JOIN product_note_files pf ON pf.note_id = pn.id").
		Where("pn.deleted_at IS NULL AND pn.product_id = ?", filter.ProductID)

	if filter.OnlyPinned != nil {
		if *filter.OnlyPinned {
			baseQuery = baseQuery.Where("pn.pinned_at IS NOT NULL")
		} else {
			baseQuery = baseQuery.Where("pn.pinned_at IS NULL")
		}
	}

	if filter.AuthorID != nil {
		baseQuery = baseQuery.Where("pn.author_id = ?", *filter.AuthorID)
	}

	if filter.Query != nil && *filter.Query != "" {
		baseQuery = baseQuery.Where("pn.content ILIKE ?", "%"+*filter.Query+"%")
	}

	countSubQuery := baseQuery.
		Select("pn.id").
		Group("pn.id")

	if err := r.DB.
		Table("(?) AS t", countSubQuery).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := baseQuery.
		Select(`
			pn.id,
			pn.content,
			pn.pinned_at,
			pn.author_id,
			pn.created_at,

			COALESCE(
				ARRAY_AGG(DISTINCT pm.user_id)
				FILTER (WHERE pm.user_id IS NOT NULL),
				'{}'
			) AS mention_ids,

			COALESCE(
				ARRAY_AGG(DISTINCT pf.id)
				FILTER (WHERE pf.id IS NOT NULL),
				'{}'
			) AS file_ids
		`).
		Group("pn.id").
		Order("pn.pinned_at DESC NULLS LAST, pn.created_at DESC").
		Limit(filter.GetLimit()).
		Offset(filter.GetOffset()).
		Scan(&result).Error; err != nil {
		return nil, 0, err
	}

	return result, total, nil
}
