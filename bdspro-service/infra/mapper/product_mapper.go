package mapper

import (
	"bdspro/infra/client"
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"bdspro/internal/repo"
	_dto "common/domain/dto"
	_provider "common/domain/provider"
	_usecase "common/domain/usecase"
	_utils "common/utils"
	"context"
	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"
	"strings"
)

type ProductMapper struct {
	*CommonMapper
	PostMapper               *PostMapper
	ApartmentMapper          *ApartmentMapper
	ApartmentAttributeMapper *ApartmentAttributeMapper
	BuildMapper              *BuildMapper
	ProjectMapper            *ProjectMapper
	DeveloperMapper          *DeveloperMapper
	DepositeMapper           *DepositeMapper
	CodeDataUsecase          *_usecase.CodeDataUsecase
	UserClient               *client.UserClient
	CrmClient                _provider.CrmProvider
	AssetRepo                repo.AssetRepo
	PropertyMapper           *PropertyMapper
	// TransactionMapper        TransactionTransformer
}

func NewProductMapper(postMapper *PostMapper,
	apartmentMapper *ApartmentMapper,
	apartmentAttributeMapper *ApartmentAttributeMapper,
	buildMapper *BuildMapper,
	projectMapper *ProjectMapper,
	developerMapper *DeveloperMapper,
	depositeMapper *DepositeMapper,
	codeDataUsecase *_usecase.CodeDataUsecase,
	userClient *client.UserClient,
	crmClient _provider.CrmProvider,
	assetRepo repo.AssetRepo,
	// transactionMapper TransactionTransformer,
) *ProductMapper {
	result := &ProductMapper{
		CommonMapper:             &CommonMapper{},
		PostMapper:               postMapper,
		ApartmentMapper:          apartmentMapper,
		ApartmentAttributeMapper: apartmentAttributeMapper,
		BuildMapper:              buildMapper,
		ProjectMapper:            projectMapper,
		DeveloperMapper:          developerMapper,
		DepositeMapper:           depositeMapper,
		CodeDataUsecase:          codeDataUsecase,
		UserClient:               userClient,
		CrmClient:                crmClient,
		AssetRepo:                assetRepo,
		// PropertyMapper:           propertyMapper,
		// TransactionMapper:        transactionMapper,
	}
	postMapper.ProductMapper = result
	return result
}

func (m *ProductMapper) ApplyPriceUpdatePbToDTO(
	req *bdspropb.ApplyPriceUpdateRequest,
) *dto.ApplyPriceUpdateRequest {
	result := &dto.ApplyPriceUpdateRequest{
		ProductID:   req.ProductId,
		Price:       req.Price,
		PriceStatus: enums.ProductPriceStatus(req.Price),
	}

	if req.PriceStatus != 0 {
		result.PriceStatus = enums.ProductPriceStatus(req.PriceStatus)
	}

	if req.ChangeNote != nil {
		result.ChangeNote = *req.ChangeNote
	}

	return result
}

func (m *ProductMapper) ApplyPriceUpdateDTOToPb(
	resp *dto.ApplyPriceUpdateResponse,
) *bdspropb.ApplyPriceUpdateResponse {

	pb := &bdspropb.ApplyPriceUpdateResponse{
		ProductId: resp.ProductID,
	}

	if resp.OldPrice != nil {
		pb.OldPrice = *resp.OldPrice
	}
	if resp.NewPrice != nil {
		pb.NewPrice = *resp.NewPrice
	}
	if resp.PricePerM2 != nil {
		pb.PricePerM2 = *resp.PricePerM2
	}
	pb.UpdatedAt = _utils.FormatTimeToString(&resp.UpdatedAt)

	return pb
}

func (m *ProductMapper) PbToUpdateProductSourceDTO(
	req *bdspropb.UpdateProductSourceRequest,
) *dto.UpdateProductSourceRequest {

	var sourceContactId uint64
	if req.SourceContactId != 0 {
		sourceContactId = req.SourceContactId
	}

	// var assignedUserId *uint64
	// if req.AssignedUserId != 0 {
	// 	assignedUserId = &req.AssignedUserId
	// }

	return &dto.UpdateProductSourceRequest{
		ProductID:         req.ProductId,
		SourceType:        enums.ESourceType(req.SourceType),
		SourceContactId:   sourceContactId,
		SourceContactNote: req.SourceContactNote,
		SourceStatus:      enums.ESourceStatus(req.SourceStatus),
		// AssignedUserID:    assignedUserId,
		Priority:     enums.EPriority(req.Priority),
		SourceTagIDs: req.SourceTagIds,
	}
}

// cleanAddressUnit loại bỏ các từ tiền tố địa chỉ như "Thành phố", "Tỉnh", "Phường", "Thị xã", "Huyện"
func cleanAddressUnit(input string) string {
	if input == "" {
		return input
	}

	// Danh sách các tiền tố cần loại bỏ (lowercase để so sánh)
	prefixes := []string{
		"thành phố ",
		"tỉnh ",
		"phường ",
		"xã ",
		"thị xã ",
		"thị trấn ",
		"huyện ",
		"quận ",
		"tp. ",
		"tp ",
		"tt. ",
		"tt ",
	}

	// Chuyển sang lowercase để so sánh
	lowerInput := strings.ToLower(input)

	// Kiểm tra và loại bỏ tiền tố
	for _, prefix := range prefixes {
		if strings.HasPrefix(lowerInput, prefix) {
			// Lấy phần sau tiền tố từ chuỗi gốc (giữ nguyên chữ hoa/thường gốc)
			result := input[len(prefix):]
			return strings.TrimSpace(result)
		}
	}

	return strings.TrimSpace(input)
}

// cleanAddressUnitPtr loại bỏ các từ tiền tố địa chỉ cho pointer string
func cleanAddressUnitPtr(input *string) *string {
	if input == nil || *input == "" {
		return input
	}
	result := cleanAddressUnit(*input)
	return &result
}

// func (m *ProductMapper) MapProductDetailPb(product *domain.Product) *bdspropb.ProductResponse {
// 	result := m.MapProductPb(product)
// 	result.Id = product.ID
// 	result.ParentId = product.ParentId
// 	// result.DepositeId = product.DepositeID
// 	// result.Deposite           = product.Deposite
// 	result.AssetId = product.AssetId
// 	result.PropertyTypeId = product.PropertyTypeId
// 	// result.PropertyType       = product.PropertyType
// 	result.DocTypeId = product.DocTypeId
// 	// result.DocType            = product.DocType
// 	// result.Amenities          = product.AmenityIDs
// 	result.Visibility = uint32(product.Visibility)
// 	result.Code = product.Code
// 	result.CategoryId = product.CategoryID
// 	result.PositionUrl = product.PositionUrl
// 	result.Description = product.Description
// 	// result.OrganizationId     = product.OrganizationID
// 	// result.OrgStatus          = product.OrgStatus
// 	result.LastPriceId = product.LastPriceID
// 	// result.Price              = product.Price
// 	result.ApartmentId = product.ApartmentID
// 	// result.Apartment          = product.Apartment
// 	// result.ApartmentAttribute = product.ApartmentAttribute
// 	// result.Build              = product.Build
// 	// result.Project            = product.Project
// 	// result.Developer          = product.Developer
// 	result.SourceType = int32(product.SourceType)
// 	result.CustomerId = product.CustomerID
// 	// result.PrivateData        = product.PrivateData
// 	result.Note = product.Note
// 	result.TransactionType = int32(product.TransactionType)
// 	result.SaleVisibility = uint32(product.SaleVisibility)
// 	result.RentVisibility = uint32(product.RentVisibility)
// 	result.OwnerId = product.OwnerID
// 	result.OwnerType = uint32(product.OwnerType)
// 	result.Address = product.Address
// 	result.GoogleMapLink = product.GoogleMapLink
// 	// result.Attributes         = product.Attributes
// 	// result.MediaList          = product.MediaList
// 	// result.HouseInfo          = product.HouseInfo
// 	result.Archived = product.Archived
// 	// result.PriceSuggest       = product.PriceSuggest

// 	if product.Deposite != nil {
// 		result.Deposite = m.DepositeMapper.MapDepositePb(product.Deposite)
// 	}

// 	if product.PropertyType != nil {
// 		result.PropertyType = &bdspropb.PropertyType{
// 			Id:   product.PropertyType.ID,
// 			Name: product.PropertyType.Name,
// 			// Active:    product.PropertyType.Ac,
// 		}
// 	}

// 	if product.DocType != nil {
// 		result.DocType = &bdspropb.DocType{
// 			Id:   product.DocType.ID,
// 			Name: product.DocType.Name,
// 		}
// 	}

// 	// todo: map thêm category
// 	// if product.Category != nil {
// 	// 	result.Category = &bdspropb.Category{
// 	// 		Id:   product.Category.ID,
// 	// 		Name: product.Category.Name,
// 	// 	}
// 	// }

// 	if product.Apartment != nil {
// 		result.Apartment = m.ApartmentMapper.MapApartmentPb(product.Apartment)
// 	}

// 	if product.ApartmentAttr != nil {
// 		result.ApartmentAttribute = m.ApartmentAttributeMapper.MapApartmentAttributePb(product.ApartmentAttr)
// 	}

// 	if product.Build != nil {
// 		result.Build = m.BuildMapper.MapBuildPbItem(product.Build)
// 	}

// 	if product.Project != nil {
// 		result.Project = m.ProjectMapper.MapProjectPbItem(product.Project)
// 	}

// 	if product.Developer != nil {
// 		result.Developer = m.DeveloperMapper.MapDeveloperPbItem(product.Developer)
// 	}
// 	if product.SaleTransaction != nil {
// 		result.SaleTransaction = m.TransactionMapper.EntityToTransactionResponse(product.SaleTransaction)
// 	}
// 	if product.RentTransaction != nil {
// 		result.RentTransaction = m.TransactionMapper.EntityToTransactionResponse(product.RentTransaction)
// 	}

// 	// if product.Attributes != nil {
// 	// 	result.Attributes = make([]*bdspropb.AttributeProduct, len(product.Attributes))
// 	// 	for i, attribute := range product.Attributes {
// 	// 		result.Attributes[i] = &bdspropb.AttributeProduct{
// 	// 			Id:   attribute.ID,
// 	// 			Name: attribute.Name,
// 	// 		}
// 	// 	}
// 	// }

// 	return result
// }

func (m *ProductMapper) MapOriginProfilePb(op *sharepb.OriginProfile) *sharepb.OriginProfile {
	if op == nil {
		return nil
	}

	return &sharepb.OriginProfile{
		OriginId:    op.OriginId,
		DisplayName: op.DisplayName,
		Phone:       op.Phone,
		Avatar:      op.Avatar,
		OwnerOf:     uint32(op.OwnerOf),
	}
}

func (m *ProductMapper) MapDetailProductPb(ctx context.Context, product *domain.Product) *bdspropb.ProductResponse {
	result := &bdspropb.ProductResponse{
		Id:       product.ID,
		ParentId: product.ParentId,
		// DepositeId: product.DepositeID,
		// Deposite:        product.Deposite,
		AssetId:        product.AssetId,
		PropertyTypeId: product.PropertyTypeId,
		Property:       m.PropertyMapper.PropertyToPb(product.Property),
		// PropertyType:    product.PropertyType,
		DocTypeId: product.DocTypeId,
		// DocType:         product.DocType,
		// Amenities       []AmenityItem        `gorm:"many2many:product_amenity;foreignKey:ID;joinForeignKey:ProductID;References:ID;joinReferences:AmenityID" json:"amenities,omitempty"`
		Visibility:  uint32(product.Visibility),
		Name:        product.Name,
		Code:        m.CodeDataUsecase.GetProductCode(ctx, product.ID),
		CategoryId:  product.CategoryID,
		Area:        product.Area,
		PositionUrl: product.PositionUrl,
		Description: product.Description,
		// OrganizationId: *product.OrganizationID,
		// OrgStatus:   int32(product.OrgStatus),
		LastPriceId: product.LastPriceID,
		// Price:           product.Price,
		ApartmentId: product.ApartmentID,
		// Apartment:       product.Apartment,
		// ApartmentAttr:   product.ApartmentAttr,
		// Build:           product.Build,
		// Project:         product.Project,
		// Developer:       product.Developer,
		SourceType:     uint32(product.SourceType),
		SourceTypeName: enums.SourceTypeMap[product.SourceType],
		// UPDATE SOURCE RESPONSE
		SourceStatus:      uint32(product.SourceStatus),
		Priority:          uint32(product.Priority),
		SourceContactId:   product.SourceContactId,
		SourceContactNote: product.SourceContactNote,
		// CustomerId: product.CustomerID,
		Note:            product.Note,
		TransactionType: int32(product.TransactionType),
		SaleStatus:      uint32(product.SaleStatus),
		SaleVisibility:  uint32(product.SaleVisibility),
		RentStatus:      uint32(product.RentStatus),
		RentVisibility:  uint32(product.RentVisibility),
		OwnerId:         &product.OwnerID,
		OwnerType:       uint32(product.OwnerOf),
		ProvinceId:      product.ProvinceID,
		// DistrictId:      product.DistrictID,
		WardId: product.WardID,
		Address: &sharepb.AddressV3Proto{
			Detail:       product.Address,
			ProvinceId:   product.ProvinceID,
			ProvinceName: cleanAddressUnit(product.ProvinceName),
			WardId:       product.WardID,
			WardName:     cleanAddressUnit(product.WardName),
		},
		CertificateHouseNote: product.CertificateHouseNote,
		GoogleMapLink:        product.GoogleMapLink,
		// Attributes:    product.Attributes,
		// MediaList:     product.MediaList,
		// HouseInfo:     product.HouseInfo,
		Archived:  product.Archived,
		CreatedAt: _utils.FormatTimeToString(product.CreatedAt),
		UpdatedAt: _utils.FormatTimeToString(product.UpdatedAt),
		CreatedBy: product.CreatedBy,

		PostCount:        product.PostCount,
		ContactCount:     product.ContactCount,
		ShareCount:       product.ShareCount,
		AssetCount:       product.AssetCount,
		DealCount:        product.DealCount,
		AppointmentCount: product.AppointmentCount,
	}

	if result.Name == "" && product.Property != nil {
		// result.Name = product.Property.Title
	}
	// if product.Attributes != nil {
	// 	result.Attributes = make([]*bdspropb.AttributeProduct, len(product.Attributes))
	// 	for i, attribute := range product.Attributes {
	// 		result.Attributes[i] = &bdspropb.AttributeProduct{
	// 			Id:   attribute.ID,
	// 			Name: attribute.Name,
	// 		}
	// 	}
	// }
	// if product.Deposite != nil {
	// 	result.Deposite = m.DepositeMapper.MapDepositePb(product.Deposite)
	// }

	if product.PropertyType != nil {
		result.PropertyType = &bdspropb.ItemResponse{
			Id:   product.PropertyType.ID,
			Name: product.PropertyType.Name,
		}
	}

	if product.DocType != nil {
		result.DocType = &bdspropb.DocType{
			Id:   product.DocType.ID,
			Name: product.DocType.Name,
		}
	}

	if product.Amenities != nil {
		result.Amenities = make([]*bdspropb.AmenityItem, len(product.Amenities))
		for i, amenity := range product.Amenities {
			result.Amenities[i] = &bdspropb.AmenityItem{
				Id:   amenity.ID,
				Name: amenity.Name,
			}
		}
	}

	if product.Project != nil {
		result.Project = m.ProjectMapper.MapProjectPbItem(product.Project)
	}

	if product.Developer != nil {
		result.Developer = m.DeveloperMapper.MapDeveloperPbItem(product.Developer)
	}

	if product.Build != nil {
		result.Build = m.BuildMapper.MapBuildPbItem(product.Build)
	}

	if product.Apartment != nil {
		result.Apartment = m.ApartmentMapper.MapApartmentPb(product.Apartment)
	}

	if product.ApartmentAttr != nil {
		result.ApartmentAttribute = m.ApartmentAttributeMapper.MapApartmentAttributePb(product.ApartmentAttr)
	}

	if product.HouseInfo != nil {
		result.HouseInfo = m.HouseInfoToPb(product.HouseInfo)
	}

	// if product.District != nil {
	// 	result.DistrictName = cleanAddressUnit(product.District.Name)
	// }
	if product.Province != nil {
		result.ProvinceName = cleanAddressUnit(product.Province.Name)
	}
	if product.Ward != nil {
		result.WardName = cleanAddressUnit(product.Ward.Name)
	}

	if len(product.MediaList) > 0 {
		result.MediaList = make([]*bdspropb.MediaItem, len(product.MediaList))
		result.Avatar = &bdspropb.MediaItem{
			Id:        product.MediaList[0].ID,
			MediaUrl:  product.MediaList[0].MediaURL,
			MediaType: product.MediaList[0].MediaType,
		}
		for i, media := range product.MediaList {
			result.MediaList[i] = &bdspropb.MediaItem{
				Id:        media.ID,
				MediaUrl:  media.MediaURL,
				MediaType: media.MediaType,
			}
		}
	}

	if product.Price != nil {
		// Chỉ trả về PriceOwner nếu user request là chủ sở hữu của sản phẩm
		profileId := _utils.GetProfileIdWithContext(ctx)
		isOwner := product.OwnerID == profileId

		result.Price = &bdspropb.ProductPrice{
			Id:                 product.Price.ID,
			Currency:           product.Price.Currency,
			SalePrice:          product.Price.SalePrice,
			SaleCommission:     product.Price.SaleCommission,
			RentPrice:          product.Price.RentPrice,
			RentCommission:     product.Price.RentCommission,
			Deposite:           product.Price.Deposite,
			RentPaymentCycle:   product.Price.RentPaymentCycle,
			SaleCommissionType: product.Price.SaleCommissionType,
			RentCommissionType: product.Price.RentCommissionType,
			ChangeNote:         product.Price.ChangeNote,
			ChannelPrice:       product.Price.ChannelPrice,
		}

		// Chỉ set PriceOwner nếu user là owner
		if isOwner {
			result.Price.PriceOwner = product.PrivateData.ImportPrice
		}

	}

	// Map privateData nếu có (chỉ khi user là người tạo)
	if product.PrivateData != nil {
		result.PrivateData = m.PrivateDataToPb(product.PrivateData)
	}

	// if product.SaleTransaction != nil {
	// 	result.SaleTransaction = m.TransactionMapper.EntityToTransactionResponse(product.SaleTransaction)
	// }
	// if product.RentTransaction != nil {
	// 	result.RentTransaction = m.TransactionMapper.EntityToTransactionResponse(product.RentTransaction)
	// }

	// Get owner information if OwnerID is available
	if product.OwnerID != 0 && m.UserClient != nil {
		if owner := m.UserClient.GetProfileById(ctx, product.OwnerID); owner != nil {
			result.CreatedUser = owner
		}
	}

	// Get source contact information if SourceContactId is available
	if product.SourceContactId != nil && *product.SourceContactId > 0 {
		if contact, err := m.CrmClient.GetContactByOriginId(ctx, *product.SourceContactId); err == nil && contact != nil {
			result.SourceContact = contact
		}
	}

	// Add assets list - lấy tất cả assets liên kết với product qua bảng product_asset
	assets := []*bdspropb.AssetItem{}
	assetMap := make(map[uint64]bool) // Để tránh duplicate

	// Lấy assets liên kết qua bảng product_asset
	if product.ID != 0 && m.AssetRepo != nil {
		linkedAssets, err := m.AssetRepo.GetAssetsByProductID(ctx, product.ID)
		if err == nil && len(linkedAssets) > 0 {
			for _, asset := range linkedAssets {
				if !assetMap[asset.ID] {
					assetCode := m.CodeDataUsecase.GetAssetCode(ctx, asset.ID)
					assetItem := &bdspropb.AssetItem{
						Id:            asset.ID,
						Code:          assetCode,
						Name:          asset.Name,
						Area:          asset.Area,
						Address:       asset.Address,
						PurchasePrice: asset.PurchasePrice,
						LegalStatus:   uint32(asset.LegalStatus),
						RentStatus:    uint32(asset.RentStatus),
					}
					if asset.ProvinceID != nil {
						assetItem.ProvinceId = asset.ProvinceID
						if asset.Province != nil {
							assetItem.ProvinceName = asset.Province.Name
						}
					}
					if asset.WardID != nil {
						assetItem.WardId = asset.WardID
						if asset.Ward != nil {
							assetItem.WardName = asset.Ward.Name
						}
					}
					assets = append(assets, assetItem)
					assetMap[asset.ID] = true
				}
			}
		}
	}

	// Nếu có AssetId trực tiếp (không qua product_asset) và chưa được thêm
	if product.AssetId != nil && !assetMap[*product.AssetId] {
		assetCode := m.CodeDataUsecase.GetAssetCode(ctx, *product.AssetId)
		// Lấy thông tin asset nếu có
		if asset, err := m.AssetRepo.GetByID(*product.AssetId); err == nil && asset != nil {
			assetItem := &bdspropb.AssetItem{
				Id:            asset.ID,
				Code:          assetCode,
				Name:          asset.Name,
				Area:          asset.Area,
				Address:       asset.Address,
				PurchasePrice: asset.PurchasePrice,
				LegalStatus:   uint32(asset.LegalStatus),
				RentStatus:    uint32(asset.RentStatus),
			}
			if asset.ProvinceID != nil {
				assetItem.ProvinceId = asset.ProvinceID
				if asset.Province != nil {
					assetItem.ProvinceName = asset.Province.Name
				}
			}
			if asset.WardID != nil {
				assetItem.WardId = asset.WardID
				if asset.Ward != nil {
					assetItem.WardName = asset.Ward.Name
				}
			}
			assets = append(assets, assetItem)
		} else {
			// Nếu không lấy được thông tin, chỉ trả về id và code
			assets = append(assets, &bdspropb.AssetItem{
				Id:   *product.AssetId,
				Code: assetCode,
			})
		}
	}

	result.Assets = assets

	// Add posts list
	if len(product.PostIds) > 0 {
		posts := make([]*bdspropb.PostItem, len(product.PostIds))
		for i, postId := range product.PostIds {
			postCode := m.CodeDataUsecase.GetPostCode(ctx, postId)
			posts[i] = &bdspropb.PostItem{
				Id:   postId,
				Code: postCode,
			}
		}
		result.Posts = posts
	}

	return result
}

func (m *ProductMapper) MapProductPb(product *domain.Product) *bdspropb.ProductResponse {
	result := &bdspropb.ProductResponse{
		Id:   product.ID,
		Name: product.Name,
		Code: m.CodeDataUsecase.GetProductCode(nil, product.ID),
		Area: product.Area,
		Address: &sharepb.AddressV3Proto{
			Detail:       product.Address,
			ProvinceId:   product.ProvinceID,
			ProvinceName: cleanAddressUnit(product.ProvinceName),
			WardId:       product.WardID,
			WardName:     cleanAddressUnit(product.WardName),
		},
		SaleStatus: uint32(product.SaleStatus),
		RentStatus: uint32(product.RentStatus),
		AssetId:    product.AssetId,
		ImageId:    product.ImageId,

		PropertyTypeId:  product.PropertyTypeId,
		DocTypeId:       product.DocTypeId,
		Visibility:      uint32(product.Visibility),
		SourceType:      uint32(product.SourceType),
		TransactionType: int32(product.TransactionType),
		OwnerId:         &product.OwnerID,
	}

	if product.PropertyType != nil {
		result.PropertyType = &bdspropb.ItemResponse{
			Id:   product.PropertyType.ID,
			Name: product.PropertyType.Name,
		}
	}

	// if product.District != nil {
	// 	result.DistrictName = cleanAddressUnit(product.District.Name)
	// 	result.DistrictId = &product.District.ID
	// }
	if product.Province != nil {
		result.ProvinceName = cleanAddressUnit(product.Province.Name)
		result.ProvinceId = &product.Province.ID
	}
	if product.Ward != nil {
		result.WardName = cleanAddressUnit(product.Ward.Name)
		result.WardId = &product.Ward.ID
	}

	if len(product.MediaList) > 0 {
		result.MediaList = make([]*bdspropb.MediaItem, len(product.MediaList))
		for i, media := range product.MediaList {
			result.MediaList[i] = &bdspropb.MediaItem{
				Id:        media.ID,
				MediaUrl:  media.MediaURL,
				MediaType: media.MediaType,
			}
		}
	}

	if product.Price != nil {
		// MapProductPb không có context nên không trả về PriceOwner
		// Chỉ MapDetailProductPb mới trả về PriceOwner khi user là owner
		result.Price = &bdspropb.ProductPrice{
			Id:                 product.Price.ID,
			Currency:           product.Price.Currency,
			SalePrice:          product.Price.SalePrice,
			SaleCommission:     product.Price.SaleCommission,
			RentPrice:          product.Price.RentPrice,
			RentCommission:     product.Price.RentCommission,
			Deposite:           product.Price.Deposite,
			RentPaymentCycle:   product.Price.RentPaymentCycle,
			SaleCommissionType: product.Price.SaleCommissionType,
			RentCommissionType: product.Price.RentCommissionType,
		}
		// PriceOwner không được trả về ở đây vì không có context để kiểm tra owner
	}

	// if product.SaleTransaction != nil {
	// 	result.SaleTransaction = m.TransactionMapper.EntityToTransactionResponse(product.SaleTransaction)
	// }
	// if product.RentTransaction != nil {
	// 	result.RentTransaction = m.TransactionMapper.EntityToTransactionResponse(product.RentTransaction)
	// }

	// Statistics fields
	result.PublicListingCount = product.PublicListingCount

	// result.PostCount = product.PostCount
	// result.ContactCount = product.ContactCount
	// result.ShareCount = product.ShareCount
	// result.AssetCount = product.AssetCount
	// result.DealCount = product.DealCount
	result.AppointmentCount = product.AppointmentCount

	return result
}

func (m *ProductMapper) MapProductPbListResponse(products []dto.ProductListResponse) []*bdspropb.ProductResponse {
	productList := make([]*bdspropb.ProductResponse, len(products))
	for i, product := range products {

		if product.DeletedAt != nil {
			delAt := _utils.FormatTimeToString(product.DeletedAt)
			productList[i] = &bdspropb.ProductResponse{
				Id:        product.ID,
				DeletedAt: &delAt,
			}
			continue
		}
		// Clean address data nếu có
		var address *sharepb.AddressV3Proto
		if product.Address != nil {
			address = &sharepb.AddressV3Proto{
				Detail:       product.Address.Detail,
				ProvinceId:   product.Address.ProvinceId,
				ProvinceName: cleanAddressUnit(product.Address.ProvinceName),
				WardId:       product.Address.WardId,
				WardName:     cleanAddressUnit(product.Address.WardName),
			}
		}

		// Thêm numBedroom và numBathroom vào HouseInfo nếu có
		var houseInfo *bdspropb.HouseInfo
		if product.NumBedroom != nil || product.NumBathroom != nil {
			houseInfo = &bdspropb.HouseInfo{
				NumBedroom:  product.NumBedroom,
				NumBathroom: product.NumBathroom,
			}
		}

		// Map Avatar nếu có
		var avatar *bdspropb.MediaItem
		if product.AvatarURL != "" {
			avatar = &bdspropb.MediaItem{
				MediaUrl:  product.AvatarURL,
				MediaType: product.AvatarMediaType,
				Order:     int32(product.AvatarSortOrder),
			}
		}

		// Map BuildingInfo nếu có
		var buildingInfo *bdspropb.PropertyBuildingInfo
		if product.BuildingInfo != nil {
			buildingType := uint32(product.BuildingInfo.BuildingType)
			direction := uint32(product.BuildingInfo.Direction)
			balconyDirection := uint32(product.BuildingInfo.BalconyDirection)
			buildingInfo = &bdspropb.PropertyBuildingInfo{
				BuildingType: buildingType,
				// ConstructionArea: product.BuildingInfo.ConstructionArea,
				// FloorArea:        product.BuildingInfo.FloorArea,
				Floors:           product.BuildingInfo.Floors,
				Bedrooms:         product.BuildingInfo.Bedrooms,
				Bathrooms:        product.BuildingInfo.Bathrooms,
				Direction:        direction,
				BalconyDirection: balconyDirection,
			}
		}

		println("product.BuildingInfo", product.PropertyId)

		productList[i] = &bdspropb.ProductResponse{
			Id:              product.ID,
			PropertyId:      &product.PropertyId,
			Name:            product.Name,
			Code:            product.Code,
			Area:            product.Area,
			AvailableArea:   &product.AvailableArea,
			PropertyTypeId:  &product.PropertyTypeId,
			DocTypeId:       &product.DocTypeId,
			ProjectId:       &product.ProjectId,
			ProjectName:     product.ProjectName,
			ProvinceId:      &product.ProvinceId,
			ProvinceName:    cleanAddressUnit(product.ProvinceName),
			WardId:          &product.WardId,
			WardName:        cleanAddressUnit(product.WardName),
			TransactionType: product.TransactionType,
			Price: &bdspropb.ProductPrice{
				SalePrice: &product.PriceSalePrice,
				RentPrice: &product.PriceRentPrice,
			},
			PropertyType: &bdspropb.ItemResponse{
				Id:   product.PropertyTypeId,
				Name: product.PropertyTypeName,
			},
			Address:        address,
			HouseInfo:      houseInfo,
			BuildingInfo:   buildingInfo,
			ImageId:        product.ImageID,
			ImageUrl:       product.ImageURL,
			SourceType:     product.SourceType,
			SaleStatus:     product.SaleStatus,
			RentStatus:     product.RentStatus,
			SaleVisibility: product.SaleVisibility,
			RentVisibility: product.RentVisibility,
			Visibility:     product.RentVisibility,
			// Statistics fields
			PropertyTypeName:   &product.PropertyTypeName,
			PublicListingCount: product.PublicListingCount,
			ShareCount:         product.ShareCount,
			AssetCount:         product.AssetCount,
			DealCount:          product.DealCount,
			ContactCount:       product.ContactCount,
			AppointmentCount:   product.AppointmentCount,
			UpdatedAt:          _utils.FormatTimeToString(&product.UpdatedAt),
			PostCount:          product.PostCount,
			Avatar:             avatar,
		}
	}
	return productList
}

func (m *ProductMapper) MapProductPbList(products []domain.Product) []*bdspropb.ProductResponse {
	productList := make([]*bdspropb.ProductResponse, len(products))
	for i, product := range products {
		productList[i] = m.MapProductPb(&product)
	}
	return productList
}

func MapProductAttachmentPbList(ctx context.Context, codeDataUsecase *_usecase.CodeDataUsecase, products []domain.Product) []*bdspropb.ProductAttachment {
	productList := make([]*bdspropb.ProductAttachment, len(products))
	for i, product := range products {
		productList[i] = &bdspropb.ProductAttachment{
			Id:      product.ID,
			Name:    product.Name,
			Code:    codeDataUsecase.GetProductCode(ctx, product.ID),
			Area:    product.Area,
			AssetId: product.AssetId,
			// DepositeId: product.DepositeID,
			RentStatus:     uint32(product.RentStatus),
			SaleStatus:     uint32(product.SaleStatus),
			SourceType:     int32(product.SourceType),
			SourceTypeName: enums.SourceTypeMap[product.SourceType],
		}
		if product.Price != nil {
			productList[i].Price = *product.Price.SalePrice
		}
		if product.Property != nil && product.Property.PropertyInfo != nil && product.Property.PropertyInfo.Avatar != nil {
			productList[i].Image = product.Property.PropertyInfo.Avatar.MediaURL
		} else if len(product.MediaList) > 0 {
			productList[i].Image = product.MediaList[0].MediaURL
		}
	}
	return productList
}

func (m *ProductMapper) SaveProductToEntity(req *bdspropb.SaveChildRequest) *dto.ProductSaveRequest {
	result := &dto.ProductSaveRequest{
		ParentID:       &req.ParentId,
		Name:           req.Name,
		Code:           req.Code,
		AreaLand:       req.Area,
		Description:    req.Description,
		Note:           req.Note,
		PropertyTypeId: req.PropertyTypeId,
		DocTypeId:      req.DocTypeId,
		AmenityIds:     req.AmenityIds,
		ProjectId:      req.ProjectId,
		ProvinceID:     req.ProvinceId,
		// DistrictID:      req.DistrictId,
		WardID:            req.WardId,
		TransactionType:   req.TransactionType,
		GoogleMapLink:     req.GoogleMapLink,
		SourceType:        req.SourceType,
		SourceContactId:   req.SourceContactId,
		SourceContactNote: req.SourceContactNote,
		// SaleStatus:      enums.EProductStatus(req.SaleStatus),
		// RentStatus:      enums.EProductStatus(req.RentStatus),
		SaleVisibility: enums.EVisibility(req.SaleVisibility),
		RentVisibility: enums.EVisibility(req.RentVisibility),
		ImageId:        req.ImageId,
		// OwnerType:      enums.EOwnerOf(req.OwnerType),
	}

	// PostCreate * PostSaveWithProduct
	// AssetCreate * user_dto.AssetCreateWithProduct

	if req.PriceData != nil {
		result.PriceData = &domain.ProductPrice{
			Currency:           req.PriceData.Currency,
			SalePrice:          req.PriceData.SalePrice,
			SaleCommission:     req.PriceData.SaleCommission,
			RentPrice:          req.PriceData.RentPrice,
			RentCommission:     req.PriceData.RentCommission,
			Deposite:           req.PriceData.Deposite,
			SaleCommissionType: req.PriceData.SaleCommissionType,
			RentCommissionType: req.PriceData.RentCommissionType,
		}
	}

	if req.MediaItems != nil {
		result.MediaItems = make([]dto.MediaItem, len(req.MediaItems))
		for i, media := range req.MediaItems {
			result.MediaItems[i] = dto.MediaItem{
				ID:        media.Id,
				MediaURL:  media.MediaUrl,
				MediaType: media.MediaType,
				IsMain:    media.IsMain,
				Order:     int(media.Order),
			}
		}
	}

	if req.HouseInfo != nil {
		result.HouseInfo = m.HouseInfoToDomain(req.HouseInfo)
	}

	if req.ProductPrivate != nil {
		result.PrivateData = m.PrivateDataToDomain(req.ProductPrivate)
	}

	// todo: map nôt asset và post
	// if req.Post != nil {
	// 	result.PostCreate = &dto.PostSaveWithProduct{
	// 		ProductID:  req.Post.ProductId,
	// 		Title:      req.Post.Title,
	// 		Content:    req.Post.Content,
	// 		Price:      req.Post.PostPrice,
	// 		Visibility: enums.EVisibility(req.Post.Visibility),
	// 	}
	// }

	// if req.Asset != nil {
	// 	result.AssetCreate = &user_dto.AssetCreateWithProduct{
	// 		PurchasePrice: req.Asset.PurchasePrice,
	// 		PurchaseDate:  req.Asset.PurchaseDate,
	// 		LegalStatus:   req.Asset.LegalStatus,
	// 		RentStatus:    req.Asset.RentStatus,
	// 		Description:   req.Asset.Description,
	// 		LegalItems:    req.Asset.LegalItems,
	// 	}
	// }

	return result
}

func (m *ProductMapper) MapDevideChildRequest(req *bdspropb.DevideChildRequest) *dto.DevideChildRequest {
	childs := make([]dto.ProductSaveRequest, len(req.Childs))
	for i, child := range req.Childs {
		r := m.SaveProductToEntity(child)
		if r != nil {
			childs[i] = *r
		}
	}
	return &dto.DevideChildRequest{
		ParentID: &req.ParentId,
		Childs:   childs,
	}
}

func (m *ProductMapper) MapMarketProductToDTO(req *bdspropb.ProductMarketRequest) *dto.ProductSearchRequest {
	return &dto.ProductSearchRequest{
		Pagable: _dto.Pagable{
			Page: req.Page,
			Size: req.Size,
		},
		Text: req.Text,
		// ProvinceIds:      []uint64{req.ProvinceId},
		// DistrictIds:      []uint64{req.DistrictId},
		// WardIds:          []uint64{req.WardId},
		// SourceTypeIds:    []uint{req.SourceType},
		// TransactionTypes: []uint{req.TransactionType},
		SaleStatus: []uint32{uint32(req.SaleStatus)},
		RentStatus: []uint32{uint32(req.RentStatus)},
	}
}

func (m *ProductMapper) MapProductMarketToPb(ctx context.Context, product *domain.ProductMarket) *bdspropb.ProductMarket {
	result := &bdspropb.ProductMarket{
		Id:         product.ID,
		CreatedAt:  _utils.FormatTimeToString(&product.CreatedAt),
		UpdatedAt:  _utils.FormatTimeToString(&product.UpdatedAt),
		Name:       product.Name,
		Code:       m.CodeDataUsecase.GetProductCode(ctx, uint64(product.ID)),
		CategoryId: product.CategoryID,
		Area:       product.Area,
		// Description:        product.Description,
		TransactionType: product.TransactionType,
		SaleStatus:      int32(product.SaleStatus),
		SaleVisibility:  int32(product.SaleVisibility),
		RentStatus:      int32(product.RentStatus),
		RentVisibility:  int32(product.RentVisibility),
		// OwnerUserId:        product.OwnerUserID,
		// OwnerOrgId:         product.OrganizationID,
		// PropertyTypeId:     product.PropertyTypeID,
		// DocTypeId:          product.DocTypeID,
		// PositionUrl:        product.PositionUrl,
		// Address:            product.Address,
		// GoogleMapLink:      product.GoogleMapLink,
		// AssetId:            product.AssetId,
		// ApartmentId:        product.ApartmentID,
		// LastPriceId:        product.LastPriceID,
		// TransactionPrice:   product.TransactionPrice,
		// Archived:           product.Archived,
		// NumBedroom:         product.NumBedroom,
		// NumBathroom:        product.NumBathroom,
		// NumFloor:           product.NumFloor,
		// Furniture:          product.Furniture,
		// Orientation:        product.Orientation,
		// NumFront:           product.NumFront,
		// NumCarPark:         product.NumCarPark,
		// NumToilet:          product.NumToilet,
		// AmenityIds:         product.AmenityIDs,
		// MediaUrls:          product.MediaList,
		// MediaTypes:         product.MediaList,
		// MainImage:          product.MediaList[0].MediaURL,
		ProvinceName: cleanAddressUnitPtr(product.ProvinceName),
		// DistrictName:     cleanAddressUnitPtr(product.DistrictName),
		WardName:         cleanAddressUnitPtr(product.WardName),
		DocTypeName:      product.DocTypeName,
		PropertyTypeName: product.PropertyTypeName,
		SalePrice:        product.SalePrice,
		// Currency:           product.Price.Currency,
		// PriceOwner:         product.Price.PriceOwner,
		// SalePrice:          product.Price.SalePrice,
		// SaleCommission:     product.Price.SaleCommission,
		// SaleCommissionType: product.Price.SaleCommissionType,
		// RentPrice:          product.Price.RentPrice,
		// RentCommission:     product.Price.RentCommission,
		// RentCommissionType: product.Price.RentCommissionType,
		// RentPaymentCycle:   product.Price.RentPaymentCycle,
		// ImportPrice:        product.PrivateData.ImportPrice,
		// OperatingCost:      product.PrivateData.OperatingCost,
		// InternalNote:       product.PrivateData.InternalNote,
		// TargetProfit:       product.PrivateData.TargetProfit,
		// PrivateDocs:        product.PrivateData.PrivateDocs,
		// Deposite:           product.Price.Deposite,
	}
	if product.MediaURLs != nil {
		result.MediaUrls = product.MediaURLs
		result.MediaTypes = product.MediaTypes
	}
	return result
}

func (m *ProductMapper) MapProductMarketPbList(ctx context.Context, products []domain.ProductMarket) []*bdspropb.ProductMarket {
	productList := make([]*bdspropb.ProductMarket, len(products))
	for i, product := range products {
		productList[i] = m.MapProductMarketToPb(ctx, &product)
	}
	return productList
}

func (m *ProductMapper) UpdateProductPbToDTO(req *bdspropb.UpdateProductRequest) *dto.UpdateProductRequest {
	result := &dto.UpdateProductRequest{
		Name:        req.Name,
		Status:      req.Status,
		Visibility:  enums.EVisibility(req.Visibility),
		Description: req.Description,
		Note:        req.Note,

		// Property fields
		PropertyTypeID:   req.PropertyTypeId,
		AreaLand:         req.AreaLand,
		AreaTotal:        req.AreaTotal,
		Direction:        enums.EHouseOrient(req.Direction),
		BalconyDirection: enums.EHouseOrient(req.BalconyDirection),
		Bedrooms:         req.Bedrooms,
		Bathrooms:        req.Bathrooms,
		RoadWidth:        req.RoadWidth,
		Floors:           req.Floors,

		// Amenity
		AmenityIDs: req.AmenityIds,

		// Price
		Price:           req.Price,
		ImportPrice:     req.ImportPrice,
		CommissionValue: req.CommissionValue,
		CommissionType:  req.CommissionType,
		ChannelPrice:    req.ChannelPrice,
		ChangeNote:      req.ChangeNote,

		// Media
		MediaItems: m.MediaItemsToDomain(req.MediaItems),

		// Nguồn
		ResourceType:      enums.ESourceType(req.SourceType),
		SourceContactId:   req.SourceContactId,
		SourceContactNote: req.SourceContactNote,
		SourceStatus:      enums.ESourceStatus(req.SourceStatus),
		Priority:          enums.EPriority(req.Priority),
		SourceTagIDs:      req.SourceTagIds,

		ProjectID: req.ProjectId,

		// Mining
		MiningMode:        enums.EMiningMode(req.MiningMode),
		MiningScope:       req.MiningScope,
		MiningDescription: req.MiningDescription,
	}

	return result
}

// UpdateProductRequestToProductSaveRequest converts UpdateProductRequest DTO to ProductSaveRequest DTO
func (m *ProductMapper) UpdateProductRequestToProductSaveRequest(updateReq *dto.UpdateProductRequest) *dto.ProductSaveRequest {
	result := &dto.ProductSaveRequest{
		Name:              updateReq.Name,
		Status:            updateReq.Status,
		Visibility:        updateReq.Visibility,
		Description:       updateReq.Description,
		Note:              updateReq.Note,
		AreaLand:          updateReq.AreaLand,
		AreaTotal:         updateReq.AreaTotal,
		AmenityIds:        updateReq.AmenityIDs,
		SourceType:        uint32(updateReq.ResourceType),
		SourceContactNote: updateReq.SourceContactNote,
	}

	if updateReq.PropertyTypeID != nil {
		result.PropertyTypeId = updateReq.PropertyTypeID
	}

	if updateReq.Bedrooms != nil || updateReq.Bathrooms != nil || updateReq.Floors != nil || updateReq.RoadWidth != nil {
		result.HouseInfo = &dto.HouseInfoUpdate{
			OrientationHouse: updateReq.Direction,
		}
		if updateReq.Bedrooms != nil {
			bedrooms := int32(*updateReq.Bedrooms)
			result.HouseInfo.NumBedroom = &bedrooms
		}
		if updateReq.Bathrooms != nil {
			bathrooms := int32(*updateReq.Bathrooms)
			result.HouseInfo.NumBathroom = &bathrooms
		}
		if updateReq.Floors != nil {
			floors := int32(*updateReq.Floors)
			result.HouseInfo.NumFloor = &floors
		}
		if updateReq.RoadWidth != nil {
			result.HouseInfo.RoadWidth = updateReq.RoadWidth
		}
	}
	if updateReq.SourceContactId > 0 {
		result.SourceContactId = &updateReq.SourceContactId
	}
	if updateReq.ProjectID != nil {
		result.ProjectId = updateReq.ProjectID
	}
	if updateReq.Price != nil || updateReq.CommissionValue != nil {
		result.PriceData = &domain.ProductPrice{}
		if updateReq.Price != nil {
			result.PriceData.SalePrice = updateReq.Price
		}
		if updateReq.CommissionValue != nil {
			// CommissionValue và CommissionType sẽ được xử lý riêng trong usecase
		}
	}
	if updateReq.ImportPrice != nil {
		result.PrivateData = &domain.ProductPrivate{
			ImportPrice: updateReq.ImportPrice,
		}
	}
	if len(updateReq.MediaItems) > 0 {
		result.MediaItems = updateReq.MediaItems
	}

	return result
}

func (m *ProductMapper) ProductSavePbToDTO(req *bdspropb.ProductSaveRequest) *dto.ProductSaveRequest {
	result := &dto.ProductSaveRequest{
		ParentID:       req.ParentId,
		Name:           req.Name,
		Code:           req.Code,
		AreaLand:       req.Area,
		Description:    req.Description,
		Note:           req.Note,
		PropertyTypeId: req.PropertyTypeId,
		DocTypeId:      req.DocTypeId,
		AmenityIds:     req.AmenityIds,
		ProjectId:      req.ProjectId,
		ProvinceID:     req.ProvinceId,
		// DistrictID:      req.DistrictId,
		WardID:          req.WardId,
		TransactionType: req.TransactionType,
		// HouseInfo:       req.HouseInfo,
		// PrivateData:    req.PrivateData,
		GoogleMapLink:     req.GoogleMapLink,
		SourceType:        req.SourceType,
		SourceContactId:   req.SourceContactId,
		SourceContactNote: req.SourceContactNote,
		// SaleStatus:     enums.EProductStatus(req.SaleStatus),
		// RentStatus:     enums.EProductStatus(req.RentStatus),
		SaleVisibility: enums.EVisibility(req.SaleVisibility),
		RentVisibility: enums.EVisibility(req.RentVisibility),
		Visibility:     enums.EVisibility(req.Visibility),
		OwnerType:      req.OwnerType,
		OwnerId:        req.OwnerId,
		ImageId:        req.ImageId,

		CertificateHouseNote: req.CertificateHouseNote,
	}
	if req.PriceData != nil {
		result.PriceData = m.PriceDataToDomain(req.PriceData)
	}
	if req.MediaItems != nil {
		result.MediaItems = m.MediaItemsToDomain(req.MediaItems)
	}
	if req.HouseInfo != nil {
		result.HouseInfo = m.HouseInfoToDomain(req.HouseInfo)
	}
	// todo: cập nhật thông tin private hiện chưa có
	// if req.PrivateData != nil {
	// 	result.PrivateData = m.PrivateDataToDomain(req.PrivateData)
	// }
	// PostCreate:     req.PostCreate,
	// 	AssetCreate:    req.AssetCreate,
	if req.Post != nil {
		size := len(req.Post.Metadatas)
		metadatas := make([]dto.PostMediaItem, size)
		for i, item := range req.Post.Metadatas {
			metadatas[i] = dto.PostMediaItem{
				MediaURL:  item.MediaUrl,
				MediaType: item.MediaType,
				Order:     (size - i) * 10,
			}
		}

		result.PostCreate = &dto.PostSaveRequest{
			ProductID:       req.Post.ProductId,
			Title:           req.Post.Title,
			Content:         req.Post.Content,
			TransactionType: enums.TransactionType(req.Post.TransactionType),
			Visibility:      enums.EVisibility(req.Post.Visibility),
			PackageVisible:  uint(req.Post.PackageVisible),
			NumDate:         int(req.Post.NumDate),
			Hidden:          req.Post.Hidden,
			Medatadatas:     metadatas,
			// ExpiredAt:       _utils.ParseStringToTime(req.Post.ExpiredAt),
			// Price:           req.Post.PostPrice,
		}
	}

	if req.Asset != nil {
		legalItems := make([]domain.AssetLegal, len(req.Asset.LegalItems))
		for i, item := range req.Asset.LegalItems {
			legal := &domain.AssetLegal{
				DocumentURL:  item.DocumentUrl,
				DocumentName: item.DocumentName,
				DocumentType: item.DocumentType,
			}
			legalItems[i] = *legal
		}

		result.AssetCreate = &dto.AssetCreateWithProduct{
			PurchasePrice: req.Asset.PurchasePrice,
			// PurchaseDate:  req.Asset.PurchaseDate,
			LegalStatus: uint(req.Asset.LegalStatus),
			RentStatus:  enums.EAssetStatus(req.Asset.RentStatus),
			Description: req.Asset.Description,
			LegalItems:  legalItems,
		}
	}

	return result
}

func (m *ProductMapper) PriceDataToDomain(req *bdspropb.ProductPrice) *domain.ProductPrice {
	result := &domain.ProductPrice{
		ID:        req.Id,
		ProductID: &req.ProductId,
		Currency:  req.Currency,
		// PriceOwner:       req.PriceOwner,
		SalePrice:          req.SalePrice,
		SaleCommission:     req.SaleCommission,
		SaleCommissionType: req.SaleCommissionType,
		Deposite:           req.Deposite,
		RentPrice:          req.RentPrice,
		RentCommission:     req.RentCommission,
		RentCommissionType: req.RentCommissionType,
		RentPaymentCycle:   req.RentPaymentCycle,
	}
	return result
}

func (m *ProductMapper) MediaItemsToDomain(req []*bdspropb.MediaItem) []dto.MediaItem {
	result := make([]dto.MediaItem, len(req))
	for i, media := range req {
		result[i] = dto.MediaItem{
			ID:        media.Id,
			MediaURL:  media.MediaUrl,
			MediaType: media.MediaType,
		}
	}
	return result
}

func (m *ProductMapper) HouseInfoToDomain(req *bdspropb.HouseInfo) *dto.HouseInfoUpdate {
	result := &dto.HouseInfoUpdate{
		NumBedroom:  req.NumBedroom,
		NumBathroom: req.NumBathroom,
		NumFloor:    req.NumFloor,
		NumFront:    req.NumFront,
		NumCarPark:  req.NumCarPark,
		Furniture:   req.Furniture,
		Orientation: req.Orientation,
		RoadWidth:   req.RoadWidth,
		FrontWidth:  req.FrontWidth,
		BackWidth:   req.BackWidth,
		Width:       req.Width,
		Height:      req.Height,
	}

	if req.OrientationHouse != nil {
		result.OrientationHouse = enums.EHouseOrient(*req.OrientationHouse)
	}
	if req.CertificateHouse != nil {
		result.CertificateHouse = enums.EHouseCertificate(*req.CertificateHouse)
	}

	return result
}

func (m *ProductMapper) HouseInfoToPb(req *domain.HouseInfo) *bdspropb.HouseInfo {
	result := &bdspropb.HouseInfo{
		// NumBedroom:  int32(*req.NumBedroom),
		// NumBathroom: int32(*req.NumBathroom),
		// NumFloor:    int32(*req.NumFloor),
		// NumFront:    int32(*req.NumFront),
		// NumCarPark:  int32(*req.NumCarPark),
		// Furniture:   *req.Furniture,
		// Orientation: *req.Orientation,
	}

	if req.NumBedroom != nil {
		result.NumBedroom = req.NumBedroom
	}
	if req.NumBathroom != nil {
		result.NumBathroom = req.NumBathroom
	}
	if req.NumFloor != nil {
		result.NumFloor = req.NumFloor
	}
	if req.NumFront != nil {
		result.NumFront = req.NumFront
	}
	if req.NumCarPark != nil {
		result.NumCarPark = req.NumCarPark
	}
	if req.Furniture != nil {
		result.Furniture = req.Furniture
	}
	if req.Orientation != nil {
		result.Orientation = req.Orientation
	}
	if req.RoadWidth != nil {
		result.RoadWidth = req.RoadWidth
	}
	if req.FrontWidth != nil {
		result.FrontWidth = req.FrontWidth
	}
	if req.BackWidth != nil {
		result.BackWidth = req.BackWidth
	}
	if req.Width != nil {
		result.Width = req.Width
	}
	if req.Height != nil {
		result.Height = req.Height
	}

	// Map OrientationHouse và tự động set OrientationName từ enum
	orientationHouseUint32 := uint32(req.OrientationHouse)
	result.OrientationHouse = &orientationHouseUint32

	// Tự động map name từ enum
	if name, ok := enums.EHouseOrientNames[req.OrientationHouse]; ok {
		result.OrientationName = &name
	}

	// Map CertificateHouse và tự động set CertificateName từ enum
	certificateHouseUint32 := uint32(req.CertificateHouse)
	result.CertificateHouse = &certificateHouseUint32

	// Tự động map name từ enum
	if name, ok := enums.EHouseCertificateNames[req.CertificateHouse]; ok {
		result.CertificateName = &name
	}

	return result
}

func (m *ProductMapper) PrivateDataToDomain(req *bdspropb.ProductPrivate) *domain.ProductPrivate {
	return &domain.ProductPrivate{
		PrivateDocs:   req.PrivateDocs,
		ImportPrice:   &req.ImportPrice,
		OperatingCost: &req.OperatingCost,
		InternalNote:  req.InternalNote,
		TargetProfit:  &req.TargetProfit,
	}
}

func (m *ProductMapper) PrivateDataToPb(privateData *domain.ProductPrivate) *bdspropb.ProductPrivate {
	if privateData == nil {
		return nil
	}
	result := &bdspropb.ProductPrivate{
		Id:           privateData.ID,
		ProductId:    privateData.ProductId,
		InternalNote: privateData.InternalNote,
		PrivateDocs:  privateData.PrivateDocs,
	}
	if privateData.ImportPrice != nil {
		result.ImportPrice = *privateData.ImportPrice
	}
	if privateData.OperatingCost != nil {
		result.OperatingCost = *privateData.OperatingCost
	}
	if privateData.TargetProfit != nil {
		result.TargetProfit = *privateData.TargetProfit
	}
	return result
}

// cmd := exec.Command("ffmpeg",
// 			"-i", filePath,
// 			"-preset", "ultrafast",
// 			"-b:v", "8000k", // 8 Mbps video bitrate
// 			"-maxrate", "8000k", // Giữ ổn định
// 			"-bufsize", "16000k",
// 			"-profile:v", "baseline",
// 			"-level", "3.0",
// 			"-start_number", "0",
// 			"-hls_time", "5", // → ~5MB/segment
// 			"-hls_list_size", "0",
// 			"-hls_segment_filename", filepath.Join(outputDir, "segment_%03d.ts"),
// 			"-f", "hls",
// 			filepath.Join(outputDir, "index.m3u8"),
// 		)

func (m *ProductMapper) PbSearchRequestToDTO(req *bdspropb.ProductSearch) dto.ProductSearchRequest {
	result := dto.ProductSearchRequest{
		Pagable: _dto.Pagable{
			Page: req.Page,
			Size: req.Size,
			Sort: req.Sort,
		},
		Text:       req.Text,
		ParentID:   req.ParentId,
		Archived:   req.Archived,
		ProvinceId: req.ProvinceId,
		WardId:     req.WardId,
		// EndId:      endId,
	}

	// Parse string fields thành arrays (hỗ trợ query params dạng "1,2,3")
	if req.TransactionTypes != "" {
		int32Array := _utils.ParseInt32Array(req.TransactionTypes)
		result.TransactionTypes = make([]uint32, len(int32Array))
		for i, v := range int32Array {
			result.TransactionTypes[i] = uint32(v)
		}
	}
	if req.PropertyTypeIds != "" {
		result.PropertyTypeIds = _utils.ParseUint64Array(req.PropertyTypeIds)
	}
	if req.ProjectIds != "" {
		result.ProjectIds = _utils.ParseUint64Array(req.ProjectIds)
	}
	if req.SaleStatuses != "" {
		int32Array := _utils.ParseInt32Array(req.SaleStatuses)
		result.SaleStatus = make([]uint32, len(int32Array))
		for i, v := range int32Array {
			result.SaleStatus[i] = uint32(v)
		}
	}
	if req.RentStatuses != "" {
		int32Array := _utils.ParseInt32Array(req.RentStatuses)
		result.RentStatus = make([]uint32, len(int32Array))
		for i, v := range int32Array {
			result.RentStatus[i] = uint32(v)
		}
	}
	if req.SaleVisibilities != "" {
		int32Array := _utils.ParseInt32Array(req.SaleVisibilities)
		result.SaleVisibilities = make([]uint32, len(int32Array))
		for i, v := range int32Array {
			result.SaleVisibilities[i] = uint32(v)
		}
	}
	if req.RentVisibilities != "" {
		int32Array := _utils.ParseInt32Array(req.RentVisibilities)
		result.RentVisibilities = make([]uint32, len(int32Array))
		for i, v := range int32Array {
			result.RentVisibilities[i] = uint32(v)
		}
	}
	if req.SourceTypeIds != "" {
		int32Array := _utils.ParseInt32Array(req.SourceTypeIds)
		result.SourceTypeIds = make([]uint32, len(int32Array))
		for i, v := range int32Array {
			result.SourceTypeIds[i] = uint32(v)
		}
	}
	if req.Visibilities != "" {
		int32Array := _utils.ParseInt32Array(req.Visibilities)
		result.Visibilities = make([]uint32, len(int32Array))
		for i, v := range int32Array {
			result.Visibilities[i] = uint32(v)
		}
	}
	if req.ProvinceIds != "" {
		result.ProvinceIds = _utils.ParseUint64Array(req.ProvinceIds)
	}
	if req.WardIds != "" {
		result.WardIds = _utils.ParseUint64Array(req.WardIds)
	}

	// Map single fields as arrays nếu có giá trị
	if req.PropertyTypeId > 0 {
		result.PropertyTypeIds = append(result.PropertyTypeIds, req.PropertyTypeId)
	}
	if req.ProjectId > 0 {
		result.ProjectIds = append(result.ProjectIds, req.ProjectId)
	}
	if req.TransactionType > 0 {
		result.TransactionTypes = append(result.TransactionTypes, uint32(req.TransactionType))
	}

	return result
}

func (m *ProductMapper) MapAdminProductPbList(ctx context.Context, products []domain.Product) []*bdspropb.AdminProductItem {
	result := make([]*bdspropb.AdminProductItem, len(products))
	ownerIds := make(map[uint64]struct{}, len(products))
	for i, product := range products {
		ownerIds[product.OwnerID] = struct{}{}
		result[i] = m.MapAdminProductPb(ctx, &product)
	}
	owners, err := m.UserClient.GetMapByIDs(ctx, ownerIds)
	if err != nil {
		return result
	}
	for i, product := range products {
		if owner, ok := owners[product.OwnerID]; ok {
			result[i].Owner = owner
		}
	}

	return result
}

func (m *ProductMapper) MapAdminProductPb(ctx context.Context, product *domain.Product) *bdspropb.AdminProductItem {
	return &bdspropb.AdminProductItem{
		Id:              product.ID,
		Name:            product.Name,
		Area:            product.Area,
		TransactionType: uint32(product.TransactionType),
		ProvinceName:    cleanAddressUnit(product.ProvinceName),
		// DistrictName:     cleanAddressUnit(product.DistrictName),
		WardName:         cleanAddressUnit(product.WardName),
		SalePrice:        m.GetSalePrice(product.Price),
		RentPrice:        m.GetRentPrice(product.Price),
		RentStatus:       uint32(product.RentStatus),
		SaleStatus:       uint32(product.SaleStatus),
		CreatedAt:        _utils.FormatTimeToString(product.CreatedAt),
		UpdatedAt:        _utils.FormatTimeToString(product.UpdatedAt),
		AssetId:          product.AssetId,
		PostIds:          product.PostIds,
		AssetCode:        m.CodeDataUsecase.GetAssetCodePtr(ctx, product.AssetId),
		PostCodes:        m.CodeDataUsecase.GetPostCodes(ctx, product.PostIds),
		PropertyTypeName: product.PropertyTypeName,
	}
}

func (m *ProductMapper) AdminProductSearchRequestToDTO(req *bdspropb.AdminProductSearchRequest) *dto.ProductSearchRequest {
	result := &dto.ProductSearchRequest{
		Pagable: _dto.Pagable{
			Page: req.Page,
			Size: req.Size,
		},
		FromDate:         _utils.ParseStringToTime(req.FromDate),
		ToDate:           _utils.ParseStringToTime(req.ToDate),
		Name:             req.Name,
		AreaFrom:         req.AreaFrom,
		AreaTo:           req.AreaTo,
		TransactionTypes: req.TransactionTypes,
		SaleStatus:       req.SaleStatus,
		RentStatus:       req.RentStatus,
		PropertyTypeIds:  req.PropertyTypeIds,
		ProvinceIds:      req.ProvinceIds,
	}

	return result
}

// MapProductHistoryPb maps domain.ProductHistory to bdspropb.ProductHistory
func (m *ProductMapper) MapProductHistoryPb(ctx context.Context, history *domain.ProductHistory) *bdspropb.ProductHistory {
	result := &bdspropb.ProductHistory{
		Id:            history.ID,
		Action:        uint64(history.Action),
		CreatedUserId: history.CreatedUser.ProfileId,
		CreatedAt:     _utils.FormatTimeToString(history.CreatedAt),
		UpdatedAt:     _utils.FormatTimeToString(history.UpdatedAt),
	}

	if history.ParentID != nil {
		result.ParentId = *history.ParentID
	}
	if history.ProductID != nil {
		result.ProductId = *history.ProductID
	}

	// Convert ChildIds from pq.Int64Array to []uint64
	if len(history.ChildIds) > 0 {
		childIds := make([]uint64, len(history.ChildIds))
		for i, id := range history.ChildIds {
			childIds[i] = uint64(id)
		}
		result.ChildIds = childIds
	}

	// Get created user information
	if history.CreatedUser.ProfileId != 0 && m.UserClient != nil {
		if user := m.UserClient.GetProfileById(ctx, history.CreatedUser.ProfileId); user != nil {
			result.CreatedUser = user
		}
	}

	// Map child products if available
	if len(history.Childs) > 0 {
		childs := make([]*bdspropb.ProductInfo, len(history.Childs))
		for i, child := range history.Childs {
			childs[i] = &bdspropb.ProductInfo{
				Id:   child.ID,
				Name: child.Name,
				Code: m.CodeDataUsecase.GetProductCode(ctx, child.ID),
				Area: child.Area,
				// Address: child.Address,
			}
		}
		result.Childs = childs
	}

	return result
}

func (m *ProductMapper) MapTotalSearchParserToResponse(parser *dto.TotalSearchParser) *bdspropb.ProductResponse {
	response := &bdspropb.ProductResponse{
		Name:            parser.Name,
		TransactionType: int32(parser.TransactionType),
		HouseInfo: &bdspropb.HouseInfo{
			NumBedroom:  parser.Bedroom,
			NumBathroom: parser.Bathroom,
			NumFloor:    parser.Floor,
			NumFront:    parser.Frontage,
			NumCarPark:  parser.Park,
			NumToilet:   parser.Toilet,
			// Orientation: parser.Orientation,
			// Furniture:   parser.Furniture,
		},
	}

	// Map các trường từ parser
	if parser.Area != nil {
		response.Area = float64(*parser.Area)
	}

	if parser.PriceSuggest != nil {
		response.PriceSuggest = parser.PriceSuggest
	}

	if parser.Address != nil {
		response.Address = &sharepb.AddressV3Proto{
			Detail:       *parser.Address,
			ProvinceId:   &parser.Province.ID,
			ProvinceName: cleanAddressUnit(parser.Province.Name),
			WardId:       &parser.Ward.ID,
			WardName:     cleanAddressUnit(parser.Ward.Name),
		}
	}

	// Map regions
	if parser.Province != nil {
		response.ProvinceId = &parser.Province.ID
	}

	if parser.District != nil {
		response.DistrictId = &parser.District.ID
	}

	if parser.Ward != nil {
		response.WardId = &parser.Ward.ID
	}

	// Map property types (lấy property type đầu tiên nếu có)
	if len(parser.PropertyTypes) > 0 {
		response.PropertyTypeId = &parser.PropertyTypes[0].ID
	}

	// Map doc types (lấy doc type đầu tiên nếu có)
	if parser.DocTypes != nil {
		response.DocTypeId = &parser.DocTypes.ID
	}

	// Map sale price
	if parser.SalePrice != nil {
		response.Price = &bdspropb.ProductPrice{
			SalePrice: parser.SalePrice,
		}
	}

	// Map rent price
	if parser.RentPrice != nil {
		if response.Price == nil {
			response.Price = &bdspropb.ProductPrice{}
		}
		response.Price.RentPrice = parser.RentPrice
	}

	// Map amenities
	if len(parser.Amenities) > 0 {
		response.Amenities = make([]*bdspropb.AmenityItem, len(parser.Amenities))
		for i, amenity := range parser.Amenities {
			response.Amenities[i] = &bdspropb.AmenityItem{
				Id:   amenity.ID,
				Name: amenity.Name,
			}
		}
	}

	if parser.RentPrice != nil {
		response.Price.RentPrice = parser.RentPrice
	}

	return response
}

func (m *ProductMapper) MapTotalSearchParserToProductV3Proto(parser *_dto.ProductV3DTO) *sharepb.ProductV3Proto {
	response := &sharepb.ProductV3Proto{
		Name: parser.Name,
	}

	if parser.PropertyType != nil {
		response.PropertyType = &sharepb.ItemV3Proto{
			Id:   parser.PropertyType.ID,
			Name: parser.PropertyType.Name,
		}
	}

	if parser.DocType != nil {
		response.DocType = &sharepb.ItemV3Proto{
			Id:   parser.DocType.ID,
			Name: parser.DocType.Name,
		}
	}

	if len(parser.Amenities) > 0 {
		response.Amenities = make([]*sharepb.ItemV3Proto, len(parser.Amenities))
		for i, amenity := range parser.Amenities {
			response.Amenities[i] = &sharepb.ItemV3Proto{
				Id:   amenity.ID,
				Name: amenity.Name,
			}
		}
	}
	return response
}

func (m *ProductMapper) PbToArchiveProductsDTO(req *bdspropb.ArchiveProductsRequest) *dto.ArchiveProductsDTO {
	return &dto.ArchiveProductsDTO{
		IDs:      req.Ids,
		Archived: req.Archived,
	}
}
