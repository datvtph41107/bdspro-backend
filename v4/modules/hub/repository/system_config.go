package repository

import (
	"context"
	"errors"
	"fmt"

	"userv4/modules/hub/entity"

	"gorm.io/gorm"
)

// SystemConfigRepository manages system_config persistence.
type SystemConfigRepository interface {
	GetByKey(ctx context.Context, key string) (*entity.SystemConfig, error)
	GetByGroup(ctx context.Context, group entity.ESystemConfigGroup) ([]*entity.SystemConfig, error)
	GetAll(ctx context.Context) ([]*entity.SystemConfig, error)
	Create(ctx context.Context, cfg *entity.SystemConfig) (*entity.SystemConfig, error)
	Update(ctx context.Context, cfg *entity.SystemConfig) (*entity.SystemConfig, error)
	Delete(ctx context.Context, id uint64) error
	BulkUpsert(ctx context.Context, cfgs []*entity.SystemConfig) ([]*entity.SystemConfig, int, int, error)
}

type systemConfigPostgres struct{ db *gorm.DB }

func NewSystemConfigPostgres(db *gorm.DB) SystemConfigRepository {
	return &systemConfigPostgres{db: db}
}

func (r *systemConfigPostgres) GetByKey(ctx context.Context, key string) (*entity.SystemConfig, error) {
	var cfg entity.SystemConfig
	err := r.db.WithContext(ctx).Where("key = ? AND deleted_at IS NULL", key).First(&cfg).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("system config not found: %s", key)
	}
	return &cfg, err
}

func (r *systemConfigPostgres) GetByGroup(ctx context.Context, group entity.ESystemConfigGroup) ([]*entity.SystemConfig, error) {
	var cfgs []*entity.SystemConfig
	err := r.db.WithContext(ctx).Where("group_key = ? AND deleted_at IS NULL", group).Order("key ASC").Find(&cfgs).Error
	return cfgs, err
}

func (r *systemConfigPostgres) GetAll(ctx context.Context) ([]*entity.SystemConfig, error) {
	var cfgs []*entity.SystemConfig
	err := r.db.WithContext(ctx).Where("deleted_at IS NULL").Order("key ASC").Find(&cfgs).Error
	return cfgs, err
}

func (r *systemConfigPostgres) Create(ctx context.Context, cfg *entity.SystemConfig) (*entity.SystemConfig, error) {
	if err := r.db.WithContext(ctx).Create(cfg).Error; err != nil {
		return nil, err
	}
	return cfg, nil
}

func (r *systemConfigPostgres) Update(ctx context.Context, cfg *entity.SystemConfig) (*entity.SystemConfig, error) {
	updates := map[string]interface{}{
		"name": cfg.Name, "value": cfg.Value, "group_key": cfg.GroupConfig,
	}
	if err := r.db.WithContext(ctx).Model(&entity.SystemConfig{}).
		Where("id = ?", cfg.ID).Updates(updates).Error; err != nil {
		return nil, err
	}
	return r.GetByIDInternal(ctx, cfg.ID)
}

func (r *systemConfigPostgres) GetByIDInternal(ctx context.Context, id uint64) (*entity.SystemConfig, error) {
	var cfg entity.SystemConfig
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&cfg).Error
	return &cfg, err
}

func (r *systemConfigPostgres) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&entity.SystemConfig{}, id).Error
}

func (r *systemConfigPostgres) BulkUpsert(ctx context.Context, cfgs []*entity.SystemConfig) ([]*entity.SystemConfig, int, int, error) {
	created, updated := 0, 0
	results := make([]*entity.SystemConfig, 0, len(cfgs))
	for _, cfg := range cfgs {
		existing, err := r.GetByKey(ctx, cfg.Key)
		if err != nil {
			// key does not exist → create
			created++
			c, err := r.Create(ctx, cfg)
			if err != nil {
				return nil, 0, 0, err
			}
			results = append(results, c)
		} else {
			// key exists → update
			updated++
			existing.Name = cfg.Name
			existing.Value = cfg.Value
			existing.GroupConfig = cfg.GroupConfig
			u, err := r.Update(ctx, existing)
			if err != nil {
				return nil, 0, 0, err
			}
			results = append(results, u)
		}
	}
	return results, created, updated, nil
}
