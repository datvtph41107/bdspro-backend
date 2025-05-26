package u_product

import (
	"bdspro/internal/enums"
	"bdspro/v1/data"
	a_doc_type "bdspro/v1/mod/admin/doc_type"
	a_property_type "bdspro/v1/mod/admin/property_type"
	"bdspro/v1/mod/global/amenity"
	"bdspro/v1/mod/global/house_info"
	product_accecss "bdspro/v1/mod/global/product_access"
	"bdspro/v1/mod/global/product_media"
	"bdspro/v1/mod/global/product_price"
	"bdspro/v1/mod/global/product_private"
	g_region "bdspro/v1/mod/global/region"
	"bdspro/v1/providers"

	// u_asset "bdspro/v1/mod/user/asset"

	iu_asset "bdspro/v1/mod/user/asset/interface"
	user_dto "bdspro/v1/mod/user/dto/request"
	u_post_dto "bdspro/v1/mod/user/post/dto"
	iu_post "bdspro/v1/mod/user/post/interface"
	u_product_dto "bdspro/v1/mod/user/product/dto"
	"bdspro/v1/models"
	"bdspro/v1/utils"
	_db "common/db"
	_jwt "common/jwt"
	_routes "common/routes"
	_utils "common/utils"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

// @bind: bdspro/v1/mod/user/product/interface.IUserProductService
type UserProductService struct {
	ProductRepo *UserProductRepo
	// postRepo          *u_post.PostRepo
	// assetRepo         *u_asset.AssetRepo
	mediaRepo         *product_media.ProductMediaRepo
	productAccessRepo *product_accecss.ProductAccessRepo
	houseInfoRepo     *house_info.HouseInfoRepo
	priceRepo         *product_price.ProductPriceRepo
	productPrivate    *product_private.ProductPrivateRepo
	region            *g_region.RegionRepo
	docTypeRepo       *a_doc_type.DocTypeRepo
	propertyType      *a_property_type.PropertyTypeRepo
	amenity           *amenity.AmenityRepo
	memberPlan        *providers.MemberPlanProvider
	UserPostService   iu_post.IUserPostService
	UserAssetService  iu_asset.IUserAssetService

	CacheProductMarket *providers.ListCacheProvider[u_product_dto.ProductMarket, data.ProductSearchRequest]
}

func NewUserProductService(productRepo *UserProductRepo,
	// postRepo *u_post.PostRepo,
	// assetRepo *u_asset.AssetRepo,
	mediaRepo *product_media.ProductMediaRepo,
	productAccessRepo *product_accecss.ProductAccessRepo,
	houseInfoRepo *house_info.HouseInfoRepo,
	priceRepo *product_price.ProductPriceRepo,
	productPrivate *product_private.ProductPrivateRepo,
	region *g_region.RegionRepo,
	docTypeRepo *a_doc_type.DocTypeRepo,
	propertyType *a_property_type.PropertyTypeRepo,
	amenity *amenity.AmenityRepo,
	memberPlan *providers.MemberPlanProvider,
	// postProvider *u_product_provider.UserProductProvider,
) *UserProductService {
	service := &UserProductService{
		ProductRepo: productRepo,
		// postRepo:          postRepo,
		// assetRepo:         assetRepo,
		productAccessRepo:  productAccessRepo,
		mediaRepo:          mediaRepo,
		houseInfoRepo:      houseInfoRepo,
		priceRepo:          priceRepo,
		productPrivate:     productPrivate,
		region:             region,
		docTypeRepo:        docTypeRepo,
		propertyType:       propertyType,
		amenity:            amenity,
		memberPlan:         memberPlan,
		CacheProductMarket: &providers.ListCacheProvider[u_product_dto.ProductMarket, data.ProductSearchRequest]{},
		// postProvider:      postProvider,
	}
	// CacheProvider := &providers.CacheProvider[u_product_dto.ProductMarket]{}
	service.CacheProductMarket.Setup()
	return service
}

// func (s *UserProductService) UpdateProductStatusToSelling(productId uint64) error {
// 	product, err := s.productRepo.GetByID(productId)
// 	if err != nil {
// 		return err
// 	}

// 	// Nếu chưa thuộc sở hữu thì không cập nhật tài sản
// 	if product.OwnerUserID == nil && product.OwnerOrgID == nil {
// 		return errors.New("product chưa thuộc sở hữu, không thể cập nhật")
// 	}

// 	// Chuyển trạng thái sản phẩm sang "Đang bán"
// 	product.SaleStatus = enums.ProductSelling

// 	// Kiểm tra visibility và tạo tin đăng
// 	// if product.SaleVisibility == enums.ProductS {
// 	// 	post := &models.Post{
// 	// 		ProductId: product.ID,
// 	// 		Status:    enums.PostPublic, // Giả sử 1 là trạng thái "Đang hiển thị"
// 	// 	}
// 	// 	if err := s.postRepo.CreatePost(post); err != nil {
// 	// 		return err
// 	// 	}
// 	// }

// 	// Cập nhật sản phẩm
// 	return s.productRepo.UpdateProduct(product)
// }

// Chuyển trạng thái sản phẩm về "Chưa bán"
// func (s *UserProductService) UpdateProductStatusToNotSelling(productId uint64) error {
// 	product, err := s.productRepo.GetByID(productId)
// 	if err != nil {
// 		return err
// 	}

// 	if product.OwnerOrgID != nil {
// 		// Nếu công ty sở hữu SP, cập nhật trạng thái
// 		product.SaleStatus = 1 // Chưa bán
// 		if err := s.productRepo.UpdateProduct(product); err != nil {
// 			return err
// 		}
// 	}

// 	// Hủy hoặc ẩn tin đăng
// 	post, err := s.postRepo.GetPostByProductID(product.ID)
// 	if err == nil && post != nil {
// 		post.Status = 3 // Giả sử 3 là "Tạm ẩn"
// 		if err := s.postRepo.UpdatePost(post); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (s *UserProductService) UpdateProductStatus(ctx *gin.Context,
// 	productID uint64,
// 	newStatus enums.ProductStatus,
// 	buyerID *uint64,
// 	buyerOrgID *uint64,
// ) error {
// 	product, err := s.productRepo.GetByID(ctx, productID)
// 	if err != nil {
// 		return err
// 	}

// 	switch newStatus {
// 	case enums.ProductSelling:
// 		if product.SaleStatus == enums.ProductNotSold {
// 			// nếu thuộc sở hữu của công ty hoặc cá nhân thì đổi trạng thái tài sản sang đang bán
// 			if product.OwnerUserID != nil || product.OwnerOrgID != nil {
// 				if err := s.assetRepo.UpdateStatus(ctx, product.AssetId, enums.AssetStatusSelling); err != nil {
// 					return err
// 				}
// 			}

// 			if product.SaleVisibility == enums.ProductVisiblePublic {
// 				if err := s.postRepo.CreateOrUpdatePostVisibility(ctx, productID, enums.PostVisiblePublic); err != nil {
// 					return err
// 				}
// 			} else if product.SaleVisibility == enums.ProductVisibleInternal {
// 				if err := s.postRepo.CreateOrUpdatePostVisibility(ctx, productID, enums.PostVisibleInternal); err != nil {
// 					return err
// 				}
// 			}
// 			// nếu sản phẩm là riêng tư thì không tạo tin
// 		}

// 	case enums.ProductNotSold:
// 		if product.SaleStatus == enums.ProductSelling {
// 			// sản phẩm thuộc công ty -> đổi sang chưa bán
// 			if product.OwnerOrgID != nil {
// 				if err := s.assetRepo.UpdateStatus(ctx, product.AssetId, enums.AssetStatusNotSold); err != nil {
// 					return err
// 				}
// 			}
// 			// công ty không sở hữu tài sản không thay đổi

// 			post, err := s.postRepo.GetPostByProductID(product.ID)
// 			if err != nil {
// 				return err
// 			}

// 			// tin đang công khai => tạm ẩn
// 			if post.Visibility == enums.PostVisiblePublic {
// 				if err := s.postRepo.UpdatePostStatus(ctx, productID, enums.PostHide); err != nil {
// 					return err
// 				}
// 			}

// 			// tin nội bộ => dừng
// 			if post.Visibility == enums.PostVisiblePublic {
// 				if err := s.postRepo.UpdatePostStatus(ctx, productID, enums.PostStop); err != nil {
// 					return err
// 				}
// 			}
// 		}

// 	case enums.ProductSold:
// 		if buyerID != nil || buyerOrgID != nil {
// 			// Bán cho công ty hoặc người dùng
// 			product.OwnerUserID = buyerID
// 			product.OwnerOrgID = buyerOrgID
// 			if err := s.assetRepo.UpdateStatus(ctx, product.AssetId, enums.AssetStatusOwning); err != nil {
// 				return err
// 			}
// 		} else {
// 			// todo
// 			// Bán ra ngoài
// 			product.OwnerUserID = nil
// 			product.OwnerOrgID = nil
// 		}

// 		if err := s.postRepo.UpdatePostStatus(ctx, productID, enums.PostStatusCompleted); err != nil {
// 			return err
// 		}
// 		if err := s.postRepo.HideOrDeletePost(ctx, productID); err != nil {
// 			return err
// 		}
// 	}

// 	return s.productRepo.UpdateSaleStatus(ctx, productID, newStatus)
// }

func (s *UserProductService) CreateProduct(c *gin.Context, product *u_product_dto.ProductSaveRequest) (*u_product_dto.CreateResponse, error) {
	profileId := _jwt.GetProfileId(c)

	entity := &models.Product{}
	houseInfoEntity := &models.HouseInfo{}
	copier.Copy(&entity, &product)
	copier.Copy(&houseInfoEntity, &product.HouseInfoDTO)

	if err := s.UpdatePrice(c, 0, product.PriceDTO); err != nil {
		return nil, err
	}

	entity.LastPriceID = &product.PriceDTO.ID
	entity.OwnerUserID = &profileId
	_db.SaveWithAudit(c, entity)
	// err := s.productRepo.UpdateProduct(c, entity.ID, entity)
	// err := s.UpdateReference(c, entity, product)

	if err := s.priceRepo.UpdateProductID(c, &entity.ID, &product.PriceDTO.ID); err != nil {
		return nil, err
	}
	if err := s.UpdateAmenity(entity.ID, product.AmenityIds); err != nil {
		return nil, err
	}
	if err := s.UpdateMediaItems(entity.ID, product.MediaItemDTOs); err != nil {
		return nil, err
	}
	if err := s.UpdateHouseInfo(c, entity.ID, houseInfoEntity); err != nil {
		return nil, err
	}
	if err := s.UpdatePrivateData(c, entity.ID, product.PrivateDataDTO); err != nil {
		return nil, err
	}

	// access := &models.ProductAccess{
	// 	UserId:    profileId,
	// 	ProductID: entity.ID,
	// }
	// err := s.productAccessRepo.Create(c, access)
	// if err != nil {
	// 	return nil, err
	// }

	// ProductID       uint64                  `json:"productId" binding:"required"`
	// Title           string                  `json:"title"`
	// Content         string                  `json:"content" binding:"required"`
	// Price           float64                 `json:"postPrice"`
	// TransactionType uint                    `json:"transactionType" binding:"required"`
	// Visibility      enums.PostVisibility    `json:"visibility" binding:"required"`
	// PackageVisible  uint                    `json:"packageVisible"`
	// ExpiredAt       *time.Time              `json:"expiredAt"`
	// NumDate         int                     `json:"numDate"`
	// Hidden          bool                    `json:"hidden"`
	// Metadatas       models.ProductMediaItem `json:"metadatas"`
	postInfo := product.PostCreate
	var postEntity *models.Post
	var err error
	if postInfo != nil {
		postInfo.ProductID = entity.ID
		postInfo.Content = entity.Description
		postInfo.TransactionType = entity.TransactionType
		if entity.TransactionType == 10 {
			postInfo.Price = *product.PriceDTO.SalePrice
		} else {
			postInfo.Price = *product.PriceDTO.RentPrice
		}

		dto := u_post_dto.PostSaveRequest{}
		copier.Copy(&dto, &postInfo)

		postEntity, err = s.UserPostService.CreatePost(c, &dto)
		if err != nil {
			return nil, err
		}
	}

	assetInfo := product.AssetCreate
	var assetEntity *models.Asset
	if assetInfo != nil {
		// copier.Copy(&assetEntity, &assetInfo)
		assetEntity = &models.Asset{
			ProductID:      &entity.ID,
			Name:           entity.Name,
			Area:           entity.Area,
			PropertyTypeId: entity.PropertyTypeId,
			Description:    entity.Description,
			ProvinceID:     entity.ProvinceID,
			DistrictID:     entity.DistrictID,
			WardID:         entity.WardID,
			PurchasePrice:  assetInfo.PurchasePrice,
			PurchaseDate:   assetInfo.PurchaseDate,
			LegalStatus:    assetInfo.LegalStatus,
			RentStatus:     assetInfo.RentStatus,
			// Description:    assetInfo.Description,
		}

		assetEntity, err = s.UserAssetService.CreateAsset(c, assetEntity)
		if err != nil {
			return nil, err
		}
	}

	return &u_product_dto.CreateResponse{
		Product: entity,
		Post:    postEntity,
		Asset:   assetEntity,
	}, nil
}

func (s *UserProductService) CheckAccessProduct(c *gin.Context, profileId, productId uint64) (*models.Product, error) {
	existed, err := s.productAccessRepo.ExistsByUserAndProduct(profileId, productId)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	product, err := s.ProductRepo.GetByID(c, productId)
	if err != nil {
		return nil, err
	}

	if !existed && *product.OwnerUserID != profileId {
		return nil, &_routes.Except{
			Code:    401,
			Message: "Bạn không có quyền truy cập",
		}
	}

	return product, nil
}

func (s *UserProductService) SearchProductWithOwner(c *gin.Context, dto *data.ProductSearchRequest) (_routes.ResponseDTO, error) {
	profileId := _jwt.GetProfileId(c)
	results, total, err := s.ProductRepo.SearchProductWithOwner(profileId, dto)
	return _routes.ResponseDTO{
		Code:          0,
		Data:          results,
		TotalElements: &total,
	}, err
}

// var (
// 	productCache      []u_product_dto.ProductMarket
// 	totalMarlket      int64
// 	productCacheMutex sync.RWMutex
// )

func (s *UserProductService) QueryProductMarket(c *gin.Context, dto *data.ProductSearchRequest) ([]u_product_dto.ProductMarket, int64, error) {
	var results []u_product_dto.ProductMarket
	profileId := _jwt.GetProfileId(c)
	results, total, err := s.ProductRepo.SearchProductMarket(profileId, *dto, _utils.GetPage(c))
	return results, total, err
}
func (s *UserProductService) SearchProductMarket(c *gin.Context, dto *data.ProductSearchRequest) (_routes.ResponseDTO, error) {
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

	return _routes.ResponseDTO{
		Code:          0,
		Data:          results,
		TotalElements: &total,
	}, err
}

func (s *UserProductService) UpdateReference(c *gin.Context, entity *models.Product, dto *u_product_dto.ProductSaveRequest) error {
	if dto.AmenityIds != nil {
		s.ProductRepo.UpdateProductAmenityIds(entity.ID, dto.AmenityIds)
	}

	// price := entity.Price
	privateInfo := entity.PrivateData
	houseInfo := entity.HouseInfo
	// entity.Price = nil
	// entity.PrivateData = nil
	// entity.HouseInfo = nil

	// err := _db.DB.WithContext(c).Model(entity).
	// 	Where("id = ? and deleted_at is null", id).
	// 	Updates(entity).Error

	s.mediaRepo.UpdateProductMediaItems(entity.ID, dto.MediaItemDTOs)
	// var houseInfo models.HouseInfo
	if dto.HouseInfoDTO != nil {
		// copier.Copy(&houseInfo, &dto.HouseInfo)
		houseInfo.ProductID = entity.ID
		s.houseInfoRepo.UpdateHouseInfo(c, houseInfo)
	}

	err := s.productPrivate.UpdateInfo(c, entity.ID, privateInfo)

	return err
}

func (s *UserProductService) Update(c *gin.Context, id uint64, product *u_product_dto.ProductSaveRequest) (*models.Product, error) {
	profileId := _jwt.GetProfileId(c)
	entity, err := s.CheckAccessProduct(c, profileId, id)
	if err != nil {
		return nil, err
	}

	copier.Copy(&entity, &product)

	// err = s.productRepo.UpdateProduct(c, id, product)
	// s.UpdateReference(c, product, dto)
	product.PriceDTO.ProductID = &id
	if err := s.UpdatePrice(c, *entity.LastPriceID, product.PriceDTO); err != nil {
		return nil, err
	}
	if err := s.UpdateAmenity(entity.ID, product.AmenityIds); err != nil {
		return nil, err
	}
	if err := s.UpdateMediaItems(entity.ID, product.MediaItemDTOs); err != nil {
		return nil, err
	}

	houseInfoEntity := &models.HouseInfo{}
	copier.Copy(&houseInfoEntity, &product.HouseInfoDTO)
	if err := s.UpdateHouseInfo(c, entity.ID, houseInfoEntity); err != nil {
		return nil, err
	}
	if err := s.UpdatePrivateData(c, entity.ID, product.PrivateDataDTO); err != nil {
		return nil, err
	}
	entity.LastPriceID = &product.PriceDTO.ID
	err = _db.SaveWithAudit(c, &entity)

	return entity, err
}

func (s *UserProductService) Detail(c *gin.Context) (*models.Product, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return nil, &_routes.Except{
			Code:    400,
			Message: "ID không hợp lệ",
		}
	}

	entity, err := s.ProductRepo.GetByID(c, id)
	if err != nil {
		return nil, &_routes.Except{
			Code:    404,
			Message: "Không tìm thấy dữ liệu",
		}
	}

	entity.Asset, _ = s.UserAssetService.GetAssetByProductID(c, &id)

	return entity, nil
}

func (s *UserProductService) Posts(c *gin.Context, id uint64, dto user_dto.PostLinkSearch) (*[]models.PostItem, error) {
	results, err := s.UserPostService.ProductPostLink(c, id, &dto)
	return results, err
}

func (s *UserProductService) Members(c *gin.Context, id uint64, dto user_dto.ProductAccessSearch) (*[]models.ProductAccessQuery, error) {
	// _, err := s.RequiredOwnerProduct(c, id)
	// if err != nil {
	// 	return nil, err
	// }

	entities, err := s.productAccessRepo.OwnerSearchProduct(c, id, &dto)
	return entities, err
}

func (s *UserProductService) History(c *gin.Context) (*models.Product, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return nil, &_routes.Except{
			Code:    400,
			Message: "ID không hợp lệ",
		}
	}

	entity, err := s.ProductRepo.GetByID(c, id)
	if err != nil {
		return nil, &_routes.Except{
			Code:    404,
			Message: "Không tìm thấy dữ liệu",
		}
	}

	return entity, nil
}

func (s *UserProductService) Rent(c *gin.Context, id uint64, body any) (*models.Product, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return nil, &_routes.Except{
			Code:    400,
			Message: "ID không hợp lệ",
		}
	}

	entity, err := s.ProductRepo.GetByID(c, id)
	if err != nil {
		return nil, &_routes.Except{
			Code:    404,
			Message: "Không tìm thấy dữ liệu",
		}
	}

	return entity, nil
}

func (s *UserProductService) SmartSearch(c *gin.Context, text string) (_routes.ResponseDTO, error) {
	// regionIds, err := s.region.InterText(text, 100)
	// p := _utils.GetPage(c)
	searchs, err := utils.InferSearchText(text)

	regionIds, _ := s.region.InterText(text, 10)
	docTypeIds, _ := s.docTypeRepo.InterText(text, 10)
	propertyTypeIds, _ := s.propertyType.InterText(text, 10)

	searchs.RegionIds = regionIds
	searchs.DocTypeIds = docTypeIds
	searchs.PropertyTypeIds = propertyTypeIds
	// searchs.Page = p

	results, total, err := s.ProductRepo.SmartSearchProduct(0, searchs)
	log.Println("Search= ")

	return _routes.ResponseDTO{
		Code:          0,
		Data:          results,
		TotalElements: &total,
	}, err
}

func (s *UserProductService) UpdatePrice(c *gin.Context, priceID uint64, entity *models.ProductPrice) error {
	return s.priceRepo.UpdatePrice(c, priceID, entity)
}

func (s *UserProductService) UpdateAmenity(productID uint64, ids []uint64) error {
	if len(ids) > 0 {
		return s.ProductRepo.UpdateProductAmenityIds(productID, ids)
	}
	return nil
}

func (s *UserProductService) UpdateMediaItems(productID uint64, items []u_product_dto.MediaItem) error {
	return s.mediaRepo.UpdateProductMediaItems(productID, items)
}

func (s *UserProductService) UpdateHouseInfo(c *gin.Context, productId uint64, houseInfo *models.HouseInfo) error {
	if houseInfo != nil {
		// copier.Copy(&houseInfo, &dto.HouseInfo)
		houseInfo.ProductID = productId
		return s.houseInfoRepo.UpdateHouseInfo(c, houseInfo)
	}
	return nil
}

func (s *UserProductService) UpdatePrivateData(c *gin.Context, productID uint64, privateInfo *models.ProductPrivate) error {

	return s.productPrivate.UpdateInfo(c, productID, privateInfo)
}

func (s *UserProductService) InferenceString(text string) string {
	searchs, _ := utils.InferSearchText(text)

	regionIds, _ := s.region.InterText(text, 10)
	docTypeIds, _ := s.docTypeRepo.InterText(text, 10)
	propertyTypeIds, _ := s.propertyType.InterText(text, 10)

	searchs.RegionIds = regionIds
	searchs.DocTypeIds = docTypeIds
	searchs.PropertyTypeIds = propertyTypeIds

	return ""
}

func (s *UserProductService) CreateProductWithDetect(c *gin.Context, body *u_product_dto.ProductSaveDetect) (any, error) {
	// searchs, err := utils.InferSearchText(body.Content)
	// if err != nil {
	// 	return nil, err
	// }
	// if searchs, err := utils.InferSearchText(body.Content); err != nil {
	// 	return nil, err
	// }
	var searchs data.TotalSearchParser
	var err error
	if searchs, err = utils.InferSearchText(body.Content); err != nil {
		return nil, err
	}

	text := body.Content

	// searchs, _ := utils.InferSearchText(text)

	regions, _ := s.region.InterTextToItem(text, 10)
	docTypes, _ := s.docTypeRepo.InterTextToItem(text, 10)
	propertyTypes, _ := s.propertyType.InterTextToItem(text, 10)
	amenities, _ := s.amenity.InterTextToEntity(text, 10)

	searchs.Name = _utils.TrimFirstWords(text, 12)

	// tìm địa chỉ
	var ward, district, province *models.Region
	for i := 0; i < len(regions); i++ {
		r := regions[i]
		if r.Level == 1 && province == nil {
			province = &r
			continue
		}
		if r.Level == 2 && district == nil {
			district = &r
			continue
		}
		if r.Level == 3 && ward == nil {
			ward = &r
			continue
		}
	}
	if ward != nil && (district == nil || district.ID != ward.ParentID) {
		district, _ = s.region.GetByID(c, ward.ParentID)
		province, _ = s.region.GetByID(c, district.ParentID)
	}

	if district != nil && (province == nil || province.ID != district.ParentID) {
		province, _ = s.region.GetByID(c, district.ParentID)
	}
	// ---- xong ----

	searchs.DocTypes = docTypes
	searchs.PropertyTypes = propertyTypes
	searchs.Amenities = amenities
	searchs.Ward = ward
	searchs.District = district
	searchs.Province = province
	if searchs.Price != nil {
		if *searchs.Price > 100_000_000 {
			searchs.TransactionType = 1
			searchs.SalePrice = searchs.Price
		} else {
			searchs.TransactionType = 2
			searchs.RentPrice = searchs.Price
		}
	}

	return searchs, nil
}

func (s *UserProductService) Archived(c *gin.Context, dto u_product_dto.ArchivedRequest) error {
	if _, err := s.RequiredOwnerProduct(c, *dto.ProductId); err != nil {
		return err
	}

	return s.ProductRepo.UpdateArchived(c, dto)
}

func (s *UserProductService) Delete(c *gin.Context, id uint64) error {
	// sản phẩm đang cọc
	product, err := s.RequiredOwnerProduct(c, id)
	if err != nil {
		return err
	}

	// todo: thêm hoặc có tin active
	if product.SaleStatus != enums.ProductNotSold || product.RentStatus != enums.ProductNotRent {
		return &_routes.Except{
			Code:    401,
			Message: "Không thể xóa sản phẩm này",
		}
	}
	return s.ProductRepo.Delete(c, id)
}

func (s *UserProductService) RequiredOwnerProduct(c *gin.Context, id uint64) (*models.Product, error) {
	product, err := s.ProductRepo.GetByID(c, id)
	if err != nil {
		return nil, err
	}

	profileId := _jwt.GetProfileId(c)

	if *product.OwnerUserID != profileId {
		return nil, &_routes.Except{
			Code:    401,
			Message: "Bạn không có quyền truy cập",
		}
	}
	return product, nil
}

func (s *UserProductService) RentStatusUpdate(c *gin.Context, dto u_product_dto.StatusUpdate) error {
	// sản phẩm đang cọc
	_, err := s.RequiredOwnerProduct(c, *dto.ProductID)
	if err != nil {
		return err
	}
	return s.ProductRepo.UpdateRentStatus(c, *dto.ProductID, dto.Status)
}

func (s *UserProductService) SaleStatusUpdate(c *gin.Context, dto u_product_dto.StatusUpdate) error {
	// sản phẩm đang cọc
	_, err := s.RequiredOwnerProduct(c, *dto.ProductID)
	if err != nil {
		return err
	}
	return s.ProductRepo.UpdateSaleStatus(c, *dto.ProductID, dto.Status)
}

func (s *UserProductService) RentVisibilityUpdate(c *gin.Context, dto u_product_dto.VisibilityUpdate) error {
	// sản phẩm đang cọc
	_, err := s.RequiredOwnerProduct(c, *dto.ProductID)
	if err != nil {
		return err
	}
	return s.ProductRepo.UpdateRentVisibility(c, *dto.ProductID, dto.Visibility)
}

func (s *UserProductService) SaleVisibilityUpdate(c *gin.Context, dto u_product_dto.VisibilityUpdate) error {
	// sản phẩm đang cọc
	_, err := s.RequiredOwnerProduct(c, *dto.ProductID)
	if err != nil {
		return err
	}
	return s.ProductRepo.UpdateSaleVisibility(c, *dto.ProductID, dto.Visibility)
}

func (s *UserProductService) ProductDeposite(c *gin.Context, dto u_product_dto.DepositeRequest) error {
	// sản phẩm đang cọc
	product, err := s.RequiredOwnerProduct(c, dto.ProductID)
	if err != nil {
		return err
	}

	if product.DepositeID != nil {
		return &_routes.Except{
			Code:    400,
			Message: "Sản phẩm đã đặt cọc",
		}
	}

	record := models.Deposite{
		ProductID:     dto.ProductID,
		CustomerName:  dto.CustomerName,
		CustomerPhone: dto.CustomerPhone,
		Amount:        dto.Amount,
		Note:          dto.Note,
	}
	if err := s.ProductRepo.NewDeposite(c, &record); err != nil {
		return err
	}
	product.DepositeID = &record.ID

	return s.ProductRepo.Update(c, product.ID, product)
}

func (s *UserProductService) CancelDeposite(c *gin.Context, id uint64) error {
	// sản phẩm đang cọc
	product, err := s.RequiredOwnerProduct(c, id)
	if err != nil {
		return err
	}

	if product.DepositeID == nil {
		return &_routes.Except{
			Code:    400,
			Message: "Sản phẩm chưa đặt cọc",
		}
	}

	product.DepositeID = nil

	return s.ProductRepo.Update(c, product.ID, product)
}
