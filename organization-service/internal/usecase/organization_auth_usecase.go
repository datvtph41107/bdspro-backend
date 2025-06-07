package usecase

import (
	"context"

	"github.com/tuandq2112/organization-service/internal/custom_error"
	"github.com/tuandq2112/organization-service/internal/domain/repository"
)

type OrganizationAuthUsecase interface {
	HasPermission(ctx context.Context, userId uint32, organizationId uint32, permission string) (bool, error)
	HasRole(ctx context.Context, userId uint32, organizationId uint32, role string) (bool, error)
	IsMemberOrganization(ctx context.Context, userId uint32, organizationId uint32) (bool, error)
}

type organizationAuthUsecase struct {
	organizationMemberRepository repository.OrganizationMemberRepository
	organizationRoleRepository   repository.OrganizationRoleRepository
}

func NewOrganizationAuthUsecase(
	organizationMemberRepository repository.OrganizationMemberRepository, 
	organizationRoleRepository repository.OrganizationRoleRepository,
	organizationPermissionRepository repository.OrganizationPermissionRepository,
	) OrganizationAuthUsecase {
	return &organizationAuthUsecase{
		organizationMemberRepository: organizationMemberRepository,
		organizationRoleRepository:   organizationRoleRepository,
	}
}
func (u *organizationAuthUsecase) HasPermission(ctx context.Context, userId uint32, organizationId uint32, permission string) (bool, error) {
	organizationMember, err := u.organizationMemberRepository.FindByUserIdAndOrganizationId(ctx, userId, organizationId)
	if err != nil {
		return false, err
	}
	if organizationMember == nil {
		return false, custom_error.RecordNotFound("Organization member not found")
	}
	role, err := u.organizationRoleRepository.FindById(ctx, organizationMember.RoleId)
	if err != nil {
		return false, err
	}

	for _, p := range role.Permissions {
		if p.Key == permission {
			return true, nil
		}
	}
	
	return false, nil
}

func (u *organizationAuthUsecase) HasRole(ctx context.Context, userId uint32, organizationId uint32, role string) (bool, error) {
	organizationMember, err := u.organizationMemberRepository.FindByUserIdAndOrganizationId(ctx, userId, organizationId)
	if err != nil {
		return false, err
	}
	if organizationMember == nil {
		return false, custom_error.RecordNotFound("Organization member not found")
	}
	roleEntity, err := u.organizationRoleRepository.FindById(ctx, organizationMember.RoleId)
	if err != nil {
		return false, err	
	}
	return roleEntity.Key == role, nil
}

func (u *organizationAuthUsecase) IsMemberOrganization(ctx context.Context, userId uint32, organizationId uint32) (bool, error) {
	organizationMember, err := u.organizationMemberRepository.FindByUserIdAndOrganizationId(ctx, userId, organizationId)
	if err != nil {
		return false, err
	}
	return organizationMember != nil, nil
}