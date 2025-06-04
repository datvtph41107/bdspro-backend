package user_usecase

import (
	"bdspro/infra/providers"
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"bdspro/internal/repo"
	user_repo "bdspro/internal/repo/user"
	"bdspro/internal/usecases"
	shared_usecase "bdspro/internal/usecases/shared"
	_utils "common/utils"
	"context"
)

type UserProductUC struct {
	*shared_usecase.ProductUsecase
	CacheProductMarket *providers.ListCacheProvider[domain.ProductMarket, dto.ProductSearchRequest]
}

func NewUserProductUsecase(
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
	postUsecase *UserPostUC,
	// UserPostUsecase usecases.UPostUsecase,
	// UAssetUsecase usecases.UAssetUsecase,
) *UserProductUC {
	uc := &UserProductUC{
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
			enums.EOwnerTypeMember,
		),
		CacheProductMarket: &providers.ListCacheProvider[domain.ProductMarket, dto.ProductSearchRequest]{},
	}
	uc.ProductUsecase.This = uc
	// uc.CacheProductMarket.Setup()
	return uc
}

func (u *UserProductUC) GetOwnerID(c context.Context) (*uint64, error) {
	profileId := _utils.GetProfileIdWithContext(c)
	return &profileId, nil
}

func (u *UserProductUC) GetOwnerType(c context.Context) enums.EOwnerType {
	return enums.EOwnerTypeMember
}

func (u *UserProductUC) HasPermission(c context.Context, productId *uint64, permission enums.EPermission) (bool, error) {
	return true, nil
}

func (s *UserProductUC) QueryProductMarket(c context.Context, dto *dto.ProductSearchRequest) ([]domain.ProductMarket, int64, error) {
	var results []domain.ProductMarket
	profileId := _utils.GetProfileIdWithContext(c)
	results, total, err := s.ProductRepo.SearchProductMarket(profileId, *dto)
	return results, total, err
}
func (s *UserProductUC) SearchProductMarket(c context.Context, dto *dto.ProductSearchRequest) ([]domain.ProductMarket, int64, error) {
	// var results []u_product_dto.ProductMarket
	// var total int64
	// var err error
	results, total, err := s.CacheProductMarket.GetDataWithCache(c, dto, s.QueryProductMarket)
	// if len(productCache) > 0 {
	// 	results = productCache
	// 	total = totalMarlket
	// } else {
	// 	results, total, err = s.ProductRepo.SearchProductMarket(profileId, *dto, _utils.GetPage(c))
	// 	totalMarlket = total
	// 	productCache = results
	// }

	return results, total, err
}
