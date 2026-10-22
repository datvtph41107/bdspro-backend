package repo

import (
	_crud "common/domain/crud"
	_dto "common/domain/dto"
	"context"
	"mdn4/modules/user/domain/user"
)

type IProfessionRepo interface {
	_crud.ICrudRepo[user.ProfessionEntity]
	ListItemByProfileID(c context.Context, profileID uint64) ([]user.ProfessionEntity, error)
	BatchSave(c context.Context, professions *[]user.ProfessionEntity, deletedIds []uint64) error
	ListVerified(ctx context.Context, pagable _dto.IPagable) ([]user.ProfessionEntity, uint32, error)
}
