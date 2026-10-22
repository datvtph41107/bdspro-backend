package repository

import (
	"context"
	"mdn4/modules/crm/dto"
	domain "mdn4/modules/crm/entity"
	"mdn4/modules/crm/enums"
)

type PipelineRepo interface {
	Search(c context.Context, ownerId uint64, ownerType enums.EOwnerOf, dto dto.PipelineSearchDTO, withStages bool) ([]domain.PipelineEntity, int64, error)
	Create(c context.Context, entity *domain.PipelineEntity) (*domain.PipelineEntity, error)
	Update(c context.Context, entity *domain.PipelineEntity) (*domain.PipelineEntity, error)
	Delete(c context.Context, id uint64) error
	GetByID(c context.Context, id uint64) (*domain.PipelineEntity, error)
	GetDefault(c context.Context, organizationId uint64) (*domain.PipelineEntity, error)
}
