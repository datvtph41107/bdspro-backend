package facade

import (
	"context"
	"core/kernel/contract"
)

func WithAuth(ctx context.Context, a contract.Auth) context.Context {
	return with(ctx, authKey, a)
}

func Auth(ctx context.Context) contract.Auth {
	if v := get(ctx, authKey); v != nil {
		return v.(contract.Auth)
	}
	return guestAuth{}
}

// default guest
type guestAuth struct{}

func (guestAuth) ID() string  { return "" }
func (guestAuth) Guest() bool { return true }
