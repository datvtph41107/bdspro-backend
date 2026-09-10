package client

import (
	"context"

	"pb/clients"
	sharepb "pb/types/shared"
	userpb "pb/types/user"
)

// @bind: hub/internal/interface.IUserClient
type UserClient struct {
	*clients.UserGrpcClient
}

func NewUserClient(rpcClient *clients.UserGrpcClient) *UserClient {
	return &UserClient{UserGrpcClient: rpcClient}
}

// ListAllUsers lấy danh sách tất cả user với phân trang
func (c *UserClient) ListAllUsers(ctx context.Context, page, size uint32) (*userpb.ListAllUsersResponse, error) {
	if c == nil || c.AdminClient == nil {
		return nil, nil
	}

	req := &userpb.ListAllUsersRequest{
		Page:   page,
		Size:   size,
		Status: 10, // Chỉ lấy user active
	}

	return c.AdminClient.ListAllUsers(ctx, req)
}

// GetProfilesByIDs lấy map profile theo danh sách ID.
func (c *UserClient) GetProfilesByIDs(ctx context.Context, profileIDs []uint64) (map[uint64]*sharepb.ProfileItem, error) {
	if c == nil || c.UserGrpcClient == nil || len(profileIDs) == 0 {
		return map[uint64]*sharepb.ProfileItem{}, nil
	}

	idSet := make(map[uint64]struct{}, len(profileIDs))
	for _, id := range profileIDs {
		if id > 0 {
			idSet[id] = struct{}{}
		}
	}
	if len(idSet) == 0 {
		return map[uint64]*sharepb.ProfileItem{}, nil
	}

	return c.GetMapByIDs(ctx, idSet)
}
