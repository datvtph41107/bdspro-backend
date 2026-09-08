package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"map/internal/dto"
	"map/internal/interface/repo"
)

// MapPointUsecase defines the usecase for map point operations
type MapPointUsecase struct {
	mapPointRepo repo.MapPointRepository
}

// NewMapPointUsecase creates a new map point usecase
func NewMapPointUsecase(mapPointRepo repo.MapPointRepository) *MapPointUsecase {
	return &MapPointUsecase{
		mapPointRepo: mapPointRepo,
	}
}

// CreateMapPoint creates a new map point
func (u *MapPointUsecase) CreateMapPoint(ctx context.Context, req *dto.CreateMapPointRequest) (*dto.MapPointResponse, error) {
	// Validate coordinates
	if req.Lat < -90 || req.Lat > 90 {
		return nil, errors.New("invalid latitude: must be between -90 and 90")
	}
	if req.Lng < -180 || req.Lng > 180 {
		return nil, errors.New("invalid longitude: must be between -180 and 180")
	}
	if strings.TrimSpace(req.Name) == "" {
		return nil, errors.New("name cannot be empty")
	}

	// Create map point
	mapPoint, err := u.mapPointRepo.CreateMapPoint(ctx, req.Lat, req.Lng, req.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to create map point: %w", err)
	}

	return &dto.MapPointResponse{
		ID:   mapPoint.ID,
		Lat:  mapPoint.Lat,
		Lng:  mapPoint.Lng,
		Name: mapPoint.Name,
	}, nil
}

// FindNearbyLocations finds locations near a given point
func (u *MapPointUsecase) FindNearbyLocations(ctx context.Context, req *dto.NearbyLocationsRequest) (*dto.NearbyLocationsResponse, error) {
	// Validate coordinates
	if req.Lat < -90 || req.Lat > 90 {
		return nil, errors.New("invalid latitude: must be between -90 and 90")
	}
	if req.Lng < -180 || req.Lng > 180 {
		return nil, errors.New("invalid longitude: must be between -180 and 180")
	}
	if req.Radius <= 0 {
		req.Radius = 5000 // Default 5km
	}

	// Find nearby locations
	locations, err := u.mapPointRepo.FindNearbyLocations(ctx, req.Lat, req.Lng, req.Radius)
	if err != nil {
		return nil, fmt.Errorf("failed to find nearby locations: %w", err)
	}

	// Convert to response DTOs
	responseData := make([]dto.MapPointResponse, len(locations))
	for i, location := range locations {
		responseData[i] = dto.MapPointResponse{
			ID:   location.ID,
			Lat:  location.Lat,
			Lng:  location.Lng,
			Name: location.Name,
		}
	}

	return &dto.NearbyLocationsResponse{
		Data: responseData,
	}, nil
}

// FindLocationsByPolygon finds locations within a polygon
func (u *MapPointUsecase) FindLocationsByPolygon(ctx context.Context, req *dto.FindByPolygonRequest) (*dto.FindByPolygonResponse, error) {
	// Validate polygon
	if len(req.Coordinates) < 3 {
		return nil, errors.New("polygon must have at least 3 points")
	}

	// Validate coordinates
	for i, coord := range req.Coordinates {
		if coord.Lat < -90 || coord.Lat > 90 {
			return nil, fmt.Errorf("invalid latitude at point %d: must be between -90 and 90", i+1)
		}
		if coord.Lng < -180 || coord.Lng > 180 {
			return nil, fmt.Errorf("invalid longitude at point %d: must be between -180 and 180", i+1)
		}
	}

	// Convert coordinates to polygon string
	var coords []string
	for _, coord := range req.Coordinates {
		coords = append(coords, fmt.Sprintf("%f %f", coord.Lng, coord.Lat))
	}
	// Close the polygon by adding the first point again
	coords = append(coords, fmt.Sprintf("%f %f", req.Coordinates[0].Lng, req.Coordinates[0].Lat))
	polygon := fmt.Sprintf("POLYGON((%s))", strings.Join(coords, ", "))

	// Find locations in polygon
	locations, err := u.mapPointRepo.FindLocationsInPolygon(ctx, polygon)
	if err != nil {
		return nil, fmt.Errorf("failed to find locations in polygon: %w", err)
	}

	// Convert to response DTOs
	responseData := make([]dto.MapPointResponse, len(locations))
	for i, location := range locations {
		responseData[i] = dto.MapPointResponse{
			ID:   location.ID,
			Lat:  location.Lat,
			Lng:  location.Lng,
			Name: location.Name,
		}
	}

	return &dto.FindByPolygonResponse{
		Data: responseData,
	}, nil
}

// GetMapPointByID gets a map point by ID
func (u *MapPointUsecase) GetMapPointByID(ctx context.Context, id uint) (*dto.MapPointResponse, error) {
	mapPoint, err := u.mapPointRepo.GetMapPointByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get map point: %w", err)
	}

	return &dto.MapPointResponse{
		ID:   mapPoint.ID,
		Lat:  mapPoint.Lat,
		Lng:  mapPoint.Lng,
		Name: mapPoint.Name,
	}, nil
}

// UpdateMapPoint updates an existing map point
func (u *MapPointUsecase) UpdateMapPoint(ctx context.Context, id uint, req *dto.CreateMapPointRequest) (*dto.MapPointResponse, error) {
	// Validate coordinates
	if req.Lat < -90 || req.Lat > 90 {
		return nil, errors.New("invalid latitude: must be between -90 and 90")
	}
	if req.Lng < -180 || req.Lng > 180 {
		return nil, errors.New("invalid longitude: must be between -180 and 180")
	}
	if strings.TrimSpace(req.Name) == "" {
		return nil, errors.New("name cannot be empty")
	}

	// Update map point
	mapPoint, err := u.mapPointRepo.UpdateMapPoint(ctx, id, req.Lat, req.Lng, req.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to update map point: %w", err)
	}

	return &dto.MapPointResponse{
		ID:   mapPoint.ID,
		Lat:  mapPoint.Lat,
		Lng:  mapPoint.Lng,
		Name: mapPoint.Name,
	}, nil
}

// DeleteMapPoint deletes a map point by ID
func (u *MapPointUsecase) DeleteMapPoint(ctx context.Context, id uint) error {
	err := u.mapPointRepo.DeleteMapPoint(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete map point: %w", err)
	}
	return nil
}
