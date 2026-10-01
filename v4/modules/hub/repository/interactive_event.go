package repository

import (
	"context"

	"userv4/modules/hub/entity"

	"gorm.io/gorm"
)

// InteractiveEventRepository persists interactive event tracking rows.
type InteractiveEventRepository interface {
	Insert(ctx context.Context, ev *entity.InteractiveEvent) error
	InsertBatch(ctx context.Context, evs []*entity.InteractiveEvent) error
}

type interactiveEventPostgres struct{ db *gorm.DB }

func NewInteractiveEventPostgres(db *gorm.DB) InteractiveEventRepository {
	return &interactiveEventPostgres{db: db}
}

func (r *interactiveEventPostgres) Insert(ctx context.Context, ev *entity.InteractiveEvent) error {
	return r.db.WithContext(ctx).Create(ev).Error
}

func (r *interactiveEventPostgres) InsertBatch(ctx context.Context, evs []*entity.InteractiveEvent) error {
	if len(evs) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).CreateInBatches(evs, 100).Error
}
