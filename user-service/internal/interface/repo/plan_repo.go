package repo

import (
	_crud "common/domain/crud"
	"context"
	"user/internal/models"
)

type IPlanRepo interface {
	_crud.ICrudRepo[models.Plan]
	GetVisiblePlans(ctx context.Context) ([]models.Plan, error)
}
