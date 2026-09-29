package repo

import (
	"context"
	"user/internal/domain/auth"
)

// StatusRepository contains only status state operations consumed by serving usecases.
type StatusRepository interface {
	GetByID(ctx context.Context, authID uint64) (*auth.UserStatusEntity, error)
	CreateStatus(ctx context.Context, status *auth.UserStatusEntity) error
	UpdateStatus(ctx context.Context, status *auth.UserStatusEntity) error
}
