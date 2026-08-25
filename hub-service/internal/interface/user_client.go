package iprovider

import (
	"context"

	sharepb "pb/types/shared"
	userpb "pb/types/user"
)

// @bind: hub/infra/client.UserClient
type IUserClient interface {
	ListAllUsers(ctx context.Context, page, size uint32) (*userpb.ListAllUsersResponse, error)
	GetProfilesByIDs(ctx context.Context, profileIDs []uint64) (map[uint64]*sharepb.ProfileItem, error)
}
