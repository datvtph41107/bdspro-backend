package admin

import (
	"context"
	"errors"
	"strings"

	"common/identity"
)

var ErrUnauthorized = errors.New("trusted catalog actor identity is required")

const (
	PermissionPlanView    = "CATALOG_PLAN_VIEW"
	PermissionPlanManage  = "CATALOG_PLAN_MANAGE"
	PermissionPlanPublish = "CATALOG_PLAN_PUBLISH"
)

// PermissionAuthorizer evaluates the current IAM assignment for one actor.
// The Catalog consumer owns this semantic boundary; role names are not authority.
type PermissionAuthorizer interface {
	HasPermission(ctx context.Context, actorID uint64, permissionCode string) (bool, error)
}

/** ActorIDFromContext returns only a verified user access-token identity. */
func ActorIDFromContext(ctx context.Context) (uint64, error) {
	if _, ok := identity.ServiceCallerFromContext(ctx); !ok {
		return 0, ErrUnauthorized
	}
	caller, ok := identity.CallerFromContext(ctx)
	if !ok || caller.Kind != identity.CallerUser {
		return 0, ErrUnauthorized
	}
	actor, ok := identity.ActorFromContext(ctx)
	if !ok || actor.ProfileID == 0 || !strings.EqualFold(actor.TokenType, "ACCESS") {
		return 0, ErrUnauthorized
	}
	return actor.ProfileID, nil
}
