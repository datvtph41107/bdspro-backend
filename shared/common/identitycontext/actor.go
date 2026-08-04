package identitycontext

import "context"

const (
	AuthIDMetadataKey         = "authid"
	ProfileIDMetadataKey      = "profileid"
	OriginIDMetadataKey       = "originid"
	OrganizationIDMetadataKey = "organizationid"
	SessionIDMetadataKey      = "session"
	RoleMetadataKey           = "role"
	TokenTypeMetadataKey      = "type"
)

/**
 * Actor is the request identity asserted by the transport layer.
 */
type Actor struct {
	AuthID         uint64
	ProfileID      uint64
	OriginID       uint64
	OrganizationID *uint64
	SessionID      uint64
	Role           string
	TokenType      string
}

type actorContextKey struct{}

/**
 * IsZero reports whether the actor contains no identity assertion.
 */
func (a Actor) IsZero() bool {
	return a.AuthID == 0 &&
		a.ProfileID == 0 &&
		a.OriginID == 0 &&
		a.OrganizationID == nil &&
		a.SessionID == 0 &&
		a.Role == "" &&
		a.TokenType == ""
}

/**
 * WithActor stores a defensive actor copy in context.
 */
func WithActor(ctx context.Context, actor Actor) context.Context {
	if ctx == nil || actor.IsZero() {
		return ctx
	}

	actor = cloneActor(actor)
	return context.WithValue(ctx, actorContextKey{}, actor)
}

/**
 * FromContext returns a defensive actor copy.
 */
func FromContext(ctx context.Context) (Actor, bool) {
	if ctx == nil {
		return Actor{}, false
	}

	actor, ok := ctx.Value(actorContextKey{}).(Actor)
	if !ok || actor.IsZero() {
		return Actor{}, false
	}

	return cloneActor(actor), true
}

/**
 * OrganizationID *uint64 con trỏ
 * có thể bị thay đổi bên ngoài
 * khi gán WithActor
 * -> clone -> giả định thay đổi, gốc giữ nguyên
 */

func cloneActor(actor Actor) Actor {
	if actor.OrganizationID != nil {
		organizationID := *actor.OrganizationID
		actor.OrganizationID = &organizationID
	}
	return actor
}
