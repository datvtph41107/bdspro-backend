package repo

import (
	"bdspro/internal/enums"
	"context"
)

type RecordHistoryRepo interface {
	// Create(c context.Context, entity *domain.RecordHistory) error
	// Update(c context.Context, id uint64, entity *domain.RecordHistory) error
	// Delete(c context.Context, id uint64) error
	// GetAll() ([]domain.RecordHistory, error)
	// GetByID(id uint64) (*domain.AssetHistory, error)
	CreateHistory(
		c context.Context,
		action enums.EHistoryAction,
		recordID *uint64,
		parentID *uint64,
		childIDs []uint64,
		date string,
		description string,
		performedBy uint64,
	) error
}
