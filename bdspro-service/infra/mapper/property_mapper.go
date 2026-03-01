package mapper

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	_dto "common/domain/dto"
	_utils "common/utils"
	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"
	"time"

	"github.com/jinzhu/copier"
)

type PropertyMapper struct {
	ProductMapper *ProductMapper
	AssetMapper   *AssetMapper
}

func NewPropertyMapper(productMapper *ProductMapper, assetMapper *AssetMapper) *PropertyMapper {
	result := &PropertyMapper{
		ProductMapper: productMapper,
		AssetMapper:   assetMapper,
	}

	result.ProductMapper.PropertyMapper = result

	return result
}

func (m *PropertyMapper) PbToTagSearchDTO(
	req *bdspropb.TagSearchRequest,
) *dto.TagSearchDTO {
	page := req.Page
	if page == 0 {
		page = 1
	}
	size := req.Size
	if size == 0 {
		size = 20
	}
	return &dto.TagSearchDTO{
		Pagable: _dto.Pagable{
			Page: page,
			Size: size,
		},
		Keyword: req.Keyword,
		Type:    req.Type,
	}
}

func (m *PropertyMapper) TagToProto(
	item *dto.TagDTO,
) *bdspropb.TagResponse {

	if item == nil {
		return nil
	}

	return &bdspropb.TagResponse{
		Id:          item.ID,
		Name:        item.Name,
		Type:        item.Type,
		Description: item.Description,
		Icon:        item.Icon,
	}
}

func (m *PropertyMapper) TagListToProto(
	items []*dto.TagDTO,
) []*bdspropb.TagResponse {
	result := make([]*bdspropb.TagResponse, 0, len(items))
	for _, item := range items {
		result = append(result, m.TagToProto(item))
	}
	return result
}

// PbToPropertySearchDTO converts protobuf request to DTO
func (m *PropertyMapper) PbToPropertySearchDTO(req *bdspropb.PropertySearchRequest) (*dto.PropertySearchDTO, error) {
	result := &dto.PropertySearchDTO{
		Pagable: _dto.Pagable{
			Page: req.Page,
			Size: req.Size,
			Sort: req.Sort,
		},
		Text:     req.Text,
		RegionID: &req.RegionID,
	}

	if req.SourceType != nil {
		val := uint32(*req.SourceType)
		result.SourceType = &val
	}

	if req.Owned == "owned" {
		result.Owned = true
	}
	result.Owned = false

	if req.PropertyTypeId != nil {
		result.PropertyTypeID = req.PropertyTypeId
	}
	if req.LocationId != nil {
		result.LocationID = req.LocationId
	}
	if req.ProjectId != nil {
		result.ProjectID = req.ProjectId
	}
	if req.RecordStatus != nil {
		result.RecordStatus = *req.RecordStatus
	}

	// NEW: map các filter mới
	if req.Scope != nil {
		result.Scope = req.Scope
	}
	if req.Identified != nil {
		result.Identified = req.Identified
	}
	if req.NationalVerified != nil {
		result.NationalVerified = req.NationalVerified
	}
	if req.ProvinceId != nil {
		result.ProvinceID = req.ProvinceId
	}
	if req.DistrictId != nil {
		result.DistrictID = req.DistrictId
	}
	if req.WardId != nil {
		result.WardID = req.WardId
	}
	if req.HasAsset != nil {
		result.HasAsset = req.HasAsset
	}
	if req.HasProduct != nil {
		result.HasProduct = req.HasProduct
	}
	if req.HasListing != nil {
		result.HasListing = req.HasListing
	}

	// Parse thời gian
	if req.CreatedFrom != nil {
		result.CreatedFrom = _utils.ParseStringToTime(*req.CreatedFrom)
	}
	if req.CreatedTo != nil {
		result.CreatedTo = _utils.ParseStringToTime(*req.CreatedTo)
	}
	if req.UpdatedFrom != nil {
		result.UpdatedFrom = _utils.ParseStringToTime(*req.UpdatedFrom)
	}
	if req.UpdatedTo != nil {
		result.UpdatedTo = _utils.ParseStringToTime(*req.UpdatedTo)
	}

	return result, nil
}

func parseScope(s string) enums.EPropertyScope {
	switch s {
	case "private":
		return enums.PropertyScopePrivate
	case "shared":
		return enums.PropertyScopeShared
	case "public":
		return enums.PropertyScopePublic
	default:
		return enums.PropertyScopePrivate
	}
}

func scopeToString(scope enums.EPropertyScope) string {
	switch scope {
	case enums.PropertyScopePrivate:
		return "private"
	case enums.PropertyScopeShared:
		return "shared"
	case enums.PropertyScopePublic:
		return "public"
	default:
		return "private"
	}
}

func emptyThenDefault(s, def string) string {
	if s == "" {
		return def
	}
	return s
}

// PropertyToPb converts domain Property to protobuf PropertyResponse (alias for PropertyToProto)
func (m *PropertyMapper) PropertyToPb(property *domain.Property) *bdspropb.PropertyResponse {
	return m.PropertyToProto(property)
}

// PropertyWithCountsToProto converts PropertyWithCounts to protobuf PropertyResponse
func (m *PropertyMapper) PropertyToProto(property *domain.Property) *bdspropb.PropertyResponse {
	result := &bdspropb.PropertyResponse{
		Id:    property.ID,
		Title: property.Title,
		// MapUrl:           property.MapURL,
		// Level:            property.Level,
		// UnitCode:         property.UnitCode,
		// LegalStatus:      uint32(property.LegalStatus),
		// LegalNote:        property.LegalNote,
		// RecordStatus:     uint32(property.RecordStatus),
		// LegalStatusName:  property.LegalStatus.String(),
		// RecordStatusName: property.RecordStatus.String(),
		// PropertyTypeId:   property.PropertyTypeID,
		// LocationId:       property.LocationID,
		// ProjectId:        property.ProjectID,
		// Latitude:         property.Latitude,
		// Longitude:        property.Longitude,
		AreaTotal: property.AreaTotal,
		AreaLand:  property.AreaLand,
		// AreaResidential:  property.AreaResidential,
		CreatedAt:        _utils.FormatTimeToString(&property.CreatedAt),
		UpdatedAt:        _utils.FormatTimeToString(&property.UpdatedAt),
		RecordStatus:     uint32(property.RecordStatus),
		RecordStatusName: property.RecordStatus.String(),
		Visibility:       uint32(property.Visibility),
		// BuildingInfoId:   property.BuildingInfoID,

		// NEW fields
		Scope:      uint32(enums.EPropertyScope(property.Scope)),
		SourceType: uint32(property.SourceType),
		// Identifier:         property.Identifier,
		// NationalId:         property.NationalID,
		// NationalIdVerified: property.NationalIDVerified,
		// Add counts
		ProductCount: &property.ProductCount,
		AssetCount:   &property.AssetCount,
		PostCount:    &property.ListingCount,
	}

	// AddressV3Proto
	if property.AddressDetail != "" || property.ProvinceID != nil || property.WardID != nil {
		addr := &sharepb.AddressV3Proto{
			Detail:     property.AddressDetail,
			ProvinceId: property.ProvinceID,
			DistrictId: property.DistrictID,
			WardId:     property.WardID,
		}

		if property.Province != nil {
			addr.ProvinceName = property.Province.Name
		}
		if property.District != nil {
			addr.DistrictName = property.District.Name
		}
		if property.Ward != nil {
			addr.WardName = property.Ward.Name
		}

		result.Address = addr
	}

	// if property.RegionID != nil && property.RegionName != "" {
	// 	result.Region = &bdspropb.Region{
	// 		Id:   *property.RegionID,
	// 		Name: property.RegionName,
	// 		Code: property.RegionCode,
	// 	}
	// }

	// Project Tag
	if property.Project != nil && property.ProjectID != nil {
		result.Project = &sharepb.TagItem{
			Id:   *property.ProjectID,
			Name: property.Project.Name,
		}
	}

	// PropertyType Tag
	result.PropertyTypeName = property.PropertyType.Name
	if property.PropertyTypeID != nil && property.PropertyType.Name != "" {
		result.PropertyType = &sharepb.TagItem{
			Id:   *property.PropertyTypeID,
			Name: property.PropertyType.Name,
		}
	}

	if property.BuildingInfo != nil {
		result.BuildingInfo = &bdspropb.PropertyBuildingInfo{
			BuildingType:     uint32(property.BuildingInfo.BuildingType),
			ConstructionArea: property.BuildingInfo.ConstructionArea,
			FloorArea:        property.BuildingInfo.FloorArea,
			Floors:           property.BuildingInfo.Floors,
			Bedrooms:         property.BuildingInfo.Bedrooms,
			Bathrooms:        property.BuildingInfo.Bathrooms,
			Direction:        uint32(property.BuildingInfo.Direction),
			BalconyDirection: uint32(property.BuildingInfo.BalconyDirection),
		}
	}

	if property.Avatar != nil && property.Avatar.URL != "" {
		result.Avatar = &bdspropb.PropertyMediaItem{
			Url:       property.Avatar.URL,
			MediaType: property.Avatar.MediaType,
			SortOrder: property.Avatar.SortOrder,
		}
	}

	// Amenities
	if len(property.Amenities) > 0 {
		amenities := make([]*sharepb.TagItem, len(property.Amenities))
		for i, a := range property.Amenities {
			amenities[i] = &sharepb.TagItem{
				Id:   a.ID,
				Name: a.Name,
			}
		}
		result.Amenities = amenities
	}

	if len(property.MediaList) > 0 {
		mediaList := make([]*bdspropb.MediaItem, len(property.MediaList))
		for i, media := range property.MediaList {
			mediaList[i] = &bdspropb.MediaItem{
				MediaType: media.MediaType,
				MediaUrl:  media.URL,
			}
		}
		result.MediaList = mediaList
	}

	return result
}

// PropertyListToProto converts list of PropertyWithCounts to protobuf PropertyResponse list
func (m *PropertyMapper) PropertyListToProto(properties []*domain.Property) []*bdspropb.PropertyResponse {
	result := make([]*bdspropb.PropertyResponse, 0, len(properties))
	for _, property := range properties {
		result = append(result, m.PropertyToProto(property))
	}
	return result
}

// CreateProductPropertyPbToDTO map CreateProductPropertyRequest (pb) sang CreatePropertyProductRequest (DTO)
func (m *PropertyMapper) CreatePropertyPbToDTO(req *bdspropb.CreatePropertyRequest) *dto.CreatePropertyProductRequest {
	result := &dto.CreatePropertyProductRequest{
		Note: req.Note,
	}

	// Map Product
	if req.Product != nil {
		productDTO := m.ProductMapper.ProductSavePbToDTO(req.Product)

		product := &domain.Product{}
		copier.Copy(product, productDTO)

		if req.Product.PriceData != nil {
			product.Price = m.ProductMapper.PriceDataToDomain(req.Product.PriceData)
		}

		result.Product = product
	}

	// Map Property
	if req.Property != nil {
		property := &domain.Property{
			PropertyTypeID:  req.Property.PropertyTypeId,
			Title:           req.Property.Title,
			LocationID:      req.Property.LocationId,
			AddressDetail:   req.Property.AddressDetail,
			Latitude:        req.Property.Latitude,
			Longitude:       req.Property.Longitude,
			MapURL:          req.Property.MapUrl,
			ProjectID:       req.Property.ProjectId,
			Level:           req.Property.Level,
			UnitCode:        req.Property.UnitCode,
			AreaTotal:       req.Property.AreaTotal,
			AreaLand:        req.Property.AreaLand,
			AreaResidential: req.Property.AreaResidential,
			LegalNote:       req.Property.LegalNote,
			ProvinceID:      req.Property.ProvinceId,
			WardID:          req.Property.WardId,
			SourceType:      enums.EPropertySourceType(req.Property.SourceType),
		}
		// Convert LegalStatus string to enum (nếu có)
		if req.Property.LegalStatus != 0 {
			property.LegalStatus = enums.EHouseCertificate(req.Property.LegalStatus)
		}
		if req.Property.RecordStatus != 0 {
			property.RecordStatus = enums.EPropertyStatus(req.Property.RecordStatus)
		}
		result.Property = property
	}

	// Map PropertyLandInfo
	if req.LandInfo != nil {
		landInfo := &domain.PropertyLandInfo{
			// Frontage:  req.LandInfo.Frontage,
			Depth: req.LandInfo.Depth,
			// RoadWidth: req.LandInfo.RoadWidth,
			LandNote: req.LandInfo.LandNote,
		}
		result.LandInfo = landInfo
	}

	// Map PropertyBuildingInfo
	if req.BuildingInfo != nil {
		buildingType := enums.EBuildingTypeResidential
		if req.BuildingInfo.BuildingType != 0 {
			buildingType = enums.EBuildingType(req.BuildingInfo.BuildingType)
		}
		buildingInfo := &domain.PropertyBuildingInfo{
			BuildingType:     buildingType,
			ConstructionArea: req.BuildingInfo.ConstructionArea,
			FloorArea:        req.BuildingInfo.FloorArea,
			Floors:           req.BuildingInfo.Floors,
			Bedrooms:         req.BuildingInfo.Bedrooms,
			Bathrooms:        req.BuildingInfo.Bathrooms,
		}
		if req.BuildingInfo.Direction != 0 {
			buildingInfo.Direction = enums.EHouseOrient(req.BuildingInfo.Direction)
		}
		if req.BuildingInfo.BalconyDirection != 0 {
			buildingInfo.BalconyDirection = enums.EHouseOrient(req.BuildingInfo.BalconyDirection)
		}
		result.BuildingInfo = buildingInfo
	}

	// Map AssetLegal
	if len(req.LegalInfo) > 0 {
		legalInfo := &domain.AssetLegal{}
		// Lấy legal item đầu tiên (nếu có)
		if len(req.LegalInfo) > 0 {
			legal := req.LegalInfo[0]
			legalInfo.DocumentName = legal.DocumentName
			legalInfo.DocumentURL = legal.DocumentUrl
			legalInfo.DocumentType = legal.DocumentType
			legalInfo.Description = legal.Description
			legalInfo.RelatedSplitMergeID = legal.RelatedSplitMergeId
			if legal.IssuedDate != nil {
				issuedDate := time.Unix(*legal.IssuedDate, 0)
				legalInfo.IssuedDate = &issuedDate
			}
			if legal.ExpiryDate != nil {
				expiryDate := time.Unix(*legal.ExpiryDate, 0)
				legalInfo.ExpiryDate = &expiryDate
			}
		}
		result.LegalInfo = legalInfo
	}

	// Map PropertyMedia
	if len(req.MediaList) > 0 {
		mediaList := make([]domain.PropertyMedia, len(req.MediaList))
		for i, media := range req.MediaList {
			mediaList[i] = domain.PropertyMedia{
				MediaType: media.MediaType,
				URL:       media.Url,
				ThumbURL:  media.ThumbUrl,
				SortOrder: media.SortOrder,
			}
		}
		result.MediaList = mediaList
	}

	// Map TagIDs
	result.TagIDs = req.TagIds

	// Map Post
	if req.Post != nil {
		post := &domain.Post{}
		copier.Copy(post, req.Post)
		result.Post = post
	}

	// Map Asset
	if req.Asset != nil {
		asset := &domain.Asset{
			PurchasePrice: req.Asset.PurchasePrice,
			LegalStatus:   enums.EDocType(req.Asset.LegalStatus),
			RentStatus:    enums.EAssetStatus(req.Asset.RentStatus),
			Description:   req.Asset.Description,
		}
		if req.Asset.PurchaseDate > 0 {
			purchaseDate := time.Unix(int64(req.Asset.PurchaseDate), 0)
			asset.PurchaseDate = &purchaseDate
		}
		result.Asset = asset
	}

	return result
}

// UpdatePropertyPbToDTO converts protobuf UpdatePropertyRequest to domain Property, PropertyLandInfo, PropertyBuildingInfo
func (m *PropertyMapper) UpdatePropertyPbToDTO(req *bdspropb.UpdatePropertyRequest) (*domain.Property, *domain.PropertyLandInfo, *domain.PropertyBuildingInfo) {
	var property *domain.Property
	var landInfo *domain.PropertyLandInfo
	var buildingInfo *domain.PropertyBuildingInfo

	// Map Property - chỉ map các field được set (không nil)
	property = &domain.Property{
		ID:              req.Id,
		BuildingInfoID:  req.BuildingInfoId,
		Title:           req.Title,
		PropertyTypeID:  req.PropertyTypeId,
		LocationID:      req.LocationId,
		AddressDetail:   req.AddressDetail,
		Latitude:        req.Latitude,
		Longitude:       req.Longitude,
		MapURL:          req.MapUrl,
		Level:           req.Level,
		UnitCode:        req.UnitCode,
		ProjectID:       req.ProjectId,
		AreaTotal:       req.AreaTotal,
		AreaLand:        req.AreaLand,
		AreaResidential: req.AreaResidential,
		LegalNote:       req.LegalNote,
		ProvinceID:      req.ProvinceId,
		WardID:          req.WardId,
		RecordStatus:    enums.EPropertyStatus(req.RecordStatus),
		LegalStatus:     enums.EHouseCertificate(req.LegalStatus),
	}

	// Map AddressV3Proto nếu có
	if req.Address != nil {
		if req.Address.ProvinceId != nil {
			property.ProvinceID = req.Address.ProvinceId
		}
		if req.Address.WardId != nil {
			property.WardID = req.Address.WardId
		}
	}

	// Map PropertyLandInfo nếu có
	if req.LandInfo != nil {
		landInfo = &domain.PropertyLandInfo{
			PropertyID: req.Id,
		}
		// if req.LandInfo.Frontage != nil {
		// 	landInfo.Frontage = req.LandInfo.Frontage
		// }
		if req.LandInfo.Depth != nil {
			landInfo.Depth = req.LandInfo.Depth
		}
		// if req.LandInfo.RoadWidth != nil {
		// 	landInfo.RoadWidth = req.LandInfo.RoadWidth
		// }
		if req.LandInfo.LandNote != "" {
			landInfo.LandNote = req.LandInfo.LandNote
		}
	}

	// Map PropertyBuildingInfo nếu có
	if req.BuildingInfo != nil {
		buildingType := enums.EBuildingTypeResidential
		if req.BuildingInfo.BuildingType != 0 {
			buildingType = enums.EBuildingType(req.BuildingInfo.BuildingType)
		}
		buildingInfo = &domain.PropertyBuildingInfo{
			PropertyID:       req.Id,
			BuildingType:     buildingType,
			Direction:        enums.EHouseOrient(req.BuildingInfo.Direction),
			BalconyDirection: enums.EHouseOrient(req.BuildingInfo.BalconyDirection),
		}
		if req.BuildingInfo.ConstructionArea != nil {
			buildingInfo.ConstructionArea = req.BuildingInfo.ConstructionArea
		}
		if req.BuildingInfo.FloorArea != nil {
			buildingInfo.FloorArea = req.BuildingInfo.FloorArea
		}
		if req.BuildingInfo.Floors != nil {
			buildingInfo.Floors = req.BuildingInfo.Floors
		}
		if req.BuildingInfo.Bedrooms != nil {
			buildingInfo.Bedrooms = req.BuildingInfo.Bedrooms
		}
		if req.BuildingInfo.Bathrooms != nil {
			buildingInfo.Bathrooms = req.BuildingInfo.Bathrooms
		}
	}

	return property, landInfo, buildingInfo
}

// PropertyRelationToProto converts domain.PropertyRelation to proto PropertyRelation
func (m *PropertyMapper) PropertyRelationToProto(item *domain.PropertyRelation) *bdspropb.PropertyRelation {
	if item == nil {
		return nil
	}

	result := &bdspropb.PropertyRelation{
		PropertyId:   item.PropertyID,
		RelationType: uint64(item.RelationType),
	}

	// Set id = relationID nếu có
	if item.RelationID != nil {
		result.Id = *item.RelationID
		result.RelationId = *item.RelationID
	}

	// Map Asset nếu có
	if item.Asset != nil && m.AssetMapper != nil {
		result.Asset = m.AssetMapper.MapAssetPb(item.Asset)
	}

	// Map Product nếu có
	if item.Product != nil && m.ProductMapper != nil {
		result.Product = m.ProductMapper.MapProductPb(item.Product)
	}

	// Post sẽ được xử lý trong handler nếu cần (vì Post liên kết với Product, không phải Property trực tiếp)

	return result
}

// PropertyRelationListToProto converts []domain.PropertyRelation to []PropertyRelation proto
func (m *PropertyMapper) PropertyRelationListToProto(items []domain.PropertyRelation) []*bdspropb.PropertyRelation {
	result := make([]*bdspropb.PropertyRelation, len(items))
	for i := range items {
		result[i] = m.PropertyRelationToProto(&items[i])
	}
	return result
}

// PbToArchivePropertiesDTO converts protobuf ArchivePropertiesRequest to DTO
func (m *PropertyMapper) PbToArchivePropertiesDTO(req *bdspropb.ArchivePropertiesRequest) *dto.ArchivePropertiesDTO {
	return &dto.ArchivePropertiesDTO{
		IDs:      req.Ids,
		Archived: req.Archived,
	}
}

// PbToArchiveProductsDTO converts protobuf ArchiveProductsRequest to DTO
func (m *PropertyMapper) PbToArchiveProductsDTO(req *bdspropb.ArchiveProductsRequest) *dto.ArchiveProductsDTO {
	return &dto.ArchiveProductsDTO{
		IDs:      req.Ids,
		Archived: req.Archived,
	}
}
