package repo

import (
	"auth/internal/domain"
	"context"
)

type AuthMethodRepository interface {
	// Existing methods
	FindByUsername(c context.Context, username string) (*domain.AuthMethodDomain, error)
	FindByPhone(c context.Context, phone string) (*domain.AuthMethodDomain, error)
	FindByEmail(c context.Context, email string) (*domain.AuthMethodDomain, error)
	FindByOAuthID(c context.Context, oauthID, provider string) (*domain.AuthMethodDomain, error)
	Create(c context.Context, auth *domain.AuthMethodDomain) (*domain.AuthMethodDomain, error)
	Update(c context.Context, auth *domain.AuthMethodDomain) (*domain.AuthMethodDomain, error)
	Delete(c context.Context, id uint64) error
	SoftDelete(c context.Context, id uint64) error
	FindByID(c context.Context, id uint64) (*domain.AuthMethodDomain, error)

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
