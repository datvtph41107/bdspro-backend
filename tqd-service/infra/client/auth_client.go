package client

import (
	"context"
	clients "pb/clients"
	"tqd/env"
	"tqd/internal/enums"
	"tqd/internal/interface/provider"
)

// PermissionClient implements PermissionProvider interface
type PermissionClient struct {
	*clients.AuthGrpcClient
}

// @bind: tqd/internal/interface/provider.PermissionProvider
func NewPermissionClient() provider.PermissionProvider {
	return &PermissionClient{
		AuthGrpcClient: clients.MakeAuthGrpcClient(env.ENV_RUNTIME),
	}
}

func (c *PermissionClient) CheckPermission(ctx context.Context, userID uint64, permission enums.Permission) (bool, error) {
	// Return true for now, implement later when needed
	return true, nil
}

func (c *PermissionClient) CheckMultiplePermissions(ctx context.Context, userID uint64, permissions []enums.Permission) (map[enums.Permission]bool, error) {
	// Return all true for now, implement later when needed
	result := make(map[enums.Permission]bool)
	for _, perm := range permissions {
		result[perm] = true
	}
	return result, nil
}

func (c *PermissionClient) UserCanAccessTarget(ctx context.Context, userID uint64, targetID uint64, targetType string) (bool, error) {
	// Return true for now, implement later when needed
	return true, nil
}

func (c *PermissionClient) UserInOwner(ctx context.Context, userID uint64, ownerID uint64, ownerType string) (bool, error) {
	// Return true for now, implement later when needed
	return true, nil
}

func (c *PermissionClient) HasRoleWithOwner(ctx context.Context, userID uint64, ownerID uint64, ownerType string, role string) (bool, error) {
	// Return true for now, implement later when needed
	return true, nil
}

func (c *PermissionClient) GetUserPermissions(ctx context.Context, userID uint64) ([]enums.Permission, error) {
	// Return empty for now, implement later when needed
	return []enums.Permission{}, nil
}

func (c *PermissionClient) GetUserRoles(ctx context.Context, userID uint64) ([]string, error) {
	// Return empty for now, implement later when needed
	return []string{}, nil
}

func (c *PermissionClient) IsAdmin(ctx context.Context, userID uint64) (bool, error) {
	// Return false for now, implement later when needed
	return false, nil
}

func (c *PermissionClient) IsSuperAdmin(ctx context.Context, userID uint64) (bool, error) {
	// Return false for now, implement later when needed
	return false, nil
}
