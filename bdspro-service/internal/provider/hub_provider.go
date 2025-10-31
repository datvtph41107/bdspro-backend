package provider

import (
	"context"
	hubpb "pb/types/hub"
)

type HubProvider interface {
	GetLocationsByIds(ctx context.Context, ids []uint64) ([]*hubpb.Location, error)
	GetLocationsByIdsMap(ctx context.Context, ids []uint64) (map[uint64]*hubpb.Location, error)
}
