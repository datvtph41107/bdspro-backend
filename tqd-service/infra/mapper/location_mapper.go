// internal/mapper/location_mapper.go
package mapper

import (
	_utils "common/utils"
	tqdpb "pb/types/tqd"
	"tqd/internal/dto"
)

type LocationMapper struct{}

func NewLocationMapper() *LocationMapper {
	return &LocationMapper{}
}

// ==================== DTO -> PROTO ====================

func (m *LocationMapper) ToProtoProvince(dto *dto.ProvinceDTO) *tqdpb.Province {
	if dto == nil {
		return nil
	}
	return &tqdpb.Province{
		Id:        dto.ID,
		Code:      dto.Code,
		FullName:  dto.FullName,
		ShortName: dto.ShortName,
		Latitude:  dto.Lat,
		Longitude: dto.Lng,
		WardCount: dto.WardCount,
		CreatedAt: _utils.FormatTimeToString(dto.CreatedAt),
		UpdatedAt: _utils.FormatTimeToString(dto.UpdatedAt),
	}
}

func (m *LocationMapper) ToProtoWard(dto *dto.WardDTO) *tqdpb.Ward {
	if dto == nil {
		return nil
	}
	return &tqdpb.Ward{
		Id:           dto.ID,
		Code:         dto.Code,
		FullName:     dto.FullName,
		ShortName:    dto.ShortName,
		Latitude:     dto.Lat,
		Longitude:    dto.Lng,
		ProvinceId:   dto.ProvinceID,
		ProvinceCode: dto.ProvinceCode,
		CreatedAt:    _utils.FormatTimeToString(dto.CreatedAt),
		UpdatedAt:    _utils.FormatTimeToString(dto.UpdatedAt),
	}
}

func (m *LocationMapper) ToProtoLocationResponse(dto *dto.LocationResponse) *tqdpb.LocationResponse {
	if dto == nil {
		return nil
	}
	return &tqdpb.LocationResponse{
		Province:         m.ToProtoProvince(dto.Province),
		Ward:             m.ToProtoWard(dto.Ward),
		FullAddress:      dto.FullAddress,
		DistanceKm:       dto.DistanceKm,
		Confidence:       dto.Confidence,
		ProcessingTimeMs: dto.ProcessingTimeMs,
	}
}

func (m *LocationMapper) ToProtoBatchLocationResponse(dto *dto.BatchLocationResponse) *tqdpb.BatchLocationResponse {
	if dto == nil {
		return nil
	}

	locations := make([]*tqdpb.LocationResponse, len(dto.Locations))
	for i, loc := range dto.Locations {
		locations[i] = m.ToProtoLocationResponse(loc)
	}

	return &tqdpb.BatchLocationResponse{
		Locations:             locations,
		TotalProcessingTimeMs: dto.TotalProcessingTimeMs,
	}
}

func (m *LocationMapper) ToProtoSearchResult(dto *dto.SearchResult) *tqdpb.SearchResult {
	if dto == nil {
		return nil
	}
	return &tqdpb.SearchResult{
		Type:           dto.Type,
		Province:       m.ToProtoProvince(dto.Province),
		Ward:           m.ToProtoWard(dto.Ward),
		FullAddress:    dto.FullAddress,
		RelevanceScore: dto.RelevanceScore,
	}
}

func (m *LocationMapper) ToProtoSearchResponse(dto *dto.SearchLocationsResponse) *tqdpb.SearchLocationsResponse {
	if dto == nil {
		return nil
	}

	results := make([]*tqdpb.SearchResult, len(dto.Results))
	for i, res := range dto.Results {
		results[i] = m.ToProtoSearchResult(res)
	}

	return &tqdpb.SearchLocationsResponse{
		Results: results,
		Total:   dto.Total,
	}
}

func (m *LocationMapper) ToProtoListWardsResponse(dto *dto.ListWardsResponse) *tqdpb.ListWardsResponse {
	if dto == nil {
		return nil
	}

	wards := make([]*tqdpb.Ward, len(dto.Wards))
	for i, w := range dto.Wards {
		wards[i] = m.ToProtoWard(w)
	}

	return &tqdpb.ListWardsResponse{
		Wards:    wards,
		Total:    dto.Total,
		Page:     dto.Page,
		PageSize: dto.PageSize,
	}
}

func (m *LocationMapper) ToProtoListProvincesResponse(response *dto.ListProvincesResponse) *tqdpb.ListProvincesResponse {
	if response == nil {
		return nil
	}

	provinces := make([]*tqdpb.Province, len(response.Provinces))
	for i, province := range response.Provinces {
		provinces[i] = m.ToProtoProvince(province)
	}

	return &tqdpb.ListProvincesResponse{Provinces: provinces}
}

// ==================== PROTO -> DTO ====================

func (m *LocationMapper) FromProtoGetNearestLocationRequest(req *tqdpb.GetNearestLocationRequest) *dto.GetNearestLocationRequest {
	if req == nil {
		return nil
	}
	return &dto.GetNearestLocationRequest{
		Latitude:      req.Latitude,
		Longitude:     req.Longitude,
		MaxDistanceKm: req.MaxDistanceKm,
		Limit:         req.Limit,
	}
}

func (m *LocationMapper) FromProtoBatchGetNearestLocationsRequest(req *tqdpb.BatchGetNearestLocationsRequest) *dto.BatchGetNearestLocationsRequest {
	if req == nil {
		return nil
	}

	coordinates := make([]*dto.Coordinate, len(req.Coordinates))
	for i, coord := range req.Coordinates {
		coordinates[i] = &dto.Coordinate{
			Latitude:  coord.Latitude,
			Longitude: coord.Longitude,
		}
	}

	return &dto.BatchGetNearestLocationsRequest{
		Coordinates:   coordinates,
		MaxDistanceKm: req.MaxDistanceKm,
		Limit:         req.Limit,
	}
}

func (m *LocationMapper) FromProtoSearchLocationsRequest(req *tqdpb.SearchLocationsRequest) *dto.SearchLocationsRequest {
	if req == nil {
		return nil
	}
	return &dto.SearchLocationsRequest{
		Query:      req.Query,
		Type:       req.Type,
		Limit:      req.Limit,
		ProvinceID: req.ProvinceId,
	}
}

func (m *LocationMapper) FromProtoGetProvinceRequest(req *tqdpb.GetProvinceRequest) *dto.GetProvinceRequest {
	if req == nil {
		return nil
	}
	return &dto.GetProvinceRequest{
		ID:   &req.Id,
		Code: req.Code,
	}
}

func (m *LocationMapper) FromProtoGetWardRequest(req *tqdpb.GetWardRequest) *dto.GetWardRequest {
	if req == nil {
		return nil
	}
	return &dto.GetWardRequest{
		ID:   &req.Id,
		Code: req.Code,
	}
}

func (m *LocationMapper) FromProtoListWardsByProvinceRequest(req *tqdpb.ListWardsByProvinceRequest) *dto.ListWardsByProvinceRequest {
	if req == nil {
		return nil
	}
	return &dto.ListWardsByProvinceRequest{
		ProvinceID: req.ProvinceId,
		Page:       req.Page,
		PageSize:   req.PageSize,
	}
}
