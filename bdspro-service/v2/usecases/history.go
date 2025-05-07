package usecases

import (
	"bdspro/v2/internal/domain"
	"bdspro/v2/repo"
	"context"
)

type HistoryUsecase struct {
	Repo repo.AssetHistoryRepo
}

func NewHistoryUsecase(repo repo.AssetHistoryRepo) *HistoryUsecase {
	return &HistoryUsecase{Repo: repo}
}

func (uc *HistoryUsecase) Create(ctx context.Context, history *domain.AssetSplitMergeHistory) (*domain.AssetSplitMergeHistory, error) {
	return uc.Repo.Create(ctx, history)
}

func (uc *HistoryUsecase) GetByAssetID(ctx context.Context, assetID uint64) ([]domain.AssetSplitMergeHistory, error) {
	return uc.Repo.GetByAssetID(ctx, assetID)
}
