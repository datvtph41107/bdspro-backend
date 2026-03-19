package validator

import (
	"errors"
	"fmt"
	"strings"

	"tqd/internal/dto"
)

type LocationValidator struct{}

func NewLocationValidator() *LocationValidator {
	return &LocationValidator{}
}

func (v *LocationValidator) ValidateCoordinate(lat, lng float64) error {
	if lat < -90 || lat > 90 {
		return errors.New("latitude must be between -90 and 90")
	}
	if lng < -180 || lng > 180 {
		return errors.New("longitude must be between -180 and 180")
	}
	return nil
}

func (v *LocationValidator) ValidateGetNearestLocationRequest(req *dto.GetNearestLocationRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}

	if err := v.ValidateCoordinate(req.Latitude, req.Longitude); err != nil {
		return err
	}

	if req.MaxDistanceKm != nil && *req.MaxDistanceKm <= 0 {
		return errors.New("maxDistanceKm must be positive")
	}

	if req.Limit != nil && *req.Limit <= 0 {
		return errors.New("limit must be positive")
	}

	return nil
}

func (v *LocationValidator) ValidateBatchGetNearestLocationsRequest(req *dto.BatchGetNearestLocationsRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}

	if len(req.Coordinates) == 0 {
		return errors.New("coordinates cannot be empty")
	}

	if len(req.Coordinates) > 1000 {
		return errors.New("too many coordinates (max 1000)")
	}

	for i, coord := range req.Coordinates {
		if err := v.ValidateCoordinate(coord.Latitude, coord.Longitude); err != nil {
			return fmt.Errorf("coordinate[%d] invalid: %w", i, err)
		}
	}

	return nil
}

func (v *LocationValidator) ValidateSearchLocationsRequest(req *dto.SearchLocationsRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}

	if strings.TrimSpace(req.Query) == "" {
		return errors.New("query cannot be empty")
	}

	if req.Type != nil {
		validTypes := map[string]bool{
			"province": true,
			"district": true,
			"all":      true,
		}

		if !validTypes[*req.Type] {
			return fmt.Errorf("invalid type '%s' (must be 'province', 'district', or 'all')", *req.Type)
		}
	}

	return nil
}

func (v *LocationValidator) ValidatePagination(page, pageSize *int32) error {
	if page != nil && *page < 1 {
		return errors.New("page must be >= 1")
	}

	if pageSize != nil {
		if *pageSize < 1 || *pageSize > 100 {
			return errors.New("pageSize must be between 1 and 100")
		}
	}

	return nil
}
