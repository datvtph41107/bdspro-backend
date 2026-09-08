package usecases

import (
	"bdspro/internal/domain"
	"bdspro/internal/repo"
	"context"
)

type AssetIncomeTypeUsecase struct {
	Repo repo.AssetIncomeTypeRepo
}

func NewAssetIncomeTypeUc(repo repo.AssetIncomeTypeRepo) AssetIncomeTypeUsecase {
	return AssetIncomeTypeUsecase{
		Repo: repo,
	}
}

func (uc AssetIncomeTypeUsecase) Create(c context.Context, data *domain.AssetIncomeType) (*domain.AssetIncomeType, error) {
	err := uc.Repo.Create(c, data)
	return data, err
}

func (uc AssetIncomeTypeUsecase) Update(c context.Context, id uint64, data *domain.AssetIncomeType) (*domain.AssetIncomeType, error) {
	err := uc.Repo.Update(c, id, data)
	return data, err
}

func (uc AssetIncomeTypeUsecase) Delete(c context.Context, id uint64) error {
	err := uc.Repo.Delete(c, id)
	return err
}

func (uc AssetIncomeTypeUsecase) GetData(c context.Context, id uint64) ([]domain.AssetIncomeType, error) {
	result, err := uc.Repo.GetData(c)
	return result, err
}

func (uc AssetIncomeTypeUsecase) GetByID(c context.Context, id uint64) (*domain.AssetIncomeType, error) {
	result, err := uc.Repo.GetByID(c, id)
	return result, err
}
