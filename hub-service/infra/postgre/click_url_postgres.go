package postgres

import (
	"context"

	_db "common/db"
	"hub/internal/domain"
	"hub/internal/repo"
)

type ClickURLPostgres struct {
	db *_db.TransactionRepo
}

// @bind: internal/repo/IClickURLRepo
func NewClickURLRepo(db *_db.TransactionRepo) repo.IClickURLRepo {
	return &ClickURLPostgres{
		db: db,
	}
}

func (r *ClickURLPostgres) Insert(ctx context.Context, event *domain.ClickURL) error {
	return r.db.GetDB(ctx).WithContext(ctx).Create(event).Error
}

func (r *ClickURLPostgres) InsertBatch(ctx context.Context, events []*domain.ClickURL) error {
	if len(events) == 0 {
		return nil
	}
	return r.db.GetDB(ctx).WithContext(ctx).CreateInBatches(events, 100).Error
}
