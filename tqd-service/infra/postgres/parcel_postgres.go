package postgres

import (
	"context"
	"fmt"
	"strings"
	"tqd/internal/domain"
	"tqd/internal/interface/repo"

	"gorm.io/gorm"
)

type ParcelPostgres struct {
	DB *gorm.DB
}

func NewParcelPostgres(db *gorm.DB) repo.IParcelRepo {
	return &ParcelPostgres{DB: db}
}

func (p *ParcelPostgres) CreateBatch(ctx context.Context, parcels []domain.Parcel) error {
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
) (*domain.Parcel, error) {

	var parcel domain.Parcel

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

func (p *ParcelPostgres) FindByPolygon(
	ctx context.Context,
	geoJSON string,
	limit, offset int,
) ([]domain.Parcel, int64, error) {
	var parcels []domain.Parcel
	var total int64

	countQuery := `
		SELECT COUNT(*)
		FROM parcels
		WHERE ST_Within(
			ST_Centroid(geometry),
			ST_SetSRID(ST_GeomFromGeoJSON(?), 4326)
		)
	`
	if err := p.DB.WithContext(ctx).Raw(countQuery, geoJSON).Scan(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count query failed: %w", err)
	}

	if total == 0 {
		return []domain.Parcel{}, 0, nil
	}

	// Lấy dữ liệu phân trang
	query := `
		SELECT 
			id,
			ST_AsGeoJSON(geometry) as geometry,
			lat,
			lon,
			ref_id,
			ref_type,
			created_at,
			updated_at
		FROM parcels
		WHERE ST_Within(
			ST_Centroid(geometry),
			ST_SetSRID(ST_GeomFromGeoJSON(?), 4326)
		)
		ORDER BY id
		LIMIT ? OFFSET ?
	`
	err := p.DB.WithContext(ctx).Raw(query, geoJSON, limit, offset).Scan(&parcels).Error
	if err != nil {
		return nil, 0, fmt.Errorf("query failed: %w", err)
	}

	return parcels, total, nil
}
