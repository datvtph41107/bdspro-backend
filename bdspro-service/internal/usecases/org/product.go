package org_usecase

import (
	"bdspro/internal/enums"
	shared_usecase "bdspro/internal/usecases/shared"
	"context"
)

type OrgProductUsecase struct {
	*shared_usecase.ProductUsecase
}

func NewOrgProductUsecase(productUsecase *shared_usecase.ProductUsecase) *OrgProductUsecase {
	uc := &OrgProductUsecase{}
	clone := *productUsecase
	newPtr := &clone

	newPtr.ProductUsecase = uc
	uc.ProductUsecase = newPtr
	return uc
}

func (s OrgProductUsecase) GetOwnerID(c context.Context) (*uint64, error) {
	// ownerId, err := _jwt.GetPrincipal(c).GetOrgID()
	// if err != nil {
	// 	return nil, err
	// }
	organizationId := uint64(1)
	return &organizationId, nil
}

func (s OrgProductUsecase) HasPermission(c context.Context, productId *uint64, permission enums.EPermission) (bool, error) {
	return true, nil
}

func (s OrgProductUsecase) GetOwnerType(c context.Context) enums.ETargetType {
	return enums.ETargetTypeOrganization
}
