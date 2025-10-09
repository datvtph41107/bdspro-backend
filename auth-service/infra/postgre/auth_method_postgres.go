package postgre

import (
	"context"
	"fmt"
	"time"

	"auth/internal/domain"
	"auth/internal/interface/repo"
	_db "common/db"
)

type AuthMethodPostgres struct {
	*_db.TransactionRepo
}

// NewAuthRepository tạo mới AuthRepository
func NewAuthRepository(db *_db.TransactionRepo) repo.AuthMethodRepository {
	return &AuthMethodPostgres{TransactionRepo: db}
}

// FindByPhone tìm kiếm auth bằng số điện thoại
func (r *AuthMethodPostgres) FindByPhone(ctx context.Context, phone string) (*domain.AuthMethod, error) {
	var auth domain.AuthMethod
	err := r.GetDB(ctx).Where("provider = ? AND oauth_id = ? AND deleted_at IS NULL", "PHONE", phone).First(&auth).Error
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

// GetByID lấy auth theo ID
func (r *AuthMethodPostgres) GetByID(ctx context.Context, id uint64) (*domain.AuthMethod, error) {
	var auth domain.AuthMethod
	err := r.GetDB(ctx).First(&auth, id).Error
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

// ExistsByOAuth kiểm tra sự tồn tại của auth theo OAuth
func (r *AuthMethodPostgres) ExistsByOAuth(ctx context.Context, phone string) (bool, error) {
	var count int64
	err := r.GetDB(ctx).Model(&domain.AuthMethod{}).Where("provider = ? AND oauth_id = ? AND deleted_at IS NULL", "PHONE", phone).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindByOAuthId tìm kiếm bằng provider và oauth_id
func (r *AuthMethodPostgres) FindByOAuthID(ctx context.Context, oauthId, provider string) (*domain.AuthMethod, error) {
	var user domain.AuthMethod
	err := r.GetDB(ctx).Where("provider = ? AND oauth_id = ? AND deleted_at IS NULL", provider, oauthId).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create tạo mới auth
func (r *AuthMethodPostgres) Create(ctx context.Context, auth *domain.AuthMethod) (*domain.AuthMethod, error) {
	err := r.GetDB(ctx).Create(auth).Error
	if err != nil {
		return nil, err
	}
	return auth, nil
}

// Update cập nhật auth
func (r *AuthMethodPostgres) Update(ctx context.Context, auth *domain.AuthMethod) (*domain.AuthMethod, error) {
	err := r.GetDB(ctx).Save(auth).Error
	if err != nil {
		return nil, err
	}
	return auth, nil
}

// DeleteProfile xóa profile
func (r *AuthMethodPostgres) DeleteProfile(ctx context.Context, profileID uint64) error {
	return r.GetDB(ctx).Model(&domain.AuthMethod{}).
		Where("user_id = ?", profileID).
		Update("deleted_at", time.Now()).Error
}

// FindByUsername tìm kiếm auth bằng username
func (r *AuthMethodPostgres) FindByUsername(ctx context.Context, username string) (*domain.AuthMethod, error) {
	var auth domain.AuthMethod
	err := r.GetDB(ctx).Where("auth_name = ? AND deleted_at IS NULL", username).First(&auth).Error
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

// FindByEmail tìm auth method theo email
func (r *AuthMethodPostgres) FindByEmail(ctx context.Context, email string) (*domain.AuthMethod, error) {
	var auth domain.AuthMethod
	err := r.GetDB(ctx).Where("email = ? AND deleted_at IS NULL", email).First(&auth).Error
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

// FindByID tìm auth method theo ID
func (r *AuthMethodPostgres) FindByID(ctx context.Context, id uint64) (*domain.AuthMethod, error) {
	var auth domain.AuthMethod
	err := r.GetDB(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&auth).Error
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

// FindAdminsByProvider tìm danh sách admin theo provider với phân trang
func (r *AuthMethodPostgres) FindAdminsByProvider(ctx context.Context, provider string, page, size int, roleId *uint64, name string) ([]*domain.AuthMethod, int64, error) {
	var admins []*domain.AuthMethod
	var total int64

	// Đếm tổng số
	db := r.GetDB(ctx).WithContext(ctx).Model(&domain.AuthMethod{}).
		Where("provider = ? AND deleted_at IS NULL", provider)

	if roleId != nil && *roleId > 0 {
		db = db.Where("role_key = ?", *roleId)
	}
	if name != "" {
		like := "%" + name + "%"
		db = db.Where("(full_name ILIKE ? OR auth_name ILIKE ?)", like, like)
	}

	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Lấy danh sách với phân trang
	offset := (page - 1) * size
	listQuery := r.GetDB(ctx).WithContext(ctx).
		Model(&domain.AuthMethod{}).
		Where("provider = ? AND deleted_at IS NULL", provider)

	if roleId != nil && *roleId > 0 {
		listQuery = listQuery.Where("role_key = ?", *roleId)
	}
	if name != "" {
		like := "%" + name + "%"
		listQuery = listQuery.Where("(full_name ILIKE ? OR auth_name ILIKE ?)", like, like)
	}

	err = listQuery.Order("created_at DESC").
		Offset(offset).
		Limit(size).
		Find(&admins).Error
	if err != nil {
		return nil, 0, err
	}

	return admins, total, nil
}

// GetByUserIdAndProvider tìm auth method theo userID và provider
func (r *AuthMethodPostgres) GetByUserIdAndProvider(ctx context.Context, userID uint64, provider string) (*domain.AuthMethod, error) {
	var auth domain.AuthMethod
	err := r.GetDB(ctx).WithContext(ctx).
		Where("user_id = ? AND provider = ? AND deleted_at IS NULL", userID, provider).
		First(&auth).Error
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

// FindByUsernameAndProvider tìm auth method theo username và provider
func (r *AuthMethodPostgres) FindByAuthNameAndProvider(ctx context.Context, username, provider string) (*domain.AuthMethod, error) {
	var auth domain.AuthMethod
	err := r.GetDB(ctx).WithContext(ctx).
		Where("username = ? AND provider = ? AND deleted_at IS NULL", username, provider).
		First(&auth).Error
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

// CountByProvider đếm số lượng auth method theo provider
func (r *AuthMethodPostgres) CountByProvider(ctx context.Context, provider string) (int64, error) {
	var count int64
	err := r.GetDB(ctx).WithContext(ctx).Model(&domain.AuthMethod{}).
		Where("provider = ? AND deleted_at IS NULL", provider).
		Count(&count).Error
	return count, err
}

func (r *AuthMethodPostgres) Delete(ctx context.Context, id uint64) error {
	return r.GetDB(ctx).WithContext(ctx).Model(&domain.AuthMethod{}).
		Where("id = ?", id).
		Delete(&domain.AuthMethod{}).Error
}

// SoftDelete soft delete auth method
func (r *AuthMethodPostgres) SoftDelete(ctx context.Context, id uint64) error {
	now := time.Now()
	return r.GetDB(ctx).WithContext(ctx).Model(&domain.AuthMethod{}).Where("id = ?", id).Update("deleted_at", now).Error
}

// RestoreAccount khôi phục tài khoản đã bị soft delete
func (r *AuthMethodPostgres) RestoreAccount(c context.Context, id uint64) error {
	return r.GetDB(c).WithContext(c).Model(&domain.AuthMethod{}).
		Where("id = ?", id).
		Update("deleted_at", nil).Error
}

// FindDeletedByPhone tìm auth method đã bị xóa theo phone
func (r *AuthMethodPostgres) FindDeletedByPhone(c context.Context, phone string) (*domain.AuthMethod, error) {
	var auth domain.AuthMethod
	err := r.GetDB(c).
		WithContext(c).
		Where("phone = ? AND deleted_at IS NOT NULL", phone).
		First(&auth).Error
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

// FindDeletedByEmail tìm auth method đã bị xóa theo email
func (r *AuthMethodPostgres) FindDeletedByEmail(c context.Context, email string) (*domain.AuthMethod, error) {
	var auth domain.AuthMethod
	err := r.GetDB(c).
		WithContext(c).
		Where("email = ? AND deleted_at IS NOT NULL", email).
		First(&auth).Error
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

// FindDeletedByUsername tìm auth method đã bị xóa theo username
func (r *AuthMethodPostgres) FindDeletedByUsername(c context.Context, username string) (*domain.AuthMethod, error) {
	var auth domain.AuthMethod
	err := r.GetDB(c).
		WithContext(c).
		Where("username = ? AND deleted_at IS NOT NULL", username).
		First(&auth).
		Error
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

// LockAccount khóa tài khoản (tạm thời hoặc vĩnh viễn)
func (r *AuthMethodPostgres) LockAccount(ctx context.Context, authID uint64, lockType string, duration int, reason string, lockedBy uint64) error {
	now := time.Now()
	var lockedUntil *time.Time

	// Xác định trạng thái và thời gian khóa
	var status uint8
	if lockType == "permanent" {
		status = uint8(domain.StatusPermanentlyLocked)
	} else if lockType == "temporary" {
		status = uint8(domain.StatusTemporarilyLocked)
		// Tính thời gian khóa đến
		expiryTime := now.Add(time.Duration(duration) * time.Hour)
		lockedUntil = &expiryTime
	} else {
		return fmt.Errorf("invalid lock type: %s", lockType)
	}

	// Cập nhật tài khoản
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

	err := r.GetDB(ctx).Model(&domain.AuthMethod{}).
		Where("id = ? AND deleted_at IS NULL", authID).
		Updates(updates).Error

	if err != nil {
		return fmt.Errorf("failed to lock account: %w", err)
	}

	return nil
}

// UnlockAccount mở khóa tài khoản
func (r *AuthMethodPostgres) UnlockAccount(ctx context.Context, authID uint64, reason string, unlockedBy uint64) error {
	now := time.Now()

	// Cập nhật tài khoản về trạng thái active
	updates := map[string]interface{}{
		"status":       uint8(domain.StatusActive),
		"locked_at":    nil,
		"locked_until": nil,
		"lock_reason":  "",
		"locked_by":    nil,
		"updated_at":   now,
	}

	err := r.GetDB(ctx).Model(&domain.AuthMethod{}).
		Where("id = ? AND deleted_at IS NULL", authID).
		Updates(updates).Error

	if err != nil {
		return fmt.Errorf("failed to unlock account: %w", err)
	}

	return nil
}

// GetAccountStatus lấy trạng thái tài khoản
func (r *AuthMethodPostgres) GetAccountStatus(ctx context.Context, authID uint64) (*domain.AuthMethod, error) {
	var auth domain.AuthMethod
	err := r.GetDB(ctx).Where("id = ? AND deleted_at IS NULL", authID).First(&auth).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get account status: %w", err)
	}
	return &auth, nil
}

// CheckLockExpiration kiểm tra và tự động mở khóa các tài khoản hết hạn
func (r *AuthMethodPostgres) CheckLockExpiration(ctx context.Context) error {
	now := time.Now()

	// Tìm các tài khoản khóa tạm thời đã hết hạn
	var expiredAccounts []domain.AuthMethod
	err := r.GetDB(ctx).Where("status = ? AND locked_until IS NOT NULL AND locked_until < ?",
		uint8(domain.StatusTemporarilyLocked), now).Find(&expiredAccounts).Error

	if err != nil {
		return fmt.Errorf("failed to find expired accounts: %w", err)
	}

	// Mở khóa các tài khoản hết hạn
	if len(expiredAccounts) > 0 {
		var authIDs []uint64
		for _, account := range expiredAccounts {
			authIDs = append(authIDs, account.ID)
		}

		err = r.GetDB(ctx).Model(&domain.AuthMethod{}).
			Where("id IN ?", authIDs).
			Updates(map[string]interface{}{
				"status":       uint8(domain.StatusActive),
				"locked_at":    nil,
				"locked_until": nil,
				"lock_reason":  "",
				"locked_by":    nil,
				"updated_at":   now,
			}).Error

		if err != nil {
			return fmt.Errorf("failed to unlock expired accounts: %w", err)
		}
	}

	return nil
}
