package provider

import (
	"context"
	hubpb "pb/types/hub"
)

type HubProvider interface {
	GetLocationsByIds(ctx context.Context, ids []uint64) ([]*hubpb.Location, error)
	GetLocationsByIdsMap(ctx context.Context, ids []uint64) (map[uint64]*hubpb.Location, error)

	GetProductStatsViewClient(ctx context.Context, productID uint64, fromTime int64) (*hubpb.ProductStatsViewResponse, error)
}
