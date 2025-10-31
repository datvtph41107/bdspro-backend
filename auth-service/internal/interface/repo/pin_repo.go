package repo

import (
	"auth/internal/domain"
	"context"
	"time"
)

// PINRepository interface cho việc tương tác với dữ liệu PIN
type PINRepository interface {
	// PIN methods
	GetByAuthID(ctx context.Context, authID uint64) (*domain.UserPINEntity, error)
	Create(ctx context.Context, pin *domain.UserPINEntity) error
	Update(ctx context.Context, pin *domain.UserPINEntity) error
	Delete(ctx context.Context, authID uint64) error

	// PIN utilities
	CheckPINExists(ctx context.Context, authID uint64) (bool, error)
	ResetCheckCounter(ctx context.Context, authID uint64) error
	IncrementCheckTime(ctx context.Context, authID uint64) error
	LockPIN(ctx context.Context, authID uint64, lockedUntil time.Time) error
	UnlockPIN(ctx context.Context, authID uint64) error
}
