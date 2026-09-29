// infra/postgres/layer_resolver_config_postgres.go
package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/interface/repo"

	"gorm.io/gorm"
)

type QHLayerResolverConfigPostgres struct {
	db *gorm.DB
}

func NewLayerResolverConfigPostgres(db *gorm.DB) repo.IQHLayerResolverConfigRepo {
	return &QHLayerResolverConfigPostgres{db: db}
}

func (r *QHLayerResolverConfigPostgres) GetByKey(ctx context.Context, key string) (json.RawMessage, error) {
	var config qh_domain.QHLayerResolverConfig
	err := r.db.WithContext(ctx).Where("config_key = ?", key).First(&config).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("get config by key failed: %w", err)
	}
	return config.ConfigValue, nil
}

func (r *QHLayerResolverConfigPostgres) GetAll(ctx context.Context) (map[string]json.RawMessage, error) {
	var configs []qh_domain.QHLayerResolverConfig
	err := r.db.WithContext(ctx).Find(&configs).Error
	if err != nil {
		return nil, fmt.Errorf("get all configs failed: %w", err)
	}
	result := make(map[string]json.RawMessage)
	for _, cfg := range configs {
		result[cfg.ConfigKey] = cfg.ConfigValue
	}
	return result, nil
}

func (r *QHLayerResolverConfigPostgres) UpdateByKey(ctx context.Context, key string, value json.RawMessage, updatedBy string) error {
	now := time.Now()
	result := r.db.WithContext(ctx).Model(&qh_domain.QHLayerResolverConfig{}).
		Where("config_key = ?", key).
		Updates(map[string]interface{}{
			"config_value": value,
			"updated_by":   updatedBy,
			"updated_at":   now,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		// Insert new record
		config := &qh_domain.QHLayerResolverConfig{
			ConfigKey:   key,
			ConfigValue: value,
			UpdatedBy:   updatedBy,
			UpdatedAt:   now,
			CreatedAt:   now,
		}
		return r.db.WithContext(ctx).Create(config).Error
	}
	return nil
}
