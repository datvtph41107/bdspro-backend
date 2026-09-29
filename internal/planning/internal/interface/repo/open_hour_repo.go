package repo

import (
	"context"

	_crud "common/domain/crud"
	"tqd/internal/domain"
	"tqd/internal/dto"
)

// IOpenHourRepo defines repository interface for OpenHour
type IOpenHourRepo interface {
	_crud.ICrudRepo[domain.OpenHour]

	// FindByPOI gets open hours by POI ID
	FindByPOI(ctx context.Context, poiID uint64) ([]domain.OpenHour, error)

	// FindByDayOfWeek gets open hours by day of week
	FindByDayOfWeek(ctx context.Context, dayOfWeek int, poiID *uint64) ([]domain.OpenHour, error)

	// FindByTimeRange finds open hours within time range
	FindByTimeRange(ctx context.Context, poiID uint64, dayOfWeek int, startTime, endTime string) ([]domain.OpenHour, error)

	// CheckTimeOverlap checks if time slot overlaps with existing ones
	CheckTimeOverlap(ctx context.Context, poiID uint64, dayOfWeek int, startTime, endTime string, excludeID *uint64) (bool, error)

	// ListWithFilter lists open hours with filters
	ListWithFilter(ctx context.Context, filter *dto.OpenHourFilter) ([]domain.OpenHour, int64, error)

	// GetGroupedByDay gets open hours grouped by day
	GetGroupedByDay(ctx context.Context, poiID uint64) (map[int][]domain.OpenHour, error)

	// BulkCreate creates multiple open hours
	BulkCreate(ctx context.Context, items []domain.OpenHour) error

	// BulkUpdate updates multiple open hours
	BulkUpdate(ctx context.Context, items []domain.OpenHour) error

	// BulkDelete deletes multiple open hours
	BulkDelete(ctx context.Context, ids []uint64) error
}
