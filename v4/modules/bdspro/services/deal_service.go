package services

import (
	"context"

	bdspropb "pb/types/bdspro"

	"userv4/modules/bdspro/dto"
)

// DealService owns deal list/detail flows migrated into modules/bdspro.
type DealService interface {
	GetDealsUser(ctx context.Context, req dto.DealSearchRequest) (*bdspropb.GetDealsResponse, error)
	GetDealsByProduct(ctx context.Context, productID uint64, req dto.SyncByIDRequest) (*bdspropb.GetOrganizationDealsResponse, error)
}
