package postgre

import (
	"context"
	"fmt"

	"map/internal/domain"
	"map/internal/interface/repo"

	"gorm.io/gorm"
)

// MapPointPostgres implements MapPointRepository interface
type MapPointPostgres struct {
	db *gorm.DB
}

// NewMapPointPostgres creates a new map point postgres repository
// @bind: map/internal/interface/repo.MapPointRepository
func NewMapPointPostgres(db *gorm.DB) repo.MapPointRepository {
	return &MapPointPostgres{
		db: db,
	}
}

// FindNearbyLocations finds locations near a given point within a radius
func (r *MapPointPostgres) FindNearbyLocations(ctx context.Context, lat, lng, radius float64) ([]domain.Location, error) {
	var points []domain.Location
	query := `
		SELECT 
			id, 
			ST_Y(geom::geometry) AS lat, 
			ST_X(geom::geometry) AS lng, 
			name,
			created_at,
			updated_at,
			ST_Distance(geom, ST_MakePoint(?, ?)::geography) AS distance
		FROM map_points
		WHERE ST_Distance(geom, ST_MakePoint(?, ?)::geography) <= ?
		ORDER BY distance
	`

	result := r.db.WithContext(ctx).Raw(query, lng, lat, lng, lat, radius).Scan(&points)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find nearby locations: %w", result.Error)
	}

	return points, nil
}

// FindLocationsInPolygon finds locations within a polygon
func (r *MapPointPostgres) FindLocationsInPolygon(ctx context.Context, polygon string) ([]domain.Location, error) {
	var points []domain.Location
	query := `
		SELECT 
			id,
			ST_Y(geom::geometry) AS lat, 
			ST_X(geom::geometry) AS lng, 
			name,
			created_at,
			updated_at
		FROM get_partitioned_map_points(ST_GeogFromText(?))
		WHERE ST_Within(geom::geometry, ST_SetSRID(?::geometry, 4326))
		LIMIT 200
	`

	result := r.db.WithContext(ctx).Raw(query, polygon, polygon).Scan(&points)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find locations in polygon: %w", result.Error)
	}

	return points, nil
}

// CreateMapPoint creates a new map point
func (r *MapPointPostgres) CreateMapPoint(ctx context.Context, lat, lng float64, name string) (*domain.MapPoint, error) {
	query := `SELECT create_point(CAST(? AS NUMERIC), CAST(? AS NUMERIC), CAST(? AS VARCHAR))`

	result := r.db.WithContext(ctx).Exec(query, lat, lng, name)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to create map point: %w", result.Error)
	}

	// Get the created point
	var point domain.MapPoint
	getQuery := `
		SELECT 
			id,
			ST_Y(geom::geometry) AS lat,
			ST_X(geom::geometry) AS lng,
			name,
			created_at,
			updated_at
		FROM map_points 
		WHERE ST_Y(geom::geometry) = ? AND ST_X(geom::geometry) = ? AND name = ?
		ORDER BY created_at DESC
		LIMIT 1
	`

	result = r.db.WithContext(ctx).Raw(getQuery, lat, lng, name).Scan(&point)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get created map point: %w", result.Error)
	}

	return &point, nil
}

// GetMapPointByID gets a map point by ID
func (r *MapPointPostgres) GetMapPointByID(ctx context.Context, id uint) (*domain.MapPoint, error) {
	var point domain.MapPoint
	query := `
		SELECT 
			id,
			ST_Y(geom::geometry) AS lat,
			ST_X(geom::geometry) AS lng,
			name,
			created_at,
			updated_at
		FROM map_points 
		WHERE id = ?
	`

	result := r.db.WithContext(ctx).Raw(query, id).Scan(&point)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get map point: %w", result.Error)
	}

	if point.ID == 0 {
		return nil, fmt.Errorf("map point not found")
	}

	return &point, nil
}

// UpdateMapPoint updates an existing map point
func (r *MapPointPostgres) UpdateMapPoint(ctx context.Context, id uint, lat, lng float64, name string) (*domain.MapPoint, error) {
	// First check if the point exists
	_, err := r.GetMapPointByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update the point
	query := `
		UPDATE map_points 
		SET geom = ST_SetSRID(ST_MakePoint(?, ?), 4326), name = ?, updated_at = NOW()
		WHERE id = ?
	`

	result := r.db.WithContext(ctx).Exec(query, lng, lat, name, id)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to update map point: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("map point not found")
	}

	// Get the updated point
	updatedPoint, err := r.GetMapPointByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return updatedPoint, nil
}

// DeleteMapPoint deletes a map point by ID
func (r *MapPointPostgres) DeleteMapPoint(ctx context.Context, id uint) error {
	// First check if the point exists
	_, err := r.GetMapPointByID(ctx, id)
	if err != nil {
		return err
	}

	query := `DELETE FROM map_points WHERE id = ?`

	result := r.db.WithContext(ctx).Exec(query, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete map point: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("map point not found")
	}

	return nil
}
