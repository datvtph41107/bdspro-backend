package useradmin

import (
	"context"
	"errors"
	"strings"

	"common/identity"
)

const (
	PermissionView   = "USER_ADMIN_VIEW"
	PermissionManage = "USER_ADMIN_MANAGE"
)

var (
	ErrActorRequired                  = errors.New("trusted user admin actor identity is required")
	ErrPermissionDenied               = errors.New("user admin permission denied")
	ErrPermissionAuthorityUnavailable = errors.New("user admin permission authority unavailable")
)

// PermissionAuthorizer is the User Admin consumer's narrow IAM boundary.
// Role names and UI visibility are never treated as durable authority.
type PermissionAuthorizer interface {
	HasPermission(ctx context.Context, actorID uint64, permissionCode string) (bool, error)
}

// RequirePermission binds three independent facts: a verified service hop, a
// verified human access-token actor, and the actor's current IAM assignment.
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
