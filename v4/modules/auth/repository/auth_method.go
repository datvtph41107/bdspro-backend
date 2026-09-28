package repository

import (
	"context"

	"userv4/modules/auth/entity"
)

// AuthMethodRepository is the login-facing credential store.
// Implementations may adapt legacy infra/postgres during migration.
type AuthMethodRepository interface {
	FindByPhone(ctx context.Context, phone string) (*entity.AuthMethod, error)
	FindByAuthNameAndProvider(ctx context.Context, authName, provider string) (*entity.AuthMethod, error)
	GetByUserIdAndProvider(ctx context.Context, userID uint64, provider string) (*entity.AuthMethod, error)
}
