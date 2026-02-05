package outbox

import (
	"context"

	"gorm.io/gorm"
)

type Repo struct {
	DB *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{DB: db}
}

func (r *Repo) SaveTx(ctx context.Context, tx *gorm.DB, ev *OutboxEvent) error {
	if tx == nil {
		return r.DB.WithContext(ctx).Create(ev).Error
	}
	return tx.WithContext(ctx).Create(ev).Error
}

func (r *Repo) FetchPending(ctx context.Context, limit int) ([]OutboxEvent, error) {
	var res []OutboxEvent
	// use FOR UPDATE SKIP LOCKED for concurrency safety if using Postgres
	err := r.DB.WithContext(ctx).Where("processed = FALSE").Order("id asc").Limit(limit).Find(&res).Error
	return res, err
}

func (r *Repo) MarkProcessed(ctx context.Context, id uint64) error {
	return r.DB.WithContext(ctx).Model(&OutboxEvent{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"processed":    true,
			"processed_at": gorm.Expr("now()"),
		}).Error
}

func (r *Repo) IncrementAttempts(ctx context.Context, id uint64) error {
	return r.DB.WithContext(ctx).Model(&OutboxEvent{}).Where("id = ?", id).
		UpdateColumn("attempts", gorm.Expr("attempts + 1")).Error
}
