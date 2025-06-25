package repository

import (
	"context"

	"organization/internal/domain/entity"
)

type GroupSettingRepository interface {
	Create(ctx context.Context, setting *entity.GroupSetting) (*entity.GroupSetting, error)
	Update(ctx context.Context, setting *entity.GroupSetting) (*entity.GroupSetting, error)
	Delete(ctx context.Context, id uint32) error
	GetByID(ctx context.Context, id uint32) (*entity.GroupSetting, error)
	GetByConfigKey(ctx context.Context, configKey string) (*entity.GroupSetting, error)
	GetByGroupID(ctx context.Context, groupID uint32, page, size int) ([]*entity.GroupSetting, uint32, error)
}
