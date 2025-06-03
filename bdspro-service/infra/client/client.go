package client

import (
	"context"
)

// @bind: bdspro/internal/repo.IClientRepo
type ServiceClient struct{}

func NewServiceClient() *ServiceClient {
	return &ServiceClient{}
}

func (s *ServiceClient) UserInGroup(ctx context.Context, userID uint64, groupID uint64) (bool, error) {
	return true, nil
}
