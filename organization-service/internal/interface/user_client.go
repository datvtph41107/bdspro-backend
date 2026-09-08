package iusecase

import (
	"context"
	sharepb "pb/types/shared"
)

type IUserClient interface {
	GetProfileByIds(ctx context.Context, in *sharepb.GetProfileByIdsRequest) (*sharepb.GetProfileByIdsResponse, error)
	GetMapProfileByIds(ctx context.Context, in *sharepb.GetProfileByIdsRequest) (map[uint64]*sharepb.ProfileItem, error)
}
