package usecases

import (
	"bdspro/internal/repo"
)

type AssetHistoryUsecase struct {
	Repo repo.UAssetHistoryRepo
}

func NewAssetSplitHistoryUc(repo repo.UAssetHistoryRepo) AssetHistoryUsecase {
	return AssetHistoryUsecase{
		Repo: repo,
	}
}

// func (uc AssetSplitHistoryUsecase) Create(c context.Context, data *domain.AssetHistory) (*domain.AssetHistory, error) {
// 	err := uc.Repo.Create(c, data)
// 	return data, err
// }

// func (uc AssetSplitHistoryUsecase) Update(c context.Context, id uint64, data *domain.AssetHistory) (*domain.AssetHistory, error) {
// 	err := uc.Repo.Update(c, id, data)
// 	return data, err
// }

// func (uc AssetSplitHistoryUsecase) Delete(c context.Context, id uint64) error {
// 	err := uc.Repo.Delete(c, id)
// 	return err
// }

// func (uc AssetSplitHistoryUsecase) GetData(c context.Context, id uint64) ([]domain.AssetHistory, error) {
// 	result, err := uc.Repo.GetData(c)
// 	return result, err
// }

// func (uc AssetSplitHistoryUsecase) GetByID(c context.Context, id uint64) (*domain.AssetHistory, error) {
// 	result, err := uc.Repo.GetByID(c, id)
// 	return result, err
// }
