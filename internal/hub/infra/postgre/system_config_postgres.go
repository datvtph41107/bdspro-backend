package postgres

import (
	"context"
	"errors"
	"fmt"
	"hub/internal/domain"
	"hub/internal/dto"
	"hub/internal/enums"
	"hub/internal/repo"

	"gorm.io/gorm"
)

type SystemConfigPostgres struct {
	DB *gorm.DB
}

func NewSystemConfigPostgres(db *gorm.DB) repo.ISystemConfigRepo {
	return &SystemConfigPostgres{
		DB: db,
	}
}

func (r *SystemConfigPostgres) Create(ctx context.Context, config *domain.SystemConfigEntity) (*domain.SystemConfigEntity, error) {
	if err := r.DB.WithContext(ctx).Create(config).Error; err != nil {
		return nil, fmt.Errorf("failed to create system config: %w", err)
	}
	return config, nil
}

func (r *SystemConfigPostgres) Update(ctx context.Context, config *domain.SystemConfigEntity) (*domain.SystemConfigEntity, error) {
	updates := map[string]interface{}{
		"name":      config.Name,
		"key":       config.Key,
		"value":     config.Value,
		"group_key": config.GroupConfig,
	}

	if err := r.DB.WithContext(ctx).
		Model(&domain.SystemConfigEntity{}).
		Where("id = ?", config.ID).
		Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update system config: %w", err)
	}

	return r.GetByID(ctx, config.ID)
}

func (r *SystemConfigPostgres) Delete(ctx context.Context, id uint64) error {
	if err := r.DB.WithContext(ctx).
		Delete(&domain.SystemConfigEntity{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete system config: %w", err)
	}
	return nil
}

func (r *SystemConfigPostgres) GetByID(ctx context.Context, id uint64) (*domain.SystemConfigEntity, error) {
	var config domain.SystemConfigEntity
	if err := r.DB.WithContext(ctx).
		Where("id = ?", id).
		First(&config).Error; err != nil {
		return nil, fmt.Errorf("system config not found: %w", err)
	}
	return &config, nil
}

// GetByKey normalizes only record-not-found to absence. Unrelated database
// failures stay infrastructure errors so callers cannot mistake them for a
// missing configuration.
func (r *SystemConfigPostgres) GetByKey(ctx context.Context, key string) (*domain.SystemConfigEntity, error) {
	var config domain.SystemConfigEntity
	err := r.DB.WithContext(ctx).
		Where("key = ?", key).
		First(&config).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get system config by key %s: %w", key, err)
	}
	return &config, nil
}

func (r *SystemConfigPostgres) GetList(ctx context.Context, req *dto.ListSystemConfigRequest) ([]*domain.SystemConfigEntity, int64, error) {
	var configs []*domain.SystemConfigEntity
	var total int64

	query := r.DB.WithContext(ctx).Model(&domain.SystemConfigEntity{})

	// Filter by search
	if req.Search != "" {
		searchTerm := "%" + req.Search + "%"
		query = query.Where("name ILIKE ? OR key ILIKE ?", searchTerm, searchTerm)
	}

	// Filter by key
	if req.Key != "" {
		query = query.Where("key = ?", req.Key)
	}

	// Filter by group
	if req.GroupKey != nil {
		query = query.Where("group_key = ?", *req.GroupKey)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated list
	if err := query.
		Order("created_at DESC").
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Find(&configs).Error; err != nil {
		return nil, 0, err
	}

	return configs, total, nil
}

func (r *SystemConfigPostgres) GetAll(ctx context.Context) ([]*domain.SystemConfigEntity, error) {
	var configs []*domain.SystemConfigEntity
	if err := r.DB.WithContext(ctx).
		Order("key ASC").
		Find(&configs).Error; err != nil {
		return nil, fmt.Errorf("failed to get all system configs: %w", err)
	}
	return configs, nil
}

func (r *SystemConfigPostgres) GetByGroup(ctx context.Context, group enums.ESystemConfigGroup) ([]*domain.SystemConfigEntity, error) {
	var configs []*domain.SystemConfigEntity
	if err := r.DB.WithContext(ctx).
		Where("group_key = ?", group).
		Order("key ASC").
		Find(&configs).Error; err != nil {
		return nil, fmt.Errorf("failed to get configs by group: %w", err)
	}
	return configs, nil
}

func (r *SystemConfigPostgres) BulkUpsert(ctx context.Context, configs []*domain.SystemConfigEntity) ([]*domain.SystemConfigEntity, int, int, error) {
	createdCount := 0
	updatedCount := 0
	results := make([]*domain.SystemConfigEntity, 0, len(configs))

	// Process each config
	for _, config := range configs {
		// Check if key exists. Only semantic absence may enter the create path;
		// technical lookup failures must stop before any write.
		existing, err := r.GetByKey(ctx, config.Key)
		if err != nil {
			return nil, 0, 0, fmt.Errorf("failed to lookup config %s before upsert: %w", config.Key, err)
		}

		if existing == nil {
			created, err := r.Create(ctx, config)
			if err != nil {
				return nil, 0, 0, fmt.Errorf("failed to create config %s: %w", config.Key, err)
			}
			results = append(results, created)
			createdCount++
			continue
		}

		// Key exists, update
		existing.Name = config.Name
		existing.Value = config.Value
		existing.GroupConfig = config.GroupConfig

		updated, err := r.Update(ctx, existing)
		if err != nil {
			return nil, 0, 0, fmt.Errorf("failed to update config %s: %w", config.Key, err)
		}
		results = append(results, updated)
		updatedCount++
	}

	return results, createdCount, updatedCount, nil
}
