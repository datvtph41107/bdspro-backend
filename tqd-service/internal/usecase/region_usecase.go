package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"tqd/internal/domain"
	"tqd/internal/dto"
	"tqd/internal/enums"
	"tqd/internal/geometry"
	"tqd/internal/interface/repo"
)

type RegionUsecase interface {
	Create(ctx context.Context, req *dto.CreateRegionRequest) (*domain.Region, error)
	Update(ctx context.Context, id uint64, req *dto.UpdateRegionRequest) (*domain.Region, error)
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*dto.RegionResponse, error)
	List(ctx context.Context, req *dto.ListRegionRequest) (*dto.ListRegionResponse, error)
	GetByLocation(ctx context.Context, lat, lng float64, regionType *enums.RegionType, layerID *uint64) ([]domain.Region, error)
	GetByPoint(ctx context.Context, req *dto.PointQueryRequest) ([]dto.RegionResponse, error)
	GetByBBox(ctx context.Context, req *dto.BBoxQueryRequest) ([]dto.RegionResponse, error)
	GetIntersections(ctx context.Context, req *dto.IntersectionQueryRequest) ([]dto.RegionResponse, error)
	ImportGeoJSON(ctx context.Context, req *dto.ImportGeoJSONRequest) (*dto.ImportGeoJSONResponse, error)
	CreateBuffer(ctx context.Context, sourceID uint64, distanceMeters float64, name, displayName string, regionType uint32, status uint32, layerID *uint64) (*domain.Region, error)
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

func (u *regionUsecaseImpl) Create(ctx context.Context, req *dto.CreateRegionRequest) (*domain.Region, error) {
	if req.Name == "" {
		return nil, errors.New("name is required")
	}
	if req.DisplayName == "" {
		return nil, errors.New("display_name is required")
	}

	geometryJSON, err := u.geoParser.Parse(req)
	if err != nil {
		return nil, fmt.Errorf("parse geometry failed: %w", err)
	}

	geoInfo, err := geometry.ParseGeometry(geometryJSON)
	if err != nil {
		return nil, fmt.Errorf("calculate geometry info failed: %w", err)
	}

	if req.LayerID != nil && *req.LayerID > 0 {
		layer, err := u.layerRepo.GetByID(ctx, *req.LayerID)
		if err != nil || layer == nil {
			return nil, fmt.Errorf("layer with id %d not found", *req.LayerID)
		}
	}

	if req.SourceID > 0 && req.Distance > 0 {
		geometryJSON, err = u.regionRepo.GetBufferGeometry(ctx, req.SourceID, req.Distance)
		if err != nil {
			return nil, fmt.Errorf("create buffer failed: %w", err)
		}
		geoInfo, _ = geometry.ParseGeometry(geometryJSON)
	}

	region := &domain.Region{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		Type:        enums.RegionType(req.Type),
		Status:      enums.RegionStatus(req.Status),
		Geometry:    geometryJSON,
		Centroid:    geometry.PointToGeoJSON(geoInfo.Centroid.Lat, geoInfo.Centroid.Lng),
		Area:        geoInfo.Area,
		Length:      geoInfo.Length,
		LayerID:     req.LayerID,
		Properties:  req.Properties,
	}

	if err := u.regionRepo.Create(ctx, region); err != nil {
		return nil, fmt.Errorf("create region failed: %w", err)
	}

	return region, nil
}

func (u *regionUsecaseImpl) Update(ctx context.Context, id uint64, req *dto.UpdateRegionRequest) (*domain.Region, error) {
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
	if req.Type != nil {
		existing.Type = enums.RegionType(*req.Type)
	}
	if req.Status != nil {
		existing.Status = enums.RegionStatus(*req.Status)
	}
	if req.LayerID != nil {
		existing.LayerID = req.LayerID
	}
	if req.Properties != nil {
		existing.Properties = req.Properties
	}
	if req.Geometry != nil {
		geoInfo, err := geometry.ParseGeometry(*req.Geometry)
		if err != nil {
			return nil, fmt.Errorf("parse geometry failed: %w", err)
		}
		existing.Geometry = *req.Geometry
		existing.Centroid = geometry.PointToGeoJSON(geoInfo.Centroid.Lat, geoInfo.Centroid.Lng)
		existing.Area = geoInfo.Area
		existing.Length = geoInfo.Length
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

func (u *regionUsecaseImpl) GetByID(ctx context.Context, id uint64) (*dto.RegionResponse, error) {
	region, err := u.regionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if region == nil {
		return nil, fmt.Errorf("region not found")
	}
	return u.toResponse(region), nil
}

func (u *regionUsecaseImpl) List(ctx context.Context, req *dto.ListRegionRequest) (*dto.ListRegionResponse, error) {
	filter := &repo.RegionFilter{
		Type:    req.Type,
		Status:  req.Status,
		LayerID: req.LayerID,
		Search:  req.Search,
		Page:    req.Page,
		Limit:   req.Limit,
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

func (u *regionUsecaseImpl) GetByLocation(ctx context.Context, lat, lng float64, regionType *enums.RegionType, layerID *uint64) ([]domain.Region, error) {
	var typePtr *uint32
	if regionType != nil {
		t := uint32(*regionType)
		typePtr = &t
	}
	return u.regionRepo.GetByLocation(ctx, lat, lng, typePtr, layerID)
}

func (u *regionUsecaseImpl) GetByPoint(ctx context.Context, req *dto.PointQueryRequest) ([]dto.RegionResponse, error) {
	buffer := float64(0)
	if req.Buffer != nil {
		buffer = *req.Buffer
	}

	regions, err := u.regionRepo.GetByPoint(ctx, req.Lat, req.Lng, buffer, req.Type, req.LayerID)
	if err != nil {
		return nil, err
	}

	resp := make([]dto.RegionResponse, 0, len(regions))
	for _, r := range regions {
		resp = append(resp, *u.toResponse(&r))
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

	feature := fc.Features[0]
	geomBytes, err := json.Marshal(feature.Geometry)
	if err != nil {
		return nil, fmt.Errorf("marshal geometry failed: %w", err)
	}

	geoInfo, err := geometry.ParseGeometry(string(geomBytes))
	if err != nil {
		return nil, fmt.Errorf("parse geometry failed: %w", err)
	}

	if req.LayerID != nil && *req.LayerID > 0 {
		layer, err := u.layerRepo.GetByID(ctx, *req.LayerID)
		if err != nil || layer == nil {
			return nil, fmt.Errorf("layer not found")
		}
	}

	region := &domain.Region{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		Type:        enums.RegionType(req.Type),
		Status:      enums.RegionStatus(req.Status),
		Geometry:    string(geomBytes),
		Centroid:    geometry.PointToGeoJSON(geoInfo.Centroid.Lat, geoInfo.Centroid.Lng),
		Area:        geoInfo.Area,
		Length:      geoInfo.Length,
		LayerID:     req.LayerID,
		Properties:  feature.Properties,
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

func (u *regionUsecaseImpl) CreateBuffer(ctx context.Context, sourceID uint64, distanceMeters float64, name, displayName string, regionType uint32, status uint32, layerID *uint64) (*domain.Region, error) {
	bufferGeoJSON, err := u.regionRepo.GetBufferGeometry(ctx, sourceID, distanceMeters)
	if err != nil {
		return nil, err
	}

	geoInfo, err := geometry.ParseGeometry(bufferGeoJSON)
	if err != nil {
		return nil, err
	}

	region := &domain.Region{
		Name:        name,
		DisplayName: displayName,
		Type:        enums.RegionType(regionType),
		Status:      enums.RegionStatus(status),
		Geometry:    bufferGeoJSON,
		Centroid:    geometry.PointToGeoJSON(geoInfo.Centroid.Lat, geoInfo.Centroid.Lng),
		Area:        geoInfo.Area,
		Length:      geoInfo.Length,
		LayerID:     layerID,
		Properties: map[string]interface{}{
			"source_id": sourceID,
			"distance":  distanceMeters,
		},
	}

	if err := u.regionRepo.Create(ctx, region); err != nil {
		return nil, err
	}
	return region, nil
}

func (u *regionUsecaseImpl) toResponse(region *domain.Region) *dto.RegionResponse {
	return &dto.RegionResponse{
		ID:          region.ID,
		Name:        region.Name,
		DisplayName: region.DisplayName,
		Description: region.Description,
		Type:        uint32(region.Type),
		TypeName:    region.Type.String(),
		Status:      uint32(region.Status),
		StatusName:  region.Status.String(),
		Geometry:    region.Geometry,
		Centroid:    region.Centroid,
		BBox:        region.BBox,
		Area:        region.Area,
		Length:      region.Length,
		LayerID:     region.LayerID,
		Properties:  region.Properties,
		CreatedAt:   region.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   region.UpdatedAt.Format(time.RFC3339),
	}
}
