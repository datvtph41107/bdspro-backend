package client

import (
	"context"
	"pb/clients"
	"user/env"
	"user/internal/interface/providers"
)

// @bind: user/internal/interface/providers.AuthProvider
type AuthClient struct {
	*clients.AuthGrpcClient
}

func NewAuthClient() *AuthClient {
	return &AuthClient{
		AuthGrpcClient: clients.MakeAuthGrpcClient(env.ENV_RUNTIME),
	}
}

func (c *AuthClient) GetUserStatus(ctx context.Context, profileID uint64) (*providers.UserStatusInfo, error) {
	// TODO: Implement gọi auth-service để lấy thông tin user status
	// Hiện tại trả về mock data, sau này sẽ implement thực tế

	// Mock data - sẽ được thay thế bằng gọi API thực tế
	return &providers.UserStatusInfo{
		ProfileID:   profileID,
		Verified:    true,
		VerifiedAt:  nil, // Sẽ lấy từ auth-service
		LockedUntil: nil, // Sẽ lấy từ auth-service
		Active:      true,
	}, nil
}

func (c *AuthClient) GetUserStatusBatch(ctx context.Context, profileIDs []uint64) (map[uint64]*providers.UserStatusInfo, error) {
	// TODO: Implement gọi auth-service để lấy thông tin user status cho nhiều user
	// Hiện tại trả về mock data, sau này sẽ implement thực tế

	result := make(map[uint64]*providers.UserStatusInfo)
	for _, profileID := range profileIDs {
		result[profileID] = &providers.UserStatusInfo{
			ProfileID:   profileID,
			Verified:    true,
			VerifiedAt:  nil,
			LockedUntil: nil,
			Active:      true,
		}
	}

	return result, nil
}
