package repo

import (
	_crud "common/domain/crud"
	_dto "common/domain/dto"
	"context"
	"user/internal/models"
)

type IMainAreaRepo interface {
	_crud.ICrudRepo[models.MainAreaEntity]
	ListByProfileID(ctx context.Context, profileID uint64) ([]models.MainAreaEntity, error)
	ReplaceProfileAreas(ctx context.Context, profileID uint64, areaIDs []uint64) error
	Search(ctx context.Context, text string, pagable _dto.IPagable) ([]models.MainAreaEntity, uint32, error)
	GetByNames(ctx context.Context, names []string) ([]models.MainAreaEntity, error)
}
