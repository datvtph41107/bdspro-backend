package postgres

import (
	"context"
	"fmt"
	"math"
	"strings"

	_db "common/db"
	_dto "common/domain/dto"
	_errors "common/errors"
	_provider "common/provider"
	"tqd/internal/domain"
	_repo "tqd/internal/interface/repo"
)

// POIRepo implements the POI repository interface
type POIRepo struct {
	_provider.CrudRepo[domain.POI]
}

// NewPOIRepo creates a new POI repository instance
// @bind: infra/postgres/poi_postgres.go
func NewPOIRepo(db *_db.TransactionRepo) _repo.IPOIRepo {
	repo := &POIRepo{}
	repo.Init(repo, db)
	return repo
}

// Create creates a new POI
func (r *POIRepo) Create(ctx context.Context, poi *domain.POI) error {
	return r.CrudRepo.Create(ctx, poi)
}

// Update updates an existing POI
func (r *POIRepo) Update(ctx context.Context, id uint64, poi *domain.POI) error {
	return r.CrudRepo.Update(ctx, id, poi)
}

// GetListBase overrides the base GetList method to match CrudRepo signature
func (r *POIRepo) GetListBase(ctx context.Context, pagable _dto.IPagable) ([]domain.POI, int64, error) {
	// Convert pagable to our custom request format
	dto, ok := pagable.(*_repo.ListPOIsRequest)
	if !ok {
		return nil, 0, _errors.ReturnError(400, "invalid pagable")
	}

	resp, err := r.GetListCustom(ctx, dto)
	if err != nil {
		return nil, 0, err
	}

	return resp.Data, resp.Total, nil
}

// GetByCode retrieves a POI by its code
func (r *POIRepo) GetByCode(ctx context.Context, code string) (*domain.POI, error) {
	var poi domain.POI
	err := r.GetDB(ctx).Where("code = ? AND deleted_at IS NULL", code).First(&poi).Error
	if err != nil {
		return nil, err
	}
	return &poi, nil
}

// GetList overrides the base GetList method to match CrudRepo signature
func (r *POIRepo) GetList(ctx context.Context, pagable _dto.IPagable) ([]domain.POI, int64, error) {
	// Convert pagable to our custom request format
	dto, ok := pagable.(*_repo.ListPOIsRequest)
	if !ok {
		return nil, 0, _errors.ReturnError(400, "invalid pagable")
	}

	resp, err := r.GetListWithFilters(ctx, dto)
	if err != nil {
		return nil, 0, err
	}

	return resp.Data, resp.Total, nil
}

// GetListCustom retrieves a paginated list of POIs with filters
func (r *POIRepo) GetListCustom(ctx context.Context, req *_repo.ListPOIsRequest) (*_repo.ListPOIsResponse, error) {
	return r.GetListWithFilters(ctx, req)
}

// GetListWithFiltersInterface retrieves a paginated list of POIs with filters for interface
func (r *POIRepo) GetListWithFiltersInterface(ctx context.Context, req *_repo.ListPOIsRequest) (*_repo.ListPOIsResponse, error) {
	return r.GetListWithFilters(ctx, req)
}

// GetListInterface retrieves a paginated list of POIs with filters for interface
func (r *POIRepo) GetListInterface(ctx context.Context, req *_repo.ListPOIsRequest) (*_repo.ListPOIsResponse, error) {
	return r.GetListWithFilters(ctx, req)
}

// GetListWithFilters retrieves a paginated list of POIs with filters
func (r *POIRepo) GetListWithFilters(ctx context.Context, req *_repo.ListPOIsRequest) (*_repo.ListPOIsResponse, error) {
	var pois []domain.POI
	var total int64

	query := r.GetDB(ctx).Model(&domain.POI{}).Preload("Category").Where("deleted_at IS NULL")

	// Apply filters
	if req.Text != "" {
		searchTerm := "%" + strings.ToLower(req.Text) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(address) LIKE ? OR LOWER(code) LIKE ?", searchTerm, searchTerm, searchTerm)
	}

	if req.CategoryID != nil {
		query = query.Where("category_id = ?", *req.CategoryID)
	}

	if req.IsActive != nil {
		query = query.Where("is_active = ?", *req.IsActive)
	}

	if req.MinLat != nil && req.MaxLat != nil {
		query = query.Where("lat BETWEEN ? AND ?", *req.MinLat, *req.MaxLat)
	}

	if req.MinLng != nil && req.MaxLng != nil {
		query = query.Where("lng BETWEEN ? AND ?", *req.MinLng, *req.MaxLng)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Apply pagination
	query = query.
		// Preload("Category").
		Order("created_at DESC").
		Offset(req.GetOffset()).
		Limit(req.GetLimit())

	if err := query.Find(&pois).Error; err != nil {
		return nil, err
	}

	return &_repo.ListPOIsResponse{
		Data:  pois,
		Total: total,
	}, nil
}

// GetListByCategory retrieves POIs by category
func (r *POIRepo) GetListByCategory(ctx context.Context, categoryID uint64, limit int) ([]domain.POI, error) {
	var pois []domain.POI
	query := r.GetDB(ctx).
		Where("category_id = ? AND deleted_at IS NULL AND is_active = ?", categoryID, true).
		Preload("Category").
		Order("rating_point DESC, rating_count DESC").
		Limit(limit)

	if err := query.Find(&pois).Error; err != nil {
		return nil, err
	}

	return pois, nil
}

// GetNearby retrieves POIs within a certain radius from given coordinates
func (r *POIRepo) GetNearby(ctx context.Context, lat, lng, radius float64, limit int) ([]domain.POI, error) {
	var pois []domain.POI

	// Haversine formula for distance calculation
	// Using a simplified approach with bounding box first, then filtering by distance
	latRange := radius / 111.0                               // Approximate km per degree latitude
	lngRange := radius / (111.0 * math.Cos(lat*math.Pi/180)) // Approximate km per degree longitude

	query := r.GetDB(ctx).
		Where(`deleted_at IS NULL AND is_active = ? AND 
			lat BETWEEN ? AND ? AND 
			lng BETWEEN ? AND ?`,
			true,
			lat-latRange, lat+latRange,
			lng-lngRange, lng+lngRange).
		Preload("Category").
		Order("rating_point DESC, rating_count DESC").
		Limit(limit)

	if err := query.Find(&pois).Error; err != nil {
		return nil, err
	}

	// Filter by actual distance (more precise)
	var nearbyPOIs []domain.POI
	for _, poi := range pois {
		distance := r.calculateDistance(lat, lng, poi.Lat, poi.Lng)
		if distance <= radius {
			nearbyPOIs = append(nearbyPOIs, poi)
		}
	}

	return nearbyPOIs, nil
}

// UpdateRating updates the rating information for a POI
func (r *POIRepo) UpdateRating(ctx context.Context, id uint64, ratingPoint float64, ratingCount int) error {
	return r.GetDB(ctx).
		Model(&domain.POI{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]interface{}{
			"rating_point": ratingPoint,
			"rating_count": ratingCount,
		}).Error
}

// calculateDistance calculates the distance between two coordinates using Haversine formula
func (r *POIRepo) calculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadius = 6371 // Earth's radius in kilometers

	dLat := (lat2 - lat1) * math.Pi / 180
	dLng := (lng2 - lng1) * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLng/2)*math.Sin(dLng/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

// BeforeSave generates code if not provided
func (r *POIRepo) BeforeSave(ctx context.Context, id *uint64, entity *domain.POI) error {
	if entity.Code == "" {
		// Generate POI code
		var count int64
		r.GetDB(ctx).Model(&domain.POI{}).Count(&count)
		entity.Code = fmt.Sprintf("POI%06d", count+1)
	}
	return nil
}
