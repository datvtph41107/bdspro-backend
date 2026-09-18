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

type AuthMethodPostgres struct {
	*_db.TransactionRepo
}

// NewAuthRepository tạo mới AuthRepository
func NewAuthRepository(db *_db.TransactionRepo) repo.AuthMethodRepository {
	return &AuthMethodPostgres{TransactionRepo: db}
}

// FindByPhone tìm kiếm auth bằng số điện thoại
func (r *AuthMethodPostgres) FindByPhone(ctx context.Context, phone string) (*auth.AuthMethod, error) {
	var auth auth.AuthMethod
	err := r.GetDB(ctx).Where("provider = ? AND auth_name = ? AND deleted_at IS NULL", "PHONE", phone).First(&auth).Error
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

// GetByID lấy auth theo ID
func (r *AuthMethodPostgres) GetByID(ctx context.Context, id uint64) (*auth.AuthMethod, error) {
	var auth auth.AuthMethod
	err := r.GetDB(ctx).First(&auth, id).Error
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

// ExistsByOAuth kiểm tra sự tồn tại của auth theo OAuth
func (r *AuthMethodPostgres) ExistsByOAuth(ctx context.Context, phone string) (bool, error) {
	var count int64
	err := r.GetDB(ctx).Model(&auth.AuthMethod{}).Where("provider = ? AND auth_name = ? AND deleted_at IS NULL", "PHONE", phone).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindByOAuthId tìm kiếm bằng provider và oauth_id (auth_name cho OAuth providers)
func (r *AuthMethodPostgres) FindByOAuthID(ctx context.Context, oauthId, provider string) (*auth.AuthMethod, error) {
	var user auth.AuthMethod
	err := r.GetDB(ctx).Where("provider = ? AND auth_name = ? AND deleted_at IS NULL", provider, oauthId).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create tạo mới auth
func (r *AuthMethodPostgres) Create(ctx context.Context, auth *auth.AuthMethod) (*auth.AuthMethod, error) {
	err := r.GetDB(ctx).Create(auth).Error
	if err != nil {
		return nil, err
	}
	return auth, nil
}

// Update cập nhật auth
func (r *AuthMethodPostgres) Update(ctx context.Context, auth *auth.AuthMethod) (*auth.AuthMethod, error) {
	err := r.GetDB(ctx).Save(auth).Error
	if err != nil {
		return nil, err
	}
	return auth, nil
}

// DeleteProfile xóa profile (deprecated - sử dụng DeleteByUserId)
func (r *AuthMethodPostgres) DeleteProfile(ctx context.Context, profileID uint64) error {
	return r.DeleteByUserId(ctx, profileID)
}

// DeleteByUserId xóa tất cả auth_method của 1 user (hard delete)
func (r *AuthMethodPostgres) DeleteByUserId(ctx context.Context, userId uint64) error {
	return r.GetDB(ctx).
		Unscoped().
		Where("user_id = ?", userId).
		Delete(&auth.AuthMethod{}).Error
}

// FindByUsername tìm kiếm auth bằng username
func (r *AuthMethodPostgres) FindByUsername(ctx context.Context, username string) (*auth.AuthMethod, error) {
	var auth auth.AuthMethod
	err := r.GetDB(ctx).Where("auth_name = ? AND deleted_at IS NULL", username).First(&auth).Error
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

// FindByEmail tìm auth method theo email
func (r *AuthMethodPostgres) FindByEmail(ctx context.Context, email string) (*auth.AuthMethod, error) {
	var auth auth.AuthMethod
	err := r.GetDB(ctx).Where("email = ? AND deleted_at IS NULL", email).First(&auth).Error
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

// FindByID tìm auth method theo ID
func (r *AuthMethodPostgres) FindByID(ctx context.Context, id uint64) (*auth.AuthMethod, error) {
	var auth auth.AuthMethod
	err := r.GetDB(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&auth).Error
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

// FindAdminsByProvider tìm danh sách admin theo provider với phân trang
func (r *AuthMethodPostgres) FindAdminsByProvider(ctx context.Context, provider string, page, size int, roleId *uint64, name string) ([]*auth.AuthMethod, int64, error) {
	var admins []*auth.AuthMethod
	var total int64

	// Đếm tổng số
	db := r.GetDB(ctx).WithContext(ctx).Model(&auth.AuthMethod{}).
		Where("provider = ? AND deleted_at IS NULL", provider)

	if roleId != nil && *roleId > 0 {
		db = db.Where("role_key = ?", *roleId)
	}
	if name != "" {
		like := "%" + name + "%"
		db = db.Joins("LEFT JOIN admin_profiles ap ON (ap.auth_id = auth_method.id OR ap.id = auth_method.user_id) AND ap.deleted_at IS NULL").
			Where("(auth_method.full_name ILIKE ? OR auth_method.auth_name ILIKE ? OR ap.full_name ILIKE ?)", like, like, like)
	}

	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Lấy danh sách với phân trang
	offset := (page - 1) * size
	listQuery := r.GetDB(ctx).WithContext(ctx).
		Model(&auth.AuthMethod{}).
		Where("provider = ? AND deleted_at IS NULL", provider)

	if roleId != nil && *roleId > 0 {
		listQuery = listQuery.Where("role_key = ?", *roleId)
	}
	if name != "" {
		like := "%" + name + "%"
		listQuery = listQuery.Joins("LEFT JOIN admin_profiles ap ON (ap.auth_id = auth_method.id OR ap.id = auth_method.user_id) AND ap.deleted_at IS NULL").
			Where("(auth_method.full_name ILIKE ? OR auth_method.auth_name ILIKE ? OR ap.full_name ILIKE ?)", like, like, like)
	}

	err = listQuery.Order("auth_method.created_at DESC").
		Offset(offset).
		Limit(size).
		Find(&admins).Error
	if err != nil {
		return nil, 0, err
	}

	return admins, total, nil
}

// GetByUserIdAndProvider tìm auth method theo userID và provider
func (r *AuthMethodPostgres) GetByUserIdAndProvider(ctx context.Context, userID uint64, provider string) (*auth.AuthMethod, error) {
	var auth auth.AuthMethod
	err := r.GetDB(ctx).WithContext(ctx).
		Where("user_id = ? AND provider = ? AND deleted_at IS NULL", userID, provider).
		First(&auth).Error
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

// GetFirstByUserID lấy auth method đầu tiên (ưu tiên theo thời gian tạo) của user
func (r *AuthMethodPostgres) GetFirstByUserID(ctx context.Context, userID uint64) (*auth.AuthMethod, error) {
	var auth auth.AuthMethod
	err := r.GetDB(ctx).WithContext(ctx).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("created_at ASC").
		First(&auth).Error
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

// FindByAuthNameAndProvider tìm auth method theo auth_name và provider
func (r *AuthMethodPostgres) FindByAuthNameAndProvider(ctx context.Context, username, provider string) (*auth.AuthMethod, error) {
	var auth auth.AuthMethod
	err := r.GetDB(ctx).WithContext(ctx).
		Where("auth_name = ? AND provider = ? AND deleted_at IS NULL", username, provider).
		Order("updated_at DESC").
		First(&auth).Error
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

// FindAllByAuthNameAndProvider trả về mọi auth_method trùng auth_name + provider (xử lý dữ liệu legacy trùng username)
func (r *AuthMethodPostgres) FindAllByAuthNameAndProvider(ctx context.Context, username, provider string) ([]*auth.AuthMethod, error) {
	var auths []*auth.AuthMethod
	err := r.GetDB(ctx).WithContext(ctx).
		Where("auth_name = ? AND provider = ? AND deleted_at IS NULL", username, provider).
		Order("updated_at DESC").
		Find(&auths).Error
	if err != nil {
		return nil, err
	}
	return auths, nil
}

// FindAllByEmailAndProvider trả về mọi auth_method trùng email + provider
func (r *AuthMethodPostgres) FindAllByEmailAndProvider(ctx context.Context, email, provider string) ([]*auth.AuthMethod, error) {
	var auths []*auth.AuthMethod
	err := r.GetDB(ctx).WithContext(ctx).
		Where("email = ? AND provider = ? AND deleted_at IS NULL", email, provider).
		Order("updated_at DESC").
		Find(&auths).Error
	if err != nil {
		return nil, err
	}
	return auths, nil
}

// CountByProvider đếm số lượng auth method theo provider
func (r *AuthMethodPostgres) CountByProvider(ctx context.Context, provider string) (int64, error) {
	var count int64
	err := r.GetDB(ctx).WithContext(ctx).Model(&auth.AuthMethod{}).
		Where("provider = ? AND deleted_at IS NULL", provider).
		Count(&count).Error
	return count, err
}

func (r *AuthMethodPostgres) Delete(ctx context.Context, id uint64) error {
	return r.GetDB(ctx).WithContext(ctx).
		Unscoped().
		Where("id = ?", id).
		Delete(&auth.AuthMethod{}).Error
}

// SoftDelete soft delete auth method
func (r *AuthMethodPostgres) SoftDelete(ctx context.Context, id uint64) error {
	now := time.Now()
	return r.GetDB(ctx).WithContext(ctx).Model(&auth.AuthMethod{}).Where("id = ?", id).Update("deleted_at", now).Error
}

// RestoreAccount khôi phục tài khoản đã bị soft delete
func (r *AuthMethodPostgres) RestoreAccount(c context.Context, id uint64) error {
	return r.GetDB(c).WithContext(c).Model(&auth.AuthMethod{}).
		Where("id = ?", id).
		Update("deleted_at", nil).Error
}

// FindDeletedByPhone tìm auth method đã bị xóa theo phone
func (r *AuthMethodPostgres) FindDeletedByPhone(c context.Context, phone string) (*auth.AuthMethod, error) {
	var auth auth.AuthMethod
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
func (r *AuthMethodPostgres) FindDeletedByEmail(c context.Context, email string) (*auth.AuthMethod, error) {
	var auth auth.AuthMethod
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
func (r *AuthMethodPostgres) FindDeletedByUsername(c context.Context, username string) (*auth.AuthMethod, error) {
	var auth auth.AuthMethod
	err := r.GetDB(c).
		WithContext(c).
		Where("auth_name = ? AND deleted_at IS NOT NULL", username).
		First(&auth).
		Error
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

// GetConnectedOAuthAccounts lấy danh sách các tài khoản OAuth đã liên kết của user
func (r *AuthMethodPostgres) GetConnectedOAuthAccounts(ctx context.Context, userID uint64) ([]*auth.AuthMethod, error) {
	var accounts []*auth.AuthMethod

	// Lấy tất cả auth_method của user với provider là OAuth (GOOGLE, FACEBOOK, ZALO)
	err := r.GetDB(ctx).
		Where("user_id = ? AND provider IN (?) AND deleted_at IS NULL",
			userID, []string{"GOOGLE", "FACEBOOK", "ZALO"}).
		Order("created_at DESC").
		Find(&accounts).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get connected OAuth accounts: %w", err)
	}

	return accounts, nil
}

func (r *AuthMethodPostgres) PhoneCheck(ctx context.Context, phone string) (*auth.AuthMethod, error) {
	var auth *auth.AuthMethod
	err := r.GetDB(ctx).Debug().Where("auth_name = ? AND deleted_at IS NULL AND provider = ?", phone, "PHONE").First(&auth).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return auth, nil
}
