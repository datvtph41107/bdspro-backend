package repo

import (
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/enums"
)

type PipelineRepo interface {
	Search(c context.Context, ownerId uint64, ownerType enums.EOwnerOf, dto dto.PipelineSearchDTO, withStages bool) ([]domain.PipelineEntity, int64, error)
	Create(c context.Context, entity *domain.PipelineEntity) (*domain.PipelineEntity, error)
	Update(c context.Context, entity *domain.PipelineEntity) (*domain.PipelineEntity, error)
	Delete(c context.Context, id uint64) error
	GetByID(c context.Context, id uint64) (*domain.PipelineEntity, error)
	GetDefault(c context.Context, organizationId uint64) (*domain.PipelineEntity, error)
}