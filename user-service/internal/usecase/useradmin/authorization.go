package useradmin

import (
	"context"
	"errors"
	"strings"

	_enum "common/domain/enum"
	"common/identity"
)

// User Admin permissions are owned by the shared IAM enum. Auth/IAM evaluates
// them by durable ID; User never treats role names as authority.
var (
	PermissionView   = _enum.USER_ADMIN_VIEW
	PermissionManage = _enum.USER_ADMIN_MANAGE
)

var (
	ErrActorRequired                  = errors.New("trusted user admin actor identity is required")
	ErrPermissionDenied               = errors.New("user admin permission denied")
	ErrPermissionAuthorityUnavailable = errors.New("user admin permission authority unavailable")
)

// PermissionAuthorizer is the User Admin consumer's narrow IAM boundary.
// Role names and UI visibility are never treated as durable authority.
type PermissionAuthorizer interface {
	HasPermission(ctx context.Context, actorID uint64, permissionID uint32) (bool, error)
}

// RequirePermission binds three independent facts: a verified service hop, a
// verified human access-token actor, and the actor's current IAM assignment.
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

// ActorIDFromContext returns the canonical human owner only when both the
// service hop and the access-token identity were verified by the ingress.
func ActorIDFromContext(ctx context.Context) (uint64, error) {
	if _, ok := identity.ServiceCallerFromContext(ctx); !ok {
		return 0, ErrActorRequired
	}
	caller, ok := identity.CallerFromContext(ctx)
	if !ok || caller.Kind != identity.CallerUser {
		return 0, ErrActorRequired
	}
	actor, ok := identity.ActorFromContext(ctx)
	if !ok || actor.ProfileID == 0 || !strings.EqualFold(actor.TokenType, "ACCESS") {
		return 0, ErrActorRequired
	}
	return actor.ProfileID, nil
}
