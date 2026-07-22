// internal/repo/postgres/location_repo.go
package postgres

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"sync"

	"tqd/internal/domain"
	"tqd/internal/dto"
	"tqd/internal/interface/repo"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type locationRepo struct {
	db *gorm.DB
}

func NewLocationRepository(db *gorm.DB) repo.LocationRepository {
	return &locationRepo{db: db}
}

// ==================== PROVINCE OPERATIONS ====================

func (r *locationRepo) GetProvinceByID(ctx context.Context, id string) (*domain.Province, error) {
	var province domain.Province
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&province).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &province, err
}

func (r *locationRepo) GetProvinceByCode(ctx context.Context, code string) (*domain.Province, error) {
	var province domain.Province
	err := r.db.WithContext(ctx).
		Where("code = ?", code).
		First(&province).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &province, err
}

func (r *locationRepo) GetNearestProvince(ctx context.Context, lat, lng float64, maxDistanceKm float64) (*domain.Province, float64, error) {
	var province domain.Province

	// Sử dụng raw query cho hiệu suất tốt hơn
	query := `
		SELECT id, full_name, short_name, lat, lng, code, created_at
		FROM province_v2
		WHERE (6371 * acos(cos(radians(?)) * cos(radians(lat)) * 
		      cos(radians(lng) - radians(?)) + sin(radians(?)) * 
		      sin(radians(lat)))) < ?
		ORDER BY (6371 * acos(cos(radians(?)) * cos(radians(lat)) * 
		         cos(radians(lng) - radians(?)) + sin(radians(?)) * 
		         sin(radians(lat))))
		LIMIT 1
	`

	err := r.db.WithContext(ctx).Raw(query, lat, lng, lat, maxDistanceKm, lat, lng, lat).Scan(&province).Error
	if err != nil {
		return nil, 0, fmt.Errorf("query nearest province: %w", err)
	}

	if province.ID == "" {
		return nil, 0, nil
	}

	distance := r.haversine(lat, lng, province.Lat, province.Lng)
	return &province, distance, nil
}

func (r *locationRepo) SearchProvinces(ctx context.Context, query string, limit int) ([]*domain.Province, error) {
	var provinces []*domain.Province
	searchPattern := "%" + strings.ToLower(query) + "%"

	err := r.db.WithContext(ctx).
		Where("LOWER(full_name) LIKE ? OR LOWER(short_name) LIKE ? OR code = ?",
			searchPattern, searchPattern, query).
		Order(clause.Expr{
			SQL: `CASE
				WHEN LOWER(full_name) = LOWER(?) THEN 1
				WHEN LOWER(short_name) = LOWER(?) THEN 2
				WHEN LOWER(full_name) LIKE LOWER(?) THEN 3
				WHEN LOWER(short_name) LIKE LOWER(?) THEN 4
				ELSE 5 END, LENGTH(full_name)`,
			Vars:               []interface{}{query, query, query + "%", query + "%"},
			WithoutParentheses: true,
		}).
		Limit(limit).
		Find(&provinces).Error

	return provinces, err
}

func (r *locationRepo) ListAllProvinces(ctx context.Context) ([]*domain.Province, error) {
	var provinces []*domain.Province
	err := r.db.WithContext(ctx).
		Order("full_name").
		Find(&provinces).Error
	return provinces, err
}

// ==================== WARD OPERATIONS ====================

func (r *locationRepo) GetWardByID(ctx context.Context, id string) (*domain.Ward, error) {
	var ward domain.Ward
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&ward).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &ward, err
}

func (r *locationRepo) GetWardByCode(ctx context.Context, code string) (*domain.Ward, error) {
	var ward domain.Ward
	err := r.db.WithContext(ctx).
		Where("code = ?", code).
		First(&ward).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &ward, err
}

func (r *locationRepo) GetNearestWard(ctx context.Context, lat, lng float64, provinceID string, maxDistanceKm float64) (*domain.Ward, float64, error) {
	var ward domain.Ward

	query := `
		SELECT id, full_name, short_name, lat, lng, code, province_id
		FROM ward_v2
		WHERE province_id = ?
		  AND (6371 * acos(cos(radians(?)) * cos(radians(lat)) * 
		       cos(radians(lng) - radians(?)) + sin(radians(?)) * 
		       sin(radians(lat)))) < ?
		ORDER BY (6371 * acos(cos(radians(?)) * cos(radians(lat)) * 
		         cos(radians(lng) - radians(?)) + sin(radians(?)) * 
		         sin(radians(lat))))
		LIMIT 1
	`

	err := r.db.WithContext(ctx).Raw(query, provinceID, lat, lng, lat, maxDistanceKm, lat, lng, lat).Scan(&ward).Error
	if err != nil {
		return nil, 0, fmt.Errorf("query nearest ward: %w", err)
	}

	if ward.ID == "" {
		return nil, 0, nil
	}

	distance := r.haversine(lat, lng, ward.Lat, ward.Lng)
	return &ward, distance, nil
}

func (r *locationRepo) SearchWards(ctx context.Context, query string, provinceID string, limit int) ([]*domain.Ward, error) {
	var wards []*domain.Ward
	searchPattern := "%" + strings.ToLower(query) + "%"

	db := r.db.WithContext(ctx).
		Where("LOWER(full_name) LIKE ? OR LOWER(short_name) LIKE ?",
			searchPattern, searchPattern)

	if provinceID != "" {
		db = db.Where("province_id = ?", provinceID)
	}

	err := db.Order(clause.Expr{
		SQL: `CASE
				WHEN LOWER(full_name) = LOWER(?) THEN 1
				WHEN LOWER(short_name) = LOWER(?) THEN 2
				WHEN LOWER(full_name) LIKE LOWER(?) THEN 3
				WHEN LOWER(short_name) LIKE LOWER(?) THEN 4
				ELSE 5 END, LENGTH(full_name)`,
		Vars:               []interface{}{query, query, query + "%", query + "%"},
		WithoutParentheses: true,
	}).
		Limit(limit).
		Find(&wards).Error

	return wards, err
}

func (r *locationRepo) ListWardsByProvince(ctx context.Context, provinceID string, page, pageSize int) ([]*domain.Ward, int64, error) {
	var wards []*domain.Ward
	var total int64

	// Count total
	err := r.db.WithContext(ctx).
		Model(&domain.Ward{}).
		Where("province_id = ?", provinceID).
		Count(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("count wards: %w", err)
	}

	// Get paginated data
	offset := (page - 1) * pageSize
	err = r.db.WithContext(ctx).
		Where("province_id = ?", provinceID).
		Order("full_name").
		Offset(offset).
		Limit(pageSize).
		Find(&wards).Error

	return wards, total, err
}

// ==================== BATCH OPERATIONS ====================

func (r *locationRepo) GetNearestLocations(ctx context.Context, coordinates []*dto.Coordinate, maxDistanceKm float64) ([]*dto.LocationResult, error) {
	results := make([]*dto.LocationResult, len(coordinates))
	var wg sync.WaitGroup
	errCh := make(chan error, len(coordinates))

	for i, coord := range coordinates {
		if coord == nil {
			continue
		}

		wg.Add(1)
		go func(idx int, lat, lng float64) {
			defer wg.Done()

			result, err := r.getNearestLocationForCoordinate(ctx, lat, lng, maxDistanceKm)
			if err != nil {
				errCh <- fmt.Errorf("coordinate[%d]: %w", idx, err)
				return
			}
			results[idx] = result
		}(i, coord.Latitude, coord.Longitude)
	}

	// Wait for all goroutines
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
		close(errCh)
	}()

	select {
	case <-done:
		// Check for errors
		for err := range errCh {
			if err != nil {
				return nil, err
			}
		}
		return results, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (r *locationRepo) SearchTxtClient(ctx context.Context, query string, limit int) ([]*domain.SearchItem, error) {
	var searchItems []*domain.SearchItem
	err := r.db.WithContext(ctx).
		Where("search_text LIKE ?", "%"+strings.ToLower(query)+"%").
		Limit(limit).
		Find(&searchItems).Error
	return searchItems, err
}

// ==================== PRIVATE HELPERS ====================

func (r *locationRepo) getNearestLocationForCoordinate(ctx context.Context, lat, lng, maxDistanceKm float64) (*dto.LocationResult, error) {
	// Tìm province gần nhất
	province, provDist, err := r.GetNearestProvince(ctx, lat, lng, maxDistanceKm)
	if err != nil {
		return nil, fmt.Errorf("get nearest province: %w", err)
	}
	if province == nil {
		return nil, nil
	}

	// Tìm ward gần nhất trong province đó
	ward, wardDist, _ := r.GetNearestWard(ctx, lat, lng, province.ID, maxDistanceKm)

	// Tính confidence
	confidence := r.calculateConfidence(provDist, wardDist)

	return &dto.LocationResult{
		Province:   province,
		Ward:       ward,
		Distance:   provDist,
		Confidence: confidence,
	}, nil
}

func (r *locationRepo) haversine(lat1, lng1, lat2, lng2 float64) float64 {
	const R = 6371 // Earth radius in km

	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	dlat := (lat2 - lat1) * math.Pi / 180
	dlng := (lng2 - lng1) * math.Pi / 180

	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dlng/2)*math.Sin(dlng/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

func (r *locationRepo) calculateConfidence(provDist, wardDist float64) float64 {
	confidence := 0.5

	switch {
	case provDist < 10:
		confidence += 0.4
	case provDist < 30:
		confidence += 0.3
	case provDist < 50:
		confidence += 0.2
	case provDist < 100:
		confidence += 0.1
	default:
		confidence -= 0.2
	}

	switch {
	case wardDist > 0 && wardDist < 5:
		confidence += 0.2
	case wardDist > 0 && wardDist < 15:
		confidence += 0.1
	}

	return r.clamp(confidence, 0, 1)
}

func (r *locationRepo) clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
