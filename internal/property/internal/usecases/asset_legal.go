package usecases

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/repo"
	"context"
)

type AssetLegalUsecase struct {
	Repo repo.AssetLegalRepo
}

func NewAssetLegalUc(repo repo.AssetLegalRepo) *AssetLegalUsecase {
	return &AssetLegalUsecase{
		Repo: repo,
	}
}

func (uc AssetLegalUsecase) Create(c context.Context, data *domain.AssetLegal) (*domain.AssetLegal, error) {
	err := uc.Repo.Create(c, data)
	return data, err
}

func (uc AssetLegalUsecase) Update(c context.Context, id uint64, data *domain.AssetLegal) (*domain.AssetLegal, error) {
	err := uc.Repo.Update(c, id, data)
	return data, err
}

func (uc AssetLegalUsecase) Delete(c context.Context, id uint64) error {
	err := uc.Repo.Delete(c, id)
	return err
}

func (uc AssetLegalUsecase) GetData(c context.Context, dto dto.AssetLegalGetDTO) ([]domain.AssetLegal, int64, error) {
	result, total, err := uc.Repo.GetData(c, dto)
	return result, total, err
}

func (uc AssetLegalUsecase) GetByID(c context.Context, id uint64) (*domain.AssetLegal, error) {
	result, err := uc.Repo.GetByID(c, id)
	return result, err
}
