package admin

import (
	"context"
	"errors"
	"strings"
)

var (
	ErrPermissionDenied               = errors.New("catalog permission denied")
	ErrPermissionAuthorityUnavailable = errors.New("catalog permission authority unavailable")
)

// RequirePermission combines trusted actor identity with the current durable IAM assignment.
// It never infers authority from a role name and fails closed if the authority is unavailable.
func RequirePermission(
	ctx context.Context,
	authorizer PermissionAuthorizer,
	permissionCode string,
) (uint64, error) {
	actorID, err := ActorIDFromContext(ctx)
	if err != nil {
		return 0, err
	}
	permissionCode = strings.TrimSpace(permissionCode)
	if authorizer == nil || permissionCode == "" {
		return 0, ErrPermissionAuthorityUnavailable
	}
	allowed, err := authorizer.HasPermission(ctx, actorID, permissionCode)
	if err != nil {
		return 0, errors.Join(ErrPermissionAuthorityUnavailable, err)
	}
	if !allowed {
		return 0, ErrPermissionDenied
	}
	return actorID, nil
}
