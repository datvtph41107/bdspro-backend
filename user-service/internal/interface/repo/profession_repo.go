package repo

import (
	"context"
	"user/internal/models"
)

type IProfessionRepo interface {
	ListItemByProfileID(c context.Context, profileID uint64) ([]models.ProfessionEntity, error)
	BatchSave(c context.Context, professions *[]models.ProfessionEntity, deletedIds []uint64) error
}
