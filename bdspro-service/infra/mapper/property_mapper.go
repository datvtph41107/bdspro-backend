package mapper

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	internalutils "bdspro/internal/utils"
	_dto "common/domain/dto"
	_enum "common/domain/enum"
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
func (m *PropertyMapper) PropertyToPb(property *domain.PropertyLineage) *bdspropb.PropertyResponse {
	return m.PropertyToProto(property)
}

// PropertyWithCountsToProto converts PropertyWithCounts to protobuf PropertyResponse
func (m *PropertyMapper) PropertyToProto(property *domain.PropertyLineage) *bdspropb.PropertyResponse {
	if property == nil {
		return nil
	}
	result := &bdspropb.PropertyResponse{
		Id: property.ID,
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
		// AreaTotal: property.AreaTotal,
		// AreaLand:  property.AreaLand,
		// AreaResidential:  property.AreaResidential,

		// CreatedBy:    property.CreatedBy,
		// RecordStatus: uint32(property.RecordStatus),
		// Visibility:   uint32(property.Visibility),
		// BuildingInfoId:   property.BuildingInfoID,

		// NEW fields
		// Scope:      uint32(enums.EPropertyScope(property.Scope)),
		// SourceType: uint32(property.SourceType),
		// Identifier:         property.Identifier,
		// NationalId:         property.NationalID,
		// NationalIdVerified: property.NationalIDVerified,
		// Add counts

	}

	if property.Statistic != nil {
		result.ProductCount = &property.Statistic.ProductCount
		result.AssetCount = &property.Statistic.AssetCount
		result.PostCount = &property.Statistic.ListingCount
	}

	// if property.RecordStatus.IsValid() {
	// 	result.RecordStatusName = property.RecordStatus.String()
	// }
	// if !property.CreatedAt.IsZero() {
	// 	result.CreatedAt = _utils.FormatTimeToString(&property.CreatedAt)
	// }
	if property.UpdatedAt != nil {
		result.UpdatedAt = _utils.FormatTimeToString(property.UpdatedAt)
	}

	// Info: title, project (ListMe format)
	if property.PropertyInfo != nil {
		if property.PropertyInfo.Title != "" {
			result.Title = property.PropertyInfo.Title
		}
		if property.PropertyInfo.Project != nil {
			result.ProjectName = &property.PropertyInfo.Project.Name
			result.Project = &sharepb.TagItem{
				Id:   property.PropertyInfo.Project.ID,
				Name: property.PropertyInfo.Project.Name,
			}
		}
		if property.PropertyInfo.ProjectID != nil {
			result.ProjectId = property.PropertyInfo.ProjectID
		}
	}

	// LandInfo: areaTotal (ListMe format)
	if property.LandInfo != nil && property.LandInfo.AreaTotal != nil {
		result.AreaTotal = property.LandInfo.AreaTotal
	}

	// Location: address (ListMe format)
	if property.Location != nil {
		result.AddressDetail = property.Location.AddressDetail
		if property.Location.AddressDetail != "" || property.Location.ProvinceID != nil || property.Location.WardID != nil {
			addr := &sharepb.AddressV3Proto{
				Detail:     property.Location.AddressDetail,
				ProvinceId: property.Location.ProvinceID,
				// DistrictId: property.Location.DistrictID,
				WardId: property.Location.WardID,
			}
			if property.Location.Province != nil {
				addr.ProvinceName = property.Location.Province.Name
			}
			// if property.Location.District != nil {
			// 	addr.DistrictName = property.Location.District.Name
			// }
			if property.Location.Ward != nil {
				addr.WardName = property.Location.Ward.Name
			}
			result.Address = addr
		}
	}

	// AddressV3Proto
	// if property.AddressDetail != "" || property.ProvinceID != nil || property.WardID != nil {
	// 	addr := &sharepb.AddressV3Proto{
	// 		Detail:     property.AddressDetail,
	// 		ProvinceId: property.ProvinceID,
	// 		DistrictId: property.DistrictID,
	// 		WardId:     property.WardID,
	// 	}

	// 	if property.Province != nil {
	// 		addr.ProvinceName = property.Province.Name
	// 	}
	// 	if property.District != nil {
	// 		addr.DistrictName = property.District.Name
	// 	}
	// 	if property.Ward != nil {
	// 		addr.WardName = property.Ward.Name
	// 	}

	// 	result.Address = addr
	// }

	// if property.RegionID != nil && property.RegionName != "" {
	// 	result.Region = &bdspropb.Region{
	// 		Id:   *property.RegionID,
	// 		Name: property.RegionName,
	// 		Code: property.RegionCode,
	// 	}
	// }

	// Project Tag
	// if property.Project != nil && property.ProjectID != nil {
	// 	result.Project = &sharepb.TagItem{
	// 		Id:   *property.ProjectID,
	// 		Name: property.Project.Name,
	// 	}
	// }

	// // PropertyType Tag
	// result.PropertyTypeName = property.PropertyType.Name
	// if property.PropertyTypeID != nil && property.PropertyType.Name != "" {
	// 	result.PropertyType = &sharepb.TagItem{
	// 		Id:   *property.PropertyTypeID,
	// 		Name: property.PropertyType.Name,
	// 	}
	// }

	if property.BuildingInfo != nil {
		bi := property.BuildingInfo
		result.BuildingInfo = &bdspropb.PropertyBuildingInfo{
			BuildingType:     uint32(bi.BuildingType),
			AreaActual:       bi.AreaActual,
			AreaFloor:        bi.AreaFloor,
			AreaConstruction: bi.AreaConstruction,
			Note:             bi.Note,
			Floors:           bi.Floors,
			RoomNumber:       bi.RoomNumber,
			Bedrooms:         bi.Bedrooms,
			Bathrooms:        bi.Bathrooms,
			BuildStatus:      uint32(bi.BuildStatus),
			Direction:        uint32(bi.Direction),
			BalconyDirection: uint32(bi.BalconyDirection),
		}
	}

	// Avatar (ListMe format - từ MediaList[0])
	if len(property.MediaList) > 0 {
		media := &property.MediaList[0]
		result.Avatar = &bdspropb.PropertyMediaItem{
			Url:       media.MediaURL,
			ThumbUrl:  media.ThumbURL,
			MediaType: media.MediaType,
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
				MediaUrl:  media.MediaURL,
			}
		}
		result.MediaList = mediaList
	}

	return result
}

// PropertyListToProto converts list to protobuf PropertyResponse list (legacy)
func (m *PropertyMapper) PropertyListToProto(properties []*domain.PropertyLineage) []*bdspropb.PropertyResponse {
	result := make([]*bdspropb.PropertyResponse, 0, len(properties))
	for _, property := range properties {
		result = append(result, m.PropertyToProto(property))
	}
	return result
}

// PropertyListToDetailProto converts list to PropertyDetailResponse list (format giống GetDetail)
func (m *PropertyMapper) PropertyListToDetailProto(properties []*domain.PropertyLineage) []*bdspropb.PropertyDetailResponse {
	result := make([]*bdspropb.PropertyDetailResponse, 0, len(properties))
	for _, property := range properties {
		result = append(result, m.PropertyDetailToProto(property, _enum.EDataModePartial))
	}
	return result
}

// PropertyDetailToProto converts PropertyLineage to PropertyDetailResponse
func (m *PropertyMapper) PropertyDetailToProto(property *domain.PropertyLineage, mode _enum.EDataMode) *bdspropb.PropertyDetailResponse {
	if property == nil {
		return nil
	}
	res := &bdspropb.PropertyDetailResponse{
		Id: property.ID,
		// Lineage: &bdspropb.PropertyLineageInfo{
		// 	NationalId: property.NationalID,
		// },
	}
	if property.CreatedAt != nil {
		res.CreatedAt = _utils.FormatTimeToString(property.CreatedAt)
	}
	if property.UpdatedAt != nil {
		res.UpdatedAt = _utils.FormatTimeToString(property.UpdatedAt)
	}
	if property.VerifiedNationalAt != nil {
		ms := property.VerifiedNationalAt.UnixMilli()
		res.Lineage.VerifiedNationalAt = &ms
	}

	if property.PropertyIdentify != nil && property.PropertyIdentify.ID != 0 {
		res.Identifier = &bdspropb.PropertyIdentify{
			Id:      property.PropertyIdentify.ID,
			Pid:     property.PropertyIdentify.PID,
			Version: property.PropertyIdentify.Version,
			Code:    internalutils.ConvertToBDS(property.PropertyIdentify.PID),
		}
		if property.PropertyIdentify.LineageID != nil {
			res.Identifier.LineageId = *property.PropertyIdentify.LineageID
		}
	}

	if property.Statistic != nil {
		res.Statistic = &bdspropb.PropertyStatistic{
			Id:           property.Statistic.ID,
			ProductCount: &property.Statistic.ProductCount,
			AssetCount:   &property.Statistic.AssetCount,
			ListingCount: &property.Statistic.ListingCount,
		}
	}

	if property.PropertyInfo != nil {
		if property.PropertyInfo.Avatar != nil {
			res.Avatar = &bdspropb.MediaItem{
				Id:        property.PropertyInfo.Avatar.ID,
				MediaType: property.PropertyInfo.Avatar.MediaType,
				MediaUrl:  property.PropertyInfo.Avatar.MediaURL,
			}
		}
		if property.PropertyInfo.Project != nil {
			res.Project = &bdspropb.Project{
				Id:   property.PropertyInfo.Project.ID,
				Name: property.PropertyInfo.Project.Name,
			}
		}
		if property.PropertyInfo.PropertyType != nil {
			res.PropertyType = &bdspropb.PropertyType{
				Id:   property.PropertyInfo.PropertyType.ID,
				Name: property.PropertyInfo.PropertyType.Name,
			}
		}
	}
	if property.Developer != nil {
		res.Developer = &bdspropb.Developer{
			Id:      property.Developer.ID,
			Name:    property.Developer.Name,
			Slug:    property.Developer.Slug,
			LogoUrl: property.Developer.LogoUrl,
			Phone:   property.Developer.Phone,
			Email:   property.Developer.Email,
			Website: property.Developer.Website,
		}
	}
	if property.Personalization != nil {
		res.Personalization = &bdspropb.PropertyPersonalization{
			Id:                property.Personalization.ID,
			OwnerOriginId:     property.Personalization.OwnerOriginID,
			RoleId:            property.Personalization.RoleID,
			PropertyLineageId: property.Personalization.PropertyLineageID,
			OriginProfileId:   property.Personalization.OriginProfileID,
		}
		if !property.Personalization.OwnerAt.IsZero() {
			ownerAt := _utils.FormatTimeToString(&property.Personalization.OwnerAt)
			res.Personalization.OwnerAt = &ownerAt
		}
		if property.Personalization.ArchivedAt != nil {
			archivedAt := _utils.FormatTimeToString(property.Personalization.ArchivedAt)
			res.Personalization.ArchivedAt = &archivedAt
		}
		if property.Personalization.HiddenAt != nil {
			hiddenAt := _utils.FormatTimeToString(property.Personalization.HiddenAt)
			res.Personalization.HiddenAt = &hiddenAt
		}
		recordStatus := uint32(property.Personalization.ResolveRecordStatus())
		res.Personalization.RecordStatus = &recordStatus
	}
	// areaTotal từ landInfo, areaActual từ buildingInfo
	// if property.LandInfo != nil && property.LandInfo.AreaTotal != nil {
	// 	res.AreaTotal = property.LandInfo.AreaTotal
	// }
	// if property.BuildingInfo != nil && property.BuildingInfo.AreaActual != nil {
	// 	res.AreaActual = property.BuildingInfo.AreaActual
	// }

	// BasicInfo (PropertyIdentify)
	if property.PropertyInfo != nil {
		pi := property.PropertyInfo
		res.Info = &bdspropb.PropertyInfo{
			PropertyIdentifyId: &pi.ID,
			PropertyTypeId:     pi.PropertyTypeID,
			ProjectId:          pi.ProjectID,
			Note:               &pi.Note,
			UnitCode:           &pi.UnitCode,
			Identifier:         &pi.Identifier,
			Level:              &pi.Level,
			Title:              &pi.Title,
		}
		if pi.PropertyTypeID != nil {
			res.Info.PropertyTypeId = pi.PropertyTypeID
		}
		if pi.ProjectID != nil {
			res.Info.ProjectId = pi.ProjectID
		}
		if pi.LegalStatus != 0 {
			v := uint32(pi.LegalStatus)
			res.Info.LegalStatus = &v
		}
		// if pi.PrivacyLevel != 0 {
		// 	v := uint32(pi.PrivacyLevel)
		// 	res.Info.PrivacyLevel = &v
		// }
		if property.Personalization != nil {
			v := uint32(property.Personalization.ResolveRecordStatus())
			res.Info.RecordStatus = &v
		}
		if pi.Note != "" {
			res.Info.Note = &pi.Note
		}
		if pi.UnitCode != "" {
			res.Info.UnitCode = &pi.UnitCode
		}
		if pi.Identifier != "" {
			res.Info.Identifier = &pi.Identifier
		}
		if pi.Level != "" {
			res.Info.Level = &pi.Level
		}
		if pi.Title != "" {
			res.Info.Title = &pi.Title
		}

		// if property.PropertyInfo.Avatar != nil {
		// 	res.Info.Avatar = &bdspropb.MediaItem{
		// 		Id:        property.PropertyInfo.Avatar.ID,
		// 		MediaType: property.PropertyInfo.Avatar.MediaType,
		// 		MediaUrl:  property.PropertyInfo.Avatar.MediaURL,
		// 		Order:     int32(property.PropertyInfo.Avatar.SortOrder),
		// 	}
		// }

		if property.PropertyInfo.PropertyType != nil {
			res.PropertyType = &bdspropb.PropertyType{
				Id:   property.PropertyInfo.PropertyType.ID,
				Name: property.PropertyInfo.PropertyType.Name,
			}
		}
	}

	// Location
	if property.Location != nil {
		loc := property.Location
		location := &bdspropb.PropertyLocationInfo{
			Id:         loc.ID,
			Detail:     loc.AddressDetail,
			MapUrl:     loc.MapURL,
			RegionId:   loc.RegionID,
			ProvinceId: loc.ProvinceID,
			// DistrictId: loc.DistrictID,
			WardId:    loc.WardID,
			Latitude:  loc.Latitude,
			Longitude: loc.Longitude,
		}
		if loc.Province != nil {
			location.ProvinceName = loc.Province.Name
		}
		// if loc.District != nil {
		// 	location.DistrictName = loc.District.Name
		// }
		if loc.Ward != nil {
			location.WardName = loc.Ward.Name
		}
		// if loc.Province != nil || loc.District != nil || loc.Ward != nil {
		// 	addr := &sharepb.AddressV3Proto{
		// 		Detail:     loc.AddressDetail,
		// 		ProvinceId: loc.ProvinceID,
		// 		DistrictId: loc.DistrictID,
		// 		WardId:     loc.WardID,
		// 	}

		// 	info.Address = addr
		// }
		res.Location = location
	}

	// LandInfo
	if property.LandInfo != nil {
		li := property.LandInfo
		land := &bdspropb.PropertyLandInfoDetail{
			Id:           li.ID,
			DocumentNo:   li.DocumentNo,
			IssuringAuth: li.IssuringAuth,
			Note:         li.Note,
			LandNote:     li.LandNote,
			Plot:         &li.Plot,
			Sheet:        &li.Sheet,
			DocumentType: uint32(li.DocumentType),
			AreaTotal:    li.AreaTotal,
			AreaLand:     li.AreaLand,
			AreaPlant:    li.AreaPlant,
			FrontWidth:   li.FrontWidth,
			Depth:        li.Depth,
			StreetWidth:  li.StreetWidth,
		}
		if li.PurposeUsed != 0 {
			v := uint32(li.PurposeUsed)
			land.PurposeUsed = &v
		}
		if li.ExpiredLand != nil {
			ms := _utils.FormatTimeToString(li.ExpiredLand)
			land.ExpiredLand = &ms
		}
		if li.ExpiredPlant != nil {
			ms := _utils.FormatTimeToString(li.ExpiredPlant)
			land.ExpiredPlant = &ms
		}
		res.LandInfo = land
	}

	// BuildingInfo
	if property.BuildingInfo != nil {
		bi := property.BuildingInfo
		res.BuildingInfo = &bdspropb.PropertyBuildingInfo{
			BuildingType:     uint32(bi.BuildingType),
			AreaActual:       bi.AreaActual,
			AreaFloor:        bi.AreaFloor,
			AreaConstruction: bi.AreaConstruction,
			Floors:           bi.Floors,
			RoomNumber:       bi.RoomNumber,
			Bedrooms:         bi.Bedrooms,
			Bathrooms:        bi.Bathrooms,
			BuildStatus:      uint32(bi.BuildStatus),
			Direction:        uint32(bi.Direction),
			BalconyDirection: uint32(bi.BalconyDirection),
			Note:             bi.Note,
		}
	}

	// Edvidence
	if property.Edvidence != nil {
		ed := property.Edvidence
		res.Edvidence = &bdspropb.PropertyEdvidenceInfo{
			Id:          ed.ID,
			Title:       ed.Title,
			FileId:      ed.FileID,
			Description: ed.Description,
		}
	}

	// MediaList
	if len(property.MediaList) > 0 {
		mediaList := make([]*bdspropb.MediaItem, len(property.MediaList))
		for i, media := range property.MediaList {
			mediaList[i] = &bdspropb.MediaItem{
				Id:        media.ID,
				MediaType: media.MediaType,
				MediaUrl:  media.MediaURL,
			}
		}
		res.MediaList = mediaList
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
		res.Amenities = amenities
	}

	// AreaRegions
	if len(property.AreaRegions) > 0 {
		areaRegions := make([]*sharepb.TagItem, len(property.AreaRegions))
		for i, ar := range property.AreaRegions {
			areaRegions[i] = &sharepb.TagItem{
				Id:   ar.ID,
				Name: ar.Name,
			}
		}
		res.AreaRegions = areaRegions
	}

	if property.UpdatedAt != nil {
		res.UpdatedAt = _utils.FormatTimeToString(property.UpdatedAt)
	}
	if property.CreatedAt != nil {
		res.CreatedAt = _utils.FormatTimeToString(property.CreatedAt)
	}
	if property.VerifiedNationalAt != nil {
		ms := property.VerifiedNationalAt.UnixMilli()
		res.Lineage.VerifiedNationalAt = &ms
	}

	return res
}

// CreatePropertyPbToDTO map CreatePropertyRequest (pb) sang CreatePropertyProductRequest (DTO)
func (m *PropertyMapper) CreatePropertyPbToDTO(req *bdspropb.CreatePropertyRequest) *dto.CreatePropertyProductRequest {
	result := &dto.CreatePropertyProductRequest{
		Note:       req.Note,
		TagIDs:     req.TagIds,
		AmenityIds: req.AmenityIds,
	}

	// BasicInfo (PropertyIdentify)
	if req.GetInfo() != nil {
		bi := req.GetInfo()
		result.Info = &domain.PropertyInfo{}
		if bi.SourceType != nil {
			result.Info.SourceType = enums.EPropertySourceType(*bi.SourceType)
		}
		if bi.PropertyTypeId != nil {
			result.Info.PropertyTypeID = bi.PropertyTypeId
		}
		if bi.ProjectId != nil {
			result.Info.ProjectID = bi.ProjectId
		}
		if bi.LegalStatus != nil {
			result.Info.LegalStatus = enums.EHouseCertificate(*bi.LegalStatus)
		}
		// if bi.PrivacyLevel != nil {
		// 	result.Info.PrivacyLevel = enums.EVisibility(*bi.PrivacyLevel)
		// }
		if bi.Note != nil {
			result.Info.Note = *bi.Note
		}
		if bi.UnitCode != nil {
			result.Info.UnitCode = *bi.UnitCode
		}
		if bi.Identifier != nil {
			result.Info.Identifier = *bi.Identifier
		}
		if bi.Level != nil {
			result.Info.Level = *bi.Level
		}
		if bi.Title != nil {
			result.Info.Title = *bi.Title
		}
	}

	// Location
	if req.GetLocation() != nil {
		loc := req.GetLocation()
		result.Location = &domain.PropertyLocation{
			AddressDetail: loc.AddressDetail,
			MapURL:        loc.MapUrl,
		}
		result.Location.RegionID = loc.RegionId
		result.Location.ProvinceID = loc.ProvinceId
		result.Location.WardID = loc.WardId
		result.Location.Latitude = loc.Latitude
		result.Location.Longitude = loc.Longitude
	}

	// LandInfo
	if req.GetLandInfo() != nil {
		li := req.GetLandInfo()
		landInfo := &domain.PropertyLandInfo{
			DocumentNo:   li.DocumentNo,
			IssuringAuth: li.IssuringAuth,
			Note:         li.Note,
			LandNote:     li.LandNote,
		}
		if li.DocumentType != nil {
			landInfo.DocumentType = enums.EDocUnknown.Parse(li.DocumentType)
		}
		if li.Plot != nil {
			landInfo.Plot = *li.Plot
		}
		if li.Sheet != nil {
			landInfo.Sheet = *li.Sheet
		}
		landInfo.AreaTotal = li.AreaTotal
		landInfo.AreaLand = li.AreaLand
		landInfo.AreaPlant = li.AreaPlant
		landInfo.FrontWidth = li.FrontWidth
		landInfo.Depth = li.Depth
		landInfo.StreetWidth = li.StreetWidth
		if li.PurposeUsed != nil {
			landInfo.PurposeUsed = enums.EPurposeUsed(*li.PurposeUsed)
		}
		if li.ExpiredLand != nil {
			landInfo.ExpiredLand = _utils.ParseStringToTime(*li.ExpiredLand)
		}
		if li.ExpiredPlant != nil {
			landInfo.ExpiredPlant = _utils.ParseStringToTime(*li.ExpiredPlant)
		}
		result.LandInfo = landInfo
	}

	// BuildingInfo
	if req.GetBuildingInfo() != nil {
		bi := req.GetBuildingInfo()
		buildingType := enums.EBuildingTypeResidential
		if bi.BuildingType != nil && *bi.BuildingType != 0 {
			buildingType = enums.EBuildingType(*bi.BuildingType)
		}
		buildingInfo := &domain.PropertyBuildingInfo{
			BuildingType: buildingType,
			Note:         bi.Note,
		}
		buildingInfo.AreaActual = bi.AreaActual
		buildingInfo.AreaFloor = bi.AreaFloor
		buildingInfo.AreaConstruction = bi.AreaConstruction
		buildingInfo.Floors = bi.Floors
		buildingInfo.RoomNumber = bi.RoomNumber
		buildingInfo.Bedrooms = bi.Bedrooms
		buildingInfo.Bathrooms = bi.Bathrooms
		if bi.BuildStatus != nil {
			buildingInfo.BuildStatus = enums.EBuildStatus(*bi.BuildStatus)
		}
		if bi.Direction != nil && *bi.Direction != 0 {
			buildingInfo.Direction = enums.EHouseOrient(*bi.Direction)
		}
		if bi.BalconyDirection != nil && *bi.BalconyDirection != 0 {
			buildingInfo.BalconyDirection = enums.EHouseOrient(*bi.BalconyDirection)
		}
		result.BuildingInfo = buildingInfo
	}

	// ExternalRef
	if req.GetExternalRef() != nil {
		er := req.GetExternalRef()
		extRef := &domain.PropertyExternalRef{
			SourceCode: er.SourceCode,
			Note:       er.Note,
		}
		if er.ExternalRefId != nil {
			extRef.ExternalRefID = *er.ExternalRefId
		}
		if er.SourceSystem != nil {
			extRef.SourceSystem = enums.ERefSourceSys(*er.SourceSystem)
		}
		if er.Confidence != nil {
			extRef.Confidence = *er.Confidence
		}
		if er.SyncStatus != nil {
			extRef.SyncStatus = enums.ESyncStatus(*er.SyncStatus)
		}
		result.ExternalRef = extRef
	}

	// Edvidence
	if req.GetEdvidence() != nil {
		ev := req.GetEdvidence()
		result.Edvidence = &domain.PropertyEdvidence{
			Title:       ev.Title,
			Description: ev.Description,
		}
		result.Edvidence.FileID = ev.FileId
	}

	// Lineage
	if req.GetLineage() != nil {
		ln := req.GetLineage()
		result.Lineage = &domain.PropertyLineage{NationalID: ln.GetNationalId()}
		if ln.VerifiedNationalAt != nil {
			t := time.UnixMilli(*ln.VerifiedNationalAt)
			result.Lineage.VerifiedNationalAt = &t
		}
	}

	// MediaList
	for _, media := range req.GetMediaList() {
		result.MediaList = append(result.MediaList, domain.PropertyMedia{
			MediaType: media.MediaType,
			MediaURL:  media.Url,
			ThumbURL:  media.ThumbUrl,
			SortOrder: media.SortOrder,
		})
	}

	// LegalInfo
	if len(req.GetLegalInfo()) > 0 {
		legal := req.GetLegalInfo()[0]
		legalInfo := &domain.AssetLegal{
			DocumentName: legal.DocumentName,
			DocumentURL:  legal.DocumentUrl,
			DocumentType: legal.DocumentType,
			Description:  legal.Description,
		}
		if legal.RelatedSplitMergeId != nil {
			legalInfo.RelatedSplitMergeID = legal.RelatedSplitMergeId
		}
		if legal.IssuedDate != nil {
			t := time.UnixMilli(*legal.IssuedDate)
			legalInfo.IssuedDate = &t
		}
		if legal.ExpiryDate != nil {
			t := time.UnixMilli(*legal.ExpiryDate)
			legalInfo.ExpiryDate = &t
		}
		result.LegalInfo = legalInfo
	}

	// Product
	if req.GetProduct() != nil {
		productDTO := m.ProductMapper.ProductSavePbToDTO(req.GetProduct())
		product := &domain.Product{}
		copier.Copy(product, productDTO)
		if req.GetProduct().PriceData != nil {
			product.Price = m.ProductMapper.PriceDataToDomain(req.GetProduct().PriceData)
		}
		result.Product = product
	}

	// Post
	if req.GetPost() != nil {
		post := &domain.Post{}
		copier.Copy(post, req.GetPost())
		result.Post = post
	}

	// Asset
	if req.GetAsset() != nil {
		a := req.GetAsset()
		asset := &domain.Asset{
			PurchasePrice: a.PurchasePrice,
			LegalStatus:   enums.EDocType(a.LegalStatus),
			RentStatus:    enums.EAssetStatus(a.RentStatus),
			Description:   a.Description,
		}
		if a.PurchaseDate > 0 {
			t := time.Unix(int64(a.PurchaseDate), 0)
			asset.PurchaseDate = &t
		}
		result.Asset = asset
	}

	return result
}

// UpdatePropertyPbToDTO converts protobuf UpdatePropertyRequest to domain Property, PropertyLandInfo, PropertyBuildingInfo
func (m *PropertyMapper) UpdatePropertyPbToDTO(req *bdspropb.UpdatePropertyRequest) (*domain.PropertyLineage, *domain.PropertyLandInfo, *domain.PropertyBuildingInfo) {
	var property *domain.PropertyLineage
	var landInfo *domain.PropertyLandInfo
	var buildingInfo *domain.PropertyBuildingInfo

	// Map Property - chỉ map các field được set (không nil)
	property = &domain.PropertyLineage{
		// ID:             req.Id,
		// BuildingInfoID: req.BuildingInfoId,
		// Title:          req.Title,
		// PropertyTypeID: req.PropertyTypeId,
		// LocationID:     req.LocationId,
		// AddressDetail:  req.AddressDetail,
		// Latitude:       req.Latitude,
		// Longitude:      req.Longitude,
		// MapURL:         req.MapUrl,
		// Level:          req.Level,
		// UnitCode:       req.UnitCode,
		// ProjectID:      req.ProjectId,
		// AreaTotal:      req.AreaTotal,
		// ProvinceID:     req.ProvinceId,
		// WardID:         req.WardId,
		// RecordStatus:   enums.EPropertyStatus(req.RecordStatus),
		// AreaLand:        req.AreaLand,
		// AreaResidential: req.AreaResidential,
		// LegalNote:       req.LegalNote,
		// LegalStatus:     enums.EHouseCertificate(req.LegalStatus),
	}

	// Map AddressV3Proto nếu có
	// if req.Address != nil {
	// 	if req.Address.ProvinceId != nil {
	// 		property.ProvinceID = req.Address.ProvinceId
	// 	}
	// 	if req.Address.WardId != nil {
	// 		property.WardID = req.Address.WardId
	// 	}
	// }

	// Map PropertyLandInfo nếu có
	if req.LandInfo != nil {
		landInfo = &domain.PropertyLandInfo{
			// PropertyID: req.Id,
		}
		// if req.LandInfo.Frontage != nil {
		// 	landInfo.Frontage = req.LandInfo.Frontage
		// }
		// if req.LandInfo.Depth != nil {
		// 	landInfo.Depth = req.LandInfo.Depth
		// }
		// if req.LandInfo.RoadWidth != nil {
		// 	landInfo.RoadWidth = req.LandInfo.RoadWidth
		// }
		if req.LandInfo.LandNote != nil {
			landInfo.LandNote = *req.LandInfo.LandNote
		}
	}

	// Map PropertyBuildingInfo nếu có
	if req.BuildingInfo != nil {
		buildingType := enums.EBuildingTypeResidential
		if req.BuildingInfo.BuildingType != nil {
			buildingType = enums.EBuildingType(*req.BuildingInfo.BuildingType)
		}
		buildingInfo = &domain.PropertyBuildingInfo{
			// PropertyID:       req.Id,
			BuildingType:     buildingType,
			Direction:        enums.EHouseOrient(*req.BuildingInfo.Direction),
			BalconyDirection: enums.EHouseOrient(*req.BuildingInfo.BalconyDirection),
		}
		// if req.BuildingInfo.ConstructionArea != nil {
		// 	buildingInfo.ConstructionArea = req.BuildingInfo.ConstructionArea
		// }
		// if req.BuildingInfo.FloorArea != nil {
		// 	buildingInfo.FloorArea = req.BuildingInfo.FloorArea
		// }
		// if req.BuildingInfo.Floors != nil {
		// 	buildingInfo.Floors = req.BuildingInfo.Floors
		// }
		// if req.BuildingInfo.Bedrooms != nil {
		// 	buildingInfo.Bedrooms = req.BuildingInfo.Bedrooms
		// }
		// if req.BuildingInfo.Bathrooms != nil {
		// 	buildingInfo.Bathrooms = req.BuildingInfo.Bathrooms
		// }
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

// PbToHiddenPropertiesDTO converts protobuf HiddenPropertiesRequest to DTO
func (m *PropertyMapper) PbToHiddenPropertiesDTO(req *bdspropb.HiddenPropertiesRequest) *dto.HiddenPropertiesDTO {
	return &dto.HiddenPropertiesDTO{
		IDs:    req.Ids,
		Hidden: req.Hidden,
	}
}

// PbToArchiveProductsDTO converts protobuf ArchiveProductsRequest to DTO
func (m *PropertyMapper) PbToArchiveProductsDTO(req *bdspropb.ArchiveProductsRequest) *dto.ArchiveProductsDTO {
	return &dto.ArchiveProductsDTO{
		IDs:      req.Ids,
		Archived: req.Archived,
	}
}
