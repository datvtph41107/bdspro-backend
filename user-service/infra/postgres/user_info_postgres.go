package postgres

import (
	"context"
	"fmt"
	"time"

	_db "common/db"
	"user/internal/domain/auth"
	"user/internal/interface/repo"

	"gorm.io/gorm"
)

type UserInfoPostgres struct {
	*_db.TransactionRepo
}

// NewUserInfoRepository tạo mới UserInfoRepository
func NewUserInfoRepository(db *_db.TransactionRepo) repo.UserInfoRepository {
	return &UserInfoPostgres{TransactionRepo: db}
}

// Create tạo mới UserInfo
func (r *UserInfoPostgres) Create(ctx context.Context, userInfo *auth.AuthUser) (*auth.AuthUser, error) {
	err := r.GetDB(ctx).WithContext(ctx).Create(userInfo).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create user info: %w", err)
	}
	return userInfo, nil
}

// Update cập nhật UserInfo
func (r *UserInfoPostgres) Update(ctx context.Context, userInfo *auth.AuthUser) (*auth.AuthUser, error) {
	err := r.GetDB(ctx).WithContext(ctx).Save(userInfo).Error
	if err != nil {
		return nil, fmt.Errorf("failed to update user info: %w", err)
	}
	return userInfo, nil
}

// GetByProfileID lấy UserInfo theo profileID
func (r *UserInfoPostgres) GetByProfileID(ctx context.Context, profileID uint64) (*auth.AuthUser, error) {
	var userInfo auth.AuthUser
	err := r.GetDB(ctx).WithContext(ctx).Where("profile_id = ?", profileID).First(&userInfo).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	return &userInfo, nil
}

// Delete xóa UserInfo
func (r *UserInfoPostgres) Delete(ctx context.Context, profileID uint64) error {
	err := r.GetDB(ctx).WithContext(ctx).Where("profile_id = ?", profileID).Delete(&auth.AuthUser{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete user info: %w", err)
	}
	return nil
}

// LockUser khóa user (tạm thời hoặc vĩnh viễn)
func (r *UserInfoPostgres) LockUser(ctx context.Context, profileID uint64, lockType auth.AccountStatus, duration int, reason string, lockedBy uint64) error {
	now := time.Now()
	var lockedUntil *time.Time

	// Xác định trạng thái và thời gian khóa
	var status auth.AccountStatus
	if lockType == auth.StatusPermanentlyLocked {
		status = auth.StatusPermanentlyLocked
	} else if lockType == auth.StatusTemporarilyLocked {
		status = auth.StatusTemporarilyLocked
		// Tính thời gian khóa đến
		expiryTime := now.Add(time.Duration(duration) * time.Hour)
		lockedUntil = &expiryTime
	} else {
		return fmt.Errorf("invalid lock type: %d", lockType)
	}

	// Kiểm tra user info có tồn tại không
	var userInfo auth.AuthUser
	err := r.GetDB(ctx).WithContext(ctx).Where("profile_id = ? and deleted_at is null", profileID).First(&userInfo).Error
	if err != nil {
		// Nếu không tồn tại, tạo mới
		userInfo = auth.AuthUser{
			ProfileID:   profileID,
			Status:      status,
			LockedAt:    &now,
			LockedUntil: lockedUntil,
			LockReason:  reason,
			LockedBy:    &lockedBy,
		}
		err = r.GetDB(ctx).WithContext(ctx).Create(&userInfo).Error
	} else {
		// Nếu đã tồn tại, cập nhật
		updates := map[string]interface{}{
			"status":      status,
			"locked_at":   now,
			"lock_reason": reason,
			"locked_by":   lockedBy,
			"updated_at":  now,
		}

		if lockedUntil != nil {
			updates["locked_until"] = lockedUntil
		}

		err = r.GetDB(ctx).Model(&auth.AuthUser{}).
			Where("profile_id = ? and deleted_at is null", profileID).
			Updates(updates).Error
	}

	if err != nil {
		return fmt.Errorf("failed to lock user: %w", err)
	}

	return nil
}

// UnlockUser mở khóa user
func (r *UserInfoPostgres) UnlockUser(ctx context.Context, profileID uint64, reason string, unlockedBy uint64) error {
	now := time.Now()

	// Cập nhật user về trạng thái active
	updates := map[string]interface{}{
		"status":       uint32(auth.StatusActive),
		"locked_at":    nil,
		"locked_until": nil,
		"lock_reason":  "",
		"locked_by":    nil,
		"updated_at":   now,
	}

	err := r.GetDB(ctx).WithContext(ctx).Model(&auth.AuthUser{}).
		Where("profile_id = ? and deleted_at is null", profileID).
		Updates(updates).Error

	if err != nil {
		return fmt.Errorf("failed to unlock user: %w", err)
	}

	return nil
}

// GetUserStatus lấy trạng thái user
func (r *UserInfoPostgres) GetUserStatus(ctx context.Context, profileID uint64) (*auth.AuthUser, error) {
	var userInfo auth.AuthUser
	err := r.GetDB(ctx).WithContext(ctx).Where("profile_id = ? and deleted_at is null", profileID).First(&userInfo).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			newUser := auth.AuthUser{
				ProfileID: profileID,
				Status:    auth.StatusActive,
			}
			err = r.
				GetDB(ctx).
				WithContext(ctx).
				Create(&newUser).Error
			if err != nil {
				return nil, fmt.Errorf("failed to create user info: %w", err)
			}
			return &newUser, nil
		}
		// Nếu không tìm thấy, trả về user với trạng thái active mặc định
		return nil, err
	}
	return &userInfo, nil
}

// CheckLockExpiration kiểm tra và tự động mở khóa các user hết hạn
func (r *UserInfoPostgres) CheckLockExpiration(ctx context.Context) error {
	now := time.Now()

	// Tìm các user khóa tạm thời đã hết hạn
	var expiredUsers []auth.AuthUser
	err := r.GetDB(ctx).WithContext(ctx).Where("status = ? AND locked_until IS NOT NULL AND locked_until < ?",
		uint32(auth.StatusTemporarilyLocked), now).Find(&expiredUsers).Error

	if err != nil {
		return fmt.Errorf("failed to find expired users: %w", err)
	}

	// Mở khóa các user hết hạn
	if len(expiredUsers) > 0 {
		var profileIDs []uint64
		for _, user := range expiredUsers {
			profileIDs = append(profileIDs, user.ProfileID)
		}

		err = r.GetDB(ctx).WithContext(ctx).Model(&auth.AuthUser{}).
			Where("profile_id IN ? and deleted_at is null", profileIDs).
			Updates(map[string]interface{}{
				"status":       auth.StatusActive,
				"locked_at":    nil,
				"locked_until": nil,
				"lock_reason":  "",
				"locked_by":    nil,
				"updated_at":   now,
			}).Error

		if err != nil {
			return fmt.Errorf("failed to unlock expired users: %w", err)
		}
	}

	return nil
}

// GetUsersStatus lấy trạng thái nhiều user
func (r *UserInfoPostgres) GetUsersByProfileIDs(ctx context.Context, profileIDs []uint64) ([]auth.AuthUser, error) {
	var users []auth.AuthUser
	err := r.GetDB(ctx).
		// Model().
		Preload("Role").
		Preload("Role.Color").
		Where("profile_id IN ? and deleted_at is null", profileIDs).
		Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserInfoPostgres) GetRoleIdsByProfileId(ctx context.Context, profileID uint64) ([]uint64, error) {
	roleIDSet := make(map[uint64]struct{})

	var userInfo auth.AuthUser
	err := r.GetDB(ctx).
		Where("profile_id = ? AND deleted_at IS NULL", profileID).
		First(&userInfo).Error
	if err == nil && userInfo.RoleID != nil && *userInfo.RoleID > 0 {
		roleIDSet[*userInfo.RoleID] = struct{}{}
	}

	var profileRoleIDs []uint64
	err = r.GetDB(ctx).
		Table("role_profiles").
		Select("role_id").
		Where("profile_id = ? AND deleted_at IS NULL", profileID).
		Pluck("role_id", &profileRoleIDs).Error
	if err != nil {
		return nil, err
	}
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

// LockMultipleUsers khóa nhiều user cùng lúc
func (r *UserInfoPostgres) LockMultipleUsers(ctx context.Context, profileIDs []uint64, lockType auth.AccountStatus, duration int, reason string, lockedBy uint64) error {
	now := time.Now()
	var lockedUntil *time.Time

	// Xác định trạng thái và thời gian khóa
	var status auth.AccountStatus
	if lockType == auth.StatusPermanentlyLocked {
		status = auth.StatusPermanentlyLocked
	} else if lockType == auth.StatusTemporarilyLocked {
		status = auth.StatusTemporarilyLocked
		expiryTime := now.Add(time.Duration(duration) * time.Hour)
		lockedUntil = &expiryTime
	} else {
		return fmt.Errorf("invalid lock type: %d", lockType)
	}

	// Sử dụng transaction để đảm bảo tính nhất quán
	return r.GetDB(ctx).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, profileID := range profileIDs {
			var userInfo auth.AuthUser
			err := tx.Where("profile_id = ? and deleted_at is null", profileID).First(&userInfo).Error

			if err != nil {
				// Tạo mới nếu không tồn tại
				userInfo = auth.AuthUser{
					ProfileID:   profileID,
					Status:      status,
					LockedAt:    &now,
					LockedUntil: lockedUntil,
					LockReason:  reason,
					LockedBy:    &lockedBy,
				}
				err = tx.Create(&userInfo).Error
			} else {
				// Cập nhật nếu đã tồn tại
				updates := map[string]interface{}{
					"status":      status,
					"locked_at":   now,
					"lock_reason": reason,
					"locked_by":   lockedBy,
					"updated_at":  now,
				}

				if lockedUntil != nil {
					updates["locked_until"] = lockedUntil
				}

				err = tx.Model(&auth.AuthUser{}).
					Where("profile_id = ? and deleted_at is null", profileID).
					Updates(updates).Error
			}

			if err != nil {
				return fmt.Errorf("failed to lock user %d: %w", profileID, err)
			}
		}
		return nil
	})
}

// UnlockMultipleUsers mở khóa nhiều user cùng lúc
func (r *UserInfoPostgres) UnlockMultipleUsers(ctx context.Context, profileIDs []uint64, reason string, unlockedBy uint64) error {
	now := time.Now()

	updates := map[string]interface{}{
		"status":       uint32(auth.StatusActive),
		"locked_at":    nil,
		"locked_until": nil,
		"lock_reason":  "",
		"locked_by":    nil,
		"updated_at":   now,
	}

	err := r.GetDB(ctx).WithContext(ctx).Model(&auth.AuthUser{}).
		Where("profile_id IN ? and deleted_at is null", profileIDs).
		Updates(updates).Error

	if err != nil {
		return fmt.Errorf("failed to unlock users: %w", err)
	}

	return nil
}
