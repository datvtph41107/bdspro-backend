package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"tqd/internal/domain"
	"tqd/internal/interface/repo"
)

type FeatureUsecase struct {
	repo repo.FeatureRepository
}

func NewFeatureUsecase(repo repo.FeatureRepository) *FeatureUsecase {
	return &FeatureUsecase{repo: repo}
}

type PointRadiusRequest struct {
	Lat    float64
	Lng    float64
	Radius int
	Limit  int
	Offset int
}

type PolygonRequest struct {
	PolygonGeoJSON []byte
	Limit          int
	Offset         int
}

func (u *FeatureUsecase) GetFeaturesByPointRadius(ctx context.Context, req *PointRadiusRequest) ([]*domain.Feature, int64, error) {
	if req.Radius <= 0 {
		return nil, 0, fmt.Errorf("radius must be positive")
	}
	if req.Lat < -90 || req.Lat > 90 || req.Lng < -180 || req.Lng > 180 {
		return nil, 0, fmt.Errorf("invalid coordinates")
	}
	if req.Limit <= 0 {
		req.Limit = 1000 // default
	}
	if req.Offset < 0 {
		req.Offset = 0
	}
	return u.repo.GetByPointRadius(ctx, req.Lat, req.Lng, req.Radius, req.Limit, req.Offset)
}

func (u *FeatureUsecase) GetFeaturesByPolygon(ctx context.Context, polygonGeoJSON []byte, limit, offset int) ([]*domain.Feature, int64, error) {
	// Validate GeoJSON
	var geom map[string]interface{}
	if err := json.Unmarshal(polygonGeoJSON, &geom); err != nil {
		return nil, 0, fmt.Errorf("invalid GeoJSON: %w", err)
	}
	geomType, ok := geom["type"].(string)
	if !ok {
		return nil, 0, fmt.Errorf("missing geometry type")
	}
	if geomType != "Polygon" && geomType != "MultiPolygon" {
		return nil, 0, fmt.Errorf("geometry type must be Polygon or MultiPolygon, got %s", geomType)
	}
	if limit <= 0 {
		limit = 1000
	}
	if offset < 0 {
		offset = 0
	}
	return u.repo.GetByPolygon(ctx, polygonGeoJSON, limit, offset)
}
