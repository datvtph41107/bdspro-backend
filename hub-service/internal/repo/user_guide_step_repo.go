package repo

import (
	"context"
	"hub/internal/domain"
)

type IUserGuideStepRepo interface {
	GetByUserGuideID(ctx context.Context, userGuideID uint64) ([]*domain.UserGuideStepEntity, error)
	CreateBatch(ctx context.Context, steps []*domain.UserGuideStepEntity) error
	DeleteByUserGuideID(ctx context.Context, userGuideID uint64) error
}
