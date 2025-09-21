package mapper

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	_dto "common/dto"
	_models "common/models"
	_utils "common/utils"
	pb_bdspro "pb/types/bdspro"
)

type ProductMapper struct {
	PostMapper               *PostMapper
	ApartmentMapper          *ApartmentMapper
	ApartmentAttributeMapper *ApartmentAttributeMapper
	BuildMapper              *BuildMapper
	ProjectMapper            *ProjectMapper
	DeveloperMapper          *DeveloperMapper
	DepositeMapper           *DepositeMapper
	// TransactionMapper        TransactionTransformer
}

func NewProductMapper(postMapper *PostMapper,
	apartmentMapper *ApartmentMapper,
	apartmentAttributeMapper *ApartmentAttributeMapper,
	buildMapper *BuildMapper,
	projectMapper *ProjectMapper,
	developerMapper *DeveloperMapper,
	depositeMapper *DepositeMapper,
	// transactionMapper TransactionTransformer,
) *ProductMapper {
	result := &ProductMapper{
		PostMapper:               postMapper,
		ApartmentMapper:          apartmentMapper,
		ApartmentAttributeMapper: apartmentAttributeMapper,
		BuildMapper:              buildMapper,
		ProjectMapper:            projectMapper,
		DeveloperMapper:          developerMapper,
		DepositeMapper:           depositeMapper,
		// TransactionMapper:        transactionMapper,
	}
	postMapper.ProductMapper = result
	return result
}

// func (m *ProductMapper) MapProductDetailPb(product *domain.Product) *pb_bdspro.ProductResponse {
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
// 		result.PropertyType = &pb_bdspro.PropertyType{
// 			Id:   product.PropertyType.ID,
// 			Name: product.PropertyType.Name,
// 			// Active:    product.PropertyType.Ac,
// 		}
// 	}

// 	if product.DocType != nil {
// 		result.DocType = &pb_bdspro.DocType{
// 			Id:   product.DocType.ID,
// 			Name: product.DocType.Name,
// 		}
// 	}

// 	// todo: map thêm category
// 	// if product.Category != nil {
// 	// 	result.Category = &pb_bdspro.Category{
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
// 	// 	result.Attributes = make([]*pb_bdspro.AttributeProduct, len(product.Attributes))
// 	// 	for i, attribute := range product.Attributes {
// 	// 		result.Attributes[i] = &pb_bdspro.AttributeProduct{
// 	// 			Id:   attribute.ID,
// 	// 			Name: attribute.Name,
// 	// 		}
// 	// 	}
// 	// }

// 	return result
// }

func (m *ProductMapper) MapDetailProductPb(product *domain.Product) *pb_bdspro.ProductResponse {
	result := &pb_bdspro.ProductResponse{
		Id:       product.ID,
		ParentId: product.ParentId,
		// DepositeId: product.DepositeID,
		// Deposite:        product.Deposite,
		AssetId:        product.AssetId,
		PropertyTypeId: product.PropertyTypeId,
		// PropertyType:    product.PropertyType,
		DocTypeId: product.DocTypeId,
		// DocType:         product.DocType,
		// Amenities       []AmenityItem        `gorm:"many2many:product_amenity;foreignKey:ID;joinForeignKey:ProductID;References:ID;joinReferences:AmenityID" json:"amenities,omitempty"`
		Visibility:  uint32(product.Visibility),
		Name:        product.Name,
		Code:        product.Code,
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
		SourceType: int32(product.SourceType),
		// CustomerId: product.CustomerID,
		// PrivateData:     product.PrivateData,
		Note:            product.Note,
		TransactionType: int32(product.TransactionType),
		SaleStatus:      uint32(product.SaleStatus),
		SaleVisibility:  uint32(product.SaleVisibility),
		RentStatus:      uint32(product.RentStatus),
		RentVisibility:  uint32(product.RentVisibility),
		OwnerId:         &product.OwnerID,
		OwnerType:       uint32(product.OwnerOf),
		ProvinceId:      product.ProvinceID,
		DistrictId:      product.DistrictID,
		WardId:          product.WardID,
		Address:         product.Address,
		GoogleMapLink:   product.GoogleMapLink,
		// Attributes:    product.Attributes,
		// MediaList:     product.MediaList,
		// HouseInfo:     product.HouseInfo,
		Archived:  product.Archived,
		CreatedAt: _utils.FormatTimeToString(product.CreatedAt),
		UpdatedAt: _utils.FormatTimeToString(product.UpdatedAt),
	}

	// if product.Attributes != nil {
	// 	result.Attributes = make([]*pb_bdspro.AttributeProduct, len(product.Attributes))
	// 	for i, attribute := range product.Attributes {
	// 		result.Attributes[i] = &pb_bdspro.AttributeProduct{
	// 			Id:   attribute.ID,
	// 			Name: attribute.Name,
	// 		}
	// 	}
	// }
	// if product.Deposite != nil {
	// 	result.Deposite = m.DepositeMapper.MapDepositePb(product.Deposite)
	// }

	if product.PropertyType != nil {
		result.PropertyType = &pb_bdspro.ItemResponse{
			Id:   product.PropertyType.ID,
			Name: product.PropertyType.Name,
		}
	}

	if product.DocType != nil {
		result.DocType = &pb_bdspro.DocType{
			Id:   product.DocType.ID,
			Name: product.DocType.Name,
		}
	}

	if product.Amenities != nil {
		result.Amenities = make([]*pb_bdspro.AmenityItem, len(product.Amenities))
		for i, amenity := range product.Amenities {
			result.Amenities[i] = &pb_bdspro.AmenityItem{
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

	if product.District != nil {
		result.DistrictName = product.District.Name
	}
	if product.Province != nil {
		result.ProvinceName = product.Province.Name
	}
	if product.Ward != nil {
		result.WardName = product.Ward.Name
	}

	if len(product.MediaList) > 0 {
		result.MediaList = make([]*pb_bdspro.MediaItem, len(product.MediaList))
		for i, media := range product.MediaList {
			result.MediaList[i] = &pb_bdspro.MediaItem{
				Id:        media.ID,
				MediaUrl:  media.MediaURL,
				MediaType: media.MediaType,
			}
		}
	}

	if product.Price != nil {
		result.Price = &pb_bdspro.ProductPrice{
			// PriceOwner:     product.Price.PriceOwner,
			Id:               product.Price.ID,
			Currency:         product.Price.Currency,
			SalePrice:        product.Price.SalePrice,
			SaleCommission:   product.Price.SaleCommission,
			RentPrice:        product.Price.RentPrice,
			RentCommission:   product.Price.RentCommission,
			Deposite:         product.Price.Deposite,
			RentPaymentCycle: product.Price.RentPaymentCycle,
		}

	}

	// if product.SaleTransaction != nil {
	// 	result.SaleTransaction = m.TransactionMapper.EntityToTransactionResponse(product.SaleTransaction)
	// }
	// if product.RentTransaction != nil {
	// 	result.RentTransaction = m.TransactionMapper.EntityToTransactionResponse(product.RentTransaction)
	// }

	return result
}

func (m *ProductMapper) MapProductPb(product *domain.Product) *pb_bdspro.ProductResponse {
	result := &pb_bdspro.ProductResponse{
		Id:         product.ID,
		Name:       product.Name,
		Area:       product.Area,
		Address:    product.Address,
		SaleStatus: uint32(product.SaleStatus),
		RentStatus: uint32(product.RentStatus),
		AssetId:    product.AssetId,
	}

	if product.District != nil {
		result.DistrictName = product.District.Name
		result.DistrictId = &product.District.ID
	}
	if product.Province != nil {
		result.ProvinceName = product.Province.Name
		result.ProvinceId = &product.Province.ID
	}
	if product.Ward != nil {
		result.WardName = product.Ward.Name
		result.WardId = &product.Ward.ID
	}

	if len(product.MediaList) > 0 {
		result.MediaList = make([]*pb_bdspro.MediaItem, len(product.MediaList))
		for i, media := range product.MediaList {
			result.MediaList[i] = &pb_bdspro.MediaItem{
				Id:        media.ID,
				MediaUrl:  media.MediaURL,
				MediaType: media.MediaType,
			}
		}
	}

	if product.Price != nil {
		result.Price = &pb_bdspro.ProductPrice{
			// PriceOwner:     product.Price.PriceOwner,
			Id:               product.Price.ID,
			Currency:         product.Price.Currency,
			SalePrice:        product.Price.SalePrice,
			SaleCommission:   product.Price.SaleCommission,
			RentPrice:        product.Price.RentPrice,
			RentCommission:   product.Price.RentCommission,
			Deposite:         product.Price.Deposite,
			RentPaymentCycle: product.Price.RentPaymentCycle,
		}
	}

	// if product.SaleTransaction != nil {
	// 	result.SaleTransaction = m.TransactionMapper.EntityToTransactionResponse(product.SaleTransaction)
	// }
	// if product.RentTransaction != nil {
	// 	result.RentTransaction = m.TransactionMapper.EntityToTransactionResponse(product.RentTransaction)
	// }

	return result
}

func (m *ProductMapper) MapProductPbList(products []domain.Product) []*pb_bdspro.ProductResponse {
	productList := make([]*pb_bdspro.ProductResponse, len(products))
	for i, product := range products {
		productList[i] = m.MapProductPb(&product)
	}
	return productList
}

func MapProductAttachmentPbList(products []domain.Product) []*pb_bdspro.ProductAttachment {
	productList := make([]*pb_bdspro.ProductAttachment, len(products))
	for i, product := range products {
		productList[i] = &pb_bdspro.ProductAttachment{
			Id:      product.ID,
			Name:    product.Name,
			Code:    product.Code,
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
		if product.MediaList != nil && len(product.MediaList) > 0 {
			productList[i].Image = product.MediaList[0].MediaURL
		}
	}
	return productList
}

func (m *ProductMapper) SaveProductToEntity(req *pb_bdspro.SaveChildRequest) *dto.ProductSaveRequest {
	result := &dto.ProductSaveRequest{
		ParentID:        &req.ParentId,
		Name:            req.Name,
		Code:            req.Code,
		Area:            req.Area,
		Description:     req.Description,
		Note:            req.Note,
		PropertyTypeId:  req.PropertyTypeId,
		DocTypeId:       req.DocTypeId,
		AmenityIds:      req.AmenityIds,
		ProjectId:       req.ProjectId,
		ProvinceID:      req.ProvinceId,
		DistrictID:      req.DistrictId,
		WardID:          req.WardId,
		TransactionType: req.TransactionType,
		GoogleMapLink:   req.GoogleMapLink,
		SourceType:      req.SourceType,
		// SaleStatus:      enums.EProductStatus(req.SaleStatus),
		// RentStatus:      enums.EProductStatus(req.RentStatus),
		SaleVisibility: enums.EVisibility(req.SaleVisibility),
		RentVisibility: enums.EVisibility(req.RentVisibility),
		// OwnerType:      enums.EOwnerOf(req.OwnerType),
	}

	// PostCreate * PostSaveWithProduct
	// AssetCreate * user_dto.AssetCreateWithProduct

	if req.PriceData != nil {
		result.PriceData = &domain.ProductPrice{
			// PriceOwner:     req.PriceData.PriceOwner,
			Currency:       req.PriceData.Currency,
			SalePrice:      req.PriceData.SalePrice,
			SaleCommission: req.PriceData.SaleCommission,
			RentPrice:      req.PriceData.RentPrice,
			RentCommission: req.PriceData.RentCommission,
			Deposite:       req.PriceData.Deposite,
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

func (m *ProductMapper) MapDevideChildRequest(req *pb_bdspro.DevideChildRequest) *dto.DevideChildRequest {
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

func (m *ProductMapper) MapMarketProductToDTO(req *pb_bdspro.ProductMarketRequest) *dto.ProductSearchRequest {
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

func (m *ProductMapper) MapProductMarketToPb(product *domain.ProductMarket) *pb_bdspro.ProductMarket {
	result := &pb_bdspro.ProductMarket{
		Id:         product.ID,
		CreatedAt:  _utils.FormatTimeToString(&product.CreatedAt),
		UpdatedAt:  _utils.FormatTimeToString(&product.UpdatedAt),
		Name:       product.Name,
		Code:       product.Code,
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
		ProvinceName:     product.ProvinceName,
		DistrictName:     product.DistrictName,
		WardName:         product.WardName,
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

func (m *ProductMapper) MapProductMarketPbList(products []domain.ProductMarket) []*pb_bdspro.ProductMarket {
	productList := make([]*pb_bdspro.ProductMarket, len(products))
	for i, product := range products {
		productList[i] = m.MapProductMarketToPb(&product)
	}
	return productList
}

func (m *ProductMapper) ProductSavePbToDTO(req *pb_bdspro.ProductSaveRequest) *dto.ProductSaveRequest {
	result := &dto.ProductSaveRequest{
		ParentID:        req.ParentId,
		Name:            req.Name,
		Code:            req.Code,
		Area:            req.Area,
		Description:     req.Description,
		Note:            req.Note,
		PropertyTypeId:  req.PropertyTypeId,
		DocTypeId:       req.DocTypeId,
		AmenityIds:      req.AmenityIds,
		ProjectId:       req.ProjectId,
		ProvinceID:      req.ProvinceId,
		DistrictID:      req.DistrictId,
		WardID:          req.WardId,
		TransactionType: req.TransactionType,
		// HouseInfo:       req.HouseInfo,
		// PrivateData:    req.PrivateData,
		GoogleMapLink: req.GoogleMapLink,
		SourceType:    req.SourceType,
		// SaleStatus:     enums.EProductStatus(req.SaleStatus),
		// RentStatus:     enums.EProductStatus(req.RentStatus),
		SaleVisibility: enums.EVisibility(req.SaleVisibility),
		RentVisibility: enums.EVisibility(req.RentVisibility),
		OwnerType:      req.OwnerType,
		OwnerId:        req.OwnerId,
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

func (m *ProductMapper) PriceDataToDomain(req *pb_bdspro.ProductPrice) *domain.ProductPrice {
	result := &domain.ProductPrice{
		BaseEntity: _models.BaseEntity{
			ID: req.Id,
		},
		ProductID: &req.ProductId,
		Currency:  req.Currency,
		// PriceOwner:       req.PriceOwner,
		SalePrice:        req.SalePrice,
		SaleCommission:   req.SaleCommission,
		Deposite:         req.Deposite,
		RentPrice:        req.RentPrice,
		RentCommission:   req.RentCommission,
		RentPaymentCycle: req.RentPaymentCycle,
	}
	return result
}

func (m *ProductMapper) MediaItemsToDomain(req []*pb_bdspro.MediaItem) []dto.MediaItem {
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

func (m *ProductMapper) HouseInfoToDomain(req *pb_bdspro.HouseInfo) *dto.HouseInfoUpdate {
	return &dto.HouseInfoUpdate{
		NumBedroom:  req.NumBedroom,
		NumBathroom: req.NumBathroom,
		NumFloor:    req.NumFloor,
		NumFront:    req.NumFront,
		NumCarPark:  req.NumCarPark,
		Furniture:   req.Furniture,
		Orientation: req.Orientation,
	}
}

func (m *ProductMapper) HouseInfoToPb(req *domain.HouseInfo) *pb_bdspro.HouseInfo {
	result := &pb_bdspro.HouseInfo{
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

	return result
}

func (m *ProductMapper) PrivateDataToDomain(req *pb_bdspro.ProductPrivate) *domain.ProductPrivate {
	return &domain.ProductPrivate{
		PrivateDocs:   req.PrivateDocs,
		ImportPrice:   &req.ImportPrice,
		OperatingCost: &req.OperatingCost,
		InternalNote:  req.InternalNote,
		TargetProfit:  &req.TargetProfit,
	}
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

func (m *ProductMapper) PbSearchRequestToDTO(req *pb_bdspro.ProductSearch) *dto.ProductSearchRequest {
	return &dto.ProductSearchRequest{
		Pagable: _dto.Pagable{
			Page: req.Page,
			Size: req.Size,
		},
	}
}
