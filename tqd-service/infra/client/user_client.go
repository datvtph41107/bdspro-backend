package client

import (
	_utils "common/utils"
	"context"
	clients "pb/clients"
	"tqd/env"
	"tqd/internal/dto"
	"tqd/internal/interface/provider"
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
	if userID == 0 || c.Client == nil {
		return nil, nil
	}
	resp, err := c.GetProfileByIds(ctx, []uint64{userID})
	if err != nil || resp == nil || len(resp.Profiles) == 0 {
		return nil, err
	}
	p := resp.Profiles[0]
	return &dto.UserProfileDTO{
		ID:       p.GetId(),
		FullName: p.GetFullName(),
		Avatar:   p.GetAvatar(),
		Phone:    p.GetPhone(),
	}, nil
}

func (c *UserClient) GetProfilesByIDs(ctx context.Context, userIDs []uint64) (map[uint64]*dto.UserProfileDTO, error) {
	out := make(map[uint64]*dto.UserProfileDTO)
	if len(userIDs) == 0 || c.Client == nil {
		return out, nil
	}
	seen := make(map[uint64]struct{}, len(userIDs))
	ids := make([]uint64, 0, len(userIDs))
	for _, id := range userIDs {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return out, nil
	}
	resp, err := c.GetProfileByIds(ctx, ids)
	if err != nil {
		return out, err
	}
	for _, p := range resp.GetProfiles() {
		if p == nil || p.GetId() == 0 {
			continue
		}
		out[p.GetId()] = &dto.UserProfileDTO{
			ID:       p.GetId(),
			FullName: p.GetFullName(),
			Avatar:   p.GetAvatar(),
			Phone:    p.GetPhone(),
		}
	}
	return out, nil
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
