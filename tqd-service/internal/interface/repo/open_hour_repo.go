package repo

import (
	_crud "common/domain/crud"
	"context"
	"tqd/internal/domain"
)

type IOpenHourRepo interface {
	_crud.ICrudRepo[domain.OpenHour]
	GetGroupedByDayOfWeek(ctx context.Context, poiID *uint64) (map[int][]domain.OpenHour, error)
}
