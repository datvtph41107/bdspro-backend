package mapper

import (
	_dto "common/domain/dto"
	_utils "common/utils"
	tqdpb "pb/types/tqd"
	"tqd/internal/domain"
	"tqd/internal/dto"
	"tqd/internal/interface/repo"
)

// PoiMapper handles mapping between domain and proto
type PoiMapper struct{}

// NewPoiMapper creates a new PoiMapper
func NewPoiMapper() *PoiMapper {
	return &PoiMapper{}
}

// POIToDTO converts POI domain entity to POI DTO
func POIToDTO(poi *domain.POI) *dto.POIDTO {
	if poi == nil {
		return nil
	}

	poiDTO := &dto.POIDTO{
		ID:          poi.ID,
		Code:        poi.Code,
		CategoryID:  poi.CategoryID,
		Name:        poi.Name,
		Address:     poi.Address,
		Lat:         poi.Lat,
		Lng:         poi.Lng,
		GeoJSON:     poi.GeoJSON,
		Phone:       poi.Phone,
		Website:     poi.Website,
		RatingPoint: poi.RatingPoint,
		RatingCount: poi.RatingCount,
		Source:      poi.Source,
		Confidence:  poi.Confidence,
		IsActive:    poi.IsActive,
		CreatedAt:   _utils.FormatTimeToString(poi.CreatedAt),
		UpdatedAt:   _utils.FormatTimeToString(poi.UpdatedAt),
	}

	// Handle optional time fields
	if poi.LastVerifiedAt != nil {
		lastVerifiedStr := _utils.FormatTimeToString(poi.LastVerifiedAt)
		poiDTO.LastVerifiedAt = &lastVerifiedStr
	}

	if poi.NextReviewAt != nil {
		nextReviewStr := _utils.FormatTimeToString(poi.NextReviewAt)
		poiDTO.NextReviewAt = &nextReviewStr
	}

	// Map category if exists
	if poi.Category != nil {
		poiDTO.Category = &dto.PoiCategoryDTO{
			ID:   poi.Category.ID,
			Code: poi.Category.Code,
			Name: poi.Category.Name,
			Icon: poi.Category.Icon,
		}
	}

	return poiDTO
}

// POIListToDTO converts a list of POI domain entities to DTO list
func POIListToDTO(pois []domain.POI) []dto.POIDTO {
	poiDTOs := make([]dto.POIDTO, len(pois))
	for i, poi := range pois {
		poiDTO := POIToDTO(&poi)
		if poiDTO != nil {
			poiDTOs[i] = *poiDTO
		}
	}
	return poiDTOs
}

// CreatePOIRequestToDomain converts CreatePOIRequestDTO to POI domain entity
func CreatePOIRequestToDomain(req *dto.CreatePOIRequestDTO) *domain.POI {
	if req == nil {
		return nil
	}

	poi := &domain.POI{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Address:     req.Address,
		Lat:         req.Lat,
		Lng:         req.Lng,
		GeoJSON:     req.GeoJSON,
		Phone:       req.Phone,
		Website:     req.Website,
		RatingPoint: req.RatingPoint,
		RatingCount: req.RatingCount,
		Source:      req.Source,
		Confidence:  req.Confidence,
		IsActive:    req.IsActive,
	}

	// Parse optional time fields
	if req.LastVerifiedAt != nil {
		if parsedTime := _utils.ParseStringToTime(*req.LastVerifiedAt); parsedTime != nil {
			poi.LastVerifiedAt = parsedTime
		}
	}

	if req.NextReviewAt != nil {
		if parsedTime := _utils.ParseStringToTime(*req.NextReviewAt); parsedTime != nil {
			poi.NextReviewAt = parsedTime
		}
	}

	return poi
}

// UpdatePOIRequestToDomain converts UpdatePOIRequestDTO to POI domain entity
func UpdatePOIRequestToDomain(req *dto.UpdatePOIRequestDTO, existingPOI *domain.POI) *domain.POI {
	if req == nil || existingPOI == nil {
		return nil
	}

	// Update fields from request
	existingPOI.CategoryID = req.CategoryID
	existingPOI.Name = req.Name
	existingPOI.Address = req.Address
	existingPOI.Lat = req.Lat
	existingPOI.Lng = req.Lng
	existingPOI.GeoJSON = req.GeoJSON
	existingPOI.Phone = req.Phone
	existingPOI.Website = req.Website
	existingPOI.RatingPoint = req.RatingPoint
	existingPOI.RatingCount = req.RatingCount
	existingPOI.Source = req.Source
	existingPOI.Confidence = req.Confidence
	existingPOI.IsActive = req.IsActive

	// Parse optional time fields
	if req.LastVerifiedAt != nil {
		if parsedTime := _utils.ParseStringToTime(*req.LastVerifiedAt); parsedTime != nil {
			existingPOI.LastVerifiedAt = parsedTime
		}
	} else {
		existingPOI.LastVerifiedAt = nil
	}

	if req.NextReviewAt != nil {
		if parsedTime := _utils.ParseStringToTime(*req.NextReviewAt); parsedTime != nil {
			existingPOI.NextReviewAt = parsedTime
		}
	} else {
		existingPOI.NextReviewAt = nil
	}

	return existingPOI
}

// ListPOIsRequestToRepoRequest converts ListPOIsRequestDTO to repository request
func ListPOIsRequestToRepoRequest(req *dto.ListPOIsRequestDTO) *repo.ListPOIsRequest {
	if req == nil {
		return nil
	}

	return &repo.ListPOIsRequest{
		Text:       req.Text,
		CategoryID: req.CategoryID,
		IsActive:   req.IsActive,
		MinLat:     req.MinLat,
		MaxLat:     req.MaxLat,
		MinLng:     req.MinLng,
		MaxLng:     req.MaxLng,
		Pagable: _dto.Pagable{
			Page: req.Page,
			Size: req.Size,
		},
	}
}

// ToProto converts domain POI to proto Poi
func (m *PoiMapper) ToProto(poi *domain.POI) *tqdpb.Poi {
	if poi == nil {
		return nil
	}

	protoPoi := &tqdpb.Poi{
		Id:          poi.ID,
		Code:        poi.Code,
		Name:        poi.Name,
		Address:     poi.Address,
		Latitude:    poi.Lat,
		Longitude:   poi.Lng,
		CategoryId:  *poi.CategoryID,
		Phone:       poi.Phone,
		Website:     poi.Website,
		Rating:      float64(poi.RatingPoint),
		ReviewCount: uint32(poi.RatingCount),
		IsActive:    poi.IsActive,
		CreatedAt:   _utils.FormatTimeToString(poi.CreatedAt),
		UpdatedAt:   _utils.FormatTimeToString(poi.UpdatedAt),
	}

	return protoPoi
}

// ToProtoList converts a list of POI domain entities to proto list
func (m *PoiMapper) ToProtoList(pois []domain.POI) []*tqdpb.Poi {
	protoPoints := make([]*tqdpb.Poi, len(pois))
	for i, poi := range pois {
		protoPoints[i] = m.ToProto(&poi)
	}
	return protoPoints
}
