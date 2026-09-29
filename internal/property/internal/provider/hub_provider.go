package provider

import (
	"context"
	hubpb "pb/types/hub"
)

type HubProvider interface {
	GetLocationsByIds(ctx context.Context, ids []uint64) ([]*hubpb.Location, error)
	GetLocationsByIdsMap(ctx context.Context, ids []uint64) (map[uint64]*hubpb.Location, error)
	GetRegionsByIdsMap(ctx context.Context, ids []uint64) (map[uint64]*hubpb.Region, error)

	GetProductStatsViewClient(ctx context.Context, productID uint64, fromTime int64) (*hubpb.ProductStatsViewResponse, error)
	PutUpdate(ctx context.Context, resource string, resourceId uint64, ownerIds []uint64, updatedAt int64) error
	DelUpdate(ctx context.Context, resource string, id uint64, ownerId uint64) error
}
