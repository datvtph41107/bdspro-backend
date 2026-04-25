package postgres

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/dto"
	"tqd/internal/interface/repo"

	"gorm.io/gorm"
)

type ParcelPostgres struct {
	DB *gorm.DB
}

func NewParcelPostgres(db *gorm.DB) repo.IParcelRepo {
	return &ParcelPostgres{DB: db}
}

func (p *ParcelPostgres) CreateBatch(ctx context.Context, parcels []qh_domain.Parcel) error {
	const batchSize = 500

	return p.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i := 0; i < len(parcels); i += batchSize {
			end := i + batchSize
			if end > len(parcels) {
				end = len(parcels)
			}

			var values []string
			var args []interface{}

			for _, parcel := range parcels[i:end] {
				values = append(values, `(ST_SetSRID(ST_GeomFromGeoJSON(?),4326), ?, ?, ?, ?)`)

				args = append(args,
					parcel.Geometry,
					parcel.Lat,
					parcel.Lng,
					parcel.RefID,
					parcel.RefType,
				)
			}

			query := `
				INSERT INTO parcels (geometry, lat, lon, ref_id, ref_type)
				VALUES ` + strings.Join(values, ",")

			if err := tx.Exec(query, args...).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (p *ParcelPostgres) FindByLocation(
	ctx context.Context,
	lat, lng float64,
) (*qh_domain.Parcel, error) {

	var parcel qh_domain.Parcel

	err := p.DB.WithContext(ctx).
		Raw(`
			SELECT 
				id,
				ST_AsGeoJSON(geometry) as geometry,
				lat,
				lon,
				ref_id,
				ref_type
			FROM parcels
			WHERE ST_Intersects(
				geometry,
				ST_SetSRID(ST_Point(?, ?), 4326)
			)
			LIMIT 1
		`, lng, lat).
		Scan(&parcel).Error

	if err != nil {
		return nil, err
	}

	return &parcel, nil
}

func computeBBox(points []dto.Point) (minLon, minLat, maxLon, maxLat float64) {
	minLon, minLat = points[0].Lng, points[0].Lat
	maxLon, maxLat = points[0].Lng, points[0].Lat

	for _, p := range points {
		if p.Lng < minLon {
			minLon = p.Lng
		}
		if p.Lng > maxLon {
			maxLon = p.Lng
		}
		if p.Lat < minLat {
			minLat = p.Lat
		}
		if p.Lat > maxLat {
			maxLat = p.Lat
		}
	}
	return
}

func (p *ParcelPostgres) FindIntersectPolygon(
	ctx context.Context,
	points []dto.Point,
	limit, offset int,
) ([]qh_domain.Parcel, int64, error) {

	if len(points) < 3 {
		return nil, 0, fmt.Errorf("invalid polygon")
	}

	polygon := buildPolygonWKT(points)

	var parcels []qh_domain.Parcel
	var total int64

	countQuery := `
		WITH poly AS (
			SELECT ST_GeomFromText(?, 4326) AS geom
		)
		SELECT COUNT(*)
		FROM parcels p, poly
		WHERE 
			p.geometry && poly.geom
			AND ST_Intersects(p.geometry, poly.geom)
	`

	if err := p.DB.WithContext(ctx).
		Raw(countQuery, polygon).
		Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []qh_domain.Parcel{}, 0, nil
	}

	query := `
		WITH poly AS (
			SELECT ST_GeomFromText(?, 4326) AS geom
		)
		SELECT 
			p.id,
			ST_AsGeoJSON(p.geometry) as geometry,
			p.lat,
			p.lon,
			p.ref_id,
			p.ref_type,
			p.created_at,
			p.updated_at
		FROM parcels p, poly
		WHERE 
			p.geometry && poly.geom
			AND ST_Intersects(p.geometry, poly.geom)
		ORDER BY p.id
		LIMIT ? OFFSET ?
	`

	err := p.DB.WithContext(ctx).
		Raw(query, polygon, limit, offset).
		Scan(&parcels).Error

	return parcels, total, err
}

func (p *ParcelPostgres) FindWithinPolygon(
	ctx context.Context,
	points []dto.Point,
	limit, offset int,
) ([]qh_domain.Parcel, int64, error) {
	if len(points) < 3 {
		return nil, 0, fmt.Errorf("invalid polygon")
	}

	polygon := buildPolygonWKT(points)

	var parcels []qh_domain.Parcel
	var total int64

	countQuery := `
		WITH poly AS (
			SELECT ST_GeomFromText(?, 4326) AS geom
		)
		SELECT COUNT(*)
		FROM parcels p, poly
		WHERE ST_Contains(poly.geom, p.geometry)
	`

	if err := p.DB.WithContext(ctx).
		Raw(countQuery, polygon).
		Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	query := `
		WITH poly AS (
			SELECT ST_GeomFromText(?, 4326) AS geom
		)
		SELECT 
			p.id,
			ST_AsGeoJSON(p.geometry) as geometry,
			p.lat,
			p.lon,
			p.ref_id,
			p.ref_type,
			p.created_at,
			p.updated_at
		FROM parcels p, poly
		WHERE ST_Contains(poly.geom, p.geometry)
		ORDER BY p.id
		LIMIT ? OFFSET ?
	`

	err := p.DB.WithContext(ctx).
		Raw(query, polygon, limit, offset).
		Scan(&parcels).Error

	return parcels, total, err
}

func buildPolygonWKT(points []dto.Point) string {
	var sb strings.Builder
	sb.Grow(len(points) * 32)

	sb.WriteString("POLYGON((")

	for i := 0; i < len(points); i++ {
		p := points[i]

		sb.WriteString(strconv.FormatFloat(p.Lng, 'f', 6, 64))
		sb.WriteByte(' ')
		sb.WriteString(strconv.FormatFloat(p.Lat, 'f', 6, 64))

		if i < len(points)-1 {
			sb.WriteByte(',')
		}
	}

	sb.WriteByte(',')
	sb.WriteString(strconv.FormatFloat(points[0].Lng, 'f', 6, 64))
	sb.WriteByte(' ')
	sb.WriteString(strconv.FormatFloat(points[0].Lat, 'f', 6, 64))

	sb.WriteString("))")

	return sb.String()
}
