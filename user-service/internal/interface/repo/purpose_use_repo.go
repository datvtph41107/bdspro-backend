package repo

import (
	_crud "common/domain/crud"
	"context"
	"user/internal/models"
)

type IPurposeUseRepo interface {
	_crud.ICrudRepo[models.PurposeUseEntity]
	ListByProfileID(ctx context.Context, profileID uint64) ([]models.PurposeUseEntity, error)
	AssignToProfile(ctx context.Context, profileID uint64, purposeUseIDs []uint64) error
}
