package kernel

import (
	"context"
	"core/kernel/facade"
	"time"

	"github.com/google/uuid"
)

// ===== Auth =====
type authImpl struct {
	userID string
}

func (a authImpl) ID() string  { return a.userID }
func (a authImpl) Guest() bool { return a.userID == "" }

// ===== Clock =====
type systemClock struct{}

func (systemClock) Now() time.Time { return time.Now() }

// ===== UUID =====
type uuidGen struct{}

func (uuidGen) New() string { return uuid.NewString() }

func InjectFacades(ctx context.Context, userID string) context.Context {
	ctx = facade.WithAuth(ctx, authImpl{userID: userID})
	ctx = facade.WithClock(ctx, systemClock{})
	ctx = facade.WithUUID(ctx, uuidGen{})
	return ctx
}
