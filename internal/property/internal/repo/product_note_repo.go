package repo

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"context"
	"time"
)

type ProductNoteRepository interface {
	//note
	CreateNote(ctx context.Context, note *domain.ProductNote) error
	GetByID(ctx context.Context, id uint64) (*domain.ProductNote, error)
	GetListByProductID(ctx context.Context, filter *dto.GetProductNotesFilter) ([]*dto.ProductNoteRow, int64, error)
	UpdateContent(ctx context.Context, noteID uint64, content string) error
	UpdatePinnedAt(ctx context.Context, noteID uint64, pinnedAt *time.Time) error
	DeleteNote(ctx context.Context, noteID uint64) error

	//mention
	CreateMention(ctx context.Context, mentions []domain.ProductNoteMention) error
	DeleteMentionByNoteID(ctx context.Context, noteID uint64) error

	//attach
	CreateFiles(ctx context.Context, files []domain.ProductNoteFile) error
	ListByNoteIDs(ctx context.Context, noteIDs []uint64) ([]*domain.ProductNoteFile, error)
	DeleteFilesByNoteID(ctx context.Context, noteID uint64) error

	ListMentionsByNoteIDs(ctx context.Context, noteIDs []uint64) ([]*dto.ProductNoteMentionItem, error)
}
