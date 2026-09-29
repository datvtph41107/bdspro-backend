package providers

import (
	"context"
	hubpb "pb/types/hub"
)

type HubProvider interface {
	GetLocationsByIds(ctx context.Context, ids []uint64) (map[uint64]*hubpb.Location, error)
	GetProvinceAndWardNames(ctx context.Context, provinceId, wardId *uint64) (provinceName, wardName string)
}
