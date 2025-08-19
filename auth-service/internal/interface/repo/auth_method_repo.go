package repo

import (
	"auth/internal/domain"
	"context"
)

// AuthMethodRepository interface cho việc tương tác với dữ liệu authentication
type AuthMethodRepository interface {
	// UserAuth methods
	FindByPhone(ctx context.Context, phone string) (*domain.AuthMethodDomain, error)
	GetByID(ctx context.Context, id uint64) (*domain.AuthMethodDomain, error)
	ExistsByOAuth(ctx context.Context, phone string) (bool, error)
	Create(ctx context.Context, auth *domain.AuthMethodDomain) error
	Update(ctx context.Context, auth *domain.AuthMethodDomain) error
	DeleteProfile(ctx context.Context, profileID uint64) error
	FindByOAuthId(ctx context.Context, oauthId, provider string) (*domain.AuthMethodDomain, error)
}
