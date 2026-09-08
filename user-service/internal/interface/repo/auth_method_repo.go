package repo

import (
	"context"
	"user/internal/domain/auth"
)

type AuthMethodRepository interface {
	// Existing methods
	FindByUsername(ctx context.Context, username string) (*auth.AuthMethod, error)
	FindByPhone(ctx context.Context, phone string) (*auth.AuthMethod, error)
	FindByEmail(ctx context.Context, email string) (*auth.AuthMethod, error)
	FindByOAuthID(ctx context.Context, oauthID, provider string) (*auth.AuthMethod, error)
	Create(ctx context.Context, auth *auth.AuthMethod) (*auth.AuthMethod, error)
	Update(ctx context.Context, auth *auth.AuthMethod) (*auth.AuthMethod, error)
	Delete(ctx context.Context, id uint64) error
	SoftDelete(ctx context.Context, id uint64) error
	DeleteByUserId(ctx context.Context, userId uint64) error // Xóa tất cả auth_method của 1 user
	FindByID(ctx context.Context, id uint64) (*auth.AuthMethod, error)

	// Admin management methods
	FindAdminsByProvider(c context.Context, provider string, page, size int, roleId *uint64, name string) ([]*auth.AuthMethod, int64, error)
	FindByAuthNameAndProvider(c context.Context, authName, provider string) (*auth.AuthMethod, error)
	FindAllByAuthNameAndProvider(c context.Context, authName, provider string) ([]*auth.AuthMethod, error)
	FindAllByEmailAndProvider(c context.Context, email, provider string) ([]*auth.AuthMethod, error)
	GetByUserIdAndProvider(c context.Context, userID uint64, provider string) (*auth.AuthMethod, error)
	GetFirstByUserID(ctx context.Context, userID uint64) (*auth.AuthMethod, error)
	CountByProvider(c context.Context, provider string) (int64, error)

	// Account restoration methods
	RestoreAccount(c context.Context, id uint64) error
	FindDeletedByPhone(c context.Context, phone string) (*auth.AuthMethod, error)
	FindDeletedByEmail(c context.Context, email string) (*auth.AuthMethod, error)
	FindDeletedByUsername(c context.Context, username string) (*auth.AuthMethod, error)

	// Account lock management methods
	LockAccount(c context.Context, authID uint64, lockType string, duration int, reason string, lockedBy uint64) error
	UnlockAccount(c context.Context, authID uint64, reason string, unlockedBy uint64) error
	GetAccountStatus(c context.Context, authID uint64) (*auth.AuthMethod, error)
	CheckLockExpiration(c context.Context) error // Kiểm tra và tự động mở khóa các tài khoản hết hạn

	// OAuth connected accounts
	GetConnectedOAuthAccounts(c context.Context, userID uint64) ([]*auth.AuthMethod, error)
	PhoneCheck(ctx context.Context, phone string) (*auth.AuthMethod, error)
}
