package u_product

import (
	"bdspro/v1/data"
	"bdspro/v1/enums"
	a_doc_type "bdspro/v1/mod/admin/doc_type"
	a_property_type "bdspro/v1/mod/admin/property_type"
	"bdspro/v1/mod/global/amenity"
	"bdspro/v1/mod/global/house_info"
	product_accecss "bdspro/v1/mod/global/product_access"
	"bdspro/v1/mod/global/product_media"
	"bdspro/v1/mod/global/product_price"
	"bdspro/v1/mod/global/product_private"
	g_region "bdspro/v1/mod/global/region"
	u_asset "bdspro/v1/mod/user/asset"
	u_post "bdspro/v1/mod/user/post"
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

type ProductService struct {
	productRepo       *UserProductRepo
	postRepo          *u_post.PostRepo
	assetRepo         *u_asset.AssetRepo
	mediaRepo         *product_media.ProductMediaRepo
	productAccessRepo *product_accecss.ProductAccessRepo
	houseInfoRepo     *house_info.HouseInfoRepo
	priceRepo         *product_price.ProductPriceRepo
	productPrivate    *product_private.ProductPrivateRepo
	region            *g_region.RegionRepo
	docTypeRepo       *a_doc_type.DocTypeRepo
	propertyType      *a_property_type.PropertyTypeRepo
	amenity           *amenity.AmenityRepo
}

func NewProductService(productRepo *UserProductRepo,
	postRepo *u_post.PostRepo,
	assetRepo *u_asset.AssetRepo,
	mediaRepo *product_media.ProductMediaRepo,
	productAccessRepo *product_accecss.ProductAccessRepo,
	houseInfoRepo *house_info.HouseInfoRepo,
	priceRepo *product_price.ProductPriceRepo,
	productPrivate *product_private.ProductPrivateRepo,
	region *g_region.RegionRepo,
	docTypeRepo *a_doc_type.DocTypeRepo,
	propertyType *a_property_type.PropertyTypeRepo,
	amenity *amenity.AmenityRepo,
) *ProductService {
	return &ProductService{
		productRepo:       productRepo,
		postRepo:          postRepo,
		assetRepo:         assetRepo,
		productAccessRepo: productAccessRepo,
		mediaRepo:         mediaRepo,
		houseInfoRepo:     houseInfoRepo,
		priceRepo:         priceRepo,
		productPrivate:    productPrivate,
		region:            region,
		docTypeRepo:       docTypeRepo,
		propertyType:      propertyType,
		amenity:           amenity,
	}
}

// func (s *ProductService) UpdateProductStatusToSelling(productId uint64) error {
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
// func (s *ProductService) UpdateProductStatusToNotSelling(productId uint64) error {
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

func (s *ProductService) UpdateProductStatus(ctx *gin.Context,
	productID uint64,
	newStatus enums.ProductStatus,
	buyerID *uint64,
	buyerOrgID *uint64,
) error {
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return err
	}

	switch newStatus {
	case enums.ProductSelling:
		if product.SaleStatus == enums.ProductNotSold {
			// nếu thuộc sở hữu của công ty hoặc cá nhân thì đổi trạng thái tài sản sang đang bán
			if product.OwnerUserID != nil || product.OwnerOrgID != nil {
				if err := s.assetRepo.UpdateStatus(ctx, product.AssetId, enums.AssetStatusSelling); err != nil {
					return err
				}
			}

			if product.SaleVisibility == enums.ProductVisiblePublic {
				if err := s.postRepo.CreateOrUpdatePostVisibility(ctx, productID, enums.PostVisiblePublic); err != nil {
					return err
				}
			} else if product.SaleVisibility == enums.ProductVisibleInternal {
				if err := s.postRepo.CreateOrUpdatePostVisibility(ctx, productID, enums.PostVisibleInternal); err != nil {
					return err
				}
			}
			// nếu sản phẩm là riêng tư thì không tạo tin
		}

	case enums.ProductNotSold:
		if product.SaleStatus == enums.ProductSelling {
			// sản phẩm thuộc công ty -> đổi sang chưa bán
			if product.OwnerOrgID != nil {
				if err := s.assetRepo.UpdateStatus(ctx, product.AssetId, enums.AssetStatusNotSold); err != nil {
					return err
				}
			}
			// công ty không sở hữu tài sản không thay đổi

			post, err := s.postRepo.GetPostByProductID(product.ID)
			if err != nil {
				return err
			}

			// tin đang công khai => tạm ẩn
			if post.Visibility == enums.PostVisiblePublic {
				if err := s.postRepo.UpdatePostStatus(ctx, productID, enums.PostHide); err != nil {
					return err
				}
			}

			// tin nội bộ => dừng
			if post.Visibility == enums.PostVisiblePublic {
				if err := s.postRepo.UpdatePostStatus(ctx, productID, enums.PostStop); err != nil {
					return err
				}
			}
		}

	case enums.ProductSold:
		if buyerID != nil || buyerOrgID != nil {
			// Bán cho công ty hoặc người dùng
			product.OwnerUserID = buyerID
			product.OwnerOrgID = buyerOrgID
			if err := s.assetRepo.UpdateStatus(ctx, product.AssetId, enums.AssetStatusOwning); err != nil {
				return err
			}
		} else {
			// todo
			// Bán ra ngoài
			product.OwnerUserID = nil
			product.OwnerOrgID = nil
		}

		if err := s.postRepo.UpdatePostStatus(ctx, productID, enums.PostStatusCompleted); err != nil {
			return err
		}
		if err := s.postRepo.HideOrDeletePost(ctx, productID); err != nil {
			return err
		}
	}

	return s.productRepo.UpdateStatus(ctx, productID, newStatus)
}

func (s *ProductService) CreateProduct(c *gin.Context, product *data.ProductSaveRequest) (*models.Product, error) {
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

	return entity, nil
}

func (s *ProductService) CheckAccessProduct(c *gin.Context, profileId, productId uint64) (*models.Product, error) {
	existed, err := s.productAccessRepo.ExistsByUserAndProduct(profileId, productId)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	product, err := s.productRepo.GetByID(c, productId)
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

func (s *ProductService) SearchProductWithOwner(c *gin.Context, dto *data.ProductSearchRequest) (_routes.ResponseDTO, error) {
	profileId := _jwt.GetProfileId(c)
	results, total, err := s.productRepo.SearchProductWithOwner(profileId, *dto, _utils.GetPage(c))
	return _routes.ResponseDTO{
		Code:          0,
		Data:          results,
		TotalElements: &total,
	}, err
}

func (s *ProductService) SearchProductMarket(c *gin.Context, dto *data.ProductSearchRequest) (_routes.ResponseDTO, error) {
	profileId := _jwt.GetProfileId(c)
	results, total, err := s.productRepo.SearchProductMarket(profileId, *dto, _utils.GetPage(c))
	return _routes.ResponseDTO{
		Code:          0,
		Data:          results,
		TotalElements: &total,
	}, err
}

func (s *ProductService) UpdateReference(c *gin.Context, entity *models.Product, dto *data.ProductSaveRequest) error {
	if dto.AmenityIds != nil {
		s.productRepo.UpdateProductAmenityIds(entity.ID, dto.AmenityIds)
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

func (s *ProductService) Update(c *gin.Context, id uint64, product *data.ProductSaveRequest) (*models.Product, error) {
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

func (s *ProductService) Detail(c *gin.Context) (*models.Product, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return nil, &_routes.Except{
			Code:    400,
			Message: "ID không hợp lệ",
		}
	}

	entity, err := s.productRepo.GetByID(c, id)
	if err != nil {
		return nil, &_routes.Except{
			Code:    404,
			Message: "Không tìm thấy dữ liệu",
		}
	}

	return entity, nil
}

func (s *ProductService) Rent(c *gin.Context, id uint64, body any) (*models.Product, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return nil, &_routes.Except{
			Code:    400,
			Message: "ID không hợp lệ",
		}
	}

	entity, err := s.productRepo.GetByID(c, id)
	if err != nil {
		return nil, &_routes.Except{
			Code:    404,
			Message: "Không tìm thấy dữ liệu",
		}
	}

	return entity, nil
}

func (s *ProductService) SmartSearch(c *gin.Context, text string) (_routes.ResponseDTO, error) {
	// regionIds, err := s.region.InterText(text, 100)
	p := _utils.GetPage(c)
	searchs, err := utils.InferSearchText(text)

	regionIds, _ := s.region.InterText(text, 10)
	docTypeIds, _ := s.docTypeRepo.InterText(text, 10)
	propertyTypeIds, _ := s.propertyType.InterText(text, 10)

	searchs.RegionIds = regionIds
	searchs.DocTypeIds = docTypeIds
	searchs.PropertyTypeIds = propertyTypeIds

	results, total, err := s.productRepo.SmartSearchProduct(0, searchs, p)
	log.Println("Search= ")

	return _routes.ResponseDTO{
		Code:          0,
		Data:          results,
		TotalElements: &total,
	}, err
}

func (s *ProductService) UpdatePrice(c *gin.Context, priceID uint64, entity *models.ProductPrice) error {
	return s.priceRepo.UpdatePrice(c, priceID, entity)
}

func (s *ProductService) UpdateAmenity(productID uint64, ids []uint64) error {
	if len(ids) > 0 {
		return s.productRepo.UpdateProductAmenityIds(productID, ids)
	}
	return nil
}

func (s *ProductService) UpdateMediaItems(productID uint64, items []data.MediaItem) error {
	return s.mediaRepo.UpdateProductMediaItems(productID, items)
}

func (s *ProductService) UpdateHouseInfo(c *gin.Context, productId uint64, houseInfo *models.HouseInfo) error {
	if houseInfo != nil {
		// copier.Copy(&houseInfo, &dto.HouseInfo)
		houseInfo.ProductID = productId
		return s.houseInfoRepo.UpdateHouseInfo(c, houseInfo)
	}
	return nil
}

func (s *ProductService) UpdatePrivateData(c *gin.Context, productID uint64, privateInfo *models.ProductPrivate) error {
	return s.productPrivate.UpdateInfo(c, productID, privateInfo)
}

func (s *ProductService) InferenceString(text string) string {
	searchs, _ := utils.InferSearchText(text)

	regionIds, _ := s.region.InterText(text, 10)
	docTypeIds, _ := s.docTypeRepo.InterText(text, 10)
	propertyTypeIds, _ := s.propertyType.InterText(text, 10)

	searchs.RegionIds = regionIds
	searchs.DocTypeIds = docTypeIds
	searchs.PropertyTypeIds = propertyTypeIds

	return ""
}

func (s *ProductService) CreateProductWithDetect(c *gin.Context, body *data.ProductSaveDetect) (any, error) {
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
