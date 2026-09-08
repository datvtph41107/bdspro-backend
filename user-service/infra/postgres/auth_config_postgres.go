package postgres

import (
	"context"

	"gorm.io/gorm"
	"user/internal/domain/auth"
)

// AuthConfigPostgres implements the dynamic auth configuration persistence used
// by integration providers such as ZNS.
// @bind: user/internal/interface/repo.IAuthConfigRepo
type AuthConfigPostgres struct {
	db *gorm.DB
}

func NewAuthConfigRepository(db *gorm.DB) *AuthConfigPostgres {
	return &AuthConfigPostgres{db: db}
}

func (r *AuthConfigPostgres) GetByKey(ctx context.Context, key string) (*auth.AuthConfig, error) {
	var config auth.AuthConfig
	err := r.db.WithContext(ctx).
		Where("config_key = ? AND is_active = ?", key, true).
		First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *AuthConfigPostgres) UpdateValue(ctx context.Context, key string, value string) error {
	return r.db.WithContext(ctx).
		Model(&auth.AuthConfig{}).
		Where("config_key = ?", key).
		Update("config_value", value).Error
}
