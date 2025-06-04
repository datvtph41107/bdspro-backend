package org_usecase

import (
	"bdspro/internal/enums"
	iusecase "bdspro/internal/usecases/interface"
	shared_usecase "bdspro/internal/usecases/shared"
	"context"
)

type OrgAssetUsecase struct {
	*shared_usecase.AssetUsecase
	ClientRepo iusecase.IClientUC
}

func NewOrgAssetUsecase(assetUsecase *shared_usecase.AssetUsecase,
	clientRepo iusecase.IClientUC) *OrgAssetUsecase {
	uc := &OrgAssetUsecase{}
	clone := *assetUsecase
	newPtr := &clone

	newPtr.AssetUsecase = uc
	uc.AssetUsecase = newPtr
	uc.ClientRepo = clientRepo
	return uc
}

func (u *OrgAssetUsecase) GetOwnerID(c context.Context) (*uint64, error) {
	organizationId := uint64(1)
	return &organizationId, nil
}

func (u *OrgAssetUsecase) GetOwnerType(c context.Context) enums.EOwnerType {
	return enums.EOwnerTypeOrgnization
}

func (u *OrgAssetUsecase) HasPermission(c context.Context, assetId *uint64, permission enums.EPermission) (bool, error) {
	return true, nil
}
