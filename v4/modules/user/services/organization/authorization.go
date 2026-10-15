package organization

import (
	"context"
	"errors"
	"strings"

	"common/identity"
)

var ErrUnauthorized = errors.New("trusted organization admin identity is required")

type PermissionAuthorizer interface {
	HasPermission(context.Context, uint64, string) (bool, error)
}

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
