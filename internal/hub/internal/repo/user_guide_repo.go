package repo

import (
	_crud "common/domain/crud"
	_dto "common/domain/dto"
	"context"
	"hub/internal/domain"
)

type IUserGuideRepo interface {
	_crud.ICrudRepo[domain.UserGuideEntity]
	GetListWithFilter(ctx context.Context, title string, groupKey string, pagable _dto.IPagable) ([]*domain.UserGuideEntity, int64, error)
	GetGroups(ctx context.Context) ([]map[string]interface{}, error)
	GetSimpleListWithText(ctx context.Context, text string, groupKey string, pagable _dto.IPagable) ([]*domain.UserGuideEntity, int64, error)
	GetByKey(ctx context.Context, key string) (*domain.UserGuideEntity, error)
	// GetAllWithKey lấy toàn bộ user guide có key != '' (dùng cho API sync by key).
	GetAllWithKey(ctx context.Context) ([]*domain.UserGuideEntity, error)
}
