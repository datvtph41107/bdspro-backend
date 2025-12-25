package shared_usecase

import (
	"bdspro/infra/providers"
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"bdspro/internal/provider"
	"bdspro/internal/repo"
	"bdspro/internal/usecases"
	"bdspro/internal/utils"
	_dto "common/domain/dto"
	_enum "common/domain/enum"
	_routes "common/routes"
	_utils "common/utils"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/jinzhu/copier"
)

type IProductUsecase interface {
	GetOwnerID(c context.Context) (*uint64, error)
	HasPermission(c context.Context, productId *uint64, permission enums.EPermission) (bool, error)
	GetOwnerType(c context.Context) enums.EOwnerOf
	// CheckPermission(c context.Context, productID uint64, dto *dto.ProductSearchRequest) error
}

type ProductUsecase struct {
	This               IProductUsecase
	ProductRepo        repo.ProductRepo
	ProductChildRepo   repo.UProductChildRepo
	AssetRepo          repo.AssetRepo
	MediaRepo          repo.ProductMediaRepo
	ProductAccessRepo  repo.SharingAccessRepo
	HouseInfoRepo      repo.HouseInfoRepo
	PriceRepo          repo.ProductPriceRepo
	ProductPrivate     repo.ProductPrivateRepo
	Region             repo.RegionRepo
	DocTypeRepo        repo.DocTypeRepo
	PropertyType       repo.PropertyTypeRepo
	Amenity            repo.AmenityRepo
	ProductHistoryRepo repo.ProductHistoryRepo
	MemberPlan         usecases.PlanUsecase
	PostUsecase        *PostUsecase
	OwnerType          enums.EOwnerOf
	CacheProductMarket *providers.ListCacheProvider[domain.ProductMarket, dto.ProductSearchRequest]
	NotificationClient provider.NotificationProvider
	Transaction        provider.TransactionProvider
	AssetUsecase       *AssetUsecase
	DeepseekClient     provider.DeepseekProvider

	ProductOrganizationRepo repo.ProductOrganizationRepo
	ProductUserRepo         repo.ProductUserRepo
	BdsDomainUsecase        *usecases.BdsDomainUsecase
	DealRepo                repo.DealRepository
}

func NewProductUsecase(
	productRepo repo.ProductRepo,
	ProductHistoryRepo repo.ProductHistoryRepo,
	assetRepo repo.AssetRepo,
	mediaRepo repo.ProductMediaRepo,
	productAccessRepo repo.SharingAccessRepo,
	houseInfoRepo repo.HouseInfoRepo,
	priceRepo repo.ProductPriceRepo,
	productPrivate repo.ProductPrivateRepo,
	region repo.RegionRepo,
	docTypeRepo repo.DocTypeRepo,
	propertyType repo.PropertyTypeRepo,
	amenity repo.AmenityRepo,
	memberPlan usecases.PlanUsecase,
	PostUsecase *PostUsecase,
	productChildRepo repo.UProductChildRepo,
	notificationClient provider.NotificationProvider,
	transaction provider.TransactionProvider,
	assetUsecase *AssetUsecase,
	productOrganizationRepo repo.ProductOrganizationRepo,
	productUserRepo repo.ProductUserRepo,
	deepseekClient provider.DeepseekProvider,
	bdsDomainUsecase *usecases.BdsDomainUsecase,
	dealRepo repo.DealRepository,
) *ProductUsecase {
	result := &ProductUsecase{
		ProductRepo:        productRepo,
		AssetRepo:          assetRepo,
		ProductAccessRepo:  productAccessRepo,
		MediaRepo:          mediaRepo,
		HouseInfoRepo:      houseInfoRepo,
		PriceRepo:          priceRepo,
		ProductPrivate:     productPrivate,
		Region:             region,
		DocTypeRepo:        docTypeRepo,
		PropertyType:       propertyType,
		Amenity:            amenity,
		MemberPlan:         memberPlan,
		ProductHistoryRepo: ProductHistoryRepo,
		PostUsecase:        PostUsecase,
		ProductChildRepo:   productChildRepo,
		CacheProductMarket: &providers.ListCacheProvider[domain.ProductMarket, dto.ProductSearchRequest]{},
		NotificationClient: notificationClient,
		Transaction:        transaction,
		AssetUsecase:       assetUsecase,
		DeepseekClient:     deepseekClient,

		ProductOrganizationRepo: productOrganizationRepo,
		ProductUserRepo:         productUserRepo,
		BdsDomainUsecase:        bdsDomainUsecase,
		DealRepo:                dealRepo,
		// OwnerType:    ownerType,
	}

	result.PostUsecase.ProductUsecase = result

	return result
}

// func MakeProductUsecase(
// 	productRepo repo.ProductRepo,
// 	ProductHistoryRepo repo.ProductHistoryRepo,
// 	assetRepo repo.AssetRepo,
// 	mediaRepo repo.ProductMediaRepo,
// 	productAccessRepo repo.SharingAccessRepo,
// 	houseInfoRepo repo.HouseInfoRepo,
// 	priceRepo repo.ProductPriceRepo,
// 	productPrivate repo.ProductPrivateRepo,
// 	region repo.RegionRepo,
// 	docTypeRepo repo.DocTypeRepo,
// 	propertyType repo.PropertyTypeRepo,
// 	amenity repo.AmenityRepo,
// 	memberPlan usecases.PlanUsecase,
// 	PostUsecase *PostUsecase,
// 	ownerType enums.EOwnerOf,
// ) *ProductUsecase {
// 	return &ProductUsecase{
// 		ProductRepo:        productRepo,
// 		AssetRepo:          assetRepo,
// 		ProductAccessRepo:  productAccessRepo,
// 		MediaRepo:          mediaRepo,
// 		HouseInfoRepo:      houseInfoRepo,
// 		PriceRepo:          priceRepo,
// 		ProductPrivate:     productPrivate,
// 		Region:             region,
// 		DocTypeRepo:        docTypeRepo,
// 		PropertyType:       propertyType,
// 		Amenity:            amenity,
// 		MemberPlan:         memberPlan,
// 		ProductHistoryRepo: ProductHistoryRepo,
// 		PostUsecase:        PostUsecase,
// 		OwnerType:          ownerType,
// 	}
// }

func (s *ProductUsecase) CheckPermission(c context.Context, productId *uint64, permission enums.EPermission) error {
	// hasPermission, err := s.This.HasPermission(c, productId, permission)
	// if !hasPermission {
	// 	return &_routes.Except{
	// 		Code:    401,
	// 		Message: "Bạn không có quyền tạo sản phẩm",
	// 	}
	// }
	// return err
	return nil
}

func (s *ProductUsecase) GetOwnerID(c context.Context, ownerType enums.EOwnerOf) uint64 {
	// if ownerType == enums.EOwnerOfOrganization {
	// 	return s.This.GetOwnerID(c)
	// }

	ownerId := _utils.GetProfileIdWithContext(c)
	return ownerId
}

func (s *ProductUsecase) CreateOrganizationProduct(c context.Context, product *dto.ProductSaveRequest) (*dto.CreateResponse, error) {
	organizationId := _utils.GetOrganizationIdFromContext(c)

	product.RequestOrganizationID = &organizationId
	return s.CreateProduct(c, product)
}

func (s *ProductUsecase) CreateProduct(c context.Context, product *dto.ProductSaveRequest) (*dto.CreateResponse, error) {
	var err error
	err = s.CheckPermission(c, nil, enums.PermissionProductCreate)
	if err != nil {
		return nil, err
	}

	// ownerId là id của tổ chức, hoặc của cá nhân
	ownerId := s.GetOwnerID(c, product.OwnerType)
	if err != nil {
		return nil, err
	}

	entity := &domain.Product{}
	houseInfoEntity := &domain.HouseInfo{}
	copier.Copy(&entity, &product)
	copier.Copy(&houseInfoEntity, &product.HouseInfo)
	var postEntity *domain.Post
	var assetEntity *domain.Asset

	err = s.Transaction.WithTransaction(c, func(c context.Context) error {
		price := product.PriceData
		if err := s.PriceRepo.Create(c, price); err != nil {
			return err
		}
		entity.LastPriceID = &price.ID
		entity.OwnerID = ownerId
		entity.OwnerOf = product.OwnerType
		// _db.SaveWithContext(c, entity)
		if err := s.ProductRepo.Create(c, entity); err != nil {
			return err
		}
		// err := s.productRepo.UpdateProduct(c, entity.ID, entity)
		// err := s.UpdateReference(c, entity, product)

		if product.PriceData != nil {
			if err := s.PriceRepo.UpdateProductID(c, &entity.ID, &product.PriceData.ID); err != nil {
				return err
			}
		}
		if err := s.UpdateAmenity(c, entity.ID, product.AmenityIds); err != nil {
			return err
		}
		mediaIDs, err := s.UpdateMediaItems(c, entity.ID, product.MediaItems)
		if err != nil {
			return err
		}
		// Gán ID của ảnh đầu tiên vào ImageId
		if len(mediaIDs) > 0 {
			entity.ImageId = &mediaIDs[0]
			if err := s.ProductRepo.UpdateImageID(c, entity.ID, entity.ImageId); err != nil {
				return err
			}
		}
		if err := s.UpdateHouseInfo(c, entity.ID, houseInfoEntity); err != nil {
			return err
		}
		if err := s.UpdatePrivateData(c, entity.ID, product.PrivateData); err != nil {
			return err
		}

		postInfo := product.PostCreate
		if postInfo != nil && s.PostUsecase != nil {
			postInfo.ProductID = entity.ID
			postInfo.Content = entity.Description
			postInfo.TransactionType = enums.TransactionType(entity.TransactionType)
			// if entity.TransactionType == 10 {
			// 	postInfo.Price = *product.PriceData.SalePrice
			// } else {
			// 	postInfo.Price = *product.PriceData.RentPrice
			// }

			// dto := dto.PostSaveRequest{}
			// copier.Copy(&dto, &postInfo)

			postEntity, err = s.PostUsecase.CreatePost(c, postInfo)
			if err != nil {
				return err
			}
		}

		assetInfo := product.AssetCreate
		if assetInfo != nil {
			// copier.Copy(&assetEntity, &assetInfo)
			assetEntity = &domain.Asset{
				ProductID:      &entity.ID,
				Name:           entity.Name,
				Area:           entity.Area,
				PropertyTypeId: entity.PropertyTypeId,
				Description:    entity.Description,
				ProvinceID:     entity.ProvinceID,
				WardID:         entity.WardID,
				PurchasePrice:  assetInfo.PurchasePrice,
				PurchaseDate:   assetInfo.PurchaseDate,
				LegalStatus:    enums.EDocType(assetInfo.LegalStatus),
				LegalItems:     assetInfo.LegalItems,
				RentStatus:     assetInfo.RentStatus,
				OwnerOf:        entity.OwnerOf,
				// Description:    assetInfo.Description,
			}

			assetEntity, err = s.AssetUsecase.CreateAsset(c, assetEntity)
			if err != nil {
				return err
			}

			err = s.ProductRepo.UpdateAssetID(c, entity.ID, assetEntity.ID)
			if err != nil {
				return err
			}
		}

		// sản phẩm của tổ chức
		var profileId uint64
		if product.RequestOrganizationID != nil {
			_, err = s.ProductOrganizationRepo.CreateOwner(c, entity.ID, *product.RequestOrganizationID)
			if err != nil {
				return err
			}
		} else {
			// sản phẩm của người dùng
			profileId = _utils.GetProfileIdWithContext(c)
			_, err = s.ProductUserRepo.CreateOwner(c, entity.ID, profileId)
			if err != nil {
				return err
			}
		}

		// Tạo BDSDomain tương ứng cho Product
		if s.BdsDomainUsecase != nil {
			// Nếu chưa có profileId (tổ chức), lấy từ context
			if profileId == 0 {
				profileId = _utils.GetProfileIdWithContext(c)
			}
			// Set ProductID cho houseInfoEntity nếu chưa có
			if houseInfoEntity != nil {
				houseInfoEntity.ProductID = entity.ID
				err = s.BdsDomainUsecase.CreateBdsDomainFromProduct(c, entity, houseInfoEntity, profileId)
				// Log error nhưng không fail transaction nếu tạo BDSDomain thất bại
				if err != nil {
					log.Printf("Failed to create BDSDomain for product %d: %v", entity.ID, err)
				}
			} else {
				// Tạo BDSDomain ngay cả khi không có HouseInfo
				err = s.BdsDomainUsecase.CreateBdsDomainFromProduct(c, entity, nil, profileId)
				if err != nil {
					log.Printf("Failed to create BDSDomain for product %d: %v", entity.ID, err)
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &dto.CreateResponse{
		Product: entity,
		Post:    postEntity,
		Asset:   assetEntity,
	}, nil
}

func (s *ProductUsecase) updateProduct(contextTx context.Context, id uint64, product *dto.ProductSaveRequest) error {
	entity, err := s.ProductRepo.GetByIDContext(contextTx, id)
	if err != nil {
		return err
	}

	copier.Copy(&entity, &product)

	// err = s.productRepo.UpdateProduct(c, id, product)
	// s.UpdateReference(c, product, dto)
	if product.PriceData != nil {
		product.PriceData.ProductID = &id
	}
	if err := s.UpdatePrice(contextTx, entity.LastPriceID, product.PriceData); err != nil {
		return err
	}
	if err := s.UpdateAmenity(contextTx, entity.ID, product.AmenityIds); err != nil {
		return err
	}
	mediaIDs, err := s.UpdateMediaItems(contextTx, entity.ID, product.MediaItems)
	if err != nil {
		return err
	}
	// Gán ID của ảnh đầu tiên vào ImageId
	if len(mediaIDs) > 0 {
		entity.ImageId = &mediaIDs[0]
		if err := s.ProductRepo.UpdateImageID(contextTx, entity.ID, entity.ImageId); err != nil {
			return err
		}
	}

	houseInfoEntity := &domain.HouseInfo{}
	copier.Copy(&houseInfoEntity, &product.HouseInfo)
	if err := s.UpdateHouseInfo(contextTx, entity.ID, houseInfoEntity); err != nil {
		return err
	}
	if err := s.UpdatePrivateData(contextTx, entity.ID, product.PrivateData); err != nil {
		return err
	}
	entity.LastPriceID = &product.PriceData.ID
	if err = s.ProductRepo.UpdateProduct(contextTx, id, entity); err != nil {
		return err
	}
	return nil
}

func (s *ProductUsecase) Update(c context.Context, id uint64, product *dto.ProductSaveRequest) (*domain.Product, error) {
	// profileId := _jwt.GetProfileId(c)
	err := s.CheckPermission(c, &id, enums.PermissionProductUpdate)
	if err != nil {
		return nil, err
	}
	// ownerId, err := s.ProductUsecase.GetOwnerID(c)
	// if err != nil {
	// 	return nil, err
	// }
	// entity, err := s.CheckAccessProduct(c, ownerId, id)
	// if err != nil {
	// 	return nil, err
	// }

	contextTx := s.Transaction.StartTransaction(c)
	if err = s.updateProduct(contextTx, id, product); err != nil {
		s.Transaction.RollbackTransaction(contextTx)
		return nil, err
	}

	if err = s.Transaction.CommitTransaction(contextTx); err != nil {
		return nil, err
	}

	return nil, nil
}

// func (s *ProductUsecase) CheckAccessProduct(c context.Context, profileId, productId uint64) (*domain.Product, error) {
// 	existed, err := s.ProductAccessRepo.ExistsByUserAndProduct(profileId, productId)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if err != nil {
// 		return nil, err
// 	}

// 	product, err := s.ProductRepo.GetByID(productId)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if !existed && *product.OwnerUserID != profileId {
// 		return nil, &_routes.Except{
// 			Code:    401,
// 			Message: "Bạn không có quyền truy cập",
// 		}
// 	}

// 	return product, nil
// }

func (s *ProductUsecase) GetProductOfGroup(c context.Context, groupId uint64, dto dto.ProductSearchRequest) ([]dto.ProductListResponse, int64, error) {
	dto.RequestGroupId = &groupId
	return s.searchProduct(c, dto)
}

func (s *ProductUsecase) GetProductOfCurrentUser(c context.Context, dto dto.ProductSearchRequest) ([]dto.ProductListResponse, int64, error) {
	profileId := _utils.GetProfileIdWithContext(c)
	dto.RequestUserId = &profileId
	return s.searchProduct(c, dto)
}

func (s *ProductUsecase) GetProductOfCurrentOrganization(c context.Context, dto dto.ProductSearchRequest) ([]dto.ProductListResponse, int64, error) {
	organizationId := _utils.GetOrganizationIdFromContext(c)
	dto.RequestOrganizationId = &organizationId
	return s.searchProduct(c, dto)
}

func (s *ProductUsecase) searchProduct(c context.Context, payload dto.ProductSearchRequest) ([]dto.ProductListResponse, int64, error) {
	results, total, err := s.ProductRepo.GetProducts(c, payload)
	if err != nil {
		return nil, 0, err
	}
	for i := 0; i < len(results); i++ {
		s.MapProductStatusList(c, &results[i])
	}
	return results, total, nil
}

func (s *ProductUsecase) SearchExport(c context.Context, dto *dto.ProductSearchRequest) ([]domain.Product, error) {
	ownerID, err := s.This.GetOwnerID(c)
	if err != nil {
		return nil, err
	}
	results, err := s.ProductRepo.SearchExport(ownerID, s.This.GetOwnerType(c), dto)
	return results, err
}

// func (s *ProductUsecase) UpdateReference(c context.Context, entity *domain.Product, dto *dto.ProductSaveRequest) error {
// 	if dto.AmenityIds != nil {
// 		s.ProductRepo.UpdateProductAmenityIds(entity.ID, dto.AmenityIds)
// 	}

// 	// price := entity.Price
// 	privateInfo := entity.PrivateData
// 	houseInfo := entity.HouseInfo
// 	// entity.Price = nil
// 	// entity.PrivateData = nil
// 	// entity.HouseInfo = nil

// 	// err := _db.DB.WithContext(c).Model(entity).
// 	// 	Where("id = ? and deleted_at is null", id).
// 	// 	Updates(entity).Error

// 	s.MediaRepo.UpdateProductMediaItems(entity.ID, dto.MediaItemDTOs)
// 	// var houseInfo domain.HouseInfo
// 	if dto.HouseInfoDTO != nil {
// 		// copier.Copy(&houseInfo, &dto.HouseInfo)
// 		houseInfo.ProductID = entity.ID
// 		s.HouseInfoRepo.UpdateHouseInfo(c, houseInfo)
// 	}

// 	err := s.ProductPrivate.UpdateInfo(c, entity.ID, privateInfo)

// 	return err
// }

func (s *ProductUsecase) Detail(c context.Context, id uint64) (*domain.Product, error) {
	err := s.CheckPermission(c, &id, enums.PermissionProductRead)
	if err != nil {
		return nil, err
	}

	entity, err := s.ProductRepo.GetByIDContext(c, id)
	if err != nil {
		return nil, &_routes.Except{
			Code:    404,
			Message: "Không tìm thấy dữ liệu",
		}
	}

	s.MapProductStatus(c, entity)

	// Kiểm tra xem user có phải là người tạo không để load privateData
	profileID := _utils.GetProfileIdWithContext(c)
	if profileID > 0 && entity.CreatedBy != nil && *entity.CreatedBy == profileID {
		// Lấy privateData từ repo và gán vào entity
		privateData := s.ProductPrivate.FindByProductId(c, &entity.ID)
		if privateData != nil {
			entity.PrivateData = privateData
		}
	}

	return entity, nil
}

func (s *ProductUsecase) SmartSearch(c context.Context, text string, dto *dto.TextSearchRequest) ([]domain.Product, int64, error) {
	searchs, err := utils.InferSearchText(text)

	regionIds, _ := s.Region.InterText(text, 10)
	docTypeIds, _ := s.DocTypeRepo.InterText(c, text, 10)
	propertyTypeIds, _ := s.PropertyType.InterText(text, 10)

	searchs.RegionIds = regionIds
	searchs.DocTypeIds = docTypeIds
	searchs.PropertyTypeIds = propertyTypeIds
	searchs.Pagable = _dto.Pagable{
		Page: dto.GetPage(),
		Size: dto.GetSize(),
	}

	ownerID, err := s.This.GetOwnerID(c)
	if err != nil {
		return nil, 0, err
	}

	results, total, err := s.ProductRepo.SmartSearch(ownerID, s.This.GetOwnerType(c), &searchs)
	if err != nil {
		return nil, 0, err
	}
	for i := 0; i < len(results); i++ {
		s.MapProductStatus(c, &results[i])
	}
	return results, total, nil
}

func (s *ProductUsecase) UpdatePrice(c context.Context, priceID *uint64, entity *domain.ProductPrice) error {
	return s.PriceRepo.UpdatePrice(c, priceID, entity)
}

func (s *ProductUsecase) UpdateAmenity(ctx context.Context, productID uint64, ids []uint64) error {
	if len(ids) > 0 {
		return s.ProductRepo.UpdateProductAmenityIds(ctx, productID, ids)
	}
	return nil
}

func (s *ProductUsecase) UpdateMediaItems(ctx context.Context, productID uint64, items []dto.MediaItem) ([]uint64, error) {
	return s.MediaRepo.UpdateProductMediaItems(ctx, productID, items)
}

func (s *ProductUsecase) UpdateHouseInfo(ctx context.Context, productId uint64, houseInfo *domain.HouseInfo) error {
	if houseInfo != nil {
		// copier.Copy(&houseInfo, &dto.HouseInfo)
		houseInfo.ProductID = productId
		return s.HouseInfoRepo.UpdateHouseInfo(ctx, houseInfo)
	}
	return nil
}

func (s *ProductUsecase) UpdatePrivateData(ctx context.Context, productID uint64, privateInfo *domain.ProductPrivate) error {
	return s.ProductPrivate.UpdateInfo(ctx, productID, privateInfo)
}

// func (s *ProductUsecase) InferenceString(text string) string {
// 	searchs, _ := utils.InferSearchText(text)

// 	regionIds, _ := s.Region.InterText(text, 10)
// 	docTypeIds, _ := s.DocTypeRepo.InterText(text, 10)
// 	propertyTypeIds, _ := s.PropertyType.InterText(text, 10)

// 	searchs.RegionIds = regionIds
// 	searchs.DocTypeIds = docTypeIds
// 	searchs.PropertyTypeIds = propertyTypeIds

// 	return ""
// }

func (s *ProductUsecase) SuggestFields(c context.Context, text string) (*_dto.ProductV3DTO, error) {
	// searchs, err := utils.InferSearchText(body.Content)
	// if err != nil {
	// 	return nil, err
	// }
	// if searchs, err := utils.InferSearchText(body.Content); err != nil {
	// 	return nil, err
	// }
	searchs := _dto.ProductV3DTO{}
	// var err error
	// if searchs, err = utils.InferSearchText2(text); err != nil {
	// 	return nil, err
	// }

	// searchs, _ := utils.InferSearchText(text)

	// regions, _ := s.Region.InterTextToItem(text, 10)
	docType, _ := s.DocTypeRepo.InterTextToItem(c, text, 10)
	propertyType, _ := s.PropertyType.InterTextToItem(text, 10)
	amenities, _ := s.Amenity.InterTextToItem(text, 10)

	searchs.Name = _utils.TrimFirstWords(text, 12)

	// tìm địa chỉ
	// var ward, district, province *domain.Region
	// for i := 0; i < len(regions); i++ {
	// 	r := regions[i]
	// 	if r.Level == 1 && province == nil {
	// 		province = &r
	// 		continue
	// 	}
	// 	if r.Level == 2 && district == nil {
	// 		district = &r
	// 		continue
	// 	}
	// 	if r.Level == 3 && ward == nil {
	// 		ward = &r
	// 		continue
	// 	}
	// }
	// if ward != nil && ward.ParentID != nil && (district == nil || district.ID != *ward.ParentID) {
	// 	district, _ = s.Region.GetByID(c, *ward.ParentID)
	// 	province, _ = s.Region.GetByID(c, *district.ParentID)
	// }

	// if district != nil && district.ParentID != nil && (province == nil || province.ID != *district.ParentID) {
	// 	province, _ = s.Region.GetByID(c, *district.ParentID)
	// }
	// ---- xong ----

	searchs.DocType = docType
	searchs.PropertyType = propertyType
	searchs.Amenities = amenities
	// searchs.Ward = ward
	// searchs.District = district
	// searchs.Province = province
	// if searchs.PriceSuggest != nil {
	// 	if *searchs.PriceSuggest > 100_000_000 {
	// 		searchs.TransactionType = 10
	// 		searchs.SalePrice = searchs.PriceSuggest
	// 	} else {
	// 		searchs.TransactionType = 20
	// 		searchs.RentPrice = searchs.PriceSuggest
	// 	}
	// }

	return &searchs, nil
}

// SuggestFieldsWithDeepseek - Gọi DeepSeek AI để phân tích và gợi ý thông tin sản phẩm
func (s *ProductUsecase) SuggestFieldsWithDeepseek(c context.Context, text string) (*dto.TotalSearchParser, error) {
	// Gọi DeepSeek client để phân tích
	deepseekResult, err := s.DeepseekClient.SuggestProductInfo(c, text)
	if err != nil {
		return nil, &_routes.Except{
			Code:    500,
			Message: fmt.Sprintf("Lỗi khi gọi DeepSeek AI: %v", err),
		}
	}

	// Khởi tạo result parser
	result := &dto.TotalSearchParser{}

	// Map các trường cơ bản
	if deepseekResult.Name != nil {
		result.Name = *deepseekResult.Name
	}
	if deepseekResult.Area != nil {
		area := int32(*deepseekResult.Area)
		result.Area = &area
	}
	if deepseekResult.Address != nil {
		result.Address = deepseekResult.Address
	}

	// Map transaction type
	if deepseekResult.TransactionType != nil {
		result.TransactionType = uint32(*deepseekResult.TransactionType)
	}

	// Map giá
	if deepseekResult.SalePrice != nil {
		result.SalePrice = deepseekResult.SalePrice
	}
	if deepseekResult.RentPrice != nil {
		result.RentPrice = deepseekResult.RentPrice
	}

	// Map thông tin nhà
	if deepseekResult.Bedroom != nil {
		result.Bedroom = deepseekResult.Bedroom
	}
	if deepseekResult.Bathroom != nil {
		result.Bathroom = deepseekResult.Bathroom
	}
	if deepseekResult.Floor != nil {
		result.Floor = deepseekResult.Floor
	}
	if deepseekResult.Frontage != nil {
		frontage := int32(*deepseekResult.Frontage)
		result.Frontage = &frontage
	}

	// Tìm region (province, district, ward) từ database
	if deepseekResult.Province != nil {
		provinces, _ := s.Region.InterTextToItem(*deepseekResult.Province, 1)
		if len(provinces) > 0 {
			result.Province = &provinces[0]
		}
	}

	if deepseekResult.District != nil {
		districts, _ := s.Region.InterTextToItem(*deepseekResult.District, 1)
		if len(districts) > 0 {
			result.District = &districts[0]
		}
	}

	if deepseekResult.Ward != nil {
		wards, _ := s.Region.InterTextToItem(*deepseekResult.Ward, 1)
		if len(wards) > 0 {
			result.Ward = &wards[0]
		}
	}

	// Tìm property type từ database
	// if deepseekResult.PropertyType != nil {
	// 	propertyTypes, _ := s.PropertyType.InterTextToItem(*deepseekResult.PropertyType, 1)
	// 	if len(propertyTypes) > 0 {
	// 		result.PropertyTypes = propertyTypes
	// 	}
	// }

	// Tìm doc type từ database
	// if deepseekResult.LegalDoc != nil {
	// 	docTypes, _ := s.DocTypeRepo.InterTextToItem(c, *deepseekResult.LegalDoc, 1)
	// 	if len(docTypes) > 0 {
	// 		result.DocTypes = docTypes
	// 	}
	// }

	// Fallback: Nếu không có province từ DeepSeek, thử tìm từ text gốc
	if result.Province == nil || result.District == nil {
		regions, _ := s.Region.InterTextToItem(text, 10)
		var ward, district, province *domain.Region
		for i := 0; i < len(regions); i++ {
			r := regions[i]
			if r.Level == 1 && province == nil {
				province = &r
			}
			if r.Level == 2 && district == nil {
				district = &r
			}
			if r.Level == 3 && ward == nil {
				ward = &r
			}
		}
		if result.Province == nil {
			result.Province = province
		}
		if result.District == nil {
			result.District = district
		}
		if result.Ward == nil {
			result.Ward = ward
		}
	}

	// Xử lý quan hệ parent-child của region
	if result.Ward != nil && result.Ward.ParentID != nil && (result.District == nil || result.District.ID != *result.Ward.ParentID) {
		result.District, _ = s.Region.GetByID(c, *result.Ward.ParentID)
		if result.District != nil && result.District.ParentID != nil {
			result.Province, _ = s.Region.GetByID(c, *result.District.ParentID)
		}
	}

	if result.District != nil && result.District.ParentID != nil && (result.Province == nil || result.Province.ID != *result.District.ParentID) {
		result.Province, _ = s.Region.GetByID(c, *result.District.ParentID)
	}

	return result, nil
}

func (s *ProductUsecase) RequiredOwner(c context.Context, id uint64) (*domain.Product, error) {
	product, err := s.ProductRepo.GetByIDContext(c, id)
	if err != nil {
		return nil, err
	}

	// profileId := _jwt.GetProfileId(c)
	profileId := _utils.GetProfileIdWithContext(c)

	if product.OwnerID != profileId {
		return nil, &_routes.Except{
			Code:    401,
			Message: "Bạn không có quyền truy cập",
		}
	}
	return product, nil
}

func (s *ProductUsecase) Archived(c context.Context, dto dto.ArchivedRequest) error {
	// if _, err := s.RequiredOwner(c, *dto.ProductId); err != nil {
	// 	return err
	// }
	err := s.CheckPermission(c, dto.ProductId, enums.PermissionProductArchive)
	if err != nil {
		return err
	}

	return s.ProductRepo.Archived(c, dto.ProductId, dto.Archived)
}

func (s *ProductUsecase) Delete(c context.Context, id uint64) error {
	// sản phẩm đang cọc
	err := s.CheckPermission(c, &id, enums.PermissionProductDelete)
	if err != nil {
		return err
	}
	product, err := s.ProductRepo.GetByID(id)
	if err != nil {
		return err
	}

	// todo: thêm hoặc có tin active
	if product.SaleTransactionID != nil || product.RentTransactionID != nil {
		return &_routes.Except{
			Code:    401,
			Message: "Không thể xóa sản phẩm này",
		}
	}
	return s.ProductRepo.Delete(c, id)
}

// func (s *ProductUsecase) GetProductStatus(c context.Context, trans *domain.Transaction) enums.EProductStatus {
// 	if trans.TransactionType == enums.TransactionTypeSale {
// 		return enums.EProductStatus(trans.TransactionStatus)
// 	}
// 	if product.RentTransactionID != nil {
// 		return enums.EProductStatus(product.RentTransaction.TransactionStatus)
// 	}
// 	return enums.EProductStatus(product.TransactionStatus)
// }

// func (s *ProductUsecase) UpdateTransactionStatus(c context.Context,
// 	newStatus enums.TransactionStatus,
// 	transactionType enums.TransactionType,
// 	productId *uint64,
// ) error {
// 	// todo: không thể cập nhật trạng thái sản phẩm thành cọc ở api này. dung api deposite riêng
// 	if newStatus == enums.TransactionStatusDeposite {
// 		return &_routes.Except{
// 			Code:    400,
// 			Message: "Không thể cập nhật trạng thái sản phẩm thành cọc",
// 		}
// 	}
// 	err := s.CheckPermission(c, productId, enums.PermissionProductUpdateStatus)
// 	if err != nil {
// 		return err
// 	}
// 	product, err := s.ProductRepo.GetByID(*productId)
// 	if err != nil {
// 		return err
// 	}
// 	var trans *domain.Transaction
// 	// tạo transaction mới
// 	isNewTransaction := false
// 	if transactionType == enums.TransactionTypeSale && product.SaleTransactionID != nil {
// 		trans, err = s.TransactionRepo.GetByID(c, product.SaleTransactionID)
// 	} else if transactionType == enums.TransactionTypeRent && product.RentTransactionID != nil {
// 		trans, err = s.TransactionRepo.GetByID(c, product.RentTransactionID)
// 	} else {
// 		isNewTransaction = true
// 		trans = s.EmptyTransaction(product)
// 	}
// 	if err != nil {
// 		return err
// 	}

// 	if trans.TransactionStatus == newStatus && trans.ID != 0 {
// 		return &_routes.Except{
// 			Code:    400,
// 			Message: "Trạng thái trùng với trạng thái cũ",
// 		}
// 	}

// 	contextTx := s.Transaction.StartTransaction(c)

// 	// nếu chuyển thành chưa bắt đầu/hủy thì hủy giao dịch cũ đồng thời tạo giao dịch mới
// 	if (newStatus == enums.TransactionStatusDraft ||
// 		newStatus == enums.TransactionStatusCanceled ||
// 		newStatus < trans.TransactionStatus) && trans.ID != 0 {
// 		trans.TransactionStatus = enums.TransactionStatusCanceled
// 		_, err = s.TransactionRepo.Update(contextTx, trans)
// 		if err != nil {
// 			s.Transaction.RollbackTransaction(contextTx)
// 			return err
// 		}
// 		trans = s.EmptyTransaction(product)
// 		isNewTransaction = true
// 	} else {
// 		trans.TransactionStatus = newStatus
// 	}

// 	trans, err = s.TransactionRepo.Update(contextTx, trans)
// 	if err != nil {
// 		s.Transaction.RollbackTransaction(contextTx)
// 		return err
// 	}

// 	if isNewTransaction {
// 		if transactionType == enums.TransactionTypeSale {
// 			product.SaleTransactionID = &trans.ID
// 		} else if transactionType == enums.TransactionTypeRent {
// 			product.RentTransactionID = &trans.ID
// 		}
// 	}
// 	if isNewTransaction {
// 		err = s.ProductRepo.UpdateProduct(contextTx, *productId, product)
// 		if err != nil {
// 			s.Transaction.RollbackTransaction(contextTx)
// 			return err
// 		}
// 	}
// 	if err := s.Transaction.CommitTransaction(contextTx); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (s *ProductUsecase) ProductDeposite(c context.Context, dto *tx_domain.Tx) (*tx_domain.Tx, error) {
// 	// sản phẩm đang cọc
// 	err := s.CheckPermission(c, dto.ProductID, enums.PermissionProductDeposite)
// 	if err != nil {
// 		return nil, err
// 	}

// 	product, err := s.ProductRepo.GetProductByID(c, dto.ProductID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var tx *tx_domain.Tx
// 	// todo: transaction tại đây
// 	// if dto.TransactionType == 10 && product.SaleTransactionID != nil {
// 	// 	tx, err = s.TransactionRepo.GetByID(c, product.SaleTransactionID)
// 	// } else if dto.TransactionType == 20 && product.RentTransactionID != nil {
// 	// 	tx, err = s.TransactionRepo.GetByID(c, product.RentTransactionID)
// 	// } else {
// 	// 	tx = &tx_domain.Tx{}
// 	// }
// 	if err != nil {
// 		return nil, err
// 	}

// 	// typeName := "bán"
// 	// if dto.TransactionType == 20 {
// 	// 	typeName = "cho thuê"
// 	// }

// 	// todo: nếu đặt cọc khi giao dịch không phải nháp/hủy thì sẽ hủy giao dịch cũ đồng thời tạo giao dịch mới
// 	// if tx.TransactionStatus != enums.TransactionStatusDraft &&
// 	// 	tx.TransactionStatus != enums.TransactionStatusCanceled {
// 	// 	return nil, &_routes.Except{
// 	// 		Code:    400,
// 	// 		Message: fmt.Sprintf("Chỉ đặt cọc với sản phẩm chưa %s\nhoặc đã hủy", typeName),
// 	// 	}
// 	// }

// 	// var record *domain.Transaction
// 	err = s.Transaction.WithTransaction(c, func(c context.Context) error {
// 		if tx.TransactionStatus != enums.TransactionStatusDraft &&
// 			tx.TransactionStatus != enums.TransactionStatusCanceled {
// 			tx.TransactionStatus = enums.TransactionStatusCanceled
// 			// todo: transaction
// 			// tx, err = s.TransactionRepo.Update(c, tx)
// 			// if err != nil {
// 			// 	return err
// 			// }
// 		}

// 		newTx := &tx_domain.Tx{
// 			ProductID:         dto.ProductID,
// 			ContactID:         dto.ContactID,
// 			DepositeAmount:    dto.DepositeAmount,
// 			DepositeNote:      dto.DepositeNote,
// 			TransactionType:   dto.TransactionType,
// 			TransactionStatus: enums.TransactionStatusDeposite,
// 			OwnerId:           product.OwnerID,
// 			OwnerType:         product.OwnerType,
// 		}
// 		// newTx, err = s.TransactionRepo.Create(c, newTx)
// 		// if err != nil {
// 		// 	return err
// 		// }
// 		if dto.TransactionType == 10 {
// 			product.SaleTransactionID = &newTx.ID
// 		} else {
// 			product.RentTransactionID = &newTx.ID
// 		}
// 		if err := s.ProductRepo.Save(c, product); err != nil {
// 			return err
// 		}
// 		return nil
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	// product.DepositeID = &record.ID
// 	// if err := s.ProductRepo.Update(c, product.ID, product); err != nil {
// 	// 	return nil, err
// 	// }
// 	return tx, nil
// }

func (s *ProductUsecase) RentVisibilityUpdate(c context.Context, dto dto.VisibilityUpdate) error {
	err := s.CheckPermission(c, dto.ProductID, enums.PermissionProductUpdateVisibility)
	if err != nil {
		return err
	}
	return s.ProductRepo.SetRentVisibility(c, dto.ProductID, dto.Visibility)
}

func (s *ProductUsecase) SaleVisibilityUpdate(c context.Context, dto dto.VisibilityUpdate) error {
	err := s.CheckPermission(c, dto.ProductID, enums.PermissionProductUpdateVisibility)
	if err != nil {
		return err
	}
	return s.ProductRepo.SetSaleVisibility(c, dto.ProductID, dto.Visibility)
}

func (s *ProductUsecase) CancelDeposite(c context.Context, transactionId uint64) error {
	// sản phẩm đang cọc
	err := s.CheckPermission(c, &transactionId, enums.PermissionProductDeposite)
	if err != nil {
		return err
	}

	// todo: transaction
	// transaction, err := s.TransactionRepo.GetByID(c, &transactionId)
	// if err != nil {
	// 	return err
	// }
	// product, err := s.RequiredOwner(c, deposite.ProductID)
	// if err != nil {
	// 	return err
	// }

	// product, err := s.ProductRepo.GetProductByID(c, deposite.ProductID)
	// if err != nil {
	// 	return err
	// }

	// if product.DepositeID == nil {
	// 	return &_routes.Except{
	// 		Code:    400,
	// 		Message: "Sản phẩm chưa đặt cọc",
	// 	}
	// }

	// product.DepositeID = nil

	// return s.TransactionRepo.CancelDeposite(c, transaction.ID)
	return nil
}

func (uc *ProductUsecase) Posts(c context.Context, id uint64, dto dto.PostLinkSearch) (*[]domain.PostItem, error) {
	if uc.PostUsecase != nil {
		results, err := uc.PostUsecase.ProductPostLink(c, id, &dto)
		return results, err
	}
	return nil, nil
}

func (uc *ProductUsecase) Members(c context.Context, id uint64, dto dto.SharingAccessSearch) (*[]domain.SharingAccess, error) {
	// _, err := s.RequiredOwnerProduct(c, id)
	// if err != nil {
	// 	return nil, err
	// }

	entities, err := uc.ProductAccessRepo.OwnerSearchProduct(c, id, enums.EOwnerOfMember, &dto)
	return entities, err
}

func (uc *ProductUsecase) GetDealsByProduct(c context.Context, productID uint64, pagable _dto.Pagable) ([]*domain.Deal, int64, error) {
	if uc.DealRepo == nil {
		return nil, 0, &_routes.Except{
			Code:    500,
			Message: "Deal repository not initialized",
		}
	}

	deals, total, err := uc.DealRepo.GetDealsByProductID(c, productID, pagable)
	if err != nil {
		return nil, 0, err
	}

	return deals, total, nil
}

func (uc *ProductUsecase) History(c context.Context, productId uint64, dto dto.ProductHistorySearch) (*[]domain.ProductHistory, int64, error) {
	profileId := _utils.GetProfileIdWithContext(c)
	entities, total, err := uc.ProductHistoryRepo.History(c, profileId, productId, dto)
	if err != nil {
		return nil, 0, &_routes.Except{
			Code:    404,
			Message: "Không tìm thấy dữ liệu",
		}
	}

	return entities, total, nil
}

func (s *ProductUsecase) SearchByProductAccess(c context.Context, targetID *uint64, targetType enums.EOwnerOf, dto *dto.ProductSearchRequest) ([]domain.Product, int64, error) {
	return s.ProductRepo.SearchByProductAccess(c, targetID, targetType, dto)
}

func (s *ProductUsecase) RequiredBeforeSave(c context.Context, parentID uint64) (*domain.Product, error) {
	parent, err := s.ProductRepo.GetByID(parentID)
	if err != nil {
		return nil, err
	}

	// if parent.AssetId == nil {
	// 	return nil, &_routes.Except{
	// 		Code:    400,
	// 		Message: "Sản phẩm không đủ điều kiện để tạo sản phẩm con",
	// 	}
	// }

	profileId := _utils.GetProfileIdWithContext(c)

	if parent.OwnerID != profileId {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Sản phẩm không thuộc quyền sở hữu của bạn",
		}
	}

	return parent, nil
}

func (s *ProductUsecase) InfoArea(c context.Context, id *uint64) (*dto.InfoAreaResponse, error) {
	_, err := s.RequiredBeforeSave(c, *id)
	if err != nil {
		return nil, err
	}

	profileId := _utils.GetProfileIdWithContext(c)
	result, err := s.ProductChildRepo.InfoArea(c, &profileId, id)
	return result, err
}

func (s *ProductUsecase) CreateChild(c context.Context, body *dto.ProductSaveRequest) (*dto.CreateResponse, error) {
	parent, err := s.RequiredBeforeSave(c, *body.ParentID)
	if err != nil {
		return nil, err
	}
	totalAreaDevide := s.ProductChildRepo.TotalAreaDevide(body.ParentID)
	if totalAreaDevide+body.Area > parent.Area {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Diện tích không hợp lệ",
		}
	}

	body.ParentID = &parent.ID
	body.OwnerType = parent.OwnerOf
	result, err := s.CreateProduct(c, body)

	go func() {
		ctxClone := _utils.CloneContext(c)
		profileId := _utils.GetProfileIdWithContext(ctxClone)
		s.NotificationClient.CreateHistory(ctxClone, &_dto.HistoryDTO{
			ActionType: _enum.HistoryProductChildCreate,
			TargetId:   *body.ParentID,
			OwnerID:    &profileId,
			Title:      "Tạo sản phẩm con thành công",
			Note: []string{"Tạo sản phẩm",
				utils.ConvertToBDS(result.Product.ID),
				"thành công với diện tích",
				fmt.Sprintf("%.2f", result.Product.Area),
				"m2",
			},
		})
	}()

	// s.ProductChildRepo.CreateHistory(c, enums.EProductHistory_NewChild,
	// 	body.ParentID,
	// 	*body.ParentID,
	// 	[]uint64{result.Product.ID},
	// 	[]float64{result.Product.Area},
	// )
	return result, err
}

func (s *ProductUsecase) MergeProductChild(c context.Context, body *dto.MergeProductChildRequest) (*domain.Product, error) {
	_, err := s.RequiredBeforeSave(c, *body.ParentID)
	if err != nil {
		return nil, err
	}

	validate := s.ProductChildRepo.ValidatorMerge(*body.ParentID, *body.ChildIDs)
	if !validate {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Tồn tại sản phẩm con không đạt điều kiện để gộp",
		}
	}

	go func() {
		ctxClone := _utils.CloneContext(c)
		profileId := _utils.GetProfileIdWithContext(ctxClone)
		var arr []string
		for _, r := range *body.ChildIDs {
			arr = append(arr, utils.ConvertToBDS(r))
		}
		childCodes := strings.Join(arr, ", ")
		parentId := uint64(0)
		if body.ParentID != nil {
			parentId = *body.ParentID
		}

		s.NotificationClient.CreateHistory(ctxClone, &_dto.HistoryDTO{
			ActionType: _enum.HistoryProductMerge,
			TargetId:   parentId,
			OwnerID:    &profileId,
			Title:      "Gộp sản phẩm đã chọn thành công",
			Note:       []string{"Gộp sản phẩm", childCodes, " về sản phẩm cha", utils.ConvertToBDS(parentId)},
		})
	}()

	// body.ProductSaveRequest.ParentID = body.ParentID
	return nil, s.ProductChildRepo.MergeChilds(*body.ChildIDs)
}

func (s *ProductUsecase) MergeAllProductChild(c context.Context, parentID *uint64) (*domain.Product, error) {
	_, err := s.RequiredBeforeSave(c, *parentID)
	if err != nil {
		return nil, err
	}

	validate := s.ProductChildRepo.ValidListChild(*parentID)
	if !validate {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Tồn tại sản phẩm con không đạt điều kiện để gộp",
		}
	}

	ids := s.ProductChildRepo.AllChildId(*parentID)
	err = s.ProductChildRepo.MergeAll(*parentID)
	if err != nil {
		return nil, err
	}
	// body.ProductSaveRequest.ParentID = body.ParentID

	go func() {
		ctxClone := _utils.CloneContext(c)
		profileId := _utils.GetProfileIdWithContext(ctxClone)

		var arr []string
		for _, r := range ids {
			arr = append(arr, utils.ConvertToBDS(r))
		}
		childCodes := strings.Join(arr, ", ")

		s.NotificationClient.CreateHistory(ctxClone, &_dto.HistoryDTO{
			ActionType: _enum.HistoryProductMergeAll,
			TargetId:   *parentID,
			OwnerID:    &profileId,
			Title:      "Gộp toàn bộ sản phẩm thành công",
			Note:       []string{"Gộp toàn bộ sản phẩm", childCodes, " về sản phẩm cha", utils.ConvertToBDS(*parentID)},
		})
	}()

	return nil, nil
}

func (s *ProductUsecase) DevideProductChild(c context.Context, devideData *dto.DevideChildRequest) (*[]domain.Product, error) {
	child, err := s.RequiredBeforeSave(c, *devideData.ParentID)
	if err != nil {
		return nil, err
	}

	if child.ParentId == nil {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Chia sản phẩm con chỉ áp dụng đối với sản phẩm con",
		}
	}

	totalDevide := 0.0
	if devideData.Childs != nil {
		for i := 0; i < len(devideData.Childs); i++ {
			totalDevide += devideData.Childs[i].Area
			devideData.Childs[i].ParentID = child.ParentId
		}
	}

	if child.Area <= totalDevide || totalDevide < 1 {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Diện tích không hợp lệ",
		}
	}
	validate := s.ProductChildRepo.ValidListChild(*devideData.ParentID)
	if !validate {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Sản phẩm con không đạt điều kiện để chia",
		}
	}

	// devideu_product_dto.ProductSaveRequest.ParentID = child.ParentId
	s.ProductRepo.UpdateArea(c, child.ID, child.Area-totalDevide)
	var results []domain.Product
	for i := 0; i < len(devideData.Childs); i++ {
		result, err := s.CreateProduct(c, &devideData.Childs[i])
		if err != nil {
			return nil, err
		}
		results = append(results, *result.Product)
	}

	// var ids []uint64
	// var areas []float64

	// for _, r := range results {
	// 	ids = append(ids, r.ID)
	// 	areas = append(areas, r.Area)
	// }

	go func() {
		ctxClone := _utils.CloneContext(c)
		var arr []string
		for i := 0; i < len(results); i++ {
			arr = append(arr, fmt.Sprintf("%s(%.2f m2)", utils.ConvertToBDS(results[i].ID), results[i].Area))
		}

		s.NotificationClient.CreateHistory(ctxClone, &_dto.HistoryDTO{
			Title: "Tách sản phẩm thành công",
			Note: []string{
				"Tách sản phẩm",
				utils.ConvertToBDS(*child.ParentId),
				fmt.Sprintf("(%.2f m2) thành sản phẩm: ", child.Area),
				strings.Join(arr, ", "),
			},
			TargetId:   *child.ParentId,
			ActionType: _enum.HistoryProductDevide,
			// TargetType: shared_enum.EOwnerOf_Product,
		})
	}()

	return &results, nil
}

// func (s *ProductUsecase) ChildHistory(c context.Context, id uint64, dto *dto.ProductHistorySearch) (*[]domain.ProductHistory, int64, error) {
// 	profileId := _utils.GetProfileIdWithContext(c)
// 	result, total, err := s.ProductChildRepo.History(c, profileId, id, dto)
// 	if err != nil {
// 		return nil, 0, err
// 	}
// 	return result, total, nil
// }

func (s *ProductUsecase) QueryProductMarket(c context.Context, dto *dto.ProductSearchRequest) ([]domain.ProductMarket, int64, error) {
	var results []domain.ProductMarket
	profileId := _utils.GetProfileIdWithContext(c)
	results, total, err := s.ProductRepo.SearchProductMarket(profileId, *dto)
	return results, total, err
}
func (s *ProductUsecase) SearchProductMarket(c context.Context, dto *dto.ProductSearchRequest) ([]domain.ProductMarket, int64, error) {
	// var results []u_product_dto.ProductMarket
	// var total int64
	// var err error
	if dto.Text != "" {
		results, total, err := s.QueryProductMarket(c, dto)
		return results, total, err
	}
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

func (s *ProductUsecase) SyncProduct(c context.Context) error {
	err := s.ProductRepo.SyncData(c)
	if err != nil {
		log.Printf("SyncNewsFeed error: %v", err)
		return err
	}
	log.Printf("SyncProduct success")
	return nil
}

func (s *ProductUsecase) SyncData(c context.Context) error {
	return s.ProductRepo.SyncData(c)
}

func (s *ProductUsecase) GetGroupProducts(ctx context.Context, req *dto.GroupProductSearchRequest) ([]domain.Product, int64, error) {
	return s.ProductRepo.GetGroupProducts(ctx, req)
}

func (s *ProductUsecase) GetUserProducts(ctx context.Context, req *dto.UserProductSearchRequest) ([]domain.Product, int64, error) {
	return s.ProductRepo.GetUserProducts(ctx, req)
}

// UpdateNote cập nhật ghi chú sản phẩm
func (s *ProductUsecase) UpdateNote(ctx context.Context, id uint64, note string) error {
	// Check permission
	err := s.CheckPermission(ctx, &id, enums.PermissionProductUpdate)
	if err != nil {
		return err
	}

	// Update note in repository
	return s.ProductRepo.UpdateNote(ctx, id, note)
}

// GetSupportToday lấy ngẫu nhiên 5 sản phẩm cho support today
func (s *ProductUsecase) GetSupportToday(ctx context.Context) ([]domain.Product, error) {
	products, err := s.ProductRepo.GetRandomProducts(ctx, 5)
	if err != nil {
		return nil, err
	}
	return products, nil
}

// GetProductSummary lấy thống kê tổng hợp sản phẩm của user hiện tại với filters
func (s *ProductUsecase) GetProductSummary(ctx context.Context, searchRequest dto.ProductSearchRequest) (*dto.ProductSummaryDTO, error) {
	// Lấy profileID từ context
	profileID := _utils.GetProfileIdWithContext(ctx)
	if profileID == 0 {
		return nil, fmt.Errorf("profile ID not found in context")
	}

	// Gọi repo để lấy thống kê với filters
	summary, err := s.ProductRepo.GetProductSummary(ctx, profileID, searchRequest)
	if err != nil {
		return nil, err
	}

	return summary, nil
}

// GetPriceHistory lấy lịch sử cập nhật giá của sản phẩm
func (s *ProductUsecase) GetPriceHistory(ctx context.Context, productID uint64, page, size uint32) ([]*domain.ProductPrice, int64, error) {
	// Kiểm tra quyền truy cập sản phẩm
	err := s.CheckPermission(ctx, &productID, enums.PermissionProductRead)
	if err != nil {
		return nil, 0, err
	}

	// Lấy lịch sử giá từ repo
	prices, total, err := s.PriceRepo.GetPriceHistory(ctx, productID, page, size)
	if err != nil {
		return nil, 0, err
	}

	return prices, total, nil
}
