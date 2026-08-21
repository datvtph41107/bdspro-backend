package identitycontext

import (
	"common/identity"
	"context"
)

const (
	AuthIDMetadataKey         = "authid"
	ProfileIDMetadataKey      = "profileid"
	OriginIDMetadataKey       = "originid"
	OrganizationIDMetadataKey = "organizationid"
	SessionIDMetadataKey      = "session"
	RoleMetadataKey           = "role"
	TokenTypeMetadataKey      = "type"
	MaxActorTextLength        = identity.MaxActorTextLength
)

// Deprecated: use identity.Actor. This alias remains while callers migrate.
type Actor = identity.Actor

var (
	ErrInvalidActor         = identity.ErrInvalidActor
	ErrActorContextConflict = identity.ErrActorContextConflict
)

// Deprecated: use identity.BindActor.
func BindActor(ctx context.Context, actor Actor) (context.Context, error) {
	return identity.BindActor(ctx, actor)
}

// Deprecated: use identity.ActorFromContext.
func FromContext(ctx context.Context) (Actor, bool) {
	return identity.ActorFromContext(ctx)
}
