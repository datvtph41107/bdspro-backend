package usecase

import (
	"context"

	"github.com/tuandq2112/organization-service/appconfig"
	"github.com/tuandq2112/organization-service/internal/custom_error"
	"github.com/tuandq2112/organization-service/internal/domain/entity"
	"github.com/tuandq2112/organization-service/internal/domain/repository"
)

type OrganizationRoleUsecase interface {
	CreateNewRole(ctx context.Context, currentUserId uint32, model *entity.OrganizationRole) (*entity.OrganizationRole, error)
	UpdateRole(ctx context.Context, currentUserId uint32, model *entity.OrganizationRole) (*entity.OrganizationRole, error)
	DeleteRole(ctx context.Context, currentUserId uint32, roleId uint32) error
	GetRoleByOrganizationId(ctx context.Context, currentUserId uint32, organizationId uint32) ([]*entity.OrganizationRole, error)
}

type organizationRoleUsecase struct {
	organizationRoleRepository       repository.OrganizationRoleRepository
	organizationPermissionRepository repository.OrganizationPermissionRepository
	organizationMemberRepository     repository.OrganizationMemberRepository
	organizationRepository           repository.OrganizationRepository
}

func NewOrganizationRoleUsecase(
	organizationRoleRepository repository.OrganizationRoleRepository,
	organizationPermissionRepository repository.OrganizationPermissionRepository,
	organizationMemberRepository repository.OrganizationMemberRepository,
	organizationRepository repository.OrganizationRepository,
) OrganizationRoleUsecase {
	return &organizationRoleUsecase{
		organizationRoleRepository:       organizationRoleRepository,
		organizationPermissionRepository: organizationPermissionRepository,
		organizationMemberRepository:     organizationMemberRepository,
		organizationRepository:           organizationRepository,
	}
}

func (u *organizationRoleUsecase) CreateNewRole(ctx context.Context, currentUserId uint32, model *entity.OrganizationRole) (*entity.OrganizationRole, error) {
	if err := u.isAdmin(ctx, currentUserId, model.OrganizationId); err != nil {
		return nil, err
	}
	organization, err := u.organizationRepository.FindById(ctx, model.OrganizationId)
	if err != nil {
		return nil, err
	}
	if organization == nil {
		return nil, custom_error.RecordNotFound("Organization not found")
	}
	existRole, err := u.organizationRoleRepository.FindByOrganizationIdAndKey(ctx, model.OrganizationId, model.Key)
	if err != nil {
		return nil, err
	}
	if existRole != nil {
		return nil, custom_error.RoleKeyAlreadyExists()
	}
	model.OrganizationId = organization.ID

	model, err = u.organizationRoleRepository.Create(ctx, model)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (u *organizationRoleUsecase) GetRoleByOrganizationId(ctx context.Context, currentUserId uint32, organizationId uint32) ([]*entity.OrganizationRole, error) {
	if err := u.isAdmin(ctx, currentUserId, organizationId); err != nil {
		return nil, err
	}
	roles, err := u.organizationRoleRepository.FindByOrganizationId(ctx, organizationId)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (u *organizationRoleUsecase) UpdateRole(ctx context.Context, currentUserId uint32, model *entity.OrganizationRole) (*entity.OrganizationRole, error) {
	existRole, err := u.organizationRoleRepository.FindById(ctx, model.Id)
	if err != nil {
		return nil, err
	}
	if existRole == nil {
		return nil, custom_error.RecordNotFound("Role not found")
	}
	if  existRole.Id != model.Id && existRole.Key == model.Key {
		return nil, custom_error.RoleKeyAlreadyExists()
	}

	if err := u.isAdmin(ctx, currentUserId, existRole.OrganizationId); err != nil {
		return nil, err
	}
	existRole.Name = model.Name
	existRole.Key = model.Key

	return u.organizationRoleRepository.Update(ctx, existRole)
}

func (u *organizationRoleUsecase) DeleteRole(ctx context.Context, currentUserId uint32, roleId uint32) error {
	existRole, err := u.organizationRoleRepository.FindById(ctx, roleId)
	if err != nil {
		return err
	}
	if existRole == nil {
		return custom_error.RecordNotFound("Role not found")
	}
	if err := u.isAdmin(ctx, currentUserId, existRole.OrganizationId); err != nil {
		return err
	}
	return u.organizationRoleRepository.Delete(ctx, existRole.Id)
}

func (u *organizationRoleUsecase) isAdmin(ctx context.Context, userId uint32, organizationId uint32) error {
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
	if role == nil {
		return custom_error.RecordNotFound("Role not found")
	}
	if role.Key != appconfig.ADMIN_ROLE_KEY {
		return custom_error.Forbidden("You are not authorized to access this resource")
	}
	return nil
}
