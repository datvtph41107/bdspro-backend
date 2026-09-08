package repository

import (
	"context"

	"organization/internal/domain/entity"
)

type OrganizationLogActivityRepository interface {
	Create(ctx context.Context, log *entity.OrganizationLogActivity) (*entity.OrganizationLogActivity, error)
	GetByID(ctx context.Context, id uint32) (*entity.OrganizationLogActivity, error)
	GetByOrganizationID(ctx context.Context, organizationID uint32) ([]*entity.OrganizationLogActivity, error)
	GetByActorID(ctx context.Context, actorID uint32) ([]*entity.OrganizationLogActivity, error)
	GetByOrganizationIDWithPagination(ctx context.Context, organizationID uint32, page, size int) ([]*entity.OrganizationLogActivity, uint32, error)
	GetByActorIDWithPagination(ctx context.Context, actorID uint32, page, size int) ([]*entity.OrganizationLogActivity, uint32, error)
}
