package postgres

import (
	"context"
	"fmt"
	"sort"

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

// UpdateValuesAtomically commits a related set of dynamic configuration values
// in one transaction. Sorting keys gives concurrent writers one lock order.
func (r *AuthConfigPostgres) UpdateValuesAtomically(ctx context.Context, values map[string]string) error {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, key := range keys {
			result := tx.Model(&auth.AuthConfig{}).
				Where("config_key = ?", key).
				Update("config_value", values[key])
			if result.Error != nil {
				return fmt.Errorf("update auth config %q: %w", key, result.Error)
			}
			if result.RowsAffected != 1 {
				return fmt.Errorf("update auth config %q: expected 1 row, got %d", key, result.RowsAffected)
			}
		}
		return nil
	})
}
