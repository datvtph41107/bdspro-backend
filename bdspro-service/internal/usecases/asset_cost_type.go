package usecases

import (
	"bdspro/internal/domain"
	"bdspro/internal/repo"
	"context"
)

type AssetCostTypeUsecase struct {
	Repo repo.AssetCostTypeRepo
}

func NewAssetCostTypeUc(repo repo.AssetCostTypeRepo) AssetCostTypeUsecase {
	return AssetCostTypeUsecase{
		Repo: repo,
	}
}

func (uc AssetCostTypeUsecase) Create(c context.Context, data *domain.AssetCostType) (*domain.AssetCostType, error) {
	err := uc.Repo.Create(c, data)
	return data, err
}

func (uc AssetCostTypeUsecase) Update(c context.Context, id uint64, data *domain.AssetCostType) (*domain.AssetCostType, error) {
	err := uc.Repo.Update(c, id, data)
	return data, err
}

func (uc AssetCostTypeUsecase) Delete(c context.Context, id uint64) error {
	err := uc.Repo.Delete(c, id)
	return err
}

func (uc AssetCostTypeUsecase) GetData(c context.Context, id uint64) ([]domain.AssetCostType, error) {
	result, err := uc.Repo.GetData(c)
	return result, err
}

func (uc AssetCostTypeUsecase) GetByID(c context.Context, id uint64) (*domain.AssetCostType, error) {
	result, err := uc.Repo.GetByID(c, id)
	return result, err
}
