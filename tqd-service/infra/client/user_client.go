package client

import (
	_utils "common/utils"
	"context"
	"tqd/env"
	"tqd/internal/dto"
	"tqd/internal/interface/provider"

	clients "pb/clients"
)

// UserClient implements UserProvider interface
type UserClient struct {
	*clients.UserGrpcClient
}

// @bind: tqd/internal/interface/provider.UserProvider
func NewUserClient() provider.UserProvider {
	return &UserClient{
		UserGrpcClient: clients.MakeUserGrpcClient(env.ENV_RUNTIME),
	}
}

func (c *UserClient) GetProfileByID(ctx context.Context, userID uint64) (*dto.UserProfileDTO, error) {
	// Return nil for now, implement later when needed
	return nil, nil
}

func (c *UserClient) GetProfilesByIDs(ctx context.Context, userIDs []uint64) (map[uint64]*dto.UserProfileDTO, error) {
	// Return empty map for now, implement later when needed
	return make(map[uint64]*dto.UserProfileDTO), nil
}

func (c *UserClient) GetProfileIdWithContext(ctx context.Context) uint64 {
	return _utils.GetProfileIdWithContext(ctx)
}

func (c *UserClient) GetOrganizationIdFromContext(ctx context.Context) uint64 {
	return _utils.GetOrganizationIdFromContext(ctx)
}

func (c *UserClient) GetUserInfo(ctx context.Context, userID uint64) (*dto.UserInfoDTO, error) {
	// Return nil for now, implement later when needed
	return nil, nil
}

func (c *UserClient) GetUserOrganization(ctx context.Context, userID uint64) (*dto.UserOrganizationDTO, error) {
	// Return nil for now, implement later when needed
	return nil, nil
}
