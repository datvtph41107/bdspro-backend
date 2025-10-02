package provider

import (
	"context"
	"tqd/internal/enums"
)

// PermissionProvider defines the interface for permission operations
type PermissionProvider interface {
	// CheckPermission checks if user has specific permission
	CheckPermission(ctx context.Context, userID uint64, permission enums.Permission) (bool, error)

	// CheckMultiplePermissions checks if user has multiple permissions
	CheckMultiplePermissions(ctx context.Context, userID uint64, permissions []enums.Permission) (map[enums.Permission]bool, error)

	// UserCanAccessTarget checks if user can access a specific target/resource
	UserCanAccessTarget(ctx context.Context, userID uint64, targetID uint64, targetType string) (bool, error)

	// UserInOwner checks if user belongs to an owner (organization/user)
	UserInOwner(ctx context.Context, userID uint64, ownerID uint64, ownerType string) (bool, error)

	// HasRoleWithOwner checks if user has specific role with owner
	HasRoleWithOwner(ctx context.Context, userID uint64, ownerID uint64, ownerType string, role string) (bool, error)

	// GetUserPermissions gets all permissions for a user
	GetUserPermissions(ctx context.Context, userID uint64) ([]enums.Permission, error)

	// GetUserRoles gets all roles for a user
	GetUserRoles(ctx context.Context, userID uint64) ([]string, error)

	// IsAdmin checks if user is admin
	IsAdmin(ctx context.Context, userID uint64) (bool, error)

	// IsSuperAdmin checks if user is super admin
	IsSuperAdmin(ctx context.Context, userID uint64) (bool, error)
}
