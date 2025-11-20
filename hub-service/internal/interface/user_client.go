package iprovider

import (
	"context"
	userpb "pb/types/user"
)

// @bind: hub/infra/client.UserClient
type IUserClient interface {
	ListAllUsers(ctx context.Context, page, size uint32) (*userpb.ListAllUsersResponse, error)
}
