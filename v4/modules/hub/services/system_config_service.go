package services

import (
	"context"
	"fmt"
	"strconv"

	_errors "common/errors"
	"userv4/modules/hub/dto"
	"userv4/modules/hub/entity"
	"userv4/modules/hub/repository"
)

// SystemConfigService manages system configuration CRUD and query operations.
type SystemConfigService interface {
	GetByGroup(ctx context.Context, groupKey string) ([]*dto.SystemConfigItem, error)
	GetByKey(ctx context.Context, key string) (*dto.SystemConfigItem, error)
	GetUserSettingsByGroup(ctx context.Context, groupKey string) (map[string]string, error)
	BulkUpsert(ctx context.Context, req dto.BulkUpsertSystemConfigRequest) ([]*dto.SystemConfigItem, int, int, error)
	ResetDefaultByGroup(ctx context.Context, groupKey string) ([]*dto.SystemConfigItem, int, error)
}

type LocalSystemConfigService struct {
	repo repository.SystemConfigRepository
}

func NewLocalSystemConfigService(repo repository.SystemConfigRepository) *LocalSystemConfigService {
	return &LocalSystemConfigService{repo: repo}
}

func (s *LocalSystemConfigService) GetByGroup(ctx context.Context, groupKey string) ([]*dto.SystemConfigItem, error) {
	group, ok := entity.ParseSystemConfigGroup(groupKey)
	if !ok {
		return nil, _errors.BadRequestException(fmt.Sprintf("invalid group key: %s", groupKey))
	}
	cfgs, err := s.repo.GetByGroup(ctx, group)
	if err != nil {
		return nil, err
	}
	return mapConfigItems(cfgs), nil
}

func (s *LocalSystemConfigService) GetByKey(ctx context.Context, key string) (*dto.SystemConfigItem, error) {
	if key == "" {
		return nil, _errors.BadRequestException("key is required")
	}
	cfg, err := s.repo.GetByKey(ctx, key)
	if err != nil {
		return nil, _errors.NotFoundException(fmt.Sprintf("config not found: %s", key))
	}
	return mapConfigItem(cfg), nil
}

func (s *LocalSystemConfigService) GetUserSettingsByGroup(ctx context.Context, groupKey string) (map[string]string, error) {
	group, ok := entity.ParseSystemConfigGroup(groupKey)
	if !ok {
		return nil, _errors.BadRequestException(fmt.Sprintf("invalid group key: %s", groupKey))
	}
	cfgs, err := s.repo.GetByGroup(ctx, group)
	if err != nil {
		return nil, err
	}
	settings := make(map[string]string, len(cfgs))
	for _, cfg := range cfgs {
		settings[cfg.Key] = cfg.Value
	}
	return settings, nil
}

func (s *LocalSystemConfigService) BulkUpsert(ctx context.Context, req dto.BulkUpsertSystemConfigRequest) ([]*dto.SystemConfigItem, int, int, error) {
	cfgs := make([]*entity.SystemConfig, 0, len(req.Configs))
	for _, item := range req.Configs {
		group, ok := entity.ParseSystemConfigGroup(item.GroupKey)
		if !ok {
			return nil, 0, 0, _errors.BadRequestException(fmt.Sprintf("invalid group key for config %s: %s", item.Key, item.GroupKey))
		}
		cfgs = append(cfgs, &entity.SystemConfig{Name: item.Name, Key: item.Key, Value: item.Value, GroupConfig: group})
	}
	results, created, updated, err := s.repo.BulkUpsert(ctx, cfgs)
	if err != nil {
		return nil, 0, 0, err
	}
	return mapConfigItems(results), created, updated, nil
}

func (s *LocalSystemConfigService) ResetDefaultByGroup(ctx context.Context, groupKey string) ([]*dto.SystemConfigItem, int, error) {
	// Default reset returns empty since we don't have hardcoded defaults in v4 (caller can seed via BulkUpsert).
	return nil, 0, _errors.NotFoundException("no defaults configured for group: " + groupKey)
}

func mapConfigItem(cfg *entity.SystemConfig) *dto.SystemConfigItem {
	if cfg == nil {
		return nil
	}
	item := &dto.SystemConfigItem{
		ID:        cfg.ID,
		Name:      cfg.Name,
		Key:       cfg.Key,
		Value:     cfg.Value,
		GroupKey:  entity.SystemConfigGroupKey(cfg.GroupConfig),
		GroupName: entity.SystemConfigGroupName(cfg.GroupConfig),
	}
	if cfg.CreatedAt != nil {
		item.CreatedAt = strconv.FormatInt(cfg.CreatedAt.UnixMilli(), 10)
	}
	if cfg.UpdatedAt != nil {
		item.UpdatedAt = strconv.FormatInt(cfg.UpdatedAt.UnixMilli(), 10)
	}
	return item
}

func mapConfigItems(cfgs []*entity.SystemConfig) []*dto.SystemConfigItem {
	out := make([]*dto.SystemConfigItem, 0, len(cfgs))
	for _, c := range cfgs {
		out = append(out, mapConfigItem(c))
	}
	return out
}
