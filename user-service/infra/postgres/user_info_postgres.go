package postgres

import (
	"context"
	"fmt"
	"time"

	_db "common/db"
	_enum "common/domain/enum"
	"user/internal/domain/access"
	"user/internal/domain/user"
	"user/internal/interface/repo"

	"gorm.io/gorm"
)

// ProfileAccountPostgres persists post-profile account lock/status on user_profile.
type ProfileAccountPostgres struct {
	*_db.TransactionRepo
}

func NewUserInfoRepository(db *_db.TransactionRepo) repo.ProfileAccountRepository {
	return &ProfileAccountPostgres{TransactionRepo: db}
}

func (r *ProfileAccountPostgres) LockUser(ctx context.Context, profileID uint64, lockType _enum.EUserStatus, duration int, reason string, lockedBy uint64) error {
	now := time.Now()
	var lockedUntil *time.Time

	switch lockType {
	case _enum.EUserStatusPermanentlyLocked:
	case _enum.EUserStatusTemporaryLocked:
		if duration > 0 {
			expiryTime := now.Add(time.Duration(duration) * time.Hour)
			lockedUntil = &expiryTime
		}
	default:
		return fmt.Errorf("invalid lock type: %d", lockType)
	}

	updates := map[string]interface{}{
		"status":       lockType,
		"locked_at":    now,
		"lock_reason":  reason,
		"locked_by":    lockedBy,
		"locked_until": lockedUntil,
		"updated_at":   now,
	}

	res := r.GetDB(ctx).WithContext(ctx).Model(&user.UserProfileEntity{}).
		Where("profile_id = ? AND deleted_at IS NULL", profileID).
		Updates(updates)
	if res.Error != nil {
		return fmt.Errorf("failed to lock user: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("failed to lock user: profile %d not found", profileID)
	}
	return nil
}

func (r *ProfileAccountPostgres) UnlockUser(ctx context.Context, profileID uint64, reason string, unlockedBy uint64) error {
	now := time.Now()
	_ = reason
	_ = unlockedBy

	updates := map[string]interface{}{
		"status":       _enum.EUserStatusActive,
		"locked_at":    nil,
		"locked_until": nil,
		"lock_reason":  "",
		"locked_by":    nil,
		"updated_at":   now,
	}

	res := r.GetDB(ctx).WithContext(ctx).Model(&user.UserProfileEntity{}).
		Where("profile_id = ? AND deleted_at IS NULL", profileID).
		Updates(updates)
	if res.Error != nil {
		return fmt.Errorf("failed to unlock user: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("failed to unlock user: profile %d not found", profileID)
	}
	return nil
}

func (r *ProfileAccountPostgres) GetUserStatus(ctx context.Context, profileID uint64) (*user.UserProfileEntity, error) {
	var profile user.UserProfileEntity
	err := r.GetDB(ctx).WithContext(ctx).
		Where("profile_id = ? AND deleted_at IS NULL", profileID).
		First(&profile).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get user status: %w", err)
	}
	return &profile, nil
}

func (r *ProfileAccountPostgres) CheckLockExpiration(ctx context.Context) error {
	now := time.Now()

	res := r.GetDB(ctx).WithContext(ctx).Model(&user.UserProfileEntity{}).
		Where("status = ? AND locked_until IS NOT NULL AND locked_until < ? AND deleted_at IS NULL",
			_enum.EUserStatusTemporaryLocked, now).
		Updates(map[string]interface{}{
			"status":       _enum.EUserStatusActive,
			"locked_at":    nil,
			"locked_until": nil,
			"lock_reason":  "",
			"locked_by":    nil,
			"updated_at":   now,
		})
	if res.Error != nil {
		return fmt.Errorf("failed to unlock expired users: %w", res.Error)
	}
	return nil
}

func (r *ProfileAccountPostgres) GetUsersByProfileIDs(ctx context.Context, profileIDs []uint64) ([]user.UserProfileEntity, error) {
	if len(profileIDs) == 0 {
		return nil, nil
	}

	var profiles []user.UserProfileEntity
	err := r.GetDB(ctx).WithContext(ctx).
		Where("profile_id IN ? AND deleted_at IS NULL", profileIDs).
		Find(&profiles).Error
	if err != nil {
		return nil, err
	}

	roleIDs := make([]uint64, 0)
	roleIDSet := make(map[uint64]struct{})
	for _, p := range profiles {
		if p.RoleID != nil && *p.RoleID > 0 {
			if _, ok := roleIDSet[*p.RoleID]; !ok {
				roleIDSet[*p.RoleID] = struct{}{}
				roleIDs = append(roleIDs, *p.RoleID)
			}
		}
	}

	rolesByID := map[uint64]*access.Role{}
	if len(roleIDs) > 0 {
		var roles []access.Role
		if err := r.GetDB(ctx).WithContext(ctx).
			Preload("Color").
			Where("id IN ? AND deleted_at IS NULL", roleIDs).
			Find(&roles).Error; err != nil {
			return nil, err
		}
		for i := range roles {
			rolesByID[roles[i].ID] = &roles[i]
		}
	}

	for i := range profiles {
		if profiles[i].RoleID == nil {
			continue
		}
		role := rolesByID[*profiles[i].RoleID]
		if role == nil {
			continue
		}
		profiles[i].RoleName = role.RoleName
		profiles[i].RoleKeyData = role.RoleKey
		if role.Color != nil {
			profiles[i].RoleColor = role.Color.ContentColor
			profiles[i].RoleBgColor = role.Color.BackgroundColor
		}
	}
	return profiles, nil
}

func (r *ProfileAccountPostgres) GetRoleIdsByProfileId(ctx context.Context, profileID uint64) ([]uint64, error) {
	var profileRoleIDs []uint64
	err := r.GetDB(ctx).WithContext(ctx).
		Table("role_profiles").
		Select("role_id").
		Where("profile_id = ? AND deleted_at IS NULL", profileID).
		Pluck("role_id", &profileRoleIDs).Error
	if err != nil {
		return nil, err
	}

	roleIDSet := make(map[uint64]struct{})
	for _, roleID := range profileRoleIDs {
		if roleID > 0 {
			roleIDSet[roleID] = struct{}{}
		}
	}

	roleIDs := make([]uint64, 0, len(roleIDSet))
	for roleID := range roleIDSet {
		roleIDs = append(roleIDs, roleID)
	}
	return roleIDs, nil
}

func (r *ProfileAccountPostgres) LockMultipleUsers(ctx context.Context, profileIDs []uint64, lockType _enum.EUserStatus, duration int, reason string, lockedBy uint64) error {
	now := time.Now()
	var lockedUntil *time.Time

	switch lockType {
	case _enum.EUserStatusPermanentlyLocked:
	case _enum.EUserStatusTemporaryLocked:
		if duration > 0 {
			expiryTime := now.Add(time.Duration(duration) * time.Hour)
			lockedUntil = &expiryTime
		}
	default:
		return fmt.Errorf("invalid lock type: %d", lockType)
	}

	updates := map[string]interface{}{
		"status":       lockType,
		"locked_at":    now,
		"lock_reason":  reason,
		"locked_by":    lockedBy,
		"locked_until": lockedUntil,
		"updated_at":   now,
	}

	return r.GetDB(ctx).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, profileID := range profileIDs {
			res := tx.Model(&user.UserProfileEntity{}).
				Where("profile_id = ? AND deleted_at IS NULL", profileID).
				Updates(updates)
			if res.Error != nil {
				return fmt.Errorf("failed to lock user %d: %w", profileID, res.Error)
			}
			if res.RowsAffected == 0 {
				return fmt.Errorf("failed to lock user %d: profile not found", profileID)
			}
		}
		return nil
	})
}

func (r *ProfileAccountPostgres) UnlockMultipleUsers(ctx context.Context, profileIDs []uint64, reason string, unlockedBy uint64) error {
	now := time.Now()
	_ = reason
	_ = unlockedBy

	updates := map[string]interface{}{
		"status":       _enum.EUserStatusActive,
		"locked_at":    nil,
		"locked_until": nil,
		"lock_reason":  "",
		"locked_by":    nil,
		"updated_at":   now,
	}

	err := r.GetDB(ctx).WithContext(ctx).Model(&user.UserProfileEntity{}).
		Where("profile_id IN ? AND deleted_at IS NULL", profileIDs).
		Updates(updates).Error
	if err != nil {
		return fmt.Errorf("failed to unlock users: %w", err)
	}
	return nil
}
