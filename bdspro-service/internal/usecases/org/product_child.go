package org_usecase

import (
	"bdspro/internal/repo"
	user_repo "bdspro/internal/repo/user"
	shared_usecase "bdspro/internal/usecases/shared"
)

type OrgProductChildUsecase struct {
	*shared_usecase.ProductChildUsecase
}

func NewOrgProductChildUsecase(
	productChildRepo repo.UProductChildRepo,
	productRepo user_repo.UProductRepo,
	postRepo repo.UPostRepo,
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
	ProductUsecase *OrgProductUC,
) *OrgProductChildUsecase {
	uc := &OrgProductChildUsecase{
		ProductChildUsecase: shared_usecase.MakeProductChildUsecase(
			productChildRepo,
			productRepo,
			postRepo,
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
		),
	}
	uc.ProductUsecase = ProductUsecase.ProductUsecase
	return uc
}
