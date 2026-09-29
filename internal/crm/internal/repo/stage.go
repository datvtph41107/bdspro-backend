package repo

import (
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
)

type StageRepo interface {
	Search(ctx context.Context, pipelineId uint64, dto dto.StageSearchDTO) ([]domain.StageEntity, int64, error)
	Create(ctx context.Context, entity *domain.StageEntity) (*domain.StageEntity, error)
	Update(ctx context.Context, entity *domain.StageEntity) (*domain.StageEntity, error)
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*domain.StageEntity, error)
	GetOneStageDefault(ctx context.Context, pipelineId uint64) (*domain.StageEntity, error)
	BulkUpdate(ctx context.Context, removeIds []uint64, stages []domain.StageEntity) error
	GetByPipelineID(ctx context.Context, pipelineId uint64) ([]domain.StageEntity, error)
}