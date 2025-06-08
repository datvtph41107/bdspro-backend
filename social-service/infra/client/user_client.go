package client

import "context"

// @bind: social/internal/interface.UserClient
type UserClient struct {
}

func NewUserClient() *UserClient {
	return &UserClient{}
}

func (c *UserClient) IsFriends(ctx context.Context, profileID uint64, friendIDs []uint64) (bool, error) {
	return true, nil
}
