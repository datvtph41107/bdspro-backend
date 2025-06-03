package usecase

import (
	"context"

	"github.com/tuandq2112/organization-service/appconfig"
	"github.com/tuandq2112/organization-service/internal/custom_error"
	"github.com/tuandq2112/organization-service/internal/domain/entity"
	"github.com/tuandq2112/organization-service/internal/domain/repository"
)

type OrganizationPermissionUsecase interface {
	CreateNewPermission(ctx context.Context, currentUserId uint32, permission *entity.OrganizationPermission) (*entity.OrganizationPermission, error)
	UpdatePermission(ctx context.Context, currentUserId uint32, permission *entity.OrganizationPermission) (*entity.OrganizationPermission, error)
	GetPermissionsByOrganizationId(ctx context.Context, currentUserId uint32, organizationId uint32) ([]*entity.OrganizationPermission, error)
	DeletePermission(ctx context.Context, currentUserId uint32, permissionId uint32) error

}

type organizationPermissionUsecase struct {
	organizationPermissionRepository repository.OrganizationPermissionRepository
	organizationRoleRepository repository.OrganizationRoleRepository
	organizationMemberRepository repository.OrganizationMemberRepository
	organizationRepository repository.OrganizationRepository
}

func NewOrganizationPermissionUsecase(
	organizationPermissionRepository repository.OrganizationPermissionRepository, 
	organizationRoleRepository repository.OrganizationRoleRepository,
	organizationMemberRepository repository.OrganizationMemberRepository,
	organizationRepository repository.OrganizationRepository,
) OrganizationPermissionUsecase {
	return &organizationPermissionUsecase{
		organizationPermissionRepository: organizationPermissionRepository,
		organizationRoleRepository: organizationRoleRepository,
		organizationMemberRepository: organizationMemberRepository,
		organizationRepository: organizationRepository,
	}
}

func (u *organizationPermissionUsecase) CreateNewPermission(ctx context.Context, currentUserId uint32, permission *entity.OrganizationPermission) (*entity.OrganizationPermission, error) {

	organization, err := u.organizationRepository.FindById(ctx, permission.OrganizationId)
	if err != nil {
		return nil, err
	}
	if organization == nil {
		return nil, custom_error.RecordNotFound("Organization not found")
	}
	existPermissionByKey, err := u.organizationPermissionRepository.FindByKeyAndOrganizationId(ctx, permission.Key, organization.ID)
	if err != nil {
		return nil, err
	}
	if existPermissionByKey != nil {
		return nil, custom_error.PermissionKeyAlreadyExists()
	}
	permission.OrganizationId = organization.ID
	return u.organizationPermissionRepository.Create(ctx, permission)
}

func (u *organizationPermissionUsecase) UpdatePermission(ctx context.Context, currentUserId uint32, permission *entity.OrganizationPermission) (*entity.OrganizationPermission, error) {
	existPermission, err := u.organizationPermissionRepository.FindById(ctx, permission.Id)
	if err != nil {
		return nil, err
	}
	if existPermission == nil {
		return nil, custom_error.RecordNotFound("Permission not found")
	}
	organization, err := u.organizationRepository.FindById(ctx, existPermission.OrganizationId)
	if err != nil {
		return nil, err
	}
	if organization == nil {
		return nil, custom_error.RecordNotFound("Organization not found")
	}
	if err := u.isAdmin(ctx, currentUserId, organization.ID); err != nil {
		return nil, err
	}
	existPermission.Name = permission.Name
	existPermission.Key = permission.Key
	return u.organizationPermissionRepository.Update(ctx, permission)
}


func (u *organizationPermissionUsecase) GetPermissionsByOrganizationId(ctx context.Context, currentUserId uint32, organizationId uint32) ([]*entity.OrganizationPermission, error) {
	if err := u.isAdmin(ctx, currentUserId, organizationId); err != nil {
		return nil, err
	}
	return u.organizationPermissionRepository.FindByOrganizationId(ctx, organizationId)
}

func (u *organizationPermissionUsecase) DeletePermission(ctx context.Context, currentUserId uint32, permissionId uint32) error {
	permission, err := u.organizationPermissionRepository.FindById(ctx, permissionId)
	if err != nil {
		return err
	}
	organization, err := u.organizationRepository.FindById(ctx, permission.OrganizationId)
	if err != nil {
		return err
	}
	if organization == nil {
		return custom_error.RecordNotFound("Organization not found")
	}
	if err := u.isAdmin(ctx, currentUserId, organization.ID); err != nil {
		return err
	}
	return u.organizationPermissionRepository.Delete(ctx, permissionId)
}


func (u *organizationPermissionUsecase) isAdmin(ctx context.Context, userId uint32, organizationId uint32) error {
	organizationMember, err := u.organizationMemberRepository.FindByUserIdAndOrganizationId(ctx, userId, organizationId)
	if err != nil {
		return err
	}
	if organizationMember == nil {
		return custom_error.RecordNotFound("Organization member not found")
	}
	role, err := u.organizationRoleRepository.FindById(ctx, organizationMember.RoleId)
	if err != nil {
		return err
	}
	if role.Key != appconfig.ADMIN_ROLE_KEY {
		return custom_error.Forbidden("You are not authorized to access this resource")
	}
	return nil
}
