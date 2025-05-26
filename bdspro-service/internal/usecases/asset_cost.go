package usecases

import (
	"bdspro/internal/domain"
	"bdspro/internal/repo"
	"context"
)

type AssetCostUsecase struct {
	Repo repo.AssetCostRepo
}

func NewAssetCostUc(repo repo.AssetCostRepo) AssetCostUsecase {
	return AssetCostUsecase{
		Repo: repo,
	}
}

func (uc AssetCostUsecase) Create(c context.Context, data *domain.AssetCost) (*domain.AssetCost, error) {
	err := uc.Repo.Create(c, data)
	return data, err
}

func (uc AssetCostUsecase) Update(c context.Context, id uint64, data *domain.AssetCost) (*domain.AssetCost, error) {
	err := uc.Repo.Update(c, id, data)
	return data, err
}

func (uc AssetCostUsecase) Delete(c context.Context, id uint64) error {
	err := uc.Repo.Delete(c, id)
	return err
}

func (uc AssetCostUsecase) GetData(c context.Context, id uint64) ([]domain.AssetCost, error) {
	result, err := uc.Repo.GetData(c)
	return result, err
}

func (uc AssetCostUsecase) GetByID(c context.Context, id uint64) (*domain.AssetCost, error) {
	result, err := uc.Repo.GetByID(c, id)
	return result, err
}
