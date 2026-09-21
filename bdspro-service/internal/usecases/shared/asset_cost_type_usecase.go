package shared_usecase

import (
	"bdspro/internal"
	"bdspro/internal/common"
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"bdspro/internal/repo"
	_errors "common/errors"
	_utils "common/utils"
	"context"
)

type AssetCostTypeUsecase struct {
	common.OwnerUsecase[domain.AssetCostType, *dto.AssetCostTypeSearchDTO]
	repo repo.AssetCostTypeRepo
}

func NewAssetCostTypeUsecase(repo repo.AssetCostTypeRepo) *AssetCostTypeUsecase {
	uc := &AssetCostTypeUsecase{
		OwnerUsecase: common.OwnerUsecase[domain.AssetCostType, *dto.AssetCostTypeSearchDTO]{
			OwnerRepo: repo,
		},
		repo: repo,
	}
	uc.UC = uc
	return uc
}

func (uc *AssetCostTypeUsecase) CheckPermission(c context.Context, id uint64) error {
	profileId := _utils.GetProfileIdWithContext(c)
	assetCostType, err := uc.OwnerRepo.GetByID(c, id)
	if err != nil {
		return err
	}

	if assetCostType.OwnerID != profileId || assetCostType.OwnerType != enums.EOwnerOfMember {
		return _errors.ReturnError(service.AccessDenied, _errors.WithLegacyCode(403))
	}
	return nil
}

func (uc *AssetCostTypeUsecase) Create(c context.Context, entity *domain.AssetCostType) error {
	entity.OwnerID = _utils.GetProfileIdWithContext(c)
	entity.OwnerType = enums.EOwnerOfMember
	entity.IsCustom = true
	return uc.OwnerRepo.Create(c, entity)
}

func (uc *AssetCostTypeUsecase) SearchGroup(c context.Context, ownerType enums.EOwnerOf, dto *dto.AssetCostTypeSearchDTO) ([]domain.AssetCostType, []domain.AssetCostType, error) {
	profileId := _utils.GetProfileIdWithContext(c)
	entities, _, err := uc.repo.Search(c, profileId, ownerType, dto)
	if err != nil {
		return nil, nil, err
	}
	incomes := []domain.AssetCostType{}
	spends := []domain.AssetCostType{}
	for _, entity := range entities {
		if entity.Type == enums.ECostTypeIn {
			incomes = append(incomes, entity)
		} else {
			spends = append(spends, entity)
		}
	}
	return incomes, spends, nil
}
