package admin

import (
	"context"
	"errors"
)

var (
	ErrPermissionDenied               = errors.New("catalog permission denied")
	ErrPermissionAuthorityUnavailable = errors.New("catalog permission authority unavailable")
)

// RequirePermission combines trusted actor identity with Auth PermissionService.
// It never infers authority from a role name and fails closed if Auth is unavailable.
func RequirePermission(
	ctx context.Context,
	authorizer PermissionAuthorizer,
	permissionID uint32,
) (uint64, error) {
	actorID, err := ActorIDFromContext(ctx)
	if err != nil {
		return 0, err
	}
	if authorizer == nil || permissionID == 0 {
		return 0, ErrPermissionAuthorityUnavailable
	}
	allowed, err := authorizer.HasPermission(ctx, actorID, permissionID)
	if err != nil {
		return 0, errors.Join(ErrPermissionAuthorityUnavailable, err)
	}
	if !allowed {
		return 0, ErrPermissionDenied
	}
	return actorID, nil
}
