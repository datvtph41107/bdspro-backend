package repo

import (
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/enums"
)

type RuleRepo interface {
	Search(c context.Context, ownerId uint64, ownerType enums.EOwnerType, dto dto.RuleSearchDTO) ([]domain.RuleEntity, int64, error)
	Create(c context.Context, entity *domain.RuleEntity) (*domain.RuleEntity, error)
	Update(c context.Context, entity *domain.RuleEntity) (*domain.RuleEntity, error)
	Delete(c context.Context, id uint64) error
	GetByID(c context.Context, id uint64) (*domain.RuleEntity, error)
	Active(c context.Context, id uint64, active bool) error
}
