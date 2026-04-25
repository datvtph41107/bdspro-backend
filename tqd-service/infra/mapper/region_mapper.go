// infra/mapper/region_mapper.go
package mapper

import (
	"time"

	tqdpb "pb/types/tqd"
	"tqd/internal/domain"
	"tqd/internal/dto"

	"google.golang.org/protobuf/types/known/structpb"
)

type RegionMapper struct{}

func NewRegionMapper() *RegionMapper {
	return &RegionMapper{}
}

func (m *RegionMapper) ToProtoFromDomain(region *domain.Region) *tqdpb.RegionResponse {
	if region == nil {
		return nil
	}

	resp := &tqdpb.RegionResponse{
		Id:          uint32(region.ID),
		Name:        region.Name,
		DisplayName: region.DisplayName,
		Description: region.Description,
		Type:        uint32(region.Type),
		Status:      uint32(region.Status),
		Geometry:    region.Geometry,
		Centroid:    region.Centroid,
		Area:        region.Area,
		Length:      region.Length,
		CreatedAt:   region.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   region.UpdatedAt.Format(time.RFC3339),
	}

	if region.LayerID != nil {
		layerID := uint32(*region.LayerID)
		resp.LayerId = &layerID
	}

	if region.Properties != nil {
		if props, err := structpb.NewStruct(region.Properties); err == nil {
			resp.Properties = props
		}
	}

	return resp
}

// ToProtoFromDTO - Chuyển từ dto.RegionResponse sang proto
func (m *RegionMapper) ToProtoFromDTO(region *dto.RegionResponse) *tqdpb.RegionResponse {
	if region == nil {
		return nil
	}

	resp := &tqdpb.RegionResponse{
		Id:          uint32(region.ID),
		Name:        region.Name,
		DisplayName: region.DisplayName,
		Description: region.Description,
		Type:        region.Type,
		Status:      region.Status,
		Geometry:    region.Geometry,
		Centroid:    region.Centroid,
		Area:        region.Area,
		Length:      region.Length,
		CreatedAt:   region.CreatedAt,
		UpdatedAt:   region.UpdatedAt,
	}

	if region.LayerID != nil {
		layerID := uint32(*region.LayerID)
		resp.LayerId = &layerID
	}

	if region.Properties != nil {
		if props, err := structpb.NewStruct(region.Properties); err == nil {
			resp.Properties = props
		}
	}

	return resp
}

// ToProtoListFromDTO - Chuyển danh sách từ dto sang proto
func (m *RegionMapper) ToProtoListFromDTO(regions []dto.RegionResponse) []*tqdpb.RegionResponse {
	result := make([]*tqdpb.RegionResponse, 0, len(regions))
	for i := range regions {
		result = append(result, m.ToProtoFromDTO(&regions[i]))
	}
	return result
}

// ToProtoListFromDomain - Chuyển danh sách từ domain sang proto
func (m *RegionMapper) ToProtoListFromDomain(regions []domain.Region) []*tqdpb.RegionResponse {
	result := make([]*tqdpb.RegionResponse, 0, len(regions))
	for i := range regions {
		result = append(result, m.ToProtoFromDomain(&regions[i]))
	}
	return result
}

// ToCreateRegionDTO - Chuyển từ proto CreateRegionRequest sang dto.CreateRegionRequest
func (m *RegionMapper) ToCreateRegionDTO(req *tqdpb.CreateRegionRequest) *dto.CreateRegionRequest {
	if req == nil {
		return nil
	}

	props := make(map[string]interface{})
	if req.Properties != nil {
		props = req.Properties.AsMap()
	}

	createReq := &dto.CreateRegionRequest{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		Type:        req.Type,
		Status:      req.Status,
		GeoJSON:     req.Geometry,
		Properties:  props,
	}

	if req.LayerId != nil {
		layerID := uint64(*req.LayerId)
		createReq.LayerID = &layerID
	}

	return createReq
}

// ToUpdateRegionDTO - Chuyển từ proto UpdateRegionRequest sang dto.UpdateRegionRequest
func (m *RegionMapper) ToUpdateRegionDTO(req *tqdpb.UpdateRegionRequest) *dto.UpdateRegionRequest {
	if req == nil {
		return nil
	}

	updateReq := &dto.UpdateRegionRequest{}

	if req.Name != nil {
		updateReq.Name = req.Name
	}
	if req.DisplayName != nil {
		updateReq.DisplayName = req.DisplayName
	}
	if req.Description != nil {
		updateReq.Description = req.Description
	}
	if req.Type != nil {
		t := *req.Type
		updateReq.Type = &t
	}
	if req.Status != nil {
		s := *req.Status
		updateReq.Status = &s
	}
	if req.Geometry != nil {
		updateReq.Geometry = req.Geometry
	}
	if req.LayerId != nil {
		layerID := uint64(*req.LayerId)
		updateReq.LayerID = &layerID
	}
	if req.Properties != nil {
		updateReq.Properties = req.Properties.AsMap()
	}

	return updateReq
}

// ToListRegionDTO - Chuyển từ proto ListRegionsRequest sang dto.ListRegionRequest
func (m *RegionMapper) ToListRegionDTO(req *tqdpb.ListRegionsRequest) *dto.ListRegionRequest {
	if req == nil {
		return nil
	}

	page := int32(1)
	limit := int32(20)
	if req.Page != nil && *req.Page > 0 {
		page = *req.Page
	}
	if req.Limit != nil && *req.Limit > 0 {
		limit = *req.Limit
	}

	listReq := &dto.ListRegionRequest{
		Page:  int(page),
		Limit: int(limit),
	}

	if req.Type != nil {
		t := *req.Type
		listReq.Type = &t
	}
	if req.Status != nil {
		s := *req.Status
		listReq.Status = &s
	}
	if req.LayerId != nil {
		id := uint64(*req.LayerId)
		listReq.LayerID = &id
	}
	if req.Search != nil {
		listReq.Search = req.Search
	}

	return listReq
}

// ToImportGeoJSONDTO - Chuyển từ proto ImportGeoJSONRequest sang dto.ImportGeoJSONRequest
func (m *RegionMapper) ToImportGeoJSONDTO(req *tqdpb.ImportGeoJSONRequest) *dto.ImportGeoJSONRequest {
	if req == nil {
		return nil
	}

	importReq := &dto.ImportGeoJSONRequest{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		Type:        req.Type,
		Status:      req.Status,
		GeoJSON:     req.GeoJSON,
	}

	if req.LayerId != nil {
		layerID := uint64(*req.LayerId)
		importReq.LayerID = &layerID
	}

	return importReq
}

// ToIntersectionQueryDTO - Chuyển từ proto GetRegionIntersectionsRequest sang dto.IntersectionQueryRequest
func (m *RegionMapper) ToIntersectionQueryDTO(req *tqdpb.GetRegionIntersectionsRequest) *dto.IntersectionQueryRequest {
	if req == nil {
		return nil
	}

	buffer := float64(10)
	if req.Buffer != nil && *req.Buffer > 0 {
		buffer = *req.Buffer
	}

	return &dto.IntersectionQueryRequest{
		ID:     uint64(req.Id),
		Buffer: buffer,
	}
}

// ToPointQueryDTO - Chuyển từ proto GetRegionByLocationRequest sang dto.PointQueryRequest
func (m *RegionMapper) ToPointQueryDTO(req *tqdpb.GetRegionByLocationRequest) *dto.PointQueryRequest {
	if req == nil {
		return nil
	}

	pointReq := &dto.PointQueryRequest{
		Lat: req.Lat,
		Lng: req.Lng,
	}

	if req.Type != nil && *req.Type > 0 {
		t := *req.Type
		pointReq.Type = &t
	}
	if req.LayerId != nil && *req.LayerId > 0 {
		id := uint64(*req.LayerId)
		pointReq.LayerID = &id
	}

	return pointReq
}
