package user_usecase

import (
	"bdspro/internal/common"
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	user_repo "bdspro/internal/repo/user"
	_routes "common/routes"
	_utils "common/utils"
	"context"
)

type UserAssetCostUc struct {
	common.OwnerUsecase[domain.AssetCost, *dto.AssetCostSearchDTO]
}

func NewUserAssetCostUc(repo user_repo.UserAssetCostRepo) *UserAssetCostUc {
	uc := &UserAssetCostUc{
		OwnerUsecase: common.OwnerUsecase[domain.AssetCost, *dto.AssetCostSearchDTO]{
			OwnerRepo: repo,
		},
	}
	uc.UC = uc
	return uc
}

func (uc *UserAssetCostUc) CheckPermission(c context.Context, id uint64) error {
	profileId := _utils.GetProfileIdWithContext(c)
	assetCost, err := uc.OwnerRepo.GetByID(c, id)
	if err != nil {
		return err
	}

	if assetCost.OwnerID != profileId || assetCost.OwnerType != enums.EOwnerTypeMember {
		return &_routes.Except{
			Code:    403,
			Message: "Bạn không có quyền truy cập",
		}
	}
	return nil
}

func (uc *UserAssetCostUc) Create(c context.Context, entity *domain.AssetCost) error {
	entity.OwnerID = _utils.GetProfileIdWithContext(c)
	entity.OwnerType = enums.EOwnerTypeMember
	return uc.OwnerRepo.Create(c, entity)
}
