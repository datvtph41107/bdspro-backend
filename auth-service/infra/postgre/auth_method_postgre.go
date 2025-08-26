package postgre

import (
	"auth/internal/domain"
	_db "common/db"
	"context"
	"time"
)

// @bind: auth/internal/interface/repo.AuthMethodRepository
type AuthMethodPostgreRepo struct {
}

func NewAuthMethodPostgreRepo() *AuthMethodPostgreRepo {
	return &AuthMethodPostgreRepo{}
}

// FindByUsername tìm auth method theo username
func (r *AuthMethodPostgreRepo) FindByUsername(c context.Context, username string) (*domain.AuthMethodDomain, error) {
	var auth domain.AuthMethodDomain
	err := _db.DB.WithContext(c).Where("username = ? AND deleted_at IS NULL", username).First(&auth).Error
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

// FindByPhone tìm auth method theo phone
func (r *AuthMethodPostgreRepo) FindByPhone(c context.Context, phone string) (*domain.AuthMethodDomain, error) {
	var auth domain.AuthMethodDomain
	err := _db.DB.WithContext(c).Where("phone = ? AND deleted_at IS NULL", phone).First(&auth).Error
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

// FindByEmail tìm auth method theo email
func (r *AuthMethodPostgreRepo) FindByEmail(c context.Context, email string) (*domain.AuthMethodDomain, error) {
	var auth domain.AuthMethodDomain
	err := _db.DB.WithContext(c).Where("email = ? AND deleted_at IS NULL", email).First(&auth).Error
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

// FindByOAuthID tìm auth method theo OAuth ID và provider
func (r *AuthMethodPostgreRepo) FindByOAuthID(c context.Context, oauthID, provider string) (*domain.AuthMethodDomain, error) {
	var auth domain.AuthMethodDomain
	err := _db.DB.WithContext(c).Where("oauth_id = ? AND provider = ? AND deleted_at IS NULL", oauthID, provider).First(&auth).Error
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

// Create tạo auth method mới
func (r *AuthMethodPostgreRepo) Create(c context.Context, auth *domain.AuthMethodDomain) (*domain.AuthMethodDomain, error) {
	err := _db.DB.WithContext(c).Create(auth).Error
	if err != nil {
		return nil, err
	}
	return auth, nil
}

// Update cập nhật auth method
func (r *AuthMethodPostgreRepo) Update(c context.Context, auth *domain.AuthMethodDomain) (*domain.AuthMethodDomain, error) {
	err := _db.DB.WithContext(c).Save(auth).Error
	if err != nil {
		return nil, err
	}
	return auth, nil
}

// Delete xóa auth method
func (r *AuthMethodPostgreRepo) Delete(c context.Context, id uint64) error {
	return _db.DB.WithContext(c).Delete(&domain.AuthMethodDomain{}, id).Error
}

// SoftDelete soft delete auth method
func (r *AuthMethodPostgreRepo) SoftDelete(c context.Context, id uint64) error {
	now := time.Now()
	return _db.DB.WithContext(c).Model(&domain.AuthMethodDomain{}).Where("id = ?", id).Update("deleted_at", now).Error
}

// FindByID tìm auth method theo ID
func (r *AuthMethodPostgreRepo) FindByID(c context.Context, id uint64) (*domain.AuthMethodDomain, error) {
	var auth domain.AuthMethodDomain
	err := _db.DB.WithContext(c).Where("id = ? AND deleted_at IS NULL", id).First(&auth).Error
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

// FindAdminsByProvider tìm danh sách admin theo provider với phân trang
func (r *AuthMethodPostgreRepo) FindAdminsByProvider(c context.Context, provider string, page, size int) ([]*domain.AuthMethodDomain, int64, error) {
	var admins []*domain.AuthMethodDomain
	var total int64

	// Đếm tổng số
	err := _db.DB.WithContext(c).Model(&domain.AuthMethodDomain{}).
		Where("provider = ? AND deleted_at IS NULL", provider).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Lấy danh sách với phân trang
	offset := (page - 1) * size
	err = _db.DB.WithContext(c).
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
func (r *AuthMethodPostgreRepo) FindByUsernameAndProvider(c context.Context, username, provider string) (*domain.AuthMethodDomain, error) {
	var auth domain.AuthMethodDomain
	err := _db.DB.WithContext(c).
		Where("username = ? AND provider = ? AND deleted_at IS NULL", username, provider).
		First(&auth).Error
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

// CountByProvider đếm số lượng auth method theo provider
func (r *AuthMethodPostgreRepo) CountByProvider(c context.Context, provider string) (int64, error) {
	var count int64
	err := _db.DB.WithContext(c).Model(&domain.AuthMethodDomain{}).
		Where("provider = ? AND deleted_at IS NULL", provider).
		Count(&count).Error
	return count, err
}

// RestoreAccount khôi phục tài khoản đã bị soft delete
func (r *AuthMethodPostgreRepo) RestoreAccount(c context.Context, id uint64) error {
	return _db.DB.WithContext(c).Model(&domain.AuthMethodDomain{}).
		Where("id = ?", id).
		Update("deleted_at", nil).Error
}

// FindDeletedByPhone tìm auth method đã bị xóa theo phone
func (r *AuthMethodPostgreRepo) FindDeletedByPhone(c context.Context, phone string) (*domain.AuthMethodDomain, error) {
	var auth domain.AuthMethodDomain
	err := _db.DB.WithContext(c).Where("phone = ? AND deleted_at IS NOT NULL", phone).First(&auth).Error
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

// FindDeletedByEmail tìm auth method đã bị xóa theo email
func (r *AuthMethodPostgreRepo) FindDeletedByEmail(c context.Context, email string) (*domain.AuthMethodDomain, error) {
	var auth domain.AuthMethodDomain
	err := _db.DB.WithContext(c).Where("email = ? AND deleted_at IS NOT NULL", email).First(&auth).Error
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

// FindDeletedByUsername tìm auth method đã bị xóa theo username
func (r *AuthMethodPostgreRepo) FindDeletedByUsername(c context.Context, username string) (*domain.AuthMethodDomain, error) {
	var auth domain.AuthMethodDomain
	err := _db.DB.WithContext(c).Where("username = ? AND deleted_at IS NOT NULL", username).First(&auth).Error
	if err != nil {
		return nil, err
	}
	return &auth, nil
}
