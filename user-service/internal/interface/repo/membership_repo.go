package repo

import (
	_crud "common/domain/crud"
	"context"
	"user/internal/models"
)

type IMembershipRepo interface {
	_crud.ICrudRepo[models.Membership]
	GetByProfileID(ctx context.Context, profileID uint64) (*models.Membership, error)
}
