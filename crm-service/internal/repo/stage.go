package repo

import (
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
)

type StageRepo interface {
	Search(c context.Context, pipelineId uint64, dto dto.StageSearchDTO) ([]domain.StageEntity, int64, error)
	Create(c context.Context, entity *domain.StageEntity) (*domain.StageEntity, error)
	Update(c context.Context, entity *domain.StageEntity) (*domain.StageEntity, error)
	Delete(c context.Context, id uint64) error
	GetByID(c context.Context, id uint64) (*domain.StageEntity, error)
	GetOneStageDefault(c context.Context, pipelineId uint64) (*domain.StageEntity, error)
}
