package repo

import (
	"context"
	"user/enums"
	"user/internal/dto"
	"user/internal/models"
)

type ISystemConfigRepo interface {
	Create(ctx context.Context, config *models.SystemConfigEntity) (*models.SystemConfigEntity, error)
	Update(ctx context.Context, config *models.SystemConfigEntity) (*models.SystemConfigEntity, error)
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*models.SystemConfigEntity, error)
	GetByKey(ctx context.Context, key string) (*models.SystemConfigEntity, error)
	GetList(ctx context.Context, req *dto.ListSystemConfigRequest) ([]*models.SystemConfigEntity, int64, error)
	GetAll(ctx context.Context) ([]*models.SystemConfigEntity, error)
	GetByGroup(ctx context.Context, group enums.ESystemConfigGroup) ([]*models.SystemConfigEntity, error)
	BulkUpsert(ctx context.Context, configs []*models.SystemConfigEntity) ([]*models.SystemConfigEntity, int, int, error)
}
