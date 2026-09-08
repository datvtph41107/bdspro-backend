package repo

import (
	"context"
	"user/internal/domain/auth"
)

// IAuthConfigRepo is the persistence capability required by the ZNS provider.
// Dynamic auth configuration is addressed by stable config keys; broader CRUD/type
// queries are intentionally not part of the serving contract.
type IAuthConfigRepo interface {
	GetByKey(ctx context.Context, key string) (*auth.AuthConfig, error)
	UpdateValue(ctx context.Context, key string, value string) error
}
