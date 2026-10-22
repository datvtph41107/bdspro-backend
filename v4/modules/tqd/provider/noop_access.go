package provider

import (
	"context"

	"mdn4/modules/tqd/dto"
	"mdn4/modules/tqd/enums"
)

// NoopUserProvider satisfies UserProvider without a remote user service.
type NoopUserProvider struct{}

func NewNoopUserProvider() UserProvider { return NoopUserProvider{} }

func (NoopUserProvider) GetProfilesByIDs(context.Context, []uint64) (map[uint64]*dto.UserProfileDTO, error) {
	return map[uint64]*dto.UserProfileDTO{}, nil
}

func (NoopUserProvider) GetProfileIdWithContext(context.Context) uint64 { return 0 }

func (NoopUserProvider) GetOrganizationIdFromContext(context.Context) uint64 { return 0 }

// AllowAllPermissionProvider grants every permission check (HTTP gateway auth is the gate).
type AllowAllPermissionProvider struct{}

func NewAllowAllPermissionProvider() PermissionProvider {
	return AllowAllPermissionProvider{}
}

func (AllowAllPermissionProvider) CheckPermission(context.Context, uint64, enums.Permission) (bool, error) {
	return true, nil
}

func (AllowAllPermissionProvider) CheckMultiplePermissions(_ context.Context, _ uint64, permissions []enums.Permission) (map[enums.Permission]bool, error) {
	out := make(map[enums.Permission]bool, len(permissions))
	for _, p := range permissions {
		out[p] = true
	}
	return out, nil
}

func (AllowAllPermissionProvider) UserCanAccessTarget(context.Context, uint64, uint64, string) (bool, error) {
	return true, nil
}

func (AllowAllPermissionProvider) UserInOwner(context.Context, uint64, uint64, string) (bool, error) {
	return true, nil
}

func (AllowAllPermissionProvider) HasRoleWithOwner(context.Context, uint64, uint64, string, string) (bool, error) {
	return true, nil
}

func (AllowAllPermissionProvider) GetUserPermissions(context.Context, uint64) ([]enums.Permission, error) {
	return nil, nil
}

func (AllowAllPermissionProvider) GetUserRoles(context.Context, uint64) ([]string, error) {
	return nil, nil
}

func (AllowAllPermissionProvider) IsAdmin(context.Context, uint64) (bool, error) {
	return true, nil
}

func (AllowAllPermissionProvider) IsSuperAdmin(context.Context, uint64) (bool, error) {
	return true, nil
}
