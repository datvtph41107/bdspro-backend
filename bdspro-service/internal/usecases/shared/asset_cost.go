package shared_usecase

import (
	"bdspro/internal/common"
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"bdspro/internal/repo"
	_routes "common/routes"
	_utils "common/utils"
	"context"
)

type AssetCostUsecase struct {
	common.OwnerUsecase[domain.AssetCost, *dto.AssetCostSearchDTO]
	repo repo.AssetCostRepo
}

func NewAssetCostUsecase(repo repo.AssetCostRepo) *AssetCostUsecase {
	uc := &AssetCostUsecase{
		OwnerUsecase: common.OwnerUsecase[domain.AssetCost, *dto.AssetCostSearchDTO]{
			OwnerRepo: repo,
		},
		repo: repo,
	}
	uc.UC = uc
	return uc
}

func (uc *AssetCostUsecase) CheckPermission(c context.Context, id uint64) error {
	profileId := _utils.GetProfileIdWithContext(c)
	assetCost, err := uc.OwnerRepo.GetByID(c, id)
	if err != nil {
		return err
	}

	if assetCost.OwnerID != profileId || assetCost.OwnerType != enums.EOwnerOfMember {
		return &_routes.Except{
			Code:    403,
			Message: "Bạn không có quyền truy cập",
		}
	}
	return nil
}

func (uc *AssetCostUsecase) Create(c context.Context, entity *domain.AssetCost) error {
	entity.OwnerID = _utils.GetProfileIdWithContext(c)
	entity.OwnerType = enums.EOwnerOfMember
	return uc.OwnerRepo.Create(c, entity)
}

func (uc *AssetCostUsecase) Find(c context.Context, ownerType enums.EOwnerOf, dto *dto.AssetCostSearchDTO) ([]domain.AssetCost, int64, error) {
	profileId := _utils.GetProfileIdWithContext(c)
	return uc.repo.Find(c, profileId, ownerType, dto)
}
