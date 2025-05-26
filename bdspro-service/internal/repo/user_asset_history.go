package repo

import (
	"bdspro/internal/enums"
	"context"
)

type UAssetHistoryRepo interface {
	CreateHistory(
		c context.Context,
		action enums.AssetHistory,
		assetID *uint64,
		// parentID uint64,
		BeforeData *string,
		AfterData *string,
	) error
}

// Create(c context.Context, entity *domain.AssetHistory) error
// Update(c context.Context, id uint64, entity *domain.AssetHistory) error
// Delete(c context.Context, id uint64) error
// GetData(c context.Context) ([]domain.AssetHistory, error)
// GetByID(c context.Context, id uint64) (*domain.AssetHistory, error)
