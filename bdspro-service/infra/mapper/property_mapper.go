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

// PropertyToProto converts domain Property to protobuf PropertyResponse
func (m *PropertyMapper) PropertyToProto(property *domain.Property) *bdspropb.PropertyResponse {
	if property == nil {
		return nil
	}
	// Tạo PropertyWithCounts tạm từ property (các trường tổng hợp sẽ trống)
	p := &dto.PropertyWithCounts{
		Property: *property,
	}
	// Gán các trường tên từ relations nếu có
	if property.Province != nil {
		p.ProvinceName = property.Province.Name
	}
	if property.District != nil {
		p.DistrictName = property.District.Name
	}
	if property.Ward != nil {
		p.WardName = property.Ward.Name
	}
	if property.Project != nil {
		p.ProjectName = property.Project.Name
	}
	if property.Block != nil {
		p.BlockName = property.Block.Name
	}
	if property.Building != nil {
		p.BuildingName = property.Building.Name
	}
	if property.PropertyType != nil {
		p.PropertyTypeName = property.PropertyType.Name
	}
	if property.BuildingInfo != nil {
		bt := uint32(property.BuildingInfo.BuildingType)
		p.BuildingType = &bt
		p.ConstructionArea = property.BuildingInfo.ConstructionArea
		p.FloorArea = property.BuildingInfo.FloorArea
		p.Floors = property.BuildingInfo.Floors
		p.Bedrooms = property.BuildingInfo.Bedrooms
		p.Bathrooms = property.BuildingInfo.Bathrooms
		d := uint32(property.BuildingInfo.Direction)
		p.Direction = &d
		bd := uint32(property.BuildingInfo.BalconyDirection)
		p.BalconyDirection = &bd
	}
	if property.Avatar != nil {
		p.AvatarURL = property.Avatar.URL
		p.AvatarThumbURL = property.Avatar.ThumbURL
		p.AvatarMediaType = property.Avatar.MediaType
		p.AvatarSortOrder = property.Avatar.SortOrder
	}

	return m.PropertyWithCountsToProto(p)
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

func safeUint64(v *uint64) uint64 {
	if v == nil {
		return 0
	}
	return *v
}

func safeUint32(v *uint32) uint32 {
	if v == nil {
		return 0
	}
	return *v
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
func (m *PropertyMapper) PropertyWithCountsToProto(p *dto.PropertyWithCounts) *bdspropb.PropertyResponse {
	property := &p.Property
	result := &bdspropb.PropertyResponse{
		Id:               property.ID,
		Title:            property.Title,
		AddressDetail:    property.AddressDetail,
		MapUrl:           property.MapURL,
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
		BuildingInfoId:   property.BuildingInfoID,

		// NEW fields
		Scope:              uint32(enums.EPropertyScope(property.Scope)),
		Identifier:         property.Identifier,
		NationalId:         property.NationalID,
		NationalIdVerified: property.NationalIDVerified,
		ProvinceId:         safeUint64(property.ProvinceID),
		DistrictId:         safeUint64(property.DistrictID),
		WardId:             safeUint64(property.WardID),
		BlockId:            safeUint64(property.BlockID),
		BuildingId:         safeUint64(property.BuildingID),
	}

	// Add counts
	result.ProductCount = &p.ProductCount
	result.AssetCount = &p.AssetCount
	result.PostCount = &p.PostCount

	// Địa chỉ từ các trường tên
	result.ProvinceName = p.ProvinceName
	result.DistrictName = p.DistrictName
	result.WardName = p.WardName
	result.ProjectName = &p.ProjectName
	result.BlockName = p.BlockName
	result.BuildingName = p.BuildingName

	// AddressV3Proto
	if property.AddressDetail != "" || property.ProvinceID != nil || property.WardID != nil {
		result.Address = &sharepb.AddressV3Proto{
			Detail:       property.AddressDetail,
			ProvinceName: p.ProvinceName,
			ProvinceId:   property.ProvinceID,
			DistrictName: p.DistrictName,
			DistrictId:   property.DistrictID,
			WardName:     p.WardName,
			WardId:       property.WardID,
		}
	}

	// Project Tag
	if property.ProjectID != nil && p.ProjectName != "" {
		result.Project = &sharepb.TagItem{
			Id:   *property.ProjectID,
			Name: p.ProjectName,
		}
	}

	// PropertyType Tag
	result.PropertyTypeName = p.PropertyTypeName
	if property.PropertyTypeID != nil && p.PropertyTypeName != "" {
		result.PropertyType = &sharepb.TagItem{
			Id:   *property.PropertyTypeID,
			Name: p.PropertyTypeName,
		}
	}

	// BuildingInfo
	result.BuildingInfo = &bdspropb.PropertyBuildingInfo{
		BuildingType:     safeUint32(p.BuildingType),
		ConstructionArea: p.ConstructionArea,
		FloorArea:        p.FloorArea,
		Floors:           p.Floors,
		Bedrooms:         p.Bedrooms,
		Bathrooms:        p.Bathrooms,
		Direction:        safeUint32(p.Direction),
		BalconyDirection: safeUint32(p.BalconyDirection),
	}

	// Avatar
	if p.AvatarURL != "" {
		result.Avatar = &bdspropb.PropertyMediaItem{
			Url:       p.AvatarURL,
			ThumbUrl:  p.AvatarThumbURL,
			MediaType: p.AvatarMediaType,
			SortOrder: p.AvatarSortOrder,
		}
	}

	// LandInfo - map từ domain (FrontWidth, StreetWidth) sang proto (Frontage, RoadWidth)
	if property.LandInfo != nil {
		result.LandInfo = &bdspropb.PropertyLandInfo{
			Frontage:  property.LandInfo.FrontWidth,
			Depth:     property.LandInfo.Depth,
			RoadWidth: property.LandInfo.StreetWidth,
			LandNote:  property.LandInfo.LandNote,
		}
	}

	// Amenities
	if len(p.Amenities) > 0 {
		amenities := make([]*sharepb.TagItem, len(p.Amenities))
		for i, a := range p.Amenities {
			amenities[i] = &sharepb.TagItem{
				Id:   a.ID,
				Name: a.Name,
			}
		}
		result.Amenities = amenities
	}

	if len(p.MediaList) > 0 {
		mediaList := make([]*bdspropb.MediaItem, len(p.MediaList))
		for i, media := range p.MediaList {
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
		BaseEntity: _models.BaseEntity{
			ID: req.Id,
		},
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
