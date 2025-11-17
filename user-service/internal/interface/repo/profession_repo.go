package repo

import (
	_crud "common/domain/crud"
	_dto "common/domain/dto"
	"context"
	"user/internal/models"
)

type IProfessionRepo interface {
	_crud.ICrudRepo[models.ProfessionEntity]
	ListItemByProfileID(c context.Context, profileID uint64) ([]models.ProfessionEntity, error)
	BatchSave(c context.Context, professions *[]models.ProfessionEntity, deletedIds []uint64) error
	ListVerified(ctx context.Context, pagable _dto.IPagable) ([]models.ProfessionEntity, uint32, error)
}
