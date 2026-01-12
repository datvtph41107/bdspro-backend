package usecases

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"bdspro/internal/provider"
	"bdspro/internal/repo"
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
	}
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
			if product, err := u.ProductRepo.GetByID(*relation.RelationID); err == nil && product != nil {
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
				if err := u.PropertyRepo.Update(txCtx, property.ID, property); err != nil {
					return err
				}
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
					if err := u.PropertyRepo.Update(txCtx, property.ID, property); err != nil {
						return err
					}
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
		AddressText:    product.Address,
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
			property, err := u.PropertyRepo.GetByID(relation.PropertyID)
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
		AddressText:    asset.Address,
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
		existingProperty, err := u.PropertyRepo.GetByID(propertyID)
		if err != nil {
			return err
		}
		if existingProperty == nil {
			return errors.New("property not found")
		}

		// 2. Update Property
		if property != nil {
			// Merge với existing property để giữ các field không được update
			updateProperty := *existingProperty
			if property.PropertyTypeID != nil {
				updateProperty.PropertyTypeID = property.PropertyTypeID
			}
			if property.Title != "" {
				updateProperty.Title = property.Title
			}
			if property.LocationID != nil {
				updateProperty.LocationID = property.LocationID
			}
			if property.AddressText != "" {
				updateProperty.AddressText = property.AddressText
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
