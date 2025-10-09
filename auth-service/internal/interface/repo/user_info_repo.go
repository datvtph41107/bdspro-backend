package repo

import (
	"auth/internal/domain"
	"context"
)

type UserInfoRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, userInfo *domain.AuthUser) (*domain.AuthUser, error)
	Update(ctx context.Context, userInfo *domain.AuthUser) (*domain.AuthUser, error)
	GetByProfileID(ctx context.Context, profileID uint64) (*domain.AuthUser, error)
	Delete(ctx context.Context, profileID uint64) error

	// Lock management methods
	LockUser(ctx context.Context, profileID uint64, lockType domain.AccountStatus, duration int, reason string, lockedBy uint64) error
	UnlockUser(ctx context.Context, profileID uint64, reason string, unlockedBy uint64) error
	GetUserStatus(ctx context.Context, profileID uint64) (*domain.AuthUser, error)
	CheckLockExpiration(ctx context.Context) error // Kiểm tra và tự động mở khóa các user hết hạn

	// Batch operations
	GetUsersByProfileIDs(ctx context.Context, profileIDs []uint64) ([]domain.AuthUser, error)
	LockMultipleUsers(ctx context.Context, profileIDs []uint64, lockType domain.AccountStatus, duration int, reason string, lockedBy uint64) error
	UnlockMultipleUsers(ctx context.Context, profileIDs []uint64, reason string, unlockedBy uint64) error
}
