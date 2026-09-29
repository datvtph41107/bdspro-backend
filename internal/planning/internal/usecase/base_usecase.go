package usecase

import (
	"context"
	"tqd/internal/enums"
	"tqd/internal/interface/provider"
)

// AuthorizationUsecase provides common functionality for all usecases
type AuthorizationUsecase struct {
	userProvider       provider.UserProvider
	permissionProvider provider.PermissionProvider
}

// NewBaseUsecase creates a new base usecase
func NewBaseUsecase(userProvider provider.UserProvider, permissionProvider provider.PermissionProvider) *AuthorizationUsecase {
	return &AuthorizationUsecase{
		userProvider:       userProvider,
		permissionProvider: permissionProvider,
	}
}

// GetCurrentUserID gets current user ID from context
func (b *AuthorizationUsecase) GetCurrentUserID(ctx context.Context) uint64 {
	return b.userProvider.GetProfileIdWithContext(ctx)
}

// GetCurrentOrganizationID gets current organization ID from context
func (b *AuthorizationUsecase) GetCurrentOrganizationID(ctx context.Context) uint64 {
	return b.userProvider.GetOrganizationIdFromContext(ctx)
}

// CheckPermission checks if current user has specific permission
func (b *AuthorizationUsecase) CheckPermission(ctx context.Context, permission enums.Permission) (bool, error) {
	userID := b.GetCurrentUserID(ctx)
	if userID == 0 {
		return false, nil
	}

	return b.permissionProvider.CheckPermission(ctx, userID, permission)
}

// CheckPermissionWithUserID checks if specific user has specific permission
func (b *AuthorizationUsecase) CheckPermissionWithUserID(ctx context.Context, userID uint64, permission enums.Permission) (bool, error) {
	if userID == 0 {
		return false, nil
	}

	return b.permissionProvider.CheckPermission(ctx, userID, permission)
}

// CheckMultiplePermissions checks if current user has multiple permissions
func (b *AuthorizationUsecase) CheckMultiplePermissions(ctx context.Context, permissions []enums.Permission) (map[enums.Permission]bool, error) {
	userID := b.GetCurrentUserID(ctx)
	if userID == 0 {
		result := make(map[enums.Permission]bool)
		for _, perm := range permissions {
			result[perm] = false
		}
		return result, nil
	}

	return b.permissionProvider.CheckMultiplePermissions(ctx, userID, permissions)
}

// UserCanAccessTarget checks if current user can access a specific target
func (b *AuthorizationUsecase) UserCanAccessTarget(ctx context.Context, targetID uint64, targetType string) (bool, error) {
	userID := b.GetCurrentUserID(ctx)
	if userID == 0 {
		return false, nil
	}

	return b.permissionProvider.UserCanAccessTarget(ctx, userID, targetID, targetType)
}

// UserInOwner checks if current user belongs to an owner
func (b *AuthorizationUsecase) UserInOwner(ctx context.Context, ownerID uint64, ownerType string) (bool, error) {
	userID := b.GetCurrentUserID(ctx)
	if userID == 0 {
		return false, nil
	}

	return b.permissionProvider.UserInOwner(ctx, userID, ownerID, ownerType)
}

// HasRoleWithOwner checks if current user has specific role with owner
func (b *AuthorizationUsecase) HasRoleWithOwner(ctx context.Context, ownerID uint64, ownerType string, role string) (bool, error) {
	userID := b.GetCurrentUserID(ctx)
	if userID == 0 {
		return false, nil
	}

	return b.permissionProvider.HasRoleWithOwner(ctx, userID, ownerID, ownerType, role)
}

// IsAdmin checks if current user is admin
func (b *AuthorizationUsecase) IsAdmin(ctx context.Context) (bool, error) {
	userID := b.GetCurrentUserID(ctx)
	if userID == 0 {
		return false, nil
	}

	return b.permissionProvider.IsAdmin(ctx, userID)
}

// IsSuperAdmin checks if current user is super admin
func (b *AuthorizationUsecase) IsSuperAdmin(ctx context.Context) (bool, error) {
	userID := b.GetCurrentUserID(ctx)
	if userID == 0 {
		return false, nil
	}

	return b.permissionProvider.IsSuperAdmin(ctx, userID)
}

// RequirePermission requires current user to have specific permission, returns error if not
func (b *AuthorizationUsecase) RequirePermission(ctx context.Context, permission enums.Permission) error {
	hasPermission, err := b.CheckPermission(ctx, permission)
	if err != nil {
		return err
	}

	if !hasPermission {
		return &PermissionError{
			Permission: permission.String(),
			Message:    "Insufficient permissions",
		}
	}

	return nil
}

// RequirePermissionWithUserID requires specific user to have specific permission
func (b *AuthorizationUsecase) RequirePermissionWithUserID(ctx context.Context, userID uint64, permission enums.Permission) error {
	hasPermission, err := b.CheckPermissionWithUserID(ctx, userID, permission)
	if err != nil {
		return err
	}

	if !hasPermission {
		return &PermissionError{
			Permission: permission.String(),
			Message:    "Insufficient permissions",
		}
	}

	return nil
}

// PermissionError represents a permission error
type PermissionError struct {
	Permission string
	Message    string
}

func (e *PermissionError) Error() string {
	return e.Message + ": " + e.Permission
}
