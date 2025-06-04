package user_usecase

import (
	"bdspro/internal/enums"
	iusecase "bdspro/internal/usecases/interface"
	shared_usecase "bdspro/internal/usecases/shared"
	_utils "common/utils"
	"context"
)

type UserAssetUsecase struct {
	*shared_usecase.AssetUsecase
	ClientRepo iusecase.IClientUC
}

func NewUserAssetUsecase(assetUsecase *shared_usecase.AssetUsecase,
	clientRepo iusecase.IClientUC) *UserAssetUsecase {
	uc := &UserAssetUsecase{}
	clone := *assetUsecase
	newPtr := &clone

	newPtr.AssetUsecase = uc
	uc.AssetUsecase = newPtr
	uc.ClientRepo = clientRepo
	return uc
}

func (u *UserAssetUsecase) GetOwnerID(c context.Context) (*uint64, error) {
	profileId := _utils.GetProfileIdWithContext(c)
	return &profileId, nil
}

func (u *UserAssetUsecase) GetOwnerType(c context.Context) enums.EOwnerType {
	return enums.EOwnerTypeOrgnization
}

func (u *UserAssetUsecase) HasPermission(c context.Context, assetId *uint64, permission enums.EPermission) (bool, error) {
	return true, nil
}
