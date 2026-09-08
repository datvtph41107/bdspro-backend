package usecase

import (
	"context"
	"fmt"
	"math"

	_dto "common/domain/dto"
	_utils "common/utils"
	"tqd/infra/mapper"
	"tqd/internal/dto"
	"tqd/internal/interface/repo"
)

// PoiUsecase defines business logic for POI
type PoiUsecase interface {
	// CRUD operations
	Create(ctx context.Context, req *dto.CreatePoiRequest) (*dto.PoiResponse, error)
	GetByID(ctx context.Context, id uint64) (*dto.PoiResponse, error)
	GetByCode(ctx context.Context, code string) (*dto.PoiResponse, error)
	Update(ctx context.Context, id uint64, req *dto.UpdatePoiRequest) (*dto.PoiResponse, error)
	Delete(ctx context.Context, id uint64) error

	// List operations
	List(ctx context.Context, filter *dto.PoiFilter) (*dto.ListPoisResponse, error)
	ListByCategory(ctx context.Context, categoryID uint64, limit int) ([]dto.PoiResponse, error)
	ListNearby(ctx context.Context, req *dto.NearbyPoiRequest) (*dto.ListPoisResponse, error)

	// Rating operations
	UpdateRating(ctx context.Context, id uint64, rating float64, reviewCount uint32) (*dto.PoiResponse, error)

	// View/Like operations
	IncrementViewCount(ctx context.Context, id uint64) error
	IncrementLikeCount(ctx context.Context, id uint64) error
	DecrementLikeCount(ctx context.Context, id uint64) error

	// Verification
	VerifyPOI(ctx context.Context, id uint64, verified bool) (*dto.PoiResponse, error)
}

// poiUsecase implements PoiUsecase
type poiUsecase struct {
	repo   repo.IPOIRepo
	mapper *mapper.PoiMapper
}

// NewPoiUsecase creates new POI usecase
func NewPoiUsecase(
	repo repo.IPOIRepo,
	mapper *mapper.PoiMapper,
) PoiUsecase {
	return &poiUsecase{
		repo:   repo,
		mapper: mapper,
	}
}

// Create implements PoiUsecase.Create
func (u *poiUsecase) Create(ctx context.Context, req *dto.CreatePoiRequest) (*dto.PoiResponse, error) {
	// Get user ID
	userID := _utils.GetProfileIdWithContext(ctx)
	if userID == 0 {
		return nil, fmt.Errorf("unauthorized")
	}

	// Validate required fields
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if req.Address == "" {
		return nil, fmt.Errorf("address is required")
	}
	if req.CategoryID == 0 {
		return nil, fmt.Errorf("category ID is required")
	}
	if req.Latitude < -90 || req.Latitude > 90 {
		return nil, fmt.Errorf("latitude must be between -90 and 90")
	}
	if req.Longitude < -180 || req.Longitude > 180 {
		return nil, fmt.Errorf("longitude must be between -180 and 180")
	}

	// Check if category exists
	// This would require category repository injection

	// Convert to domain
	poi, err := u.mapper.ToDomainFromCreate(req, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to map request: %v", err)
	}

	// Create in repository
	if err := u.repo.Create(ctx, poi); err != nil {
		return nil, fmt.Errorf("failed to create POI: %v", err)
	}

	// Get created POI
	created, err := u.repo.GetByID(ctx, poi.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get created POI: %v", err)
	}

	return u.mapper.ToResponse(created), nil
}

// GetByID implements PoiUsecase.GetByID
func (u *poiUsecase) GetByID(ctx context.Context, id uint64) (*dto.PoiResponse, error) {
	if id == 0 {
		return nil, fmt.Errorf("invalid POI ID")
	}

	poi, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get POI: %v", err)
	}
	if poi == nil {
		return nil, fmt.Errorf("POI not found")
	}

	// Increment view count
	u.repo.IncrementViewCount(ctx, id)

	return u.mapper.ToResponse(poi), nil
}

// GetByCode implements PoiUsecase.GetByCode
func (u *poiUsecase) GetByCode(ctx context.Context, code string) (*dto.PoiResponse, error) {
	if code == "" {
		return nil, fmt.Errorf("code is required")
	}

	poi, err := u.repo.GetByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to get POI: %v", err)
	}
	if poi == nil {
		return nil, fmt.Errorf("POI not found")
	}

	return u.mapper.ToResponse(poi), nil
}

// Update implements PoiUsecase.Update
func (u *poiUsecase) Update(ctx context.Context, id uint64, req *dto.UpdatePoiRequest) (*dto.PoiResponse, error) {
	// Get user ID
	userID := _utils.GetProfileIdWithContext(ctx)
	if userID == 0 {
		return nil, fmt.Errorf("unauthorized")
	}

	// Get existing POI
	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get POI: %v", err)
	}
	if existing == nil {
		return nil, fmt.Errorf("POI not found")
	}

	// Convert to domain
	updated, err := u.mapper.ToDomainFromUpdate(req, existing, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to map request: %v", err)
	}

	// Update in repository
	if err := u.repo.Update(ctx, id, updated); err != nil {
		return nil, fmt.Errorf("failed to update POI: %v", err)
	}

	// Get updated POI
	result, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated POI: %v", err)
	}

	return u.mapper.ToResponse(result), nil
}

// Delete implements PoiUsecase.Delete
func (u *poiUsecase) Delete(ctx context.Context, id uint64) error {
	if id == 0 {
		return fmt.Errorf("invalid POI ID")
	}

	// Check if exists
	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get POI: %v", err)
	}
	if existing == nil {
		return fmt.Errorf("POI not found")
	}

	// Delete from repository
	if err := u.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete POI: %v", err)
	}

	return nil
}

// List implements PoiUsecase.List
func (u *poiUsecase) List(ctx context.Context, filter *dto.PoiFilter) (*dto.ListPoisResponse, error) {
	// Validate pagination
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Size < 1 {
		filter.Size = 10
	}
	if filter.Size > 100 {
		filter.Size = 100
	}

	// Get from repository
	pois, total, err := u.repo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list POIs: %v", err)
	}

	return &dto.ListPoisResponse{
		Data:  u.mapper.ToResponseList(pois),
		Total: total,
		Page:  filter.Page,
		Size:  filter.Size,
	}, nil
}

// ListByCategory implements PoiUsecase.ListByCategory
func (u *poiUsecase) ListByCategory(ctx context.Context, categoryID uint64, limit int) ([]dto.PoiResponse, error) {
	if categoryID == 0 {
		return nil, fmt.Errorf("category ID is required")
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	pois, err := u.repo.ListByCategory(ctx, categoryID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list POIs by category: %v", err)
	}

	return u.mapper.ToResponseList(pois), nil
}

// ListNearby implements PoiUsecase.ListNearby
func (u *poiUsecase) ListNearby(ctx context.Context, req *dto.NearbyPoiRequest) (*dto.ListPoisResponse, error) {
	// Validate pagination
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Size < 1 {
		req.Size = 10
	}
	if req.Size > 100 {
		req.Size = 100
	}

	// Create filter
	filter := &dto.PoiFilter{
		Pagable: _dto.Pagable{
			Page: req.Page,
			Size: req.Size,
		},
		CategoryIDs: req.CategoryIDs,
	}

	// Get from repository
	pois, total, err := u.repo.ListNearby(ctx, req.Latitude, req.Longitude, req.Radius, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list nearby POIs: %v", err)
	}

	// Add distance to each POI
	responses := u.mapper.ToResponseList(pois)
	for i := range responses {
		distance := calculateDistance(req.Latitude, req.Longitude,
			responses[i].Latitude, responses[i].Longitude)
		responses[i].Distance = &distance
	}

	return &dto.ListPoisResponse{
		Data:  responses,
		Total: total,
		Page:  req.Page,
		Size:  req.Size,
	}, nil
}

// UpdateRating implements PoiUsecase.UpdateRating
func (u *poiUsecase) UpdateRating(ctx context.Context, id uint64, rating float64, reviewCount uint32) (*dto.PoiResponse, error) {
	// Validate rating
	if rating < 0 || rating > 5 {
		return nil, fmt.Errorf("rating must be between 0 and 5")
	}

	// Check if POI exists
	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get POI: %v", err)
	}
	if existing == nil {
		return nil, fmt.Errorf("POI not found")
	}

	// Update rating
	if err := u.repo.UpdateRating(ctx, id, rating, reviewCount); err != nil {
		return nil, fmt.Errorf("failed to update rating: %v", err)
	}

	// Get updated POI
	updated, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated POI: %v", err)
	}

	return u.mapper.ToResponse(updated), nil
}

// IncrementViewCount implements PoiUsecase.IncrementViewCount
func (u *poiUsecase) IncrementViewCount(ctx context.Context, id uint64) error {
	if id == 0 {
		return fmt.Errorf("invalid POI ID")
	}

	if err := u.repo.IncrementViewCount(ctx, id); err != nil {
		return fmt.Errorf("failed to increment view count: %v", err)
	}

	return nil
}

// IncrementLikeCount implements PoiUsecase.IncrementLikeCount
func (u *poiUsecase) IncrementLikeCount(ctx context.Context, id uint64) error {
	if id == 0 {
		return fmt.Errorf("invalid POI ID")
	}

	if err := u.repo.IncrementLikeCount(ctx, id); err != nil {
		return fmt.Errorf("failed to increment like count: %v", err)
	}

	return nil
}

// DecrementLikeCount implements PoiUsecase.DecrementLikeCount
func (u *poiUsecase) DecrementLikeCount(ctx context.Context, id uint64) error {
	if id == 0 {
		return fmt.Errorf("invalid POI ID")
	}

	if err := u.repo.DecrementLikeCount(ctx, id); err != nil {
		return fmt.Errorf("failed to decrement like count: %v", err)
	}

	return nil
}

// VerifyPOI implements PoiUsecase.VerifyPOI
func (u *poiUsecase) VerifyPOI(ctx context.Context, id uint64, verified bool) (*dto.PoiResponse, error) {
	// Get user ID
	userID := _utils.GetProfileIdWithContext(ctx)
	if userID == 0 {
		return nil, fmt.Errorf("unauthorized")
	}

	// Check if POI exists
	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get POI: %v", err)
	}
	if existing == nil {
		return nil, fmt.Errorf("POI not found")
	}

	// Update verification
	if err := u.repo.VerifyPOI(ctx, id, verified); err != nil {
		return nil, fmt.Errorf("failed to verify POI: %v", err)
	}

	// Get updated POI
	updated, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated POI: %v", err)
	}

	return u.mapper.ToResponse(updated), nil
}

// Helper function
func calculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadius = 6371 // km

	dLat := (lat2 - lat1) * math.Pi / 180
	dLng := (lng2 - lng1) * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLng/2)*math.Sin(dLng/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}
