package postgre

import (
	"context"

	"auth/internal/domain"
	_db "common/db"
)

// AuthConfigPostgres implementation của AuthConfig repository
// @bind: auth/internal/interface/repo.IAuthConfigRepo
type AuthConfigPostgres struct {
	// TODO: Add database connection
}

// NewAuthConfigRepository tạo mới AuthConfigRepository
func NewAuthConfigRepository() *AuthConfigPostgres {
	return &AuthConfigPostgres{}
}

// GetByKey lấy config theo key
func (r *AuthConfigPostgres) GetByKey(ctx context.Context, key string) (*domain.AuthConfig, error) {
	var config domain.AuthConfig
	err := _db.DB.Where("config_key = ? AND is_active = ?", key, true).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// GetByType lấy tất cả config theo type
func (r *AuthConfigPostgres) GetByType(ctx context.Context, configType string) ([]*domain.AuthConfig, error) {
	var configs []*domain.AuthConfig
	err := _db.DB.Where("config_type = ? AND is_active = ?", configType, true).Find(&configs).Error
	if err != nil {
		return nil, err
	}
	return configs, nil
}

// GetAll lấy tất cả config
func (r *AuthConfigPostgres) GetAll(ctx context.Context) ([]*domain.AuthConfig, error) {
	var configs []*domain.AuthConfig
	err := _db.DB.Where("is_active = ?", true).Find(&configs).Error
	if err != nil {
		return nil, err
	}
	return configs, nil
}

// Create tạo config mới
func (r *AuthConfigPostgres) Create(ctx context.Context, config *domain.AuthConfig) error {
	return _db.DB.Create(config).Error
}

// Update cập nhật config
func (r *AuthConfigPostgres) Update(ctx context.Context, config *domain.AuthConfig) error {
	return _db.DB.Save(config).Error
}

// UpdateValue cập nhật giá trị config theo key
func (r *AuthConfigPostgres) UpdateValue(ctx context.Context, key string, value string) error {
	return _db.DB.Model(&domain.AuthConfig{}).
		Where("config_key = ?", key).
		Update("config_value", value).Error
}

// Delete xóa config
func (r *AuthConfigPostgres) Delete(ctx context.Context, id uint64) error {
	return _db.DB.Delete(&domain.AuthConfig{}, id).Error
}
