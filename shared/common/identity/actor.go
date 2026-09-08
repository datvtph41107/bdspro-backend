package identity

import (
	"context"
	"errors"
	"strings"
)

const MaxActorTextLength = 128

// Actor is the authenticated human identity carried by a request.
type Actor struct {
	// AuthID nhận diện phương thức/tài khoản đăng nhập; không dùng làm owner
	// profile, subscription hay quota.
	AuthID uint64
	// ProfileID là human actor canonical dùng cho ownership và audit.
	ProfileID uint64
	// OriginID là identity tích hợp/CRM legacy; không dùng làm billing subject.
	OriginID uint64
	// OrganizationID là organization context đã được token xác thực. Việc dùng
	// organization làm billing/quota subject vẫn cần role evidence tương ứng.
	OrganizationID *uint64
	// SessionID nhận diện phiên đăng nhập; không phải business owner.
	SessionID uint64
	// Role là evidence đã xác thực tại thời điểm phát token, không phải permission
	// catalog và không được client tự gửi để nâng quyền.
	Role      string
	TokenType string
}

type actorContextKey struct{}

var (
	ErrInvalidActor         = errors.New("actor identity is invalid")
	ErrActorContextConflict = errors.New("actor identity conflicts with the existing context")
)

// IsZero reports whether the actor contains no identity assertion.
func (a Actor) IsZero() bool {
	return a.AuthID == 0 &&
		a.ProfileID == 0 &&
		a.OriginID == 0 &&
		a.OrganizationID == nil &&
		a.SessionID == 0 &&
		a.Role == "" &&
		a.TokenType == ""
}

// IsValid reports whether the actor is a usable identity assertion.
func (a Actor) IsValid() bool {
	if a.IsZero() {
		return false
	}
	if a.AuthID == 0 && a.ProfileID == 0 && a.OriginID == 0 && a.SessionID == 0 {
		return false
	}
	if a.OrganizationID != nil && *a.OrganizationID == 0 {
		return false
	}
	return validOptionalText(a.Role) && validOptionalText(a.TokenType)
}

// BindActor validates and immutably binds one canonical Actor.
// Rebinding the same value is idempotent; rebinding a different value fails.
func BindActor(ctx context.Context, actor Actor) (context.Context, error) {
	if ctx == nil {
		return nil, errors.New("context is nil")
	}
	if !actor.IsValid() {
		return ctx, ErrInvalidActor
	}

	actor = cloneActor(actor)
	if existing, ok := ActorFromContext(ctx); ok {
		if actorsEqual(existing, actor) {
			return ctx, nil
		}
		return ctx, ErrActorContextConflict
	}

	return context.WithValue(ctx, actorContextKey{}, actor), nil
}

// ActorFromContext returns a defensive copy of the canonical Actor.
func ActorFromContext(ctx context.Context) (Actor, bool) {
	if ctx == nil {
		return Actor{}, false
	}

	actor, ok := ctx.Value(actorContextKey{}).(Actor)
	if !ok || !actor.IsValid() {
		return Actor{}, false
	}

	return cloneActor(actor), true
}

func validOptionalText(value string) bool {
	if value == "" {
		return true
	}
	if strings.TrimSpace(value) != value || len(value) > MaxActorTextLength {
		return false
	}
	for _, character := range value {
		if character < 0x20 || character == 0x7f {
			return false
		}
	}
	return true
}

func cloneActor(actor Actor) Actor {
	if actor.OrganizationID != nil {
		organizationID := *actor.OrganizationID
		actor.OrganizationID = &organizationID
	}
	return actor
}

func actorsEqual(left, right Actor) bool {
	if left.AuthID != right.AuthID ||
		left.ProfileID != right.ProfileID ||
		left.OriginID != right.OriginID ||
		left.SessionID != right.SessionID ||
		left.Role != right.Role ||
		left.TokenType != right.TokenType {
		return false
	}
	if left.OrganizationID == nil || right.OrganizationID == nil {
		return left.OrganizationID == nil && right.OrganizationID == nil
	}
	return *left.OrganizationID == *right.OrganizationID
}
