package repository

import (
	"context"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
)

type OrganizationMemberRepository interface {
	ExistByUserIdAndOrganizationId(ctx context.Context, userId, organizationId uint32) (bool, error)
	FindByUserIdAndOrganizationId(ctx context.Context, userId, organizationId uint32) (*entity.OrganizationMember, error)
	CreateOrganizationMember(ctx context.Context, organizationMember *entity.OrganizationMember) (*entity.OrganizationMember, error)
	FindById(ctx context.Context, id uint32) (*entity.OrganizationMember, error)
	UpdateOrganizationMember(ctx context.Context, organizationMember *entity.OrganizationMember) (*entity.OrganizationMember, error)
	DeleteOrganizationMember(ctx context.Context, organizationMember *entity.OrganizationMember) (*entity.OrganizationMember, error)
	FindByOrganizationId(ctx context.Context, organizationId uint32) ([]*entity.OrganizationMember, error)
}
