package repository

import (
	"context"

	"userv4/modules/bdspro/entity"

	"gorm.io/gorm"
)

// DealInternalNoteRepository persists deal_internal_notes.
type DealInternalNoteRepository interface {
	Create(ctx context.Context, note *entity.DealInternalNote) (*entity.DealInternalNote, error)
}

type dealInternalNotePostgres struct {
	db *gorm.DB
}

func NewDealInternalNotePostgres(db *gorm.DB) DealInternalNoteRepository {
	return &dealInternalNotePostgres{db: db}
}

func (r *dealInternalNotePostgres) Create(ctx context.Context, note *entity.DealInternalNote) (*entity.DealInternalNote, error) {
	if r.db == nil || note == nil {
		return nil, gorm.ErrInvalidDB
	}
	row := &entity.DealInternalNote{
		DealID:     note.DealID,
		ActorID:    note.ActorID,
		Content:    note.Content,
		ActionType: note.ActionType,
	}
	if err := r.db.WithContext(ctx).Create(row).Error; err != nil {
		return nil, err
	}
	return row, nil
}
