package repo

import (
	"auth/internal/domain"
	"context"
)

type AuthMethodRepository interface {
	// Existing methods
	FindByUsername(ctx context.Context, username string) (*domain.AuthMethodDomain, error)
	FindByPhone(ctx context.Context, phone string) (*domain.AuthMethodDomain, error)
	FindByEmail(ctx context.Context, email string) (*domain.AuthMethodDomain, error)
	FindByOAuthID(ctx context.Context, oauthID, provider string) (*domain.AuthMethodDomain, error)
	Create(ctx context.Context, auth *domain.AuthMethodDomain) (*domain.AuthMethodDomain, error)
	Update(ctx context.Context, auth *domain.AuthMethodDomain) (*domain.AuthMethodDomain, error)
	Delete(ctx context.Context, id uint64) error
	SoftDelete(ctx context.Context, id uint64) error
	FindByID(ctx context.Context, id uint64) (*domain.AuthMethodDomain, error)

	// Admin management methods
	FindAdminsByProvider(c context.Context, provider string, page, size int) ([]*domain.AuthMethodDomain, int64, error)
	FindByUsernameAndProvider(c context.Context, username, provider string) (*domain.AuthMethodDomain, error)
	CountByProvider(c context.Context, provider string) (int64, error)

	// Account restoration methods
	RestoreAccount(c context.Context, id uint64) error
	FindDeletedByPhone(c context.Context, phone string) (*domain.AuthMethodDomain, error)
	FindDeletedByEmail(c context.Context, email string) (*domain.AuthMethodDomain, error)
	FindDeletedByUsername(c context.Context, username string) (*domain.AuthMethodDomain, error)
}
