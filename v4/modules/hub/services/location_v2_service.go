package services

import (
	"context"

	hubpb "pb/types/hub"

	"userv4/modules/hub/dto"
)

// LocationV2Service owns v2 location flows for modules/hub.
type LocationV2Service interface {
	GetProvincesV2(ctx context.Context, req dto.GetProvincesV2Request) (*hubpb.GetProvincesV2Response, error)
	GetProvinceV2ByCode(ctx context.Context, code int32, req dto.GetProvinceV2ByCodeRequest) (*hubpb.GetProvinceV2ByCodeResponse, error)
	GetWardsV2(ctx context.Context, req dto.GetWardsV2Request) (*hubpb.GetWardsV2Response, error)
	GetWardV2ByCode(ctx context.Context, code int32) (*hubpb.GetWardV2ByCodeResponse, error)
	SearchLocationV2(ctx context.Context, req dto.SearchLocationV2Request) (*hubpb.SearchLocationV2Response, error)
}
