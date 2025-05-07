package repo

import (
	"bdspro/v2/internal/domain"
	"context"
)

type AssetHistoryRepo interface {
	Create(ctx context.Context, history *domain.AssetSplitMergeHistory) (*domain.AssetSplitMergeHistory, error)
	GetByAssetID(ctx context.Context, assetID uint64) ([]domain.AssetSplitMergeHistory, error)
}
