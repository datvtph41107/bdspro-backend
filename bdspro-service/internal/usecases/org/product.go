package org_usecase

import (
	"bdspro/internal/enums"
	"bdspro/internal/repo"
	user_repo "bdspro/internal/repo/user"
	"bdspro/internal/usecases"
	shared_usecase "bdspro/internal/usecases/shared"
	"context"
)

type OrgProductUC struct {
	*shared_usecase.ProductUsecase
}

func NewOrgProductUsecase(
	productRepo user_repo.UProductRepo,
	ProductHistoryRepo repo.ProductHistoryRepo,
	assetRepo user_repo.UAssetRepo,
	mediaRepo repo.ProductMediaRepo,
	productAccessRepo repo.ProductAccessRepo,
	houseInfoRepo repo.HouseInfoRepo,
	priceRepo repo.ProductPriceRepo,
	productPrivate repo.ProductPrivateRepo,
	region repo.RegionRepo,
	docTypeRepo repo.DocTypeRepo,
	propertyType repo.PropertyTypeRepo,
	amenity repo.AmenityRepo,
	memberPlan usecases.PlanUsecase,
	depositeRepo repo.DepositeRepo,
	postUsecase *OrgPostUC,
	// PostUsecase OrgPostUsecase,
) *OrgProductUC {
	uc := &OrgProductUC{
		ProductUsecase: shared_usecase.MakeProductUsecase(
			productRepo,
			ProductHistoryRepo,
			assetRepo,
			mediaRepo,
			productAccessRepo,
			houseInfoRepo,
			priceRepo,
			productPrivate,
			region,
			docTypeRepo,
			propertyType,
			amenity,
			memberPlan,
			postUsecase.PostUsecase,
			depositeRepo,
			enums.EOwnerTypeOrgnization,
		),
	}
	uc.ProductUsecase.This = uc
	return uc
}

func (s OrgProductUC) GetOwnerID(c context.Context) (*uint64, error) {
	// ownerId, err := _jwt.GetPrincipal(c).GetOrgID()
	// if err != nil {
	// 	return nil, err
	// }
	organizationId := uint64(1)
	return &organizationId, nil
}

func (s OrgProductUC) HasPermission(c context.Context, productId *uint64, permission enums.EPermission) (bool, error) {
	return true, nil
}

func (s OrgProductUC) GetOwnerType(c context.Context) enums.EOwnerType {
	return enums.EOwnerTypeOrgnization
}
