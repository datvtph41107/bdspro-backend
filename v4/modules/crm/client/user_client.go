package client

import (
	"context"

	sharepb "pb/types/shared"
)

// UserClient is a thin adapter used by CRM services/mappers for profile lookups.
type UserClient struct{}

func NewUserClient() *UserClient {
	return &UserClient{}
}

func (c *UserClient) GetProfileByIds(ctx context.Context, ids []uint64) ([]*sharepb.ProfileItem, error) {
	_ = ctx
	_ = ids
	return nil, nil
}

func (c *UserClient) GetMapByIDs(ctx context.Context, profileIDSet map[uint64]struct{}) (map[uint64]*sharepb.ProfileItem, error) {
	_ = ctx
	_ = profileIDSet
	return map[uint64]*sharepb.ProfileItem{}, nil
}
