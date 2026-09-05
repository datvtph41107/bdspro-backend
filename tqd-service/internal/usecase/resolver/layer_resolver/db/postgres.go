package db

import (
	"context"
	"encoding/json"
	"time"

	qh_domain "tqd/internal/domain/qh"

	"gorm.io/gorm"
)

type ConfigPostgres struct {
	db *gorm.DB
}

func NewConfigPostgres(db *gorm.DB) IConfigRepository {
	return &ConfigPostgres{db: db}
}

func (r *ConfigPostgres) GetByKey(ctx context.Context, key string) (json.RawMessage, error) {
	var config qh_domain.QHLayerResolverConfig
	err := r.db.WithContext(ctx).Where("config_key = ?", key).First(&config).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return config.ConfigValue, nil
}

func (r *ConfigPostgres) GetAll(ctx context.Context) (map[string]json.RawMessage, error) {
	var configs []qh_domain.QHLayerResolverConfig
	err := r.db.WithContext(ctx).Find(&configs).Error
	if err != nil {
		return nil, err
	}
	result := make(map[string]json.RawMessage)
	for _, cfg := range configs {
		result[cfg.ConfigKey] = cfg.ConfigValue
	}
	return result, nil
}

func (r *ConfigPostgres) UpdateByKey(ctx context.Context, key string, value json.RawMessage, updatedBy string) error {
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
		// Insert new
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
