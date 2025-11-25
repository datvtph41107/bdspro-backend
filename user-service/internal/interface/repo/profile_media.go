package repo

import (
	"context"
	models "user/internal/models"
)

type IProfileMediaRepo interface {
	ListItemByProfileID(c context.Context, profileID uint64) ([]models.ProfileMediaEntity, error)
	BatchSave(c context.Context, medias *[]models.ProfileMediaEntity, deletedIds []uint64) error
}
