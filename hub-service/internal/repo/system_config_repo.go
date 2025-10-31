package repo

import (
	"context"
	"hub/internal/domain"
	"hub/internal/dto"
	"hub/internal/enums"
)

type ISystemConfigRepo interface {
	Create(ctx context.Context, config *domain.SystemConfigEntity) (*domain.SystemConfigEntity, error)
	Update(ctx context.Context, config *domain.SystemConfigEntity) (*domain.SystemConfigEntity, error)
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*domain.SystemConfigEntity, error)
	GetByKey(ctx context.Context, key string) (*domain.SystemConfigEntity, error)
	GetList(ctx context.Context, req *dto.ListSystemConfigRequest) ([]*domain.SystemConfigEntity, int64, error)
	GetAll(ctx context.Context) ([]*domain.SystemConfigEntity, error)
	GetByGroup(ctx context.Context, group enums.ESystemConfigGroup) ([]*domain.SystemConfigEntity, error)
	BulkUpsert(ctx context.Context, configs []*domain.SystemConfigEntity) ([]*domain.SystemConfigEntity, int, int, error)
}
