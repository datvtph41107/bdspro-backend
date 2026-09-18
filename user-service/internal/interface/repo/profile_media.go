package repo

import (
	"context"
	"user/internal/domain/user"
)

type IProfileMediaRepo interface {
	ListItemByProfileID(c context.Context, profileID uint64) ([]user.ProfileMediaEntity, error)
	BatchSave(c context.Context, medias *[]user.ProfileMediaEntity, deletedIds []uint64) error
}
