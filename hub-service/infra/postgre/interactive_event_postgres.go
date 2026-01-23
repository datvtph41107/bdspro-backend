package postgres

import (
	"context"

	"hub/internal/domain"
	"hub/internal/repo"

	_db "common/db"
)

type InteractiveEventPostgres struct {
	db *_db.TransactionRepo
}

func NewInteractiveEventRepo(db *_db.TransactionRepo) repo.IInteractiveEventRepo {
	return &InteractiveEventPostgres{
		db: db,
	}
}

func (r *InteractiveEventPostgres) Insert(ctx context.Context, event *domain.InteractiveEvent) error {
	return r.db.GetDB(ctx).WithContext(ctx).Create(event).Error
}

func (r *InteractiveEventPostgres) InsertBatch(ctx context.Context, events []*domain.InteractiveEvent) error {
	if len(events) == 0 {
		return nil
	}
	return r.db.GetDB(ctx).WithContext(ctx).CreateInBatches(events, 100).Error
}
