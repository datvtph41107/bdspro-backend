package postgres

import (
	"context"

	"hub/internal/domain"
	"hub/internal/dto"
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

func (r *InteractiveEventPostgres) GetProductStatsView(
	ctx context.Context,
	productID uint64,
	fromTime int64,
) (*dto.ViewStatsProduct, error) {
	var result dto.ViewStatsProduct

	err := r.db.GetDB(ctx).WithContext(ctx).
		Raw(`
			SELECT
				COUNT(*)                      AS view_count,
				COALESCE(SUM(duration), 0)    AS total_duration,
				COALESCE(MAX(start_time), 0)  AS max_event_time
			FROM interactive_events
			WHERE ref_id = ?
				AND event = 'view_product'
				AND start_time > ?
		`, productID, fromTime).
		Scan(&result).Error
	if err != nil {
		return nil, err
	}

	if result.ViewCount == 0 {
		return nil, nil
	}
	return &result, nil
}
