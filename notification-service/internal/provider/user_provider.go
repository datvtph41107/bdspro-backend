package provider

import (
	"context"
	sharepb "pb/types/shared"
	userpb "pb/types/user"
)

type IUserProvider interface {
	GetProfileByIds(ctx context.Context, in *sharepb.GetProfileByIdsRequest) (*sharepb.GetProfileByIdsResponse, error)
	GetMapProfileByIds(ctx context.Context, in *sharepb.GetProfileByIdsRequest) (map[uint64]*sharepb.ProfileItem, error)
	GetProfileByID(ctx context.Context, in *userpb.GetProfileByIDRequest) (*userpb.ProfileResponse, error)
}
