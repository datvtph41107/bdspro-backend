package repository

import (
	"context"

	"organization/internal/domain/entity"
)

type OrganizationMemberRepository interface {
	FindByUserIdAndOrganizationIdWithRole(ctx context.Context, userId, organizationId uint32) (*entity.OrganizationMember, error)
	FindByUserIdAndOrganizationId(ctx context.Context, userId, organizationId uint32) (*entity.OrganizationMember, error)
	CreateOrganizationMember(ctx context.Context, member *entity.OrganizationMember) (*entity.OrganizationMember, error)
	FindById(ctx context.Context, id uint32) (*entity.OrganizationMember, error)
	UpdateOrganizationMember(ctx context.Context, member *entity.OrganizationMember) (*entity.OrganizationMember, error)
	DeleteOrganizationMember(ctx context.Context, member *entity.OrganizationMember) (*entity.OrganizationMember, error)
	FindByOrganizationId(ctx context.Context, organizationId uint32) ([]*entity.OrganizationMember, error)
	CountMemberByOrganizationId(ctx context.Context, organizationId uint32) (uint32, error)
	FindByOrganizationIdWithPagination(ctx context.Context, organizationId uint32, page, size int) ([]*entity.OrganizationMember, uint32, error)
}
