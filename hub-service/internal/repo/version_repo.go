package repo

import (
	"context"

	_crud "common/domain/crud"
	_dto "common/domain/dto"
	"hub/internal/domain"
)

type IVersionRepo interface {
	_crud.ICrudRepo[domain.VersionEntity]
	GetListWithFilter(ctx context.Context, appName string, platform string, pagable _dto.IPagable) ([]*domain.VersionEntity, int64, error)
	GetLatestByAppAndPlatform(ctx context.Context, appName string, platform string, versionName string, active bool) (*domain.VersionEntity, error)
}
