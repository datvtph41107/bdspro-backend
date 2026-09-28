package admin

import (
	"context"
	"errors"
	"strings"

	_enum "common/domain/enum"
	"common/identity"
)

var ErrUnauthorized = errors.New("trusted catalog actor identity is required")

// Catalog plan permissions are owned by the shared IAM enum. Auth PermissionService
// evaluates them; User never treats role names as authority.
var (
	PermissionPlanView    = _enum.CATALOG_PLAN_VIEW
	PermissionPlanManage  = _enum.CATALOG_PLAN_MANAGE
	PermissionPlanPublish = _enum.CATALOG_PLAN_PUBLISH
)

// PermissionAuthorizer evaluates the current IAM assignment through Auth.
type PermissionAuthorizer interface {
	HasPermission(ctx context.Context, actorID uint64, permissionID uint32) (bool, error)
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
