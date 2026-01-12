package mapper

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	_dto "common/domain/dto"
	_models "common/models"
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
	return &PropertyMapper{
		ProductMapper: productMapper,
		AssetMapper:   assetMapper,
	}
}

// PbToPropertySearchDTO converts protobuf request to DTO
func (m *PropertyMapper) PbToPropertySearchDTO(req *bdspropb.PropertySearchRequest) *dto.PropertySearchDTO {
	result := &dto.PropertySearchDTO{
		Pagable: _dto.Pagable{
			Page: req.Page,
			Size: req.Size,
			Sort: req.Sort,
		},
		Text: req.Text,
	}

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

	return result
}

// PropertyToProto converts domain Property to protobuf PropertyResponse
func (m *PropertyMapper) PropertyToProto(property *domain.Property) *bdspropb.PropertyResponse {
	return m.PropertyWithCountsToProto(&dto.PropertyWithCounts{
		Property: *property,
	})
}

// PropertyToPb converts domain Property to protobuf PropertyResponse (alias for PropertyToProto)
func (m *PropertyMapper) PropertyToPb(property *domain.Property) *bdspropb.PropertyResponse {
	return m.PropertyToProto(property)
}

// PropertyWithCountsToProto converts PropertyWithCounts to protobuf PropertyResponse
func (m *PropertyMapper) PropertyWithCountsToProto(propertyWithCounts *dto.PropertyWithCounts) *bdspropb.PropertyResponse {
	property := &propertyWithCounts.Property
	result := &bdspropb.PropertyResponse{
		Id:          property.ID,
		Title:       property.Title,
		AddressText: property.AddressText,
		MapUrl:      property.MapURL,
		// BuildingBlock: property.BuildingInfoID,
		Level:            property.Level,
		UnitCode:         property.UnitCode,
		LegalStatus:      uint32(property.LegalStatus),
		LegalNote:        property.LegalNote,
		RecordStatus:     uint32(property.RecordStatus),
		LegalStatusName:  property.LegalStatus.String(),
		RecordStatusName: property.RecordStatus.String(),
		PropertyTypeId:   property.PropertyTypeID,
		LocationId:       property.LocationID,
		ProjectId:        property.ProjectID,
		Latitude:         property.Latitude,
		Longitude:        property.Longitude,
		AreaTotal:        property.AreaTotal,
		AreaLand:         property.AreaLand,
		AreaResidential:  property.AreaResidential,
		CreatedAt:        _utils.FormatTimeToString(property.CreatedAt),
		UpdatedAt:        _utils.FormatTimeToString(property.UpdatedAt),
	}

	// Add counts
	result.ProductCount = &propertyWithCounts.ProductCount
	result.AssetCount = &propertyWithCounts.AssetCount
	result.PostCount = &propertyWithCounts.PostCount

	// Populate AddressV3Proto - use cleanAddressUnit from product_mapper (same package)
	if property.ProvinceID != nil || property.WardID != nil || property.AddressText != "" {
		result.Address = &sharepb.AddressV3Proto{
			Detail:       property.AddressText,
			ProvinceId:   property.ProvinceID,
			ProvinceName: cleanAddressUnit(propertyWithCounts.ProvinceName),
			WardId:       property.WardID,
			WardName:     cleanAddressUnit(propertyWithCounts.WardName),
		}
	}

	// Populate projectName
	if propertyWithCounts.ProjectName != "" {
		result.ProjectName = &propertyWithCounts.ProjectName
	} else if property.Project != nil {
		result.ProjectName = &property.Project.Name
	}

	if property.Project != nil {
		result.Project = &sharepb.TagItem{
			Id:   property.Project.ID,
			Name: property.Project.Name,
		}
	}

	if property.PropertyType != nil {
		result.PropertyTypeName = property.PropertyType.Name
	} else {
		result.PropertyTypeName = propertyWithCounts.PropertyTypeName
	}

	if property.PropertyType != nil {
		result.PropertyType = &sharepb.TagItem{
			Id:   property.PropertyType.ID,
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
	} else {
		result.BuildingInfo = &bdspropb.PropertyBuildingInfo{
			ConstructionArea: propertyWithCounts.ConstructionArea,
			FloorArea:        propertyWithCounts.FloorArea,
			Floors:           propertyWithCounts.Floors,
			Bedrooms:         propertyWithCounts.Bedrooms,
			Bathrooms:        propertyWithCounts.Bathrooms,
		}
		if propertyWithCounts.BuildingType != nil {
			result.BuildingInfo.BuildingType = *propertyWithCounts.BuildingType
		}
		if propertyWithCounts.ConstructionArea != nil {
			result.BuildingInfo.ConstructionArea = propertyWithCounts.ConstructionArea
		}
		if propertyWithCounts.FloorArea != nil {
			result.BuildingInfo.FloorArea = propertyWithCounts.FloorArea
		}
	}

	// Populate Avatar
	if propertyWithCounts.AvatarURL != "" {
		result.Avatar = &bdspropb.PropertyMediaItem{
			Url:       propertyWithCounts.AvatarURL,
			ThumbUrl:  propertyWithCounts.AvatarThumbURL,
			MediaType: propertyWithCounts.AvatarMediaType,
			SortOrder: propertyWithCounts.AvatarSortOrder,
		}
	}

	if len(propertyWithCounts.MediaList) > 0 {
		mediaList := make([]*bdspropb.MediaItem, len(propertyWithCounts.MediaList))
		for i, media := range propertyWithCounts.MediaList {
			mediaList[i] = &bdspropb.MediaItem{
				Id:        media.ID,
				MediaUrl:  media.URL,
				MediaType: media.MediaType,
				Order:     int32(media.SortOrder),
				IsMain:    media.IsCover,
			}
		}
		result.MediaList = mediaList
	}

	if property.LandInfo != nil {
		result.LandInfo = &bdspropb.PropertyLandInfo{
			Frontage:  property.LandInfo.Frontage,
			Depth:     property.LandInfo.Depth,
			RoadWidth: property.LandInfo.RoadWidth,
			LandNote:  property.LandInfo.LandNote,
		}
	}

	if len(propertyWithCounts.Amenities) > 0 {
		result.Amenities = make([]*sharepb.TagItem, len(propertyWithCounts.Amenities))
		for i, amenity := range propertyWithCounts.Amenities {
			result.Amenities[i] = &sharepb.TagItem{
				Id:   amenity.ID,
				Name: amenity.Name,
			}
		}
	}

	return result
}

// PropertyListToProto converts list of PropertyWithCounts to protobuf PropertyResponse list
func (m *PropertyMapper) PropertyListToProto(properties []dto.PropertyWithCounts) []*bdspropb.PropertyResponse {
	result := make([]*bdspropb.PropertyResponse, 0, len(properties))
	for _, property := range properties {
		result = append(result, m.PropertyWithCountsToProto(&property))
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
		// Convert DTO to domain Product
		product := &domain.Product{}
		copier.Copy(product, productDTO)
		result.Product = product
		result.Product.Price = m.ProductMapper.PriceDataToDomain(req.Product.PriceData)
	}

	// Map Property
	if req.Property != nil {
		property := &domain.Property{
			PropertyTypeID:  req.Property.PropertyTypeId,
			Title:           req.Property.Title,
			LocationID:      req.Property.LocationId,
			AddressText:     req.Property.AddressText,
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
			Frontage:  req.LandInfo.Frontage,
			Depth:     req.LandInfo.Depth,
			RoadWidth: req.LandInfo.RoadWidth,
			LandNote:  req.LandInfo.LandNote,
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
				SortOrder: int(media.SortOrder),
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
		BaseEntity: _models.BaseEntity{
			ID: req.Id,
		},
	}

	if req.PropertyTypeId != nil {
		property.PropertyTypeID = req.PropertyTypeId
	}
	if req.Title != nil {
		property.Title = *req.Title
	}
	if req.LocationId != nil {
		property.LocationID = req.LocationId
	}
	if req.AddressText != nil {
		property.AddressText = *req.AddressText
	}
	if req.Latitude != nil {
		property.Latitude = req.Latitude
	}
	if req.Longitude != nil {
		property.Longitude = req.Longitude
	}
	if req.MapUrl != nil {
		property.MapURL = *req.MapUrl
	}
	if req.ProjectId != nil {
		property.ProjectID = req.ProjectId
	}
	if req.Level != nil {
		property.Level = *req.Level
	}
	if req.UnitCode != nil {
		property.UnitCode = *req.UnitCode
	}
	if req.AreaTotal != nil {
		property.AreaTotal = req.AreaTotal
	}
	if req.AreaLand != nil {
		property.AreaLand = req.AreaLand
	}
	if req.AreaResidential != nil {
		property.AreaResidential = req.AreaResidential
	}
	if req.LegalStatus != nil {
		property.LegalStatus = enums.EHouseCertificate(*req.LegalStatus)
	}
	if req.LegalNote != nil {
		property.LegalNote = *req.LegalNote
	}
	if req.RecordStatus != nil {
		property.RecordStatus = enums.EPropertyStatus(*req.RecordStatus)
	}
	if req.ProvinceId != nil {
		property.ProvinceID = req.ProvinceId
	}
	if req.WardId != nil {
		property.WardID = req.WardId
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
		if req.LandInfo.Frontage != nil {
			landInfo.Frontage = req.LandInfo.Frontage
		}
		if req.LandInfo.Depth != nil {
			landInfo.Depth = req.LandInfo.Depth
		}
		if req.LandInfo.RoadWidth != nil {
			landInfo.RoadWidth = req.LandInfo.RoadWidth
		}
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
