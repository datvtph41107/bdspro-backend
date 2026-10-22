package repo

import (
	_crud "common/domain/crud"
	_dto "common/domain/dto"
	"context"
	"mdn4/modules/user/entity"
)

type IMainAreaRepo interface {
	_crud.ICrudRepo[entity.MainAreaEntity]
	ListByProfileID(ctx context.Context, profileID uint64) ([]entity.MainAreaEntity, error)
	ReplaceProfileAreas(ctx context.Context, profileID uint64, areaIDs []uint64) error
	Search(ctx context.Context, text string, pagable _dto.IPagable) ([]entity.MainAreaEntity, uint32, error)
	GetByNames(ctx context.Context, names []string) ([]entity.MainAreaEntity, error)
}
