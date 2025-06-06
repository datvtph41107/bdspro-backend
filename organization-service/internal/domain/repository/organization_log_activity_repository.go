package repository

import (
	"context"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
)

type OrganizationLogActivityRepository interface {
	Create(ctx context.Context, log *entity.OrganizationLogActivity) (*entity.OrganizationLogActivity, error)
	GetByID(ctx context.Context, id uint32) (*entity.OrganizationLogActivity, error)
	GetByOrganizationID(ctx context.Context, organizationID uint32) ([]*entity.OrganizationLogActivity, error)
	GetByActorID(ctx context.Context, actorID uint32) ([]*entity.OrganizationLogActivity, error)
} 