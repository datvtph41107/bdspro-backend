package mapper

import (
	"map/internal/domain"
	"map/internal/dto"
)

// MapPointMapper handles mapping between domain entities and DTOs
type MapPointMapper struct{}

// NewMapPointMapper creates a new map point mapper
func NewMapPointMapper() *MapPointMapper {
	return &MapPointMapper{}
}

// ToMapPointResponse converts domain MapPoint to MapPointResponse DTO
func (m *MapPointMapper) ToMapPointResponse(point *domain.MapPoint) *dto.MapPointResponse {
	if point == nil {
		return nil
	}

	return &dto.MapPointResponse{
		ID:   point.ID,
		Lat:  point.Lat,
		Lng:  point.Lng,
		Name: point.Name,
	}
}

// ToMapPointResponseList converts slice of domain MapPoint to slice of MapPointResponse DTO
func (m *MapPointMapper) ToMapPointResponseList(points []domain.MapPoint) []dto.MapPointResponse {
	if points == nil {
		return nil
	}

	responses := make([]dto.MapPointResponse, len(points))
	for i, point := range points {
		responses[i] = *m.ToMapPointResponse(&point)
	}

	return responses
}

// ToCreateMapPointRequest converts CreateMapPointRequest DTO to domain values
func (m *MapPointMapper) ToCreateMapPointRequest(req *dto.CreateMapPointRequest) (lat, lng float64, name string) {
	if req == nil {
		return 0, 0, ""
	}

	return req.Lat, req.Lng, req.Name
}

// ToNearbyLocationsRequest converts NearbyLocationsRequest DTO to domain values
func (m *MapPointMapper) ToNearbyLocationsRequest(req *dto.NearbyLocationsRequest) (lat, lng, radius float64) {
	if req == nil {
		return 0, 0, 5000 // Default radius
	}

	radius = req.Radius
	if radius <= 0 {
		radius = 5000 // Default 5km
	}

	return req.Lat, req.Lng, radius
}

// ToFindByPolygonRequest converts FindByPolygonRequest DTO to domain values
func (m *MapPointMapper) ToFindByPolygonRequest(req *dto.FindByPolygonRequest) []dto.MapPointCoordinate {
	if req == nil {
		return nil
	}

	return req.Coordinates
}
