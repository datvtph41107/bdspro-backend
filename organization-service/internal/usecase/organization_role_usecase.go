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
	organizationRepository           repository.OrganizationRepository
	organizationAuthUsecase          OrganizationAuthUsecase
}

func NewOrganizationRoleUsecase(
	organizationRoleRepository repository.OrganizationRoleRepository,
	organizationPermissionRepository repository.OrganizationPermissionRepository,
	organizationRepository repository.OrganizationRepository,
	organizationAuthUsecase OrganizationAuthUsecase,
) OrganizationRoleUsecase {
	return &organizationRoleUsecase{
		organizationRoleRepository:       organizationRoleRepository,
		organizationPermissionRepository: organizationPermissionRepository,
		organizationRepository:           organizationRepository,
		organizationAuthUsecase:          organizationAuthUsecase,
	}
}

func (u *organizationRoleUsecase) CreateNewRole(ctx context.Context, currentUserId uint32, model *entity.OrganizationRole) (*entity.OrganizationRole, error) {
	hasRole, err := u.organizationAuthUsecase.HasRole(ctx, currentUserId, model.OrganizationId, appconfig.ADMIN_ROLE_KEY)
	if err != nil {
		return nil, err
	}
	if !hasRole {
		return nil, custom_error.Forbidden("You are not authorized to access this resource")
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
	hasRole, err := u.organizationAuthUsecase.HasRole(ctx, currentUserId, organizationId, appconfig.ADMIN_ROLE_KEY)
	if err != nil {
		return nil, err
	}
	if !hasRole {
		return nil, custom_error.Forbidden("You are not authorized to access this resource")
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
	if existRole.Id != model.Id && existRole.Key == model.Key {
		return nil, custom_error.RoleKeyAlreadyExists()
	}

	hasRole, err := u.organizationAuthUsecase.HasRole(ctx, currentUserId, existRole.OrganizationId, appconfig.ADMIN_ROLE_KEY)
	if err != nil {
		return nil, err
	}
	if !hasRole {
		return nil, custom_error.Forbidden("You are not authorized to access this resource")
	}

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
	hasRole, err := u.organizationAuthUsecase.HasRole(ctx, currentUserId, existRole.OrganizationId, appconfig.ADMIN_ROLE_KEY)
	if err != nil {
		return err
	}
	if !hasRole {
		return custom_error.Forbidden("You are not authorized to access this resource")
	}
	return u.organizationRoleRepository.Delete(ctx, existRole.Id)
}
