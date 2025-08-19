package repo

import (
	"auth/internal/domain"
	"context"
	"time"
)

// StatusRepository interface cho việc tương tác với dữ liệu user status
type StatusRepository interface {
	// Status methods
	GetByID(ctx context.Context, authID uint64) (*domain.UserStatusEntity, error)
	CreateStatus(ctx context.Context, status *domain.UserStatusEntity) error
	UpdateStatus(ctx context.Context, status *domain.UserStatusEntity) error
	DeleteStatus(ctx context.Context, authID uint64) error

	// Status utilities
	GetStatusByUserID(ctx context.Context, userID uint64) (*domain.UserStatusEntity, error)
	GetLockedUsers(ctx context.Context) ([]*domain.UserStatusEntity, error)
	GetInactiveUsers(ctx context.Context) ([]*domain.UserStatusEntity, error)
	UnlockUser(ctx context.Context, authID uint64) error
	ActivateUser(ctx context.Context, authID uint64) error
	DeactivateUser(ctx context.Context, authID uint64) error
	SetLockUntil(ctx context.Context, authID uint64, lockUntil time.Time) error
}
