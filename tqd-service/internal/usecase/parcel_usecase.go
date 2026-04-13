package usecase

import (
	"context"
	"fmt"

	"tqd/internal/domain"
	"tqd/internal/dto"
	"tqd/internal/interface/repo"
)

// ParcelUsecase defines business logic for POI
type ParcelUsecase interface {
	// CRUD operations
	CreateBatch(ctx context.Context, req *dto.ParcelList) (*dto.ParcelList, error)
	FindByLocation(ctx context.Context, lat, lng float64) (*domain.Parcel, error)
}

// parcelUsecase implements ParcelUsecase
type parcelUsecase struct {
	repo repo.IParcelRepo
}

// NewParcelUsecase creates new POI usecase
func NewParcelUsecase(
	repo repo.IParcelRepo,
) ParcelUsecase {
	return &parcelUsecase{
		repo: repo,
	}
}

// Create implements ParcelUsecase.Create
func (u *parcelUsecase) CreateBatch(ctx context.Context, req *dto.ParcelList) (*dto.ParcelList, error) {
	// Get user ID
	// userID := _utils.GetProfileIdWithContext(ctx)
	// if userID == 0 {
	// 	return nil, fmt.Errorf("unauthorized")
	// }
	for _, parcel := range req.Data {
		if parcel.Lat < -90 || parcel.Lat > 90 {
			return nil, fmt.Errorf("latitude must be between -90 and 90")
		}
		if parcel.Lng < -180 || parcel.Lng > 180 {
			return nil, fmt.Errorf("longitude must be between -180 and 180")
		}
	}

	// Create in repository
	if err := u.repo.CreateBatch(ctx, req.Data); err != nil {
		return nil, fmt.Errorf("failed to create POIs: %v", err)
	}

	return req, nil
}

// GetByID implements ParcelUsecase.GetByID
func (u *parcelUsecase) FindByLocation(ctx context.Context, lat, lng float64) (*domain.Parcel, error) {
	parcel, err := u.repo.FindByLocation(ctx, lat, lng)
	if err != nil {
		return nil, fmt.Errorf("failed to get POI: %v", err)
	}
	if parcel == nil {
		return nil, fmt.Errorf("Parcel not found")
	}

	return parcel, nil
}
