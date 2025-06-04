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

type UserCostTypeUc struct {
	common.OwnerUsecase[domain.AssetCostType, *dto.AssetCostTypeSearchDTO]
}

func NewUserCostTypeUc(repo user_repo.UserCostTypeRepo) *UserCostTypeUc {
	uc := &UserCostTypeUc{
		OwnerUsecase: common.OwnerUsecase[domain.AssetCostType, *dto.AssetCostTypeSearchDTO]{
			OwnerRepo: repo,
		},
	}
	uc.UC = uc
	return uc
}

func (uc *UserCostTypeUc) CheckPermission(c context.Context, id uint64) error {
	profileId := _utils.GetProfileIdWithContext(c)
	assetCostType, err := uc.OwnerRepo.GetByID(c, id)
	if err != nil {
		return err
	}

	if assetCostType.OwnerID != profileId || assetCostType.OwnerType != enums.EOwnerTypeMember {
		return &_routes.Except{
			Code:    403,
			Message: "Bạn không có quyền truy cập",
		}
	}
	return nil
}

func (uc *UserCostTypeUc) Create(c context.Context, entity *domain.AssetCostType) error {
	entity.OwnerID = _utils.GetProfileIdWithContext(c)
	entity.OwnerType = enums.EOwnerTypeMember
	entity.IsCustom = true
	return uc.OwnerRepo.Create(c, entity)
}
