package group_usecase

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	org_dto "bdspro/internal/dto/org"
	"bdspro/internal/enums"
	iusecase "bdspro/internal/usecases/interface"
	shared_usecase "bdspro/internal/usecases/shared"
	_db "common/db"
	_routes "common/routes"
	_utils "common/utils"
	"context"
)

type GroupAssetUsecase struct {
	*shared_usecase.AssetUsecase
	ClientRepo iusecase.IClientUC
}

func NewGroupAssetUsecase(
	assetUsecase *shared_usecase.AssetUsecase,
	clientRepo iusecase.IClientUC,
) *GroupAssetUsecase {
	uc := &GroupAssetUsecase{}
	clone := *assetUsecase
	newPtr := &clone

	newPtr.AssetUsecase = uc
	uc.AssetUsecase = newPtr
	uc.ClientRepo = clientRepo

	return uc
}

func (u *GroupAssetUsecase) GetOwnerID(c context.Context) (*uint64, error) {
	return nil, nil
}

func (u *GroupAssetUsecase) GetOwnerType(c context.Context) enums.EOwnerType {
	return enums.EOwnerTypeGroup
}

func (u *GroupAssetUsecase) HasPermission(c context.Context, productId *uint64, permission enums.EPermission) (bool, error) {
	return true, nil
}

func (u *GroupAssetUsecase) GetAsset(c context.Context, groupId uint64, dto dto.ProductSearchRequest) ([]domain.Product, int64, error) {
	profileId := _utils.GetProfileIdWithContext(c)
	userInGroup, err := u.ClientRepo.UserInGroup(c, profileId, groupId)
	if err != nil {
		return nil, 0, err
	}
	if !userInGroup {
		return nil, 0, &_routes.Except{
			Code:    401,
			Message: "Bạn không trong nhóm này",
		}
	}
	return u.AssetUsecase.SearchByAssetShare(c, &groupId, enums.EOwnerTypeGroup, &dto)
}

func (u *GroupAssetUsecase) ShareProduct(c context.Context, dto org_dto.ProductShareRequest) (*domain.ProductAccess, error) {
	_, err := u.RequiredOwner(c, dto.ProductID)
	if err != nil {
		return nil, err
	}

	profileId := _utils.GetProfileIdWithContext(c)
	userInGroup, err := u.ClientRepo.UserInGroup(c, profileId, dto.GroupID)
	if err != nil {
		return nil, err
	}
	if !userInGroup {
		return nil, &_routes.Except{
			Code:    401,
			Message: "Bạn không trong nhóm này",
		}
	}

	productAccess := domain.ProductAccess{
		ProductID:  dto.ProductID,
		TargetId:   dto.GroupID,
		TargetType: enums.EOwnerTypeGroup,
	}

	err = _db.SaveWithContext(c, &productAccess)
	if err != nil {
		return nil, err
	}

	return &productAccess, nil
}
