package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/dto"
	"tqd/internal/geometry"
	"tqd/internal/interface/repo"
)

type RegionUsecase interface {
	Create(ctx context.Context, req *dto.CreateRegionRequest) (*qh_domain.QHRegion, error)
	CreateBatch(ctx context.Context, req *dto.CreateRegionBatchRequest) (*dto.CreateRegionBatchResponse, error)
	Update(ctx context.Context, id uint64, req *dto.UpdateRegionRequest) (*qh_domain.QHRegion, error)
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*qh_domain.QHRegion, error)
	List(ctx context.Context, req *dto.ListRegionRequest) (*dto.ListRegionResponse, error)
	GetStatsByLayer(ctx context.Context, layerID uint64) (*repo.RegionStats, error)
	GetByBBox(ctx context.Context, req *dto.BBoxQueryRequest) ([]dto.RegionResponse, error)
	GetIntersections(ctx context.Context, req *dto.IntersectionQueryRequest) ([]dto.RegionResponse, error)
	ImportGeoJSON(ctx context.Context, req *dto.ImportGeoJSONRequest) (*dto.ImportGeoJSONResponse, error)
	CreateBuffer(ctx context.Context, sourceID uint64, distanceMeters float64, name, displayName string, regionType uint32, status uint32, layerID uint64) (*qh_domain.QHRegion, error)

	// Client
	ListClient(ctx context.Context, filter *repo.ClientRegionFilter) ([]qh_domain.QHRegion, int64, error)
	GetClientByID(ctx context.Context, id uint64) (*qh_domain.QHRegion, error)
	FindByPoint(ctx context.Context, lat, lng float64) (*qh_domain.QHRegion, error)
	FindByBBox(ctx context.Context, minLng, minLat, maxLng, maxLat float64, filter *repo.RegionFilter) ([]qh_domain.QHRegion, int64, error)
}

type regionUsecaseImpl struct {
	regionRepo repo.RegionRepository
	layerRepo  repo.LayerRepository
	geoParser  *dto.GeometryParser
}

func NewRegionUsecase(regionRepo repo.RegionRepository, layerRepo repo.LayerRepository) RegionUsecase {
	return &regionUsecaseImpl{
		regionRepo: regionRepo,
		layerRepo:  layerRepo,
		geoParser:  dto.NewGeometryParser(),
	}
}

func (u *regionUsecaseImpl) Create(ctx context.Context, req *dto.CreateRegionRequest) (*qh_domain.QHRegion, error) {
	if req.Name == "" {
		return nil, errors.New("name is required")
	}
	if req.DisplayName == "" {
		return nil, errors.New("display_name is required")
	}

	// geometryJSON, err := u.geoParser.Parse(req)
	// if err != nil {
	// 	return nil, fmt.Errorf("parse geometry failed: %w", err)
	// }

	// geoInfo, err := geometry.ParseGeometry(geometryJSON)
	// if err != nil {
	// 	return nil, fmt.Errorf("calculate geometry info failed: %w", err)
	// }

	if req.LayerID != nil && *req.LayerID > 0 {
		layer, err := u.layerRepo.GetByID(ctx, *req.LayerID)
		if err != nil || layer == nil {
			return nil, fmt.Errorf("layer with id %d not found", *req.LayerID)
		}
	} else {
		return nil, fmt.Errorf("layerId is required")
	}

	// if req.SourceID > 0 && req.Distance > 0 {
	// 	geometryJSON, err = u.regionRepo.GetBufferGeometry(ctx, req.SourceID, req.Distance)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("create buffer failed: %w", err)
	// 	}
	// 	geoInfo, _ = geometry.ParseGeometry(geometryJSON)
	// }

	region := &qh_domain.QHRegion{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		// Type:        enums.RegionType(req.Type),
		// Status:      enums.RegionStatus(req.Status),
		// Geometry: geometryJSON,
		// Centroid:    geometry.PointToGeoJSON(geoInfo.Centroid.Lat, geoInfo.Centroid.Lng),
		// Area:       geoInfo.Area,
		// Length:     geoInfo.Length,
		// LayerID:       req.LayerID,
		RawProperties: req.Properties,
	}

	if err := u.regionRepo.Create(ctx, region); err != nil {
		return nil, fmt.Errorf("create region failed: %w", err)
	}

	return region, nil
}

func (u *regionUsecaseImpl) Update(ctx context.Context, id uint64, req *dto.UpdateRegionRequest) (*qh_domain.QHRegion, error) {
	existing, err := u.regionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, fmt.Errorf("region with id %d not found", id)
	}

	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.DisplayName != nil {
		existing.DisplayName = *req.DisplayName
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	// if req.Type != nil {
	// 	existing.Type = enums.RegionType(*req.Type)
	// }
	// if req.Status != nil {
	// 	existing.Status = enums.RegionStatus(*req.Status)
	// }
	if req.LayerID != 0 {
		existing.LayerID = req.LayerID
	}
	if req.Properties != nil {
		existing.RawProperties = req.Properties
	}
	if req.Geometry != nil {
		geoInfo, err := geometry.ParseGeometry(*req.Geometry)
		if err != nil {
			return nil, fmt.Errorf("parse geometry failed: %w", err)
		}
		// existing.Geometry = req.Geometry
		// existing.Centroid = geometry.PointToGeoJSON(geoInfo.Centroid.Lat, geoInfo.Centroid.Lng)
		existing.AreaSqm = geoInfo.Area
		// existing.Length = geoInfo.Length
	}

	if err := u.regionRepo.Update(ctx, existing); err != nil {
		return nil, fmt.Errorf("update region failed: %w", err)
	}

	return existing, nil
}

func (u *regionUsecaseImpl) Delete(ctx context.Context, id uint64) error {
	existing, err := u.regionRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("region not found")
	}
	return u.regionRepo.Delete(ctx, id)
}

func (u *regionUsecaseImpl) GetByID(ctx context.Context, id uint64) (*qh_domain.QHRegion, error) {
	region, err := u.regionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if region == nil {
		return nil, fmt.Errorf("region not found")
	}
	return region, nil
}

func (u *regionUsecaseImpl) List(ctx context.Context, req *dto.ListRegionRequest) (*dto.ListRegionResponse, error) {
	filter := &repo.RegionFilter{
		// Type:    req.Type,
		Status:  req.Status,
		LayerID: req.LayerID,
		Search:  req.Search,
		// Page:    req.Page,
		// Limit:   req.Limit,
	}

	regions, total, err := u.regionRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	resp := &dto.ListRegionResponse{
		Data:  make([]dto.RegionResponse, 0, len(regions)),
		Total: total,
		Page:  req.Page,
		Limit: req.Limit,
	}

	for _, r := range regions {
		resp.Data = append(resp.Data, *u.toResponse(&r))
	}

	return resp, nil
}

func (u *regionUsecaseImpl) GetByBBox(ctx context.Context, req *dto.BBoxQueryRequest) ([]dto.RegionResponse, error) {
	regions, err := u.regionRepo.GetByBBox(ctx, req.MinX, req.MinY, req.MaxX, req.MaxY, req.Type, req.LayerID)
	if err != nil {
		return nil, err
	}

	resp := make([]dto.RegionResponse, 0, len(regions))
	for _, r := range regions {
		resp = append(resp, *u.toResponse(&r))
	}
	return resp, nil
}

func (u *regionUsecaseImpl) GetIntersections(ctx context.Context, req *dto.IntersectionQueryRequest) ([]dto.RegionResponse, error) {
	buffer := req.Buffer
	if buffer <= 0 {
		buffer = 10
	}

	regions, err := u.regionRepo.GetIntersections(ctx, req.ID, buffer)
	if err != nil {
		return nil, err
	}

	resp := make([]dto.RegionResponse, 0, len(regions))
	for _, r := range regions {
		resp = append(resp, *u.toResponse(&r))
	}
	return resp, nil
}

func (u *regionUsecaseImpl) ImportGeoJSON(ctx context.Context, req *dto.ImportGeoJSONRequest) (*dto.ImportGeoJSONResponse, error) {
	var fc struct {
		Type     string `json:"type"`
		Features []struct {
			Type       string                 `json:"type"`
			Geometry   map[string]interface{} `json:"geometry"`
			Properties map[string]interface{} `json:"properties"`
		} `json:"features"`
	}

	if err := json.Unmarshal([]byte(req.GeoJSON), &fc); err != nil {
		return nil, fmt.Errorf("invalid GeoJSON: %w", err)
	}

	if len(fc.Features) == 0 {
		return nil, errors.New("GeoJSON has no features")
	}

	// feature := fc.Features[0]
	// geomBytes, err := json.Marshal(feature.Geometry)
	// if err != nil {
	// 	return nil, fmt.Errorf("marshal geometry failed: %w", err)
	// }

	// geoInfo, err := geometry.ParseGeometry(string(geomBytes))
	// if err != nil {
	// 	return nil, fmt.Errorf("parse geometry failed: %w", err)
	// }

	if req.LayerID != nil && *req.LayerID > 0 {
		layer, err := u.layerRepo.GetByID(ctx, *req.LayerID)
		if err != nil || layer == nil {
			return nil, fmt.Errorf("layer not found")
		}
	}

	region := &qh_domain.QHRegion{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		// Type:        enums.RegionType(req.Type),
		// Status:      enums.RegionStatus(req.Status),
		// Geometry: string(geomBytes),
		// Centroid:    geometry.PointToGeoJSON(geoInfo.Centroid.Lat, geoInfo.Centroid.Lng),
		// Area:       geoInfo.Area,
		// Length:     geoInfo.Length,
		// LayerID:    req.LayerID,
		// Properties: feature.Properties,
	}

	if err := u.regionRepo.Create(ctx, region); err != nil {
		return nil, fmt.Errorf("create region failed: %w", err)
	}

	return &dto.ImportGeoJSONResponse{
		Region:       *u.toResponse(region),
		FeatureCount: len(fc.Features),
		Message:      fmt.Sprintf("Imported %d features", len(fc.Features)),
	}, nil
}

func (u *regionUsecaseImpl) CreateBuffer(ctx context.Context, sourceID uint64, distanceMeters float64, name, displayName string, regionType uint32, status uint32, layerID uint64) (*qh_domain.QHRegion, error) {
	// bufferGeoJSON, err := u.regionRepo.GetBufferGeometry(ctx, sourceID, distanceMeters)
	// if err != nil {
	// 	return nil, err
	// }

	// geoInfo, err := geometry.ParseGeometry(bufferGeoJSON)
	// if err != nil {
	// 	return nil, err
	// }

	region := &qh_domain.QHRegion{
		Name:        name,
		DisplayName: displayName,
		// Type:        enums.RegionType(regionType),
		// Status:      enums.RegionStatus(status),
		// Geometry: bufferGeoJSON,
		// Centroid:    geometry.PointToGeoJSON(geoInfo.Centroid.Lat, geoInfo.Centroid.Lng),
		// Area:    geoInfo.Area,
		// Length:  geoInfo.Length,
		LayerID: layerID,
		// Properties: map[string]interface{}{
		// 	"source_id": sourceID,
		// 	"distance":  distanceMeters,
		// },
	}

	if err := u.regionRepo.Create(ctx, region); err != nil {
		return nil, err
	}
	return region, nil
}

func (u *regionUsecaseImpl) CreateBatch(ctx context.Context, req *dto.CreateRegionBatchRequest) (*dto.CreateRegionBatchResponse, error) {
	resp := &dto.CreateRegionBatchResponse{
		Items: make([]dto.CreateRegionBatchResultItem, 0, len(req.Items)),
		Total: len(req.Items),
	}

	regions := make([]*qh_domain.QHRegion, len(req.Items))

	for i, item := range req.Items {
		result := dto.CreateRegionBatchResultItem{Index: i}

		// Validate layerId tồn tại
		layer, err := u.layerRepo.GetByID(ctx, item.LayerID)
		if err != nil || layer == nil {
			result.Success = false
			result.Error = fmt.Sprintf("layer with id %d not found", item.LayerID)
			resp.Items = append(resp.Items, result)
			resp.FailureCount++
			continue
		}

		// Parse geometry
		if item.Geometry == nil {
			result.Success = false
			result.Error = "geometry is required"
			resp.Items = append(resp.Items, result)
			resp.FailureCount++
			continue
		}
		// geoInfo, err := geometry.ParseGeometry(item.Geometry)
		if err != nil {
			result.Success = false
			result.Error = fmt.Sprintf("invalid geometry: %v", err)
			resp.Items = append(resp.Items, result)
			resp.FailureCount++
			continue
		}

		// Gán giá trị mặc định cho name nếu trống
		name := item.Name
		if name == "" {
			name = fmt.Sprintf("region_%d_%d_%d", item.LayerID, time.Now().UnixNano(), i)
		}
		displayName := item.DisplayName
		if displayName == "" {
			displayName = name
		}
		regionType := item.Type
		if regionType == 0 {
			regionType = 1 // mặc định Road
		}
		regionStatus := item.Status
		if regionStatus == 0 {
			regionStatus = 10 // mặc định Active
		}
		layerID := item.LayerID

		regions[i] = &qh_domain.QHRegion{
			Name:        name,
			DisplayName: displayName,
			// Type:        enums.RegionType(regionType),
			// Status:      enums.RegionStatus(regionStatus),
			Geometry: qh_domain.RawGeometry{
				Raw: item.Geometry,
			},
			// Centroid:   geometry.PointToGeoJSON(geoInfo.Centroid.Lat, geoInfo.Centroid.Lng),
			// Area:       geoInfo.Area,
			// Length:     geoInfo.Length,
			LayerID:       layerID,
			RawProperties: item.Properties,
		}

		// regionResp := u.toResponse(region)
		result.Success = true
		// result.Region = regionResp
		// resp.Items = append(resp.Items, result)
		resp.SuccessCount++
	}

	if err := u.regionRepo.CreateBatch(ctx, regions); err != nil {
		// result.Success = false
		// result.Error = fmt.Sprintf("create failed: %v", err)
		// resp.Items = append(resp.Items, result)
		resp.FailureCount++

	}

	return resp, nil
}

func (u *regionUsecaseImpl) toResponse(region *qh_domain.QHRegion) *dto.RegionResponse {
	return &dto.RegionResponse{
		ID:          region.ID,
		Name:        region.Name,
		DisplayName: region.DisplayName,
		Description: region.Description,
		// Type:        uint32(region.Type),
		// TypeName:    region.Type.String(),
		// Status:      uint32(region.Status),
		// StatusName:  region.Status.String(),
		Geometry: string(region.Geometry.Raw),
		// Centroid:   region.Centroid,
		// BBox:       region.BBox,
		// Area:       region.Area,
		// Length:     region.Length,
		LayerID: region.LayerID,
		// Properties: region.Properties,
		CreatedAt: region.CreatedAt.Format(time.RFC3339),
		UpdatedAt: region.UpdatedAt.Format(time.RFC3339),
	}
}

func (u *regionUsecaseImpl) GetStatsByLayer(ctx context.Context, layerID uint64) (*repo.RegionStats, error) {
	if layerID == 0 {
		return nil, errors.New("layer_id is required")
	}
	return u.regionRepo.GetStatsByLayer(ctx, layerID)
}

func (u *regionUsecaseImpl) ListClient(ctx context.Context, filter *repo.ClientRegionFilter) ([]qh_domain.QHRegion, int64, error) {
	return u.regionRepo.ListClient(ctx, filter)
}

func (u *regionUsecaseImpl) GetClientByID(ctx context.Context, id uint64) (*qh_domain.QHRegion, error) {
	region, err := u.regionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get region failed: %w", err)
	}
	if region == nil {
		return nil, fmt.Errorf("region not found")
	}
	if region.Status != 10 || region.DeletedAt != nil {
		return nil, fmt.Errorf("region not available")
	}
	return region, nil
}

func (u *regionUsecaseImpl) FindByPoint(ctx context.Context, lat, lng float64) (*qh_domain.QHRegion, error) {
	if lat < -90 || lat > 90 {
		return nil, errors.New("invalid latitude")
	}
	if lng < -180 || lng > 180 {
		return nil, errors.New("invalid longitude")
	}
	return u.regionRepo.FindByPoint(ctx, lat, lng)
}

func (u *regionUsecaseImpl) FindByBBox(ctx context.Context, minLng, minLat, maxLng, maxLat float64, filter *repo.RegionFilter) ([]qh_domain.QHRegion, int64, error) {
	return u.regionRepo.FindByBBox(ctx, minLng, minLat, maxLng, maxLat, filter)
}
