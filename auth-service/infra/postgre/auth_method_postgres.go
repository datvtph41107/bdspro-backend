package postgre

import (
	"context"
	"time"

	"auth/internal/domain"

	"gorm.io/gorm"
)

// AuthMethodPostgres implementation của auth repository
// @bind: auth/internal/interface/repo.AuthMethodRepository
type AuthMethodPostgres struct {
	db *gorm.DB
}

// NewAuthRepository tạo mới AuthRepository
func NewAuthRepository(db *gorm.DB) *AuthMethodPostgres {
	return &AuthMethodPostgres{db: db}
}

// FindByPhone tìm kiếm auth bằng số điện thoại
func (r *AuthMethodPostgres) FindByPhone(ctx context.Context, phone string) (*domain.AuthMethodDomain, error) {
	var auth domain.AuthMethodDomain
	err := r.db.Where("provider = ? AND oauth_id = ? AND deleted_at IS NULL", "PHONE", phone).First(&auth).Error
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

// GetByID lấy auth theo ID
func (r *AuthMethodPostgres) GetByID(ctx context.Context, id uint64) (*domain.AuthMethodDomain, error) {
	var auth domain.AuthMethodDomain
	err := r.db.First(&auth, id).Error
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

// ExistsByOAuth kiểm tra sự tồn tại của auth theo OAuth
func (r *AuthMethodPostgres) ExistsByOAuth(ctx context.Context, phone string) (bool, error) {
	var count int64
	err := r.db.Model(&domain.AuthMethodDomain{}).Where("provider = ? AND oauth_id = ? AND deleted_at IS NULL", "PHONE", phone).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindByOAuthId tìm kiếm bằng provider và oauth_id
func (r *AuthMethodPostgres) FindByOAuthId(ctx context.Context, oauthId, provider string) (*domain.AuthMethodDomain, error) {
	var user domain.AuthMethodDomain
	err := r.db.Where("provider = ? AND oauth_id = ? AND deleted_at IS NULL", provider, oauthId).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create tạo mới auth
func (r *AuthMethodPostgres) Create(ctx context.Context, auth *domain.AuthMethodDomain) error {
	return r.db.Create(auth).Error
}

// Update cập nhật auth
func (r *AuthMethodPostgres) Update(ctx context.Context, auth *domain.AuthMethodDomain) error {
	return r.db.Save(auth).Error
}

// DeleteProfile xóa profile
func (r *AuthMethodPostgres) DeleteProfile(ctx context.Context, profileID uint64) error {
	return r.db.Model(&domain.AuthMethodDomain{}).
		Where("user_id = ?", profileID).
		Update("deleted_at", time.Now()).Error
}
