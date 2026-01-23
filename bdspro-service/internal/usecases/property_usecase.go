package usecases

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"bdspro/internal/provider"
	"bdspro/internal/repo"
	_enum "common/domain/enum"
	_errors "common/errors"
	_utils "common/utils"
	"context"
	"errors"
)

type PropertyUsecase struct {
	PropertyRepo             repo.PropertyRepo
	PropertyRelationRepo     repo.PropertyRelationRepo
	PropertyProductRepo      repo.PropertyProductRepo
	PropertyLandInfoRepo     repo.PropertyLandInfoRepo
	PropertyBuildingInfoRepo repo.PropertyBuildingInfoRepo
	PropertyMediaRepo        repo.PropertyMediaRepo
	PropertyAmenityRepo      repo.PropertyAmenityRepo
	AssetLegalRepo           repo.AssetLegalRepo
	ProductRepo              repo.ProductRepo
	ProductUserRepo          repo.ProductUserRepo
	ProductPriceRepo         repo.ProductPriceRepo
	ProductMediaRepo         repo.ProductMediaRepo
	AssetRepo                repo.AssetRepo
	Transaction              provider.TransactionProvider
	NotifyProvider           provider.NotificationProvider
}

func NewPropertyUsecase(
	propertyRepo repo.PropertyRepo,
	propertyRelationRepo repo.PropertyRelationRepo,
	propertyProductRepo repo.PropertyProductRepo,
	propertyLandInfoRepo repo.PropertyLandInfoRepo,
	propertyBuildingInfoRepo repo.PropertyBuildingInfoRepo,
	propertyMediaRepo repo.PropertyMediaRepo,
	propertyAmenityRepo repo.PropertyAmenityRepo,
	assetLegalRepo repo.AssetLegalRepo,
	productRepo repo.ProductRepo,
	productUserRepo repo.ProductUserRepo,
	productPriceRepo repo.ProductPriceRepo,
	productMediaRepo repo.ProductMediaRepo,
	assetRepo repo.AssetRepo,
	transaction provider.TransactionProvider,
	notifyProvider provider.NotificationProvider,
) *PropertyUsecase {
	return &PropertyUsecase{
		PropertyRepo:             propertyRepo,
		PropertyRelationRepo:     propertyRelationRepo,
		PropertyProductRepo:      propertyProductRepo,
		PropertyLandInfoRepo:     propertyLandInfoRepo,
		PropertyBuildingInfoRepo: propertyBuildingInfoRepo,
		PropertyMediaRepo:        propertyMediaRepo,
		PropertyAmenityRepo:      propertyAmenityRepo,
		AssetLegalRepo:           assetLegalRepo,
		ProductRepo:              productRepo,
		ProductUserRepo:          productUserRepo,
		ProductPriceRepo:         productPriceRepo,
		ProductMediaRepo:         productMediaRepo,
		AssetRepo:                assetRepo,
		Transaction:              transaction,
		NotifyProvider:           notifyProvider,
	}
}

func (u *PropertyUsecase) GetPropertyAmenities(
	ctx context.Context,
	propertyID uint64,
) ([]*dto.PropertyAmenityDTO, error) {
	if _utils.GetProfileIdWithContext(ctx) == 0 {
		return nil, _errors.UnauthorizedException()
	}

	property, err := u.PropertyRepo.GetDetail(ctx, propertyID)
	if err != nil {
		return nil, err
	}
	if property == nil {
		return nil, _errors.NotFoundException("Property not found")
	}

	amenities, err := u.PropertyAmenityRepo.GetAll(ctx, nil)
	if err != nil {
		return nil, err
	}

	selectedIDs, err := u.PropertyAmenityRepo.GetAmenityIDsByProperty(ctx, propertyID)
	if err != nil {
		return nil, err
	}

	selectedMap := map[uint64]bool{}
	for _, id := range selectedIDs {
		selectedMap[id] = true
	}

	result := make([]*dto.PropertyAmenityDTO, 0, len(amenities))
	for _, a := range amenities {
		result = append(result, &dto.PropertyAmenityDTO{
			ID:      a.ID,
			Name:    a.Name,
			Active:  true,
			Checked: selectedMap[a.ID],
		})
	}

	return result, nil
}

func (u *PropertyUsecase) UpdatePropertyAmenities(
	ctx context.Context,
	req *dto.UpdatePropertyAmenitiesRequest,
) error {
	if _utils.GetProfileIdWithContext(ctx) == 0 {
		return _errors.UnauthorizedException()
	}

	property, err := u.PropertyRepo.GetDetail(ctx, req.PropertyID)
	if err != nil {
		return err
	}
	if property == nil {
		return _errors.NotFoundException("Property not found")
	}

	return u.PropertyAmenityRepo.UpdatePropertyAmenityIds(
		ctx,
		req.PropertyID,
		req.AmenityIDs,
	)
}

func (u *PropertyUsecase) GetPropertyTags(
	ctx context.Context,
	propertyID uint64,
) ([]*dto.PropertyTagDTO, error) {

	if _utils.GetProfileIdWithContext(ctx) == 0 {
		return nil, _errors.UnauthorizedException()
	}

	property, err := u.PropertyRepo.GetDetail(ctx, propertyID)
	if err != nil {
		return nil, err
	}
	if property == nil {
		return nil, _errors.NotFoundException("Property not found")
	}

	tagIDs, err := u.PropertyAmenityRepo.GetTagIDsByProperty(ctx, propertyID)
	if err != nil {
		return nil, err
	}

	if len(tagIDs) == 0 {
		return []*dto.PropertyTagDTO{}, nil
	}

	tags, err := u.PropertyAmenityRepo.FindByIDs(ctx, tagIDs)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.PropertyTagDTO, 0, len(tags))
	for _, tag := range tags {
		result = append(result, &dto.PropertyTagDTO{
			ID:   tag.ID,
			Name: tag.Name,
			Type: tag.Type,
			Icon: tag.Icon,
		})
	}

	return result, nil
}

func (u *PropertyUsecase) SearchTags(
	ctx context.Context,
	search *dto.TagSearchDTO,
) ([]*dto.TagDTO, int64, error) {
	if _utils.GetProfileIdWithContext(ctx) == 0 {
		return nil, 0, _errors.UnauthorizedException()
	}

	tags, total, err := u.PropertyAmenityRepo.Search(ctx, search)
	if err != nil {
		return nil, 0, err
	}

	items := make([]*dto.TagDTO, 0, len(tags))
	for _, tag := range tags {
		items = append(items, &dto.TagDTO{
			ID:          tag.ID,
			Name:        tag.Name,
			Type:        tag.Type,
			Description: tag.Description,
			Icon:        tag.Icon,
		})
	}

	return items, total, nil
}

func (u *PropertyUsecase) UpdatePropertyTags(
	ctx context.Context,
	req *dto.UpdatePropertyTagsRequest,
) error {
	if _utils.GetProfileIdWithContext(ctx) == 0 {
		return _errors.UnauthorizedException()
	}

	property, err := u.PropertyRepo.GetDetail(ctx, req.PropertyID)
	if err != nil {
		return err
	}
	if property == nil {
		return _errors.NotFoundException("Property not found")
	}

	// validate tag ids (nếu có)
	if len(req.TagIDs) > 0 {
		tags, err := u.PropertyAmenityRepo.FindByIDs(ctx, req.TagIDs)
		if err != nil {
			return err
		}
		if len(tags) != len(req.TagIDs) {
			return _errors.BadRequestException("Invalid tag ids")
		}
	}

	if err := u.PropertyAmenityRepo.ReplacePropertyTags(
		ctx,
		req.PropertyID,
		req.TagIDs,
	); err != nil {
		return err
	}

	return nil
}

// LinkProductToProperty link product với property của product cha (nếu có)
func (u *PropertyUsecase) LinkProductToProperty(ctx context.Context, productID uint64, parentProductID uint64) error {
	// Tìm property của product cha
	parentRelation, err := u.PropertyRelationRepo.GetByProductID(ctx, parentProductID)
	if err != nil {
		return err
	}

	// Nếu không có property relation cho parent, không làm gì
	if parentRelation == nil {
		return nil
	}

	// Tạo link product với property
	propertyProduct := &domain.PropertyProduct{
		PropertyID: parentRelation.PropertyID,
		ProductID:  productID,
	}

	return u.PropertyProductRepo.Create(ctx, propertyProduct)
}

// LinkAssetToProperty link asset với property của product (nếu có)
func (u *PropertyUsecase) LinkAssetToProperty(ctx context.Context, assetID uint64, productID uint64) error {
	// Tìm property của product
	productRelation, err := u.PropertyRelationRepo.GetByProductID(ctx, productID)
	if err != nil {
		return err
	}

	// Nếu không có property relation cho product, không làm gì
	if productRelation == nil {
		return nil
	}

	// Tạo relation mới với type Asset hoặc update existing relation nếu chưa có asset relation
	if productRelation.RelationType == enums.ERelationTypeProduct {
		// Nếu đã có product relation, tạo thêm asset relation mới
		assetRelation := &domain.PropertyRelation{
			PropertyID:   productRelation.PropertyID,
			RelationID:   &assetID,
			RelationType: enums.ERelationTypeAsset,
		}
		return u.PropertyRelationRepo.Create(ctx, assetRelation)
	}
	// Nếu chưa có relation, tạo mới
	productRelation.RelationID = &assetID
	productRelation.RelationType = enums.ERelationTypeAsset
	return u.PropertyRelationRepo.Update(ctx, productRelation)
}

// GetList lấy danh sách Property với phân trang
func (u *PropertyUsecase) GetList(ctx context.Context, searchDTO *dto.PropertySearchDTO) ([]dto.PropertyWithCounts, int64, error) {
	return u.PropertyRepo.Search(ctx, searchDTO)
}

// Detail lấy chi tiết Property theo ID
func (u *PropertyUsecase) Detail(ctx context.Context, id uint64) (*domain.Property, error) {
	property, err := u.PropertyRepo.GetDetail(ctx, id)
	if err != nil {
		return nil, err
	}
	return property, nil
}

// GetRelation lấy danh sách PropertyRelation với đầy đủ thông tin Asset/Product/Post
func (u *PropertyUsecase) GetRelation(ctx context.Context, propertyID uint64) ([]domain.PropertyRelation, error) {
	// Lấy tất cả PropertyRelations với propertyID này
	relations, err := u.PropertyRelationRepo.GetAllByPropertyID(ctx, propertyID)
	if err != nil {
		return nil, err
	}

	result := make([]domain.PropertyRelation, 0, len(relations))

	// Load thêm thông tin Asset/Product/Post cho mỗi relation
	for _, relation := range relations {
		item := domain.PropertyRelation{
			PropertyID:   relation.PropertyID,
			RelationID:   relation.RelationID,
			RelationType: relation.RelationType,
		}

		if relation.RelationID == nil {
			result = append(result, item)
			continue
		}

		switch relation.RelationType {
		case enums.ERelationTypeProduct:
			// Load Product nếu có
			if product, err := u.ProductRepo.GetProductByID(ctx, relation.RelationID); err == nil && product != nil {
				item.Product = product
			}
		case enums.ERelationTypeAsset:
			// Load Asset nếu có
			if asset, err := u.AssetRepo.GetByID(*relation.RelationID); err == nil && asset != nil {
				item.Asset = asset
			}
		}

		result = append(result, item)
	}

	return result, nil
}

// Create tạo Property mới cùng với Product, Asset, Post và các thông tin liên quan
// Lưu ý: Product, Asset, Post sẽ được tạo ở handler layer trước khi gọi usecase này
func (u *PropertyUsecase) Create(ctx context.Context, req *dto.CreatePropertyProductRequest) (*domain.Property, error) {
	var property *domain.Property
	var productID *uint64
	var assetID *uint64

	err := u.Transaction.WithTransaction(ctx, func(txCtx context.Context) error {
		// 1. Tạo Property trước
		if req.Property != nil {
			property = req.Property
			// Set default recordStatus nếu chưa có
			if property.RecordStatus == 0 {
				property.RecordStatus = enums.EPropertyStatusActive
			}
			if err := u.PropertyRepo.Create(txCtx, property); err != nil {
				return err
			}
		}

		// 2. Tạo Product nếu có trong request và chưa có ID
		var err error
		if req.Product != nil {
			// Gán PropertyID vào Product nếu đã tạo Property
			if property != nil {
				req.Product.PropertyID = &property.ID
			}
			productID, err = u.CreateProductFromProperty(txCtx, req.Product)
			if err != nil {
				return err
			}
		}
		// 3. Tạo Asset nếu có trong request và chưa có ID
		if req.Asset != nil {
			assetID, err = u.CreateAssetFromProperty(txCtx, req.Asset)
			if err != nil {
				return err
			}
		}

		// 4. Tạo các thông tin liên quan đến Property
		if property != nil {
			// 4.1. Tạo PropertyLandInfo nếu có
			if req.LandInfo != nil {
				landInfo := req.LandInfo
				landInfo.PropertyID = property.ID
				if err := u.PropertyLandInfoRepo.Create(txCtx, landInfo); err != nil {
					return err
				}
			}

			// 4.2. Tạo PropertyBuildingInfo nếu có
			if req.BuildingInfo != nil {
				buildingInfo := req.BuildingInfo
				buildingInfo.PropertyID = property.ID
				if err := u.PropertyBuildingInfoRepo.Create(txCtx, buildingInfo); err != nil {
					return err
				}
				// Link BuildingInfoID vào Property
				property.BuildingInfoID = &buildingInfo.ID
			}

			// 4.3. Tạo PropertyMedia nếu có
			if len(req.MediaList) > 0 {
				for i := range req.MediaList {
					req.MediaList[i].PropertyID = property.ID
				}
				if err := u.PropertyMediaRepo.CreateBatch(txCtx, req.MediaList); err != nil {
					return err
				}
				// Set ảnh đầu tiên làm background image
				if len(req.MediaList) > 0 {
					property.AvatarID = &req.MediaList[0].ID
				}
			}

			// 4.4. Tạo PropertyAmenity nếu có
			if len(req.AmenityIds) > 0 {
				if err := u.UpdatePropertyAmenity(txCtx, property.ID, req.AmenityIds); err != nil {
					return err
				}
			}

			// 4.5. Tạo PropertyRelation để liên kết với Product/Asset
			// Tạo relation với Product trước
			if productID != nil {
				productRelation := &domain.PropertyRelation{
					PropertyID:   property.ID,
					RelationID:   productID,
					RelationType: enums.ERelationTypeProduct,
				}
				if err := u.PropertyRelationRepo.Create(txCtx, productRelation); err != nil {
					return err
				}
			}
			// Tạo relation với Asset nếu có
			if assetID != nil {
				assetRelation := &domain.PropertyRelation{
					PropertyID:   property.ID,
					RelationID:   assetID,
					RelationType: enums.ERelationTypeAsset,
				}
				if err := u.PropertyRelationRepo.Create(txCtx, assetRelation); err != nil {
					return err
				}
			}

			if err := u.PropertyRepo.Update(txCtx, property.ID, property); err != nil {
				return err
			}
		}

		// 8. Tạo AssetLegal nếu có
		if req.LegalInfo != nil && assetID != nil {
			legalInfo := req.LegalInfo
			legalInfo.AssetID = *assetID
			if err := u.AssetLegalRepo.Create(txCtx, legalInfo); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Tạo property history sau khi tạo property thành công
	if property != nil && u.NotifyProvider != nil {
		u.NotifyProvider.RegistedEventProperty(ctx,
			&dto.CreatePropertyActivityDTO{
				SubjectID:   property.ID,
				SubjectType: _enum.SUBJECT_TYPE_PROPERTY,
				Action:      _enum.ACTION_PROPERTY_CREATE,
				ActorID:     _utils.GetProfileIdWithContext(ctx),
				Description: "Sản phẩm được tạo ở đây",
			})
	}

	return property, nil
}

// CreateProductFromProperty tạo Product nếu chưa có ID, trả về productID
func (u *PropertyUsecase) CreateProductFromProperty(txCtx context.Context, product *domain.Product) (*uint64, error) {
	if product == nil {
		return nil, nil
	}

	// Nếu Product đã có ID thì không cần tạo
	if product.ID != 0 {
		return &product.ID, nil
	}

	// Set OwnerID và OwnerOf nếu chưa có
	if product.OwnerID == 0 {
		ownerID := _utils.GetProfileIdWithContext(txCtx)
		product.OwnerID = ownerID
	}
	if product.OwnerOf == 0 {
		product.OwnerOf = enums.EOwnerOfMember
	}

	// Tạo Price nếu có
	if product.Price != nil {
		if err := u.ProductPriceRepo.Create(txCtx, product.Price); err != nil {
			return nil, err
		}
		product.LastPriceID = &product.Price.ID
	}

	// Tạo Product
	if err := u.ProductRepo.Create(txCtx, product); err != nil {
		return nil, err
	}

	// Update ProductID vào Price sau khi tạo Product
	if product.Price != nil {
		if err := u.ProductPriceRepo.UpdateProductID(txCtx, &product.ID, &product.Price.ID); err != nil {
			return nil, err
		}
	}

	// Tạo ProductUser để link Product với User
	// Chỉ tạo nếu OwnerOf là Member (không phải Organization)
	if product.OwnerOf == enums.EOwnerOfMember {
		profileID := _utils.GetProfileIdWithContext(txCtx)
		_, err := u.ProductUserRepo.CreateOwner(txCtx, product.ID, profileID)
		if err != nil {
			return nil, err
		}
	}

	return &product.ID, nil
}

// CreateAssetFromProperty tạo Asset nếu chưa có ID, trả về assetID
func (u *PropertyUsecase) CreateAssetFromProperty(txCtx context.Context, asset *domain.Asset) (*uint64, error) {
	if asset == nil {
		return nil, nil
	}

	// Nếu Asset đã có ID thì không cần tạo
	if asset.ID != 0 {
		return &asset.ID, nil
	}

	// Set OwnerID và OwnerOf nếu chưa có
	if asset.OwnerID == nil || *asset.OwnerID == 0 {
		ownerID := _utils.GetProfileIdWithContext(txCtx)
		asset.OwnerID = &ownerID
	}
	if asset.OwnerOf == 0 {
		asset.OwnerOf = enums.EOwnerOfMember
	}

	// Set default values nếu chưa có
	if asset.Archived == false && asset.ID == 0 {
		asset.Archived = false
	}

	// Tạo Asset
	if err := u.AssetRepo.Create(txCtx, asset); err != nil {
		return nil, err
	}

	return &asset.ID, nil
}

// CreateWithProduct tạo Property từ Product
// Nếu Product đã có AssetId thì không tạo Property mới (vì đã có PropertyRelation)
func (u *PropertyUsecase) CreateWithProduct(ctx context.Context, product *domain.Product) (*domain.Property, error) {
	// Nếu Product đã có AssetId, kiểm tra xem đã có PropertyRelation chưa
	if product.AssetId != nil {
		relation, err := u.PropertyRelationRepo.GetByAssetID(ctx, *product.AssetId)
		if err != nil {
			return nil, err
		}
		// Nếu đã có PropertyRelation, không tạo Property mới
		if relation != nil {
			// Lấy Property hiện có
			property, err := u.PropertyRepo.GetByID(relation.PropertyID)
			if err != nil {
				return nil, err
			}
			return property, nil
		}
	}

	// Tạo Property từ thông tin Product
	property := &domain.Property{
		PropertyTypeID: product.PropertyTypeId,
		Title:          product.Name,
		ProjectID:      product.ProjectID,
		AreaTotal:      &product.Area,
		AddressDetail:  product.Address,
		MapURL:         product.GoogleMapLink,
		RecordStatus:   enums.EPropertyStatusActive,
	}

	// Tạo Property
	if err := u.PropertyRepo.Create(ctx, property); err != nil {
		return nil, err
	}

	// Tạo PropertyRelation với type Product
	relation := &domain.PropertyRelation{
		PropertyID:   property.ID,
		RelationID:   &product.ID,
		RelationType: enums.ERelationTypeProduct,
	}

	if err := u.PropertyRelationRepo.Create(ctx, relation); err != nil {
		return nil, err
	}

	// Nếu product có AssetId, tạo thêm relation với Asset
	if product.AssetId != nil {
		assetRelation := &domain.PropertyRelation{
			PropertyID:   property.ID,
			RelationID:   product.AssetId,
			RelationType: enums.ERelationTypeAsset,
		}
		if err := u.PropertyRelationRepo.Create(ctx, assetRelation); err != nil {
			return nil, err
		}
	}

	// Copy amenityIds từ Product sang Property nếu có
	amenityIds, err := u.ProductRepo.GetAmenityIDsByProductID(ctx, product.ID)
	if err == nil && len(amenityIds) > 0 {
		if err := u.UpdatePropertyAmenity(ctx, property.ID, amenityIds); err != nil {
			// Log error nhưng không fail transaction
			// return nil, err
		}
	}

	return property, nil
}

// UpdatePropertyAmenity cập nhật amenityIds cho Property
func (u *PropertyUsecase) UpdatePropertyAmenity(ctx context.Context, propertyID uint64, amenityIDs []uint64) error {
	if len(amenityIDs) > 0 {
		return u.PropertyAmenityRepo.UpdatePropertyAmenityIds(ctx, propertyID, amenityIDs)
	}
	return nil
}

// GetOrCreatePropertyFromProductID lấy PropertyID từ ProductID, nếu chưa có thì tạo mới
// Trả về PropertyID và error
func (u *PropertyUsecase) GetOrCreatePropertyFromProductID(ctx context.Context, productID uint64) (uint64, error) {
	// Kiểm tra xem đã có PropertyRelation chưa
	relation, err := u.PropertyRelationRepo.GetByProductID(ctx, productID)
	if err != nil {
		return 0, err
	}

	// Nếu đã có PropertyRelation, trả về PropertyID
	if relation != nil {
		return relation.PropertyID, nil
	}

	// Nếu chưa có, lấy Product để tạo Property
	product, err := u.ProductRepo.GetByIDContext(ctx, productID)
	if err != nil {
		return 0, err
	}
	if product == nil {
		return 0, errors.New("product not found")
	}

	// Tạo Property từ thông tin Product
	property := &domain.Property{
		PropertyTypeID: product.PropertyTypeId,
		Title:          product.Name,
		ProjectID:      product.ProjectID,
		AreaTotal:      &product.Area,
		AddressDetail:  product.Address,
		MapURL:         product.GoogleMapLink,
		RecordStatus:   enums.EPropertyStatusActive,
	}

	// Tạo Property
	if err := u.PropertyRepo.Create(ctx, property); err != nil {
		return 0, err
	}

	// Tạo PropertyRelation với type Product
	relation = &domain.PropertyRelation{
		PropertyID:   property.ID,
		RelationID:   &productID,
		RelationType: enums.ERelationTypeProduct,
	}

	if err := u.PropertyRelationRepo.Create(ctx, relation); err != nil {
		return 0, err
	}

	// Nếu product có AssetId, tạo thêm relation với Asset
	if product.AssetId != nil {
		assetRelation := &domain.PropertyRelation{
			PropertyID:   property.ID,
			RelationID:   product.AssetId,
			RelationType: enums.ERelationTypeAsset,
		}
		if err := u.PropertyRelationRepo.Create(ctx, assetRelation); err != nil {
			return 0, err
		}
	}

	// Copy amenityIds từ Product sang Property nếu có
	amenityIds, err := u.ProductRepo.GetAmenityIDsByProductID(ctx, productID)
	if err == nil && len(amenityIds) > 0 {
		if err := u.UpdatePropertyAmenity(ctx, property.ID, amenityIds); err != nil {
			// Log error nhưng không fail transaction
		}
	}

	return property.ID, nil
}

// UpdatePropertyFromProduct cập nhật Property, BuildingInfo và LandInfo của Property từ ProductID
func (u *PropertyUsecase) UpdatePropertyFromProduct(ctx context.Context, productID uint64, houseInfo *dto.UpdateProductRequest) error {
	if houseInfo == nil {
		return nil
	}

	// Lấy Property từ Product
	propertyRelation, err := u.PropertyRelationRepo.GetByProductID(ctx, productID)
	if err != nil {
		return err
	}

	// Nếu không có Property, không làm gì
	if propertyRelation == nil {
		return nil
	}

	// Update Property với các field: AreaLand, AreaTotal, ProjectID, PropertyTypeID
	property, err := u.PropertyRepo.GetDetail(ctx, propertyRelation.PropertyID)
	if err != nil {
		return err
	}
	if property != nil {
		// Tạo property update chỉ với các field cần update
		updateProperty := &domain.Property{
			AreaLand:       &houseInfo.AreaLand,
			AreaTotal:      &houseInfo.AreaTotal,
			ProjectID:      houseInfo.ProjectID,
			PropertyTypeID: houseInfo.PropertyTypeID,
		}

		// Update Property nếu có thay đổi
		if err := u.PropertyRepo.UpdateOption(ctx, propertyRelation.PropertyID, updateProperty); err != nil {
			return err
		}
	}

	// Convert từ UpdateProductRequest sang PropertyBuildingInfo
	buildingInfo := &domain.PropertyBuildingInfo{
		PropertyID: propertyRelation.PropertyID,
	}

	// Map các field từ UpdateProductRequest sang PropertyBuildingInfo
	if houseInfo.Floors != nil {
		buildingInfo.Floors = houseInfo.Floors
	}
	if houseInfo.Bedrooms != nil {
		buildingInfo.Bedrooms = houseInfo.Bedrooms
	}
	if houseInfo.Bathrooms != nil {
		buildingInfo.Bathrooms = houseInfo.Bathrooms
	}
	if houseInfo.Direction != 0 {
		buildingInfo.Direction = houseInfo.Direction
	}
	if houseInfo.BalconyDirection != 0 {
		buildingInfo.BalconyDirection = houseInfo.BalconyDirection
	}
	// Lấy existing building info để update hoặc create
	existingBuildingInfo, err := u.PropertyBuildingInfoRepo.GetByPropertyID(ctx, propertyRelation.PropertyID)
	if err != nil {
		return err
	}

	if existingBuildingInfo != nil {
		// Update existing
		buildingInfo.ID = existingBuildingInfo.ID
		if err := u.PropertyBuildingInfoRepo.Update(ctx, buildingInfo); err != nil {
			return err
		}
	} else {
		// Create new
		if err := u.PropertyBuildingInfoRepo.Create(ctx, buildingInfo); err != nil {
			return err
		}
		// Link BuildingInfoID vào Property
		property, err := u.PropertyRepo.GetOneByID(ctx, propertyRelation.PropertyID)
		if err != nil {
			return err
		}
		if property != nil {
			property.BuildingInfoID = &buildingInfo.ID
			if err := u.PropertyRepo.Update(ctx, propertyRelation.PropertyID, property); err != nil {
				return err
			}
		}
	}

	// Update LandInfo nếu có RoadWidth (RoadWidth lưu vào LandInfo, không phải BuildingInfo)
	if houseInfo.RoadWidth != nil {
		landInfo := &domain.PropertyLandInfo{
			PropertyID: propertyRelation.PropertyID,
			RoadWidth:  houseInfo.RoadWidth,
		}

		// Lấy existing land info để update hoặc create
		existingLandInfo, err := u.PropertyLandInfoRepo.GetByPropertyID(ctx, propertyRelation.PropertyID)
		if err != nil {
			return err
		}

		if existingLandInfo != nil {
			// Update existing
			landInfo.ID = existingLandInfo.ID
			// Giữ nguyên các field khác nếu có
			if existingLandInfo.Frontage != nil {
				landInfo.Frontage = existingLandInfo.Frontage
			}
			if existingLandInfo.Depth != nil {
				landInfo.Depth = existingLandInfo.Depth
			}
			if existingLandInfo.LandNote != "" {
				landInfo.LandNote = existingLandInfo.LandNote
			}
			if err := u.PropertyLandInfoRepo.Update(ctx, landInfo); err != nil {
				return err
			}
		} else {
			// Create new
			if err := u.PropertyLandInfoRepo.Create(ctx, landInfo); err != nil {
				return err
			}
		}
	}

	return nil
}

// UpdateMediaItems cập nhật danh sách media items cho Property
// Dựa trên pattern từ đoạn code 411-423: xóa media cũ, tạo media mới, set ảnh đầu tiên làm avatar
func (u *PropertyUsecase) UpdateMediaItems(ctx context.Context, propertyID uint64, mediaItems []dto.MediaItem) ([]uint64, error) {
	// Xóa các media items cũ của property
	if err := u.PropertyMediaRepo.DeleteByPropertyID(ctx, propertyID); err != nil {
		return nil, err
	}

	// Nếu danh sách media items rỗng -> không cần thêm mới, chỉ cần xóa là đủ
	if len(mediaItems) == 0 {
		return nil, nil
	}

	// Convert từ dto.MediaItem sang domain.PropertyMedia và gán PropertyID
	var propertyMediaList []domain.PropertyMedia
	for i, mediaItem := range mediaItems {
		propertyMedia := domain.PropertyMedia{
			PropertyID: propertyID,
			MediaType:  mediaItem.MediaType,
			URL:        mediaItem.MediaURL,
			SortOrder:  int32(i + 1), // Sử dụng Order từ MediaItem hoặc index + 1
			IsCover:    mediaItem.IsMain,
		}
		// Sử dụng Order từ MediaItem nếu có
		if mediaItem.Order > 0 {
			propertyMedia.SortOrder = int32(mediaItem.Order)
		}
		propertyMediaList = append(propertyMediaList, propertyMedia)
	}

	// Tạo batch media items mới
	if err := u.PropertyMediaRepo.CreateBatch(ctx, propertyMediaList); err != nil {
		return nil, err
	}

	// Lấy danh sách ID đã tạo
	var mediaIDs []uint64
	var avatarID *uint64
	for _, media := range propertyMediaList {
		mediaIDs = append(mediaIDs, media.ID)
		// Tìm ảnh có IsCover = true làm avatar, hoặc lấy ảnh đầu tiên
		if avatarID == nil {
			if media.IsCover {
				avatarID = &media.ID
			}
		}
	}

	// Nếu không có ảnh IsCover, lấy ảnh đầu tiên
	if avatarID == nil && len(mediaIDs) > 0 {
		avatarID = &mediaIDs[0]
	}

	// Update AvatarID của Property
	if avatarID != nil {
		property, err := u.PropertyRepo.GetOneByID(ctx, propertyID)
		if err != nil {
			return nil, err
		}
		if property != nil {
			property.AvatarID = avatarID
			if err := u.PropertyRepo.Update(ctx, propertyID, property); err != nil {
				return nil, err
			}
		}
	} else {
		// Nếu không có media items, xóa AvatarID
		property, err := u.PropertyRepo.GetOneByID(ctx, propertyID)
		if err != nil {
			return nil, err
		}
		if property != nil {
			property.AvatarID = nil
			if err := u.PropertyRepo.Update(ctx, propertyID, property); err != nil {
				return nil, err
			}
		}
	}

	return mediaIDs, nil
}

// CreateWithAsset tạo Property từ Asset
// Nếu Asset đã có ProductID thì không tạo Property mới (vì đã có PropertyRelation)
func (u *PropertyUsecase) CreateWithAsset(ctx context.Context, asset *domain.Asset) (*domain.Property, error) {
	// Nếu Asset đã có ProductID, kiểm tra xem đã có PropertyRelation chưa
	if asset.ProductID != nil {
		relation, err := u.PropertyRelationRepo.GetByProductID(ctx, *asset.ProductID)
		if err != nil {
			return nil, err
		}
		// Nếu đã có PropertyRelation, không tạo Property mới
		if relation != nil {
			// Lấy Property hiện có
			property, err := u.PropertyRepo.GetOneByID(ctx, relation.PropertyID)
			if err != nil {
				return nil, err
			}
			return property, nil
		}
	}

	// Tạo Property từ thông tin Asset
	property := &domain.Property{
		PropertyTypeID: asset.PropertyTypeId,
		Title:          asset.Name,
		AreaTotal:      &asset.Area,
		AddressDetail:  asset.Address,
		RecordStatus:   enums.EPropertyStatusActive, // Default status
	}

	// Map legal status từ Asset - giữ nguyên dạng số hoặc convert nếu cần
	// LegalStatus trong Property là string, có thể để trống hoặc convert sau
	// property.LegalStatus = "" // Có thể map sau nếu cần

	// Tạo Property
	if err := u.PropertyRepo.Create(ctx, property); err != nil {
		return nil, err
	}

	// Tạo PropertyRelation với type Asset
	assetRelation := &domain.PropertyRelation{
		PropertyID:   property.ID,
		RelationID:   &asset.ID,
		RelationType: enums.ERelationTypeAsset,
	}

	if err := u.PropertyRelationRepo.Create(ctx, assetRelation); err != nil {
		return nil, err
	}

	// Nếu asset có ProductID, tạo thêm relation với Product
	if asset.ProductID != nil {
		productRelation := &domain.PropertyRelation{
			PropertyID:   property.ID,
			RelationID:   asset.ProductID,
			RelationType: enums.ERelationTypeProduct,
		}
		if err := u.PropertyRelationRepo.Create(ctx, productRelation); err != nil {
			return nil, err
		}
	}

	return property, nil
}

// Update cập nhật Property và các thông tin liên quan (LandInfo, BuildingInfo)
func (u *PropertyUsecase) Update(ctx context.Context, propertyID uint64, property *domain.Property, landInfo *domain.PropertyLandInfo, buildingInfo *domain.PropertyBuildingInfo, amenityIDs []uint64) error {
	return u.Transaction.WithTransaction(ctx, func(txCtx context.Context) error {
		// 1. Kiểm tra Property có tồn tại không
		existingProperty, err := u.PropertyRepo.GetOneByID(txCtx, propertyID)
		if err != nil {
			return err
		}
		if existingProperty == nil {
			return errors.New("property not found")
		}

		// 2. Update Property
		// if property.BuildingInfoID == nil {
		// 	existingProperty.PropertyTypeID = nil
		// 	if err := u.PropertyRepo.Update(txCtx, propertyID, existingProperty); err != nil {
		// 		return err
		// 	}
		// }

		// 4. Update PropertyBuildingInfo nếu có
		if buildingInfo != nil {
			existingBuildingInfo, err := u.PropertyBuildingInfoRepo.GetByPropertyID(txCtx, propertyID)
			if err != nil {
				return err
			}

			if existingBuildingInfo != nil {
				// Update existing
				buildingInfo.ID = existingBuildingInfo.ID
				buildingInfo.PropertyID = propertyID
				existingProperty.BuildingInfoID = &existingBuildingInfo.ID
				if err := u.PropertyBuildingInfoRepo.Update(txCtx, buildingInfo); err != nil {
					return err
				}
			} else {
				// Create new
				buildingInfo.PropertyID = propertyID
				if err := u.PropertyBuildingInfoRepo.Create(txCtx, buildingInfo); err != nil {
					return err
				}
				// Link BuildingInfoID vào Property
				updateProperty := *existingProperty
				updateProperty.BuildingInfoID = &buildingInfo.ID
				if err := u.PropertyRepo.Update(txCtx, propertyID, &updateProperty); err != nil {
					return err
				}
			}
		} else {
			existingProperty.BuildingInfoID = nil
		}

		if property != nil {
			// Merge với existing property để giữ các field không được update
			updateProperty := *existingProperty
			// updateProperty.BuildingInfoID = property.BuildingInfoID
			if property.PropertyTypeID != nil {
				updateProperty.PropertyTypeID = property.PropertyTypeID
			}
			if property.Title != "" {
				updateProperty.Title = property.Title
			}
			if property.LocationID != nil {
				updateProperty.LocationID = property.LocationID
			}
			if property.AddressDetail != "" {
				updateProperty.AddressDetail = property.AddressDetail
			}
			if property.Latitude != nil {
				updateProperty.Latitude = property.Latitude
			}
			if property.Longitude != nil {
				updateProperty.Longitude = property.Longitude
			}
			if property.MapURL != "" {
				updateProperty.MapURL = property.MapURL
			}
			if property.ProjectID != nil {
				updateProperty.ProjectID = property.ProjectID
			}
			if property.Level != "" {
				updateProperty.Level = property.Level
			}
			if property.UnitCode != "" {
				updateProperty.UnitCode = property.UnitCode
			}
			if property.AreaTotal != nil {
				updateProperty.AreaTotal = property.AreaTotal
			}
			if property.AreaLand != nil {
				updateProperty.AreaLand = property.AreaLand
			}
			if property.AreaResidential != nil {
				updateProperty.AreaResidential = property.AreaResidential
			}
			if property.LegalNote != "" {
				updateProperty.LegalNote = property.LegalNote
			}
			if property.ProvinceID != nil {
				updateProperty.ProvinceID = property.ProvinceID
			}
			if property.WardID != nil {
				updateProperty.WardID = property.WardID
			}
			// Update enum fields nếu đã set
			if property.LegalStatus != 0 {
				updateProperty.LegalStatus = property.LegalStatus
			}
			if property.RecordStatus != 0 {
				updateProperty.RecordStatus = property.RecordStatus
			}

			if err := u.PropertyRepo.Update(txCtx, propertyID, &updateProperty); err != nil {
				return err
			}
		}

		// 3. Update PropertyLandInfo nếu có
		if landInfo != nil {
			existingLandInfo, err := u.PropertyLandInfoRepo.GetByPropertyID(txCtx, propertyID)
			if err != nil {
				return err
			}

			if existingLandInfo != nil {
				// Update existing
				landInfo.ID = existingLandInfo.ID
				landInfo.PropertyID = propertyID
				if err := u.PropertyLandInfoRepo.Update(txCtx, landInfo); err != nil {
					return err
				}
			} else {
				// Create new
				landInfo.PropertyID = propertyID
				if err := u.PropertyLandInfoRepo.Create(txCtx, landInfo); err != nil {
					return err
				}
			}
		}

		// 5. Update PropertyAmenity nếu có
		if amenityIDs != nil {
			if err := u.UpdatePropertyAmenity(txCtx, propertyID, amenityIDs); err != nil {
				return err
			}
		}

		return nil
	})
}

// DeleteProperties xóa nhiều Property theo danh sách ID
func (u *PropertyUsecase) DeleteProperties(ctx context.Context, ids []uint64) error {
	if len(ids) == 0 {
		return nil
	}

	return u.Transaction.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := u.PropertyRepo.DeleteBatchProperties(txCtx, ids); err != nil {
			return err
		}

		return nil
	})
}
