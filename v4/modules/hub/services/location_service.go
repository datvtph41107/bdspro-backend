package services

import (
	"context"

	hubpb "pb/types/hub"

	"userv4/modules/hub/dto"
)

// LocationService owns v1 location flows for modules/hub.
type LocationService interface {
	GetProvinces(ctx context.Context, req dto.GetProvincesRequest) (*hubpb.GetProvincesResponse, error)
	GetDistricts(ctx context.Context, req dto.GetDistrictsRequest) (*hubpb.GetDistrictsResponse, error)
	GetWards(ctx context.Context, req dto.GetWardsRequest) (*hubpb.GetWardsResponse, error)
	SearchLocation(ctx context.Context, req dto.SearchLocationRequest) (*hubpb.SearchLocationResponse, error)
}
