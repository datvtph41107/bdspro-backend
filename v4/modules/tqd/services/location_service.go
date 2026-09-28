package services

import (
	"context"

	tqdpb "pb/types/tqd"
	"userv4/modules/tqd/dto"
)

// LocationService is the modules/tqd location catalog surface.
type LocationService interface {
	ListProvinces(ctx context.Context, req dto.ListProvincesRequest) (*tqdpb.ListProvincesResponse, error)
	ListWards(ctx context.Context, req dto.ListWardsRequest) (*tqdpb.ClientListWardsResponse, error)
}
