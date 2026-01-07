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
}

func NewPropertyMapper(productMapper *ProductMapper) *PropertyMapper {
	return &PropertyMapper{
		ProductMapper: productMapper,
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
	}

	if property.PropertyTypeID != nil {
		result.PropertyTypeId = property.PropertyTypeID
	}

	if property.LocationID != nil {
		result.LocationId = property.LocationID
	}

	if property.ProjectID != nil {
		result.ProjectId = property.ProjectID
	}

	if property.Latitude != nil {
		result.Latitude = property.Latitude
	}

	if property.Longitude != nil {
		result.Longitude = property.Longitude
	}

	if property.AreaTotal != nil {
		result.AreaTotal = property.AreaTotal
	}

	if property.AreaLand != nil {
		result.AreaLand = property.AreaLand
	}

	if property.AreaResidential != nil {
		result.AreaResidential = property.AreaResidential
	}

	// Convert timestamps
	if property.CreatedAt != nil && !property.CreatedAt.IsZero() {
		result.CreatedAt = _utils.FormatTimeToString(property.CreatedAt)
	}

	if property.UpdatedAt != nil && !property.UpdatedAt.IsZero() {
		result.UpdatedAt = _utils.FormatTimeToString(property.UpdatedAt)
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
	if property.PropertyType != nil {
		result.PropertyTypeName = property.PropertyType.Name
	} else {
		result.PropertyTypeName = propertyWithCounts.PropertyTypeName
	}

	// Populate BuildingInfo
	if propertyWithCounts.BuildingType != nil {
		buildingInfo := &bdspropb.PropertyBuildingInfo{
			BuildingType:     propertyWithCounts.BuildingType,
			ConstructionArea: propertyWithCounts.ConstructionArea,
			FloorArea:        propertyWithCounts.FloorArea,
			Floors:           propertyWithCounts.Floors,
			Bedrooms:         propertyWithCounts.Bedrooms,
			Bathrooms:        propertyWithCounts.Bathrooms,
			Direction:        propertyWithCounts.Direction,
			BalconyDirection: propertyWithCounts.BalconyDirection,
		}
		result.BuildingInfo = buildingInfo
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
		if req.BuildingInfo.BuildingType != nil {
			buildingType = enums.EBuildingType(*req.BuildingInfo.BuildingType)
		}
		buildingInfo := &domain.PropertyBuildingInfo{
			BuildingType:     buildingType,
			ConstructionArea: req.BuildingInfo.ConstructionArea,
			FloorArea:        req.BuildingInfo.FloorArea,
			Floors:           req.BuildingInfo.Floors,
			Bedrooms:         req.BuildingInfo.Bedrooms,
			Bathrooms:        req.BuildingInfo.Bathrooms,
		}
		if req.BuildingInfo.Direction != nil {
			buildingInfo.Direction = enums.EHouseOrient(*req.BuildingInfo.Direction)
		}
		if req.BuildingInfo.BalconyDirection != nil {
			buildingInfo.BalconyDirection = enums.EHouseOrient(*req.BuildingInfo.BalconyDirection)
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
