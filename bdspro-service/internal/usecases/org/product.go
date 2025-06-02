package org_usecase

import (
	"bdspro/internal/enums"
	"bdspro/internal/repo"
	"bdspro/internal/usecases"
	shared_usecase "bdspro/internal/usecases/shared"
	"context"
)

type OrgProductUsecase struct {
	shared_usecase.ProductUsecase
}

func NewOrgProductUsecase(
	productRepo repo.UProductRepo,
	ProductHistoryRepo repo.ProductHistoryRepo,
	assetRepo repo.AssetRepo,
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
	UserPostUsecase *usecases.UPostUsecase,
	UAssetUsecase *usecases.UAssetUsecase,
	depositeRepo repo.DepositeRepo,
) *OrgProductUsecase {
	uc := &OrgProductUsecase{}
	orgProductUsecase := shared_usecase.MakeProductUsecase(
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
		UserPostUsecase,
		UAssetUsecase,
		depositeRepo,
		enums.ETargetTypeOrganization,
		uc,
	)
	uc.ProductUsecase = orgProductUsecase
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
