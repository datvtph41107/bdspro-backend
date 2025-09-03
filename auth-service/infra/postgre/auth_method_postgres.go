package postgre

import (
	"context"
	"time"

	"auth/internal/domain"
	_db "common/db"
)

// AuthMethodPostgres implementation của auth repository
// @bind: auth/internal/interface/repo.AuthMethodRepository
type AuthMethodPostgres struct {
	*_db.TransactionRepo
}

// NewAuthRepository tạo mới AuthRepository
func NewAuthRepository(db *_db.TransactionRepo) *AuthMethodPostgres {
	return &AuthMethodPostgres{TransactionRepo: db}
}

// FindByPhone tìm kiếm auth bằng số điện thoại
func (r *AuthMethodPostgres) FindByPhone(ctx context.Context, phone string) (*domain.AuthMethodDomain, error) {
	var auth domain.AuthMethodDomain
	err := r.GetDB(ctx).Where("provider = ? AND oauth_id = ? AND deleted_at IS NULL", "PHONE", phone).First(&auth).Error
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

// GetByID lấy auth theo ID
func (r *AuthMethodPostgres) GetByID(ctx context.Context, id uint64) (*domain.AuthMethodDomain, error) {
	var auth domain.AuthMethodDomain
	err := r.GetDB(ctx).First(&auth, id).Error
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

// ExistsByOAuth kiểm tra sự tồn tại của auth theo OAuth
func (r *AuthMethodPostgres) ExistsByOAuth(ctx context.Context, phone string) (bool, error) {
	var count int64
	err := r.GetDB(ctx).Model(&domain.AuthMethodDomain{}).Where("provider = ? AND oauth_id = ? AND deleted_at IS NULL", "PHONE", phone).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindByOAuthId tìm kiếm bằng provider và oauth_id
func (r *AuthMethodPostgres) FindByOAuthID(ctx context.Context, oauthId, provider string) (*domain.AuthMethodDomain, error) {
	var user domain.AuthMethodDomain
	err := r.GetDB(ctx).Where("provider = ? AND oauth_id = ? AND deleted_at IS NULL", provider, oauthId).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create tạo mới auth
func (r *AuthMethodPostgres) Create(ctx context.Context, auth *domain.AuthMethodDomain) (*domain.AuthMethodDomain, error) {
	err := r.GetDB(ctx).Create(auth).Error
	if err != nil {
		return nil, err
	}
	return auth, nil
}

// Update cập nhật auth
func (r *AuthMethodPostgres) Update(ctx context.Context, auth *domain.AuthMethodDomain) (*domain.AuthMethodDomain, error) {
	err := r.GetDB(ctx).Save(auth).Error
	if err != nil {
		return nil, err
	}
	return auth, nil
}

// DeleteProfile xóa profile
func (r *AuthMethodPostgres) DeleteProfile(ctx context.Context, profileID uint64) error {
	return r.GetDB(ctx).Model(&domain.AuthMethodDomain{}).
		Where("user_id = ?", profileID).
		Update("deleted_at", time.Now()).Error
}

// FindByUsername tìm kiếm auth bằng username
func (r *AuthMethodPostgres) FindByUsername(ctx context.Context, username string) (*domain.AuthMethodDomain, error) {
	var auth domain.AuthMethodDomain
	err := r.GetDB(ctx).Where("username = ? AND deleted_at IS NULL", username).First(&auth).Error
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

// FindByEmail tìm auth method theo email
func (r *AuthMethodPostgres) FindByEmail(ctx context.Context, email string) (*domain.AuthMethodDomain, error) {
	var auth domain.AuthMethodDomain
	err := r.GetDB(ctx).Where("email = ? AND deleted_at IS NULL", email).First(&auth).Error
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

// FindByID tìm auth method theo ID
func (r *AuthMethodPostgres) FindByID(ctx context.Context, id uint64) (*domain.AuthMethodDomain, error) {
	var auth domain.AuthMethodDomain
	err := r.GetDB(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&auth).Error
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

// FindAdminsByProvider tìm danh sách admin theo provider với phân trang
func (r *AuthMethodPostgres) FindAdminsByProvider(ctx context.Context, provider string, page, size int) ([]*domain.AuthMethodDomain, int64, error) {
	var admins []*domain.AuthMethodDomain
	var total int64

	// Đếm tổng số
	err := r.GetDB(ctx).WithContext(ctx).Model(&domain.AuthMethodDomain{}).
		Where("provider = ? AND deleted_at IS NULL", provider).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Lấy danh sách với phân trang
	offset := (page - 1) * size
	err = r.GetDB(ctx).WithContext(ctx).
		Where("provider = ? AND deleted_at IS NULL", provider).
		Order("created_at DESC").
		Offset(offset).
		Limit(size).
		Find(&admins).Error
	if err != nil {
		return nil, 0, err
	}

	return admins, total, nil
}

// FindByUsernameAndProvider tìm auth method theo username và provider
func (r *AuthMethodPostgres) FindByUsernameAndProvider(ctx context.Context, username, provider string) (*domain.AuthMethodDomain, error) {
	var auth domain.AuthMethodDomain
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
	err := r.GetDB(ctx).WithContext(ctx).Model(&domain.AuthMethodDomain{}).
		Where("provider = ? AND deleted_at IS NULL", provider).
		Count(&count).Error
	return count, err
}

func (r *AuthMethodPostgres) Delete(ctx context.Context, id uint64) error {
	return r.GetDB(ctx).WithContext(ctx).Model(&domain.AuthMethodDomain{}).
		Where("id = ?", id).
		Delete(&domain.AuthMethodDomain{}).Error
}

// SoftDelete soft delete auth method
func (r *AuthMethodPostgres) SoftDelete(ctx context.Context, id uint64) error {
	now := time.Now()
	return r.GetDB(ctx).WithContext(ctx).Model(&domain.AuthMethodDomain{}).Where("id = ?", id).Update("deleted_at", now).Error
}

// RestoreAccount khôi phục tài khoản đã bị soft delete
func (r *AuthMethodPostgres) RestoreAccount(c context.Context, id uint64) error {
	return r.GetDB(c).WithContext(c).Model(&domain.AuthMethodDomain{}).
		Where("id = ?", id).
		Update("deleted_at", nil).Error
}

// FindDeletedByPhone tìm auth method đã bị xóa theo phone
func (r *AuthMethodPostgres) FindDeletedByPhone(c context.Context, phone string) (*domain.AuthMethodDomain, error) {
	var auth domain.AuthMethodDomain
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
func (r *AuthMethodPostgres) FindDeletedByEmail(c context.Context, email string) (*domain.AuthMethodDomain, error) {
	var auth domain.AuthMethodDomain
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
func (r *AuthMethodPostgres) FindDeletedByUsername(c context.Context, username string) (*domain.AuthMethodDomain, error) {
	var auth domain.AuthMethodDomain
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
