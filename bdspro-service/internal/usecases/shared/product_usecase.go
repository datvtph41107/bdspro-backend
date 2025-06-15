package shared_usecase

import (
	"bdspro/infra/providers"
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"bdspro/internal/repo"
	"bdspro/internal/usecases"
	"bdspro/internal/utils"
	_db "common/db"
	_dto "common/dto"
	_routes "common/routes"
	_utils "common/utils"
	"context"

	"github.com/jinzhu/copier"
)

type IProductUsecase interface {
	GetOwnerID(c context.Context) (*uint64, error)
	HasPermission(c context.Context, productId *uint64, permission enums.EPermission) (bool, error)
	GetOwnerType(c context.Context) enums.EOwnerType
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
	// UAssetUsecase      *UAssetUsecase
	DepositeRepo       repo.DepositeRepo
	OwnerType          enums.EOwnerType
	CacheProductMarket *providers.ListCacheProvider[domain.ProductMarket, dto.ProductSearchRequest]
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
	// UAssetUsecase *UAssetUsecase,
	depositeRepo repo.DepositeRepo,
	productChildRepo repo.UProductChildRepo,
) *ProductUsecase {
	return &ProductUsecase{
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
		// UAssetUsecase:      UAssetUsecase,
		DepositeRepo:       depositeRepo,
		ProductChildRepo:   productChildRepo,
		CacheProductMarket: &providers.ListCacheProvider[domain.ProductMarket, dto.ProductSearchRequest]{},
		// OwnerType:    ownerType,
	}
}

func MakeProductUsecase(
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
	// UAssetUsecase *UAssetUsecase,
	depositeRepo repo.DepositeRepo,
	ownerType enums.EOwnerType,
) *ProductUsecase {
	return &ProductUsecase{
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
		// UAssetUsecase:      UAssetUsecase,
		DepositeRepo: depositeRepo,
		OwnerType:    ownerType,
	}
}

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

func (s *ProductUsecase) GetOwnerID(c context.Context, ownerType enums.EOwnerType) uint64 {
	// if ownerType == enums.EOwnerTypeOrganization {
	// 	return s.This.GetOwnerID(c)
	// }

	ownerId := _utils.GetProfileIdWithContext(c)
	return ownerId
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

	if err := s.UpdatePrice(c, 0, product.PriceData); err != nil {
		return nil, err
	}
	if product.PriceData != nil {
		entity.LastPriceID = &product.PriceData.ID
	}

	entity.OwnerUserID = &ownerId
	entity.OwnerType = product.OwnerType
	_db.SaveWithContext(c, entity)
	// err := s.productRepo.UpdateProduct(c, entity.ID, entity)
	// err := s.UpdateReference(c, entity, product)

	if product.PriceData != nil {
		if err := s.PriceRepo.UpdateProductID(c, &entity.ID, &product.PriceData.ID); err != nil {
			return nil, err
		}
	}
	if err := s.UpdateAmenity(entity.ID, product.AmenityIds); err != nil {
		return nil, err
	}
	if err := s.UpdateMediaItems(entity.ID, product.MediaItems); err != nil {
		return nil, err
	}
	if err := s.UpdateHouseInfo(c, entity.ID, houseInfoEntity); err != nil {
		return nil, err
	}
	if err := s.UpdatePrivateData(c, entity.ID, product.PrivateData); err != nil {
		return nil, err
	}

	postInfo := product.PostCreate
	var postEntity *domain.Post
	if postInfo != nil && s.PostUsecase != nil {
		postInfo.ProductID = entity.ID
		postInfo.Content = entity.Description
		postInfo.TransactionType = entity.TransactionType
		if entity.TransactionType == 10 {
			postInfo.Price = *product.PriceData.SalePrice
		} else {
			postInfo.Price = *product.PriceData.RentPrice
		}

		dto := dto.PostSaveRequest{}
		copier.Copy(&dto, &postInfo)

		postEntity, err = s.PostUsecase.CreatePost(c, &dto)
		if err != nil {
			return nil, err
		}
	}

	assetInfo := product.AssetCreate
	var assetEntity *domain.Asset
	if assetInfo != nil {
		// copier.Copy(&assetEntity, &assetInfo)
		assetEntity = &domain.Asset{
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

		// assetEntity, err = s.UAssetUsecase.CreateAsset(c, &user_dto.AssetSaveRequest{
		// 	Asset:      *assetEntity,
		// 	LegalItems: assetInfo.LegalItems,
		// })
		if err != nil {
			return nil, err
		}
	}

	return &dto.CreateResponse{
		Product: entity,
		Post:    postEntity,
		Asset:   assetEntity,
	}, nil
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
	entity, err := s.ProductRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	copier.Copy(&entity, &product)

	// err = s.productRepo.UpdateProduct(c, id, product)
	// s.UpdateReference(c, product, dto)
	if product.PriceData != nil {
		product.PriceData.ProductID = &id
	}
	if err := s.UpdatePrice(c, *entity.LastPriceID, product.PriceData); err != nil {
		return nil, err
	}
	if err := s.UpdateAmenity(entity.ID, product.AmenityIds); err != nil {
		return nil, err
	}
	if err := s.UpdateMediaItems(entity.ID, product.MediaItems); err != nil {
		return nil, err
	}

	houseInfoEntity := &domain.HouseInfo{}
	copier.Copy(&houseInfoEntity, &product.HouseInfo)
	if err := s.UpdateHouseInfo(c, entity.ID, houseInfoEntity); err != nil {
		return nil, err
	}
	if err := s.UpdatePrivateData(c, entity.ID, product.PrivateData); err != nil {
		return nil, err
	}
	entity.LastPriceID = &product.PriceData.ID
	err = _db.SaveWithContext(c, &entity)

	return entity, err
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

func (s *ProductUsecase) SearchProductByOwner(c context.Context, dto *dto.ProductSearchRequest) ([]domain.Product, int64, error) {
	// profileId := _jwt.GetProfileId(c)
	ownerID := s.GetOwnerID(c, enums.EOwnerTypeMember)
	// if err != nil {
	// 	return nil, 0, err
	// }
	return s.ProductRepo.SearchProductWithOwner(&ownerID, *dto)
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

	return entity, nil
}

func (s *ProductUsecase) SmartSearch(c context.Context, text string, dto *dto.TextSearchRequest) ([]domain.Product, int64, error) {
	searchs, err := utils.InferSearchText(text)

	regionIds, _ := s.Region.InterText(text, 10)
	docTypeIds, _ := s.DocTypeRepo.InterText(text, 10)
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
	return s.ProductRepo.SmartSearch(ownerID, s.This.GetOwnerType(c), &searchs)
}

func (s *ProductUsecase) UpdatePrice(c context.Context, priceID uint64, entity *domain.ProductPrice) error {
	return s.PriceRepo.UpdatePrice(c, priceID, entity)
}

func (s *ProductUsecase) UpdateAmenity(productID uint64, ids []uint64) error {
	if len(ids) > 0 {
		return s.ProductRepo.UpdateProductAmenityIds(productID, ids)
	}
	return nil
}

func (s *ProductUsecase) UpdateMediaItems(productID uint64, items []dto.MediaItem) error {
	return s.MediaRepo.UpdateProductMediaItems(productID, items)
}

func (s *ProductUsecase) UpdateHouseInfo(c context.Context, productId uint64, houseInfo *domain.HouseInfo) error {
	if houseInfo != nil {
		// copier.Copy(&houseInfo, &dto.HouseInfo)
		houseInfo.ProductID = productId
		return s.HouseInfoRepo.UpdateHouseInfo(c, houseInfo)
	}
	return nil
}

func (s *ProductUsecase) UpdatePrivateData(c context.Context, productID uint64, privateInfo *domain.ProductPrivate) error {
	return s.ProductPrivate.UpdateInfo(c, productID, privateInfo)
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

func (s *ProductUsecase) SuggestFields(c context.Context, text string) (*dto.TotalSearchParser, error) {
	// searchs, err := utils.InferSearchText(body.Content)
	// if err != nil {
	// 	return nil, err
	// }
	// if searchs, err := utils.InferSearchText(body.Content); err != nil {
	// 	return nil, err
	// }
	var searchs dto.TotalSearchParser
	var err error
	if searchs, err = utils.InferSearchText(text); err != nil {
		return nil, err
	}

	// searchs, _ := utils.InferSearchText(text)

	regions, _ := s.Region.InterTextToItem(text, 10)
	docTypes, _ := s.DocTypeRepo.InterTextToItem(text, 10)
	propertyTypes, _ := s.PropertyType.InterTextToItem(text, 10)
	amenities, _ := s.Amenity.InterTextToEntity(text, 10)

	searchs.Name = _utils.TrimFirstWords(text, 12)

	// tìm địa chỉ
	var ward, district, province *domain.Region
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
		district, _ = s.Region.GetByID(c, ward.ParentID)
		province, _ = s.Region.GetByID(c, district.ParentID)
	}

	if district != nil && (province == nil || province.ID != district.ParentID) {
		province, _ = s.Region.GetByID(c, district.ParentID)
	}
	// ---- xong ----

	searchs.DocTypes = docTypes
	searchs.PropertyTypes = propertyTypes
	searchs.Amenities = amenities
	searchs.Ward = ward
	searchs.District = district
	searchs.Province = province
	if searchs.PriceSuggest != nil {
		if *searchs.PriceSuggest > 100_000_000 {
			searchs.TransactionType = 1
			searchs.SalePrice = searchs.PriceSuggest
		} else {
			searchs.TransactionType = 2
			searchs.RentPrice = searchs.PriceSuggest
		}
	}

	return &searchs, nil
}

func (s *ProductUsecase) RequiredOwner(c context.Context, id uint64) (*domain.Product, error) {
	product, err := s.ProductRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// profileId := _jwt.GetProfileId(c)
	profileId := _utils.GetProfileIdWithContext(c)

	if *product.OwnerUserID != profileId {
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
	if product.SaleStatus != enums.EProductNotSold || product.RentStatus != enums.EProductNotRent {
		return &_routes.Except{
			Code:    401,
			Message: "Không thể xóa sản phẩm này",
		}
	}
	return s.ProductRepo.Delete(c, id)
}

func (s *ProductUsecase) RentStatusUpdate(c context.Context, dto dto.StatusUpdate) error {
	err := s.CheckPermission(c, dto.ProductID, enums.PermissionProductUpdateStatus)
	if err != nil {
		return err
	}
	return s.ProductRepo.SetRentStatus(c, dto.ProductID, dto.Status)
}

func (s *ProductUsecase) SaleStatusUpdate(c context.Context, dto dto.StatusUpdate) error {
	err := s.CheckPermission(c, dto.ProductID, enums.PermissionProductUpdateStatus)
	if err != nil {
		return err
	}
	return s.ProductRepo.SetSaleStatus(c, dto.ProductID, dto.Status)
}

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

func (s *ProductUsecase) ProductDeposite(c context.Context, dto dto.DepositeRequest) (*domain.Deposite, error) {
	// sản phẩm đang cọc
	err := s.CheckPermission(c, dto.ProductID, enums.PermissionProductDeposite)
	if err != nil {
		return nil, err
	}

	product, err := s.ProductRepo.GetProductByID(c, dto.ProductID)
	if err != nil {
		return nil, err
	}

	if product.DepositeID != nil {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Sản phẩm đã đặt cọc",
		}
	}

	record := domain.Deposite{
		ProductID:     dto.ProductID,
		CustomerName:  dto.CustomerName,
		CustomerPhone: dto.CustomerPhone,
		Amount:        dto.Amount,
		Note:          dto.Note,
	}
	if err := s.ProductRepo.NewDeposite(c, &record); err != nil {
		return nil, err
	}
	product.DepositeID = &record.ID
	if err := s.ProductRepo.Update(c, product.ID, product); err != nil {
		return nil, err
	}
	return &record, nil
}

func (s *ProductUsecase) CancelDeposite(c context.Context, id uint64) error {
	// sản phẩm đang cọc
	err := s.CheckPermission(c, &id, enums.PermissionProductDeposite)
	if err != nil {
		return err
	}

	deposite, err := s.DepositeRepo.GetDepositeByID(c, id)
	if err != nil {
		return err
	}
	// product, err := s.RequiredOwner(c, deposite.ProductID)
	// if err != nil {
	// 	return err
	// }

	product, err := s.ProductRepo.GetProductByID(c, deposite.ProductID)
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

	return s.ProductRepo.CancelDeposite(c, product.ID)
}

func (uc *ProductUsecase) Posts(c context.Context, id uint64, dto dto.PostLinkSearch) (*[]domain.PostItem, error) {
	if uc.PostUsecase != nil {
		results, err := uc.PostUsecase.ProductPostLink(c, id, &dto)
		return results, err
	}
	return nil, nil
}

func (uc *ProductUsecase) Members(c context.Context, id uint64, dto dto.SharingAccessSearch) (*[]domain.ProductAccessQuery, error) {
	// _, err := s.RequiredOwnerProduct(c, id)
	// if err != nil {
	// 	return nil, err
	// }

	entities, err := uc.ProductAccessRepo.OwnerSearchProduct(c, id, enums.EOwnerTypeMember, &dto)
	return entities, err
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

func (s *ProductUsecase) SearchByProductAccess(c context.Context, targetID *uint64, targetType enums.EOwnerType, dto *dto.ProductSearchRequest) ([]domain.Product, int64, error) {
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
	// 		Message: "Sản phẩm không đủ điều kiến để tạo sản phẩm con",
	// 	}
	// }

	profileId := _utils.GetProfileIdWithContext(c)

	if *parent.OwnerUserID != profileId {
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

func (s *ProductUsecase) CreateChild(c context.Context, body *dto.SaveProductChildRequest) (*dto.CreateResponse, error) {
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

	body.ProductSaveRequest.ParentID = body.ParentID
	body.ProductSaveRequest.OwnerType = parent.OwnerType
	result, err := s.CreateProduct(c, &body.ProductSaveRequest)
	s.ProductChildRepo.CreateHistory(c, enums.EProductHistory_NewChild,
		body.ParentID,
		*body.ParentID,
		[]uint64{result.Product.ID},
		[]float64{result.Product.Area},
	)
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

	s.ProductChildRepo.CreateHistory(c,
		enums.EProductHistory_Merge,
		nil,
		*body.ParentID,
		*body.ChildIDs,
		[]float64{},
	)
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

	s.ProductChildRepo.CreateHistory(c,
		enums.EProductHistory_MergeAll,
		nil,
		*parentID, ids,
		[]float64{},
	)
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

	var ids []uint64
	var areas []float64

	for _, r := range results {
		ids = append(ids, r.ID)
		areas = append(areas, r.Area)
	}
	s.ProductChildRepo.CreateHistory(c, enums.EProductHistory_Devide,
		&child.ID,
		*child.ParentId,
		ids,
		areas,
	)
	return &results, nil
}

func (s *ProductUsecase) ChildHistory(c context.Context, id uint64, dto *dto.ProductHistorySearch) (*[]domain.ProductHistory, int64, error) {
	profileId := _utils.GetProfileIdWithContext(c)
	result, total, err := s.ProductChildRepo.History(c, profileId, id, dto)
	if err != nil {
		return nil, 0, err
	}
	return result, total, nil
}

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
