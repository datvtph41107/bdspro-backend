package repo

import (
	"context"
	"user/internal/domain/auth"
)

// AuthStatusRepository owns registration/credential lifecycle rows in auth_status.
type AuthStatusRepository interface {
	GetByID(ctx context.Context, authID uint64) (*auth.AuthStatus, error)
	CreateStatus(ctx context.Context, status *auth.AuthStatus) error
	UpdateStatus(ctx context.Context, status *auth.AuthStatus) error
}
