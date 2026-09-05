package client

import (
	"context"
	userpb "pb/types/user"
)

// @bind: hub/internal/interface.IUserClient
type UserClient struct {
	AdminClient userpb.AdminUserProfileServiceClient
}

func NewUserClient(
	adminClient userpb.AdminUserProfileServiceClient,
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
