package admin

import (
	"context"
	"errors"
	"strings"

	"common/callercontext"
	"common/identitycontext"
	"common/transporttrust"
)

var ErrUnauthorized = errors.New("trusted catalog actor identity is required")

const (
	PermissionPlanView    = "CATALOG_PLAN_VIEW"
	PermissionPlanPublish = "CATALOG_PLAN_PUBLISH"
)

// PermissionAuthorizer evaluates the current IAM assignment for one actor.
// The Catalog consumer owns this semantic boundary; role names are not authority.
type PermissionAuthorizer interface {
	HasPermission(ctx context.Context, actorID uint64, permissionCode string) (bool, error)
}

/** ActorIDFromContext returns only a verified user access-token identity. */
func ActorIDFromContext(ctx context.Context) (uint64, error) {
	if !transporttrust.IsVerified(ctx) {
		return 0, ErrUnauthorized
	}
	caller, ok := callercontext.FromContext(ctx)
	if !ok || caller.Kind != callercontext.KindUser {
		return 0, ErrUnauthorized
	}
	actor, ok := identitycontext.FromContext(ctx)
	if !ok || actor.ProfileID == 0 || !strings.EqualFold(actor.TokenType, "ACCESS") {
		return 0, ErrUnauthorized
	}
	return actor.ProfileID, nil
}
