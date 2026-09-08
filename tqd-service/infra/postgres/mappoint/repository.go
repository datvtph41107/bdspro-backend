package mappoint

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	domain "tqd/internal/domain/mappoint"
	usecase "tqd/internal/usecase/mappoint"
)

type Repository struct{ database *gorm.DB }

func NewRepository(database *gorm.DB) *Repository { return &Repository{database: database} }

func (r *Repository) Create(ctx context.Context, point domain.Point) (domain.Point, error) {
	var created domain.Point
	err := r.database.WithContext(ctx).Raw(`
		INSERT INTO map_points (name, geom)
		VALUES (?, ST_SetSRID(ST_MakePoint(?, ?), 4326))
		RETURNING id, name, ST_Y(geom) AS latitude, ST_X(geom) AS longitude, created_at, updated_at`,
		point.Name, point.Longitude, point.Latitude,
	).Scan(&created).Error
	if err != nil {
		return domain.Point{}, fmt.Errorf("create map point: %w", err)
	}
	return created, nil
}

func (r *Repository) Get(ctx context.Context, id uint64) (domain.Point, error) {
	var point domain.Point
	result := r.database.WithContext(ctx).Raw(`
		SELECT id, name, ST_Y(geom) AS latitude, ST_X(geom) AS longitude, created_at, updated_at
		FROM map_points WHERE id = ?`, id,
	).Scan(&point)
	if result.Error != nil {
		return domain.Point{}, fmt.Errorf("get map point: %w", result.Error)
	}
	if result.RowsAffected == 0 || point.ID == 0 {
		return domain.Point{}, usecase.ErrNotFound
	}
	return point, nil
}

func (r *Repository) Update(ctx context.Context, point domain.Point) (domain.Point, error) {
	var updated domain.Point
	result := r.database.WithContext(ctx).Raw(`
		UPDATE map_points
		SET name = ?, geom = ST_SetSRID(ST_MakePoint(?, ?), 4326), updated_at = NOW()
		WHERE id = ?
		RETURNING id, name, ST_Y(geom) AS latitude, ST_X(geom) AS longitude, created_at, updated_at`,
		point.Name, point.Longitude, point.Latitude, point.ID,
	).Scan(&updated)
	if result.Error != nil {
		return domain.Point{}, fmt.Errorf("update map point: %w", result.Error)
	}
	if result.RowsAffected == 0 || updated.ID == 0 {
		return domain.Point{}, usecase.ErrNotFound
	}
	return updated, nil
}

func (r *Repository) Delete(ctx context.Context, id uint64) error {
	result := r.database.WithContext(ctx).Exec(`DELETE FROM map_points WHERE id = ?`, id)
	if result.Error != nil {
		return fmt.Errorf("delete map point: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return usecase.ErrNotFound
	}
	return nil
}

func (r *Repository) FindNearby(ctx context.Context, center domain.Coordinate, radius float64, limit int) ([]domain.Point, error) {
	var points []domain.Point
	err := r.database.WithContext(ctx).Raw(`
		SELECT id, name, ST_Y(geom) AS latitude, ST_X(geom) AS longitude, created_at, updated_at
		FROM map_points
		WHERE ST_DWithin(geom::geography, ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography, ?)
		ORDER BY ST_Distance(geom::geography, ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography), id
		LIMIT ?`,
		center.Longitude, center.Latitude, radius,
		center.Longitude, center.Latitude, limit,
	).Scan(&points).Error
	if err != nil {
		return nil, fmt.Errorf("find nearby map points: %w", err)
	}
	return points, nil
}

func (r *Repository) FindWithinPolygon(ctx context.Context, coordinates []domain.Coordinate, limit int) ([]domain.Point, error) {
	parts := make([]string, 0, len(coordinates)+1)
	for _, coordinate := range coordinates {
		parts = append(parts, fmt.Sprintf("%g %g", coordinate.Longitude, coordinate.Latitude))
	}
	first := coordinates[0]
	last := coordinates[len(coordinates)-1]
	if first != last {
		parts = append(parts, fmt.Sprintf("%g %g", first.Longitude, first.Latitude))
	}
	wkt := "POLYGON((" + strings.Join(parts, ",") + "))"

	var points []domain.Point
	err := r.database.WithContext(ctx).Raw(`
		SELECT id, name, ST_Y(geom) AS latitude, ST_X(geom) AS longitude, created_at, updated_at
		FROM map_points
		WHERE ST_Covers(ST_GeomFromText(?, 4326), geom)
		ORDER BY id
		LIMIT ?`, wkt, limit,
	).Scan(&points).Error
	if err != nil {
		if errors.Is(err, gorm.ErrInvalidData) {
			return nil, fmt.Errorf("%w: invalid polygon", usecase.ErrInvalidInput)
		}
		return nil, fmt.Errorf("find map points in polygon: %w", err)
	}
	return points, nil
}

var _ usecase.Repository = (*Repository)(nil)
