package postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/repo"
	"context"

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

func (r *PostgreProductNoteRepo) GetByID(
	ctx context.Context,
	id uint64,
) (*domain.ProductNote, error) {
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

func (r *PostgreProductNoteRepo) UpdatePin(
	ctx context.Context,
	noteID uint64,
	isPinned bool,
) error {
	return r.DB.WithContext(ctx).
		Model(&domain.ProductNote{}).
		Where("id = ?", noteID).
		Update("is_pinned", isPinned).Error
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

func (r *PostgreProductNoteRepo) AttachFilesToNote(
	ctx context.Context,
	noteID uint64,
	files []dto.ProductNoteFileItem,
) error {

	if len(files) == 0 {
		return nil
	}

	records := make([]*domain.ProductNoteFile, 0, len(files))

	for _, f := range files {
		records = append(records, &domain.ProductNoteFile{
			NoteID:   noteID,
			FileName: f.FileName,
			FileURL:  f.URL,
			FileType: f.Type,
			FileSize: f.Size,
		})
	}

	return r.DB.WithContext(ctx).
		Create(&records).Error
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
) ([]*dto.ProductNoteItemRespone, int64, error) {
	var result []*dto.ProductNoteItemRespone
	db := r.DB.WithContext(ctx).
		Table("product_notes pn").
		Select(`
			pn.id,
			pn.content,
			pn.is_pinned,
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
		Joins("LEFT JOIN product_note_mentions pm ON pm.note_id = pn.id").
		Joins("LEFT JOIN product_note_files pf ON pf.note_id = pn.id").
		Where("pn.product_id = ?", filter.ProductID)

	if filter.OnlyPinned != nil {
		db = db.Where("pn.is_pinned = ?", *filter.OnlyPinned)
	}

	if filter.AuthorID != nil {
		db = db.Where("pn.author_id = ?", *filter.AuthorID)
	}

	if filter.Query != nil && *filter.Query != "" {
		db = db.Where("pn.content ILIKE ?", "%"+*filter.Query+"%")
	}

	var total int64
	if err := db.
		Distinct("pn.id").
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.
		Group("pn.id").
		Order("pn.is_pinned DESC, pn.created_at DESC").
		Limit(filter.GetLimit()).
		Offset(filter.GetOffset()).
		Scan(&result).Error; err != nil {
		return nil, 0, err
	}

	return result, total, nil
}
