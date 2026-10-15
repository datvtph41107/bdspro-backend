package repo

import (
	"context"

	_enum "common/domain/enum"
	"userv4/modules/user/domain/user"
)

// ProfileAccountRepository owns post-profile account lock/status on user_profile.
type ProfileAccountRepository interface {
	LockUser(ctx context.Context, profileID uint64, lockType _enum.EUserStatus, duration int, reason string, lockedBy uint64) error
	UnlockUser(ctx context.Context, profileID uint64, reason string, unlockedBy uint64) error
	GetUserStatus(ctx context.Context, profileID uint64) (*user.UserProfileEntity, error)
	CheckLockExpiration(ctx context.Context) error
	GetUsersByProfileIDs(ctx context.Context, profileIDs []uint64) ([]user.UserProfileEntity, error)
	LockMultipleUsers(ctx context.Context, profileIDs []uint64, lockType _enum.EUserStatus, duration int, reason string, lockedBy uint64) error
	UnlockMultipleUsers(ctx context.Context, profileIDs []uint64, reason string, unlockedBy uint64) error
	GetRoleIdsByProfileId(ctx context.Context, profileID uint64) ([]uint64, error)
}
