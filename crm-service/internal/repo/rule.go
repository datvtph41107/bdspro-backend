package repo

import (
	base_enum "base/enum"
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
)

type RuleRepo interface {
	Search(c context.Context, ownerId uint64, ownerType base_enum.EOwnerOf, dto dto.RuleSearchDTO) ([]domain.RuleEntity, int64, error)
	Create(c context.Context, entity *domain.RuleEntity) (*domain.RuleEntity, error)
	Update(c context.Context, entity *domain.RuleEntity) (*domain.RuleEntity, error)
	Delete(c context.Context, id uint64) error
	GetByID(c context.Context, id uint64) (*domain.RuleEntity, error)
	Active(c context.Context, id uint64, active bool) error
}