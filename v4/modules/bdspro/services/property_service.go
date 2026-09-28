package services

import (
	"context"

	bdspropb "pb/types/bdspro"

	"userv4/modules/bdspro/dto"
)

// PropertyService owns property list/detail flows for modules/bdspro.
type PropertyService interface {
	ListMe(ctx context.Context, req dto.PropertyListRequest) (*bdspropb.PropertyListResponse, error)
	GetDetail(ctx context.Context, req dto.PropertyListRequest) (*bdspropb.PropertyDetailResponse, error)
}
