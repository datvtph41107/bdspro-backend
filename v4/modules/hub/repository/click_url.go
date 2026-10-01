package repository

import (
	"context"

	"userv4/modules/hub/entity"

	"gorm.io/gorm"
)

// ClickURLRepository persists click-url tracking rows.
type ClickURLRepository interface {
	InsertBatch(ctx context.Context, evs []*entity.ClickURL) error
}

type clickURLPostgres struct{ db *gorm.DB }

func NewClickURLPostgres(db *gorm.DB) ClickURLRepository {
	return &clickURLPostgres{db: db}
}

func (r *clickURLPostgres) InsertBatch(ctx context.Context, evs []*entity.ClickURL) error {
	if len(evs) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).CreateInBatches(evs, 100).Error
}
