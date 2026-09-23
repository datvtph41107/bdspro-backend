// internal/mapper/location_mapper.go
package mapper

import (
	"log"
	"strings"
	"time"

	_dto "common/domain/dto"
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
	resp := &tqdpb.Province{
		Id:                    dto.ID,
		Code:                  dto.Code,
		FullName:              dto.FullName,
		ShortName:             dto.ShortName,
		Latitude:              dto.Lat,
		Longitude:             dto.Lng,
		DivisionType:          dto.DivisionType,
		WardCount:             dto.WardCount,
		OfficialName:          dto.OfficialName,
		OfficialOperationDate: _utils.FormatTimeToString(dto.OfficialOperationDate),
		NaturalAreaKm2:        dto.NaturalAreaKm2,
		Population:            dto.Population,
		PopulationDensity:     dto.PopulationDensity,
		Description:           dto.Description,
		UnitType:              dto.UnitType,
		UnitTypeName:          dto.UnitTypeName,
		DataSource:            dto.DataSource,
		DataSourceName:        dto.DataSourceName,
		OfficialRole:          dto.OfficialRole,
		OfficialAddress:       dto.OfficialAddress,
		CompositeId:           dto.CompositeID,
		CreatedAt:             _utils.FormatTimeToString(dto.CreatedAt),
		UpdatedAt:             _utils.FormatTimeToString(dto.UpdatedAt),
		Geometry:              dto.Geometry,
	}
	if dto.SeoID != nil && *dto.SeoID > 0 {
		seoID := *dto.SeoID
		resp.SeoId = &seoID
	}
	if dto.PlanDocID != nil && *dto.PlanDocID > 0 {
		planDocID := *dto.PlanDocID
		resp.PlanDocId = &planDocID
	}
	return resp
}

func (m *LocationMapper) ToProtoWardListItem(dto *dto.WardDTO) *tqdpb.Ward {
	if dto == nil {
		return nil
	}
	resp := &tqdpb.Ward{
		Id: dto.ID,
		// Code: dto.Code,
		// FullName:     dto.FullName,
		ShortName: dto.ShortName,
		Latitude:  dto.Lat,
		Longitude: dto.Lng,
		// DivisionType: dto.DivisionType,
		ProvinceId: dto.ProvinceID,
		// ProvinceCode:          dto.ProvinceCode,
		// OfficialName:          dto.OfficialName,
		// OfficialOperationDate: _utils.FormatTimeToString(dto.OfficialOperationDate),
		// NaturalAreaKm2:    dto.NaturalAreaKm2,
		// Population:        dto.Population,
		// PopulationDensity: dto.PopulationDensity,
		// Summary:           dto.Summary,
		// Description:       dto.Description,
		CompositeId:  dto.CompositeID,
		UnitType:     dto.UnitType,
		UnitTypeName: dto.UnitTypeName,
		// DataSource:     dto.DataSource,
		// DataSourceName: dto.DataSourceName,
		// OfficialRole:          dto.OfficialRole,
		// OfficialAddress:       dto.OfficialAddress,
		// CreatedAt:             _utils.FormatTimeToString(dto.CreatedAt),
		// UpdatedAt:             _utils.FormatTimeToString(dto.UpdatedAt),
		Geometry: dto.Geometry,
	}
	if dto.SeoID != nil && *dto.SeoID > 0 {
		seoID := *dto.SeoID
		resp.SeoId = &seoID
	}
	return resp
}

func (m *LocationMapper) ToProtoWard(dto *dto.WardDTO) *tqdpb.Ward {
	if dto == nil {
		return nil
	}
	resp := &tqdpb.Ward{
		Id:                    dto.ID,
		Code:                  dto.Code,
		FullName:              dto.FullName,
		ShortName:             dto.ShortName,
		Latitude:              dto.Lat,
		Longitude:             dto.Lng,
		DivisionType:          dto.DivisionType,
		ProvinceId:            dto.ProvinceID,
		ProvinceCode:          dto.ProvinceCode,
		OfficialName:          dto.OfficialName,
		OfficialOperationDate: _utils.FormatTimeToString(dto.OfficialOperationDate),
		NaturalAreaKm2:        dto.NaturalAreaKm2,
		Population:            dto.Population,
		PopulationDensity:     dto.PopulationDensity,
		Summary:               dto.Summary,
		Description:           dto.Description,
		CompositeId:           dto.CompositeID,
		UnitType:              dto.UnitType,
		UnitTypeName:          dto.UnitTypeName,
		DataSource:            dto.DataSource,
		DataSourceName:        dto.DataSourceName,
		OfficialRole:          dto.OfficialRole,
		OfficialAddress:       dto.OfficialAddress,
		CreatedAt:             _utils.FormatTimeToString(dto.CreatedAt),
		UpdatedAt:             _utils.FormatTimeToString(dto.UpdatedAt),
		Geometry:              dto.Geometry,
	}
	if dto.SeoID != nil && *dto.SeoID > 0 {
		seoID := *dto.SeoID
		resp.SeoId = &seoID
	}
	return resp
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
		CompositeId:    dto.CompositeID,
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
		wards[i] = m.ToProtoWardListItem(w)
	}

	return &tqdpb.ListWardsResponse{
		Wards:    wards,
		Total:    dto.Total,
		Page:     dto.Page,
		PageSize: dto.PageSize,
	}
}

func (m *LocationMapper) ToProtoClientListWardsResponse(resp *dto.ClientListWardsResponse) *tqdpb.ClientListWardsResponse {
	if resp == nil {
		return nil
	}

	data := make([]*tqdpb.Ward, len(resp.Data))
	for i, w := range resp.Data {
		data[i] = m.ToProtoWardListItem(w)
	}

	out := &tqdpb.ClientListWardsResponse{Data: data}
	if resp.Pagination != nil {
		log.Println("pagination", resp.Pagination.PageSize)
		out.Pagination = &tqdpb.LocationPagination{
			Page:  resp.Pagination.Page,
			Size:  resp.Pagination.PageSize,
			Total: resp.Pagination.Total,
		}
	}
	return out
}

func (m *LocationMapper) ToProtoAdminListProvincesResponse(dto *dto.AdminListProvincesResponse) *tqdpb.AdminListProvincesResponse {
	if dto == nil {
		return nil
	}

	data := make([]*tqdpb.Province, len(dto.Data))
	for i, p := range dto.Data {
		data[i] = m.ToProtoProvince(p)
	}

	return &tqdpb.AdminListProvincesResponse{
		Data:  data,
		Total: dto.Total,
	}
}

func (m *LocationMapper) ToProtoAdminListWardsResponse(dto *dto.AdminListWardsResponse) *tqdpb.AdminListWardsResponse {
	if dto == nil {
		return nil
	}

	data := make([]*tqdpb.Ward, len(dto.Data))
	for i, w := range dto.Data {
		data[i] = m.ToProtoWard(w)
	}

	return &tqdpb.AdminListWardsResponse{
		Data:  data,
		Total: dto.Total,
	}
}

func (m *LocationMapper) ToProtoListAdministrativeUnitsResponse(dto *dto.ListAdministrativeUnitsResponse) *tqdpb.ListAdministrativeUnitsResponse {
	if dto == nil {
		return nil
	}
	data := make([]*tqdpb.AdministrativeUnitListItem, 0, len(dto.Data))
	for _, item := range dto.Data {
		if item == nil {
			continue
		}
		data = append(data, &tqdpb.AdministrativeUnitListItem{
			Province: m.ToProtoProvince(item.Province),
			Ward:     m.ToProtoWard(item.Ward),
		})
	}
	return &tqdpb.ListAdministrativeUnitsResponse{
		Data:  data,
		Total: dto.Total,
	}
}

func (m *LocationMapper) FromProtoAdminListProvincesRequest(req *tqdpb.AdminListProvincesRequest) *dto.AdminListProvincesRequest {
	if req == nil {
		return &dto.AdminListProvincesRequest{}
	}
	return &dto.AdminListProvincesRequest{
		Pagable: _dto.Pagable{
			Page: req.Page,
			Size: req.Size,
			Sort: req.Sort,
		},
		Search: req.Search,
	}
}

func (m *LocationMapper) FromProtoAdminListWardsRequest(req *tqdpb.AdminListWardsRequest) *dto.AdminListWardsRequest {
	if req == nil {
		return &dto.AdminListWardsRequest{}
	}
	return &dto.AdminListWardsRequest{
		Pagable: _dto.Pagable{
			Page: req.Page,
			Size: req.Size,
			Sort: req.Sort,
		},
		Search:     req.Search,
		ProvinceID: req.ProvinceId,
	}
}

func (m *LocationMapper) FromProtoAdminListAdministrativeUnitsRequest(req *tqdpb.AdminListAdministrativeUnitsRequest) *dto.ListAdministrativeUnitsRequest {
	if req == nil {
		return &dto.ListAdministrativeUnitsRequest{}
	}
	return &dto.ListAdministrativeUnitsRequest{
		Pagable: _dto.Pagable{
			Page: req.Page,
			Size: req.Size,
			Sort: req.Sort,
		},
		Search:            req.Search,
		PlanningProjectID: req.PlanningProjectId,
	}
}

func (m *LocationMapper) FromProtoClientListAdministrativeUnitsRequest(req *tqdpb.ClientListAdministrativeUnitsRequest) *dto.ListAdministrativeUnitsRequest {
	if req == nil {
		return &dto.ListAdministrativeUnitsRequest{}
	}
	return &dto.ListAdministrativeUnitsRequest{
		Pagable: _dto.Pagable{
			Page: req.Page,
			Size: req.Size,
			Sort: req.Sort,
		},
		Search:            req.Search,
		PlanningProjectID: req.PlanningProjectId,
		CompositeIDs:      _utils.ParseUint64Array(req.GetCompositeIds()),
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
	out := &dto.GetProvinceRequest{
		Code: req.Code,
	}
	if req.Id != 0 {
		id := req.Id
		out.ID = &id
	}
	return out
}

func (m *LocationMapper) FromProtoGetWardRequest(req *tqdpb.GetWardRequest) *dto.GetWardRequest {
	if req == nil {
		return nil
	}
	out := &dto.GetWardRequest{
		Code: req.Code,
	}
	if req.Id != 0 {
		id := req.Id
		out.ID = &id
	}
	return out
}

func (m *LocationMapper) FromProtoListWardsRequest(req *tqdpb.ListWardsRequest) *dto.ListWardsRequest {
	if req == nil {
		return &dto.ListWardsRequest{}
	}
	return &dto.ListWardsRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
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

func (m *LocationMapper) FromProtoProvince(req *tqdpb.Province) *dto.ProvinceDTO {
	if req == nil {
		return nil
	}
	resp := &dto.ProvinceDTO{
		ID:                    req.Id,
		Code:                  req.Code,
		FullName:              req.FullName,
		ShortName:             req.ShortName,
		Lat:                   req.Latitude,
		Lng:                   req.Longitude,
		DivisionType:          req.DivisionType,
		WardCount:             req.WardCount,
		OfficialName:          req.OfficialName,
		OfficialOperationDate: parseOfficialOperationDate(req.OfficialOperationDate),
		NaturalAreaKm2:        req.NaturalAreaKm2,
		Population:            req.Population,
		PopulationDensity:     req.PopulationDensity,
		Description:           req.Description,
		UnitType:              req.UnitType,
		DataSource:            req.DataSource,
		OfficialRole:          req.OfficialRole,
		OfficialAddress:       req.OfficialAddress,
	}
	if req.SeoId != nil {
		seoID := req.GetSeoId()
		resp.SeoID = &seoID
	}
	if req.PlanDocId != nil {
		planDocID := req.GetPlanDocId()
		resp.PlanDocID = &planDocID
	}
	return resp
}

func (m *LocationMapper) FromProtoWard(req *tqdpb.Ward) *dto.WardDTO {
	if req == nil {
		return nil
	}
	return &dto.WardDTO{
		ID:                    req.Id,
		Code:                  req.Code,
		FullName:              req.FullName,
		ShortName:             req.ShortName,
		Lat:                   req.Latitude,
		Lng:                   req.Longitude,
		DivisionType:          req.DivisionType,
		ProvinceID:            req.ProvinceId,
		ProvinceCode:          req.ProvinceCode,
		OfficialName:          req.OfficialName,
		OfficialOperationDate: parseOfficialOperationDate(req.OfficialOperationDate),
		NaturalAreaKm2:        req.NaturalAreaKm2,
		Population:            req.Population,
		PopulationDensity:     req.PopulationDensity,
		Summary:               req.Summary,
		Description:           req.Description,
		UnitType:              req.UnitType,
		DataSource:            req.DataSource,
		OfficialRole:          req.OfficialRole,
		OfficialAddress:       req.OfficialAddress,
	}
}

func parseOfficialOperationDate(value string) *time.Time {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	if t := _utils.ParseStringToTime(value); t != nil {
		return t
	}
	for _, layout := range []string{"02/01/2006", "2/1/2006", "02-01-2006"} {
		t, err := time.Parse(layout, value)
		if err == nil {
			return &t
		}
	}
	return nil
}
