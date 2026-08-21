package client

import (
	"common/identity"
	"context"
	"errors"
	"fmt"

	clients "pb/clients"
	"tqd/internal/enums"
	"tqd/internal/interface/provider"
)

var (
	ErrAuthorizationSubjectMissing  = errors.New("authorization actor profile is missing")
	ErrAuthorizationSubjectMismatch = errors.New("authorization subject does not match request actor")
	ErrAuthorizationUnsupported     = errors.New("authorization operation is not implemented by the canonical Auth permission contract")
)

type permissionChecker interface {
	HasPermissions(context.Context, []string) error
}

// PermissionClient is a thin TQD adapter over the Auth authorization authority.
// It never grants access locally: permission decisions come from Auth using the
// verified Actor already carried by the request context.
type PermissionClient struct {
	auth permissionChecker
}

// @bind: tqd/internal/interface/provider.PermissionProvider
func NewPermissionClient(authClient *clients.AuthGrpcClient) provider.PermissionProvider {
	return &PermissionClient{auth: authClient}
}

func (c *PermissionClient) CheckPermission(
	ctx context.Context,
	userID uint64,
	permission enums.Permission,
) (bool, error) {
	if err := validateAuthorizationSubject(ctx, userID); err != nil {
		return false, err
	}
	code := permission.String()
	if code == "unknown" {
		return false, fmt.Errorf("unknown TQD permission value: %d", permission)
	}
	if c == nil || c.auth == nil {
		return false, fmt.Errorf("auth permission client is unavailable")
	}
	if err := c.auth.HasPermissions(ctx, []string{code}); err != nil {
		return false, err
	}
	return true, nil
}

func (c *PermissionClient) CheckMultiplePermissions(
	ctx context.Context,
	userID uint64,
	permissions []enums.Permission,
) (map[enums.Permission]bool, error) {
	if err := validateAuthorizationSubject(ctx, userID); err != nil {
		return nil, err
	}
	result := make(map[enums.Permission]bool, len(permissions))
	for _, permission := range permissions {
		allowed, err := c.CheckPermission(ctx, userID, permission)
		if err != nil {
			return nil, err
		}
		result[permission] = allowed
	}
	return result, nil
}

func validateAuthorizationSubject(ctx context.Context, userID uint64) error {
	actor, ok := identity.ActorFromContext(ctx)
	if !ok || actor.ProfileID == 0 || userID == 0 {
		return ErrAuthorizationSubjectMissing
	}
	if actor.ProfileID != userID {
		return fmt.Errorf(
			"%w: actor_profile=%d requested_profile=%d",
			ErrAuthorizationSubjectMismatch,
			actor.ProfileID,
			userID,
		)
	}
	return nil
}

// The remaining historical PermissionProvider methods have no equivalent in
// the current Auth permission-code contract. They fail closed until a real
// authority/contract is introduced; returning true here would be an
// authorization bypass.
func (c *PermissionClient) UserCanAccessTarget(context.Context, uint64, uint64, string) (bool, error) {
	return false, ErrAuthorizationUnsupported
}

func (c *PermissionClient) UserInOwner(context.Context, uint64, uint64, string) (bool, error) {
	return false, ErrAuthorizationUnsupported
}

func (c *PermissionClient) HasRoleWithOwner(context.Context, uint64, uint64, string, string) (bool, error) {
	return false, ErrAuthorizationUnsupported
}

func (c *PermissionClient) GetUserPermissions(context.Context, uint64) ([]enums.Permission, error) {
	return nil, ErrAuthorizationUnsupported
}

func (c *PermissionClient) GetUserRoles(context.Context, uint64) ([]string, error) {
	return nil, ErrAuthorizationUnsupported
}

func (c *PermissionClient) IsAdmin(context.Context, uint64) (bool, error) {
	return false, ErrAuthorizationUnsupported
}

func (c *PermissionClient) IsSuperAdmin(context.Context, uint64) (bool, error) {
	return false, ErrAuthorizationUnsupported
}
