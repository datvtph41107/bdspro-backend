package client

import (
	"context"

	hubrpc "hub/infra/rpc"
	userpb "pb/types/user"
)

// @bind: hub/internal/interface.IUserClient
type UserClient struct {
	AdminClient *hubrpc.AdminUserProfileClient
}

func NewUserClient(
	adminClient *hubrpc.AdminUserProfileClient,
) *UserClient {
	return &UserClient{
		AdminClient: adminClient,
	}
}

// ListAllUsers lấy danh sách tất cả user với phân trang
func (c *UserClient) ListAllUsers(ctx context.Context, page, size uint32) (*userpb.ListAllUsersResponse, error) {
	if c.AdminClient == nil {
		return nil, nil
	}

	req := &userpb.ListAllUsersRequest{
		Page:   page,
		Size:   size,
		Status: 10, // Chỉ lấy user active
	}

	return c.AdminClient.ListAllUsers(ctx, req)
}
