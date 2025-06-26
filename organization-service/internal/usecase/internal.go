package usecase

import (
	"context"
	"organization/internal/domain/entity"
)

type IMembershipClient interface {
	GetPlan(ctx context.Context, id uint64) (*entity.PlanEntity, error)
}
