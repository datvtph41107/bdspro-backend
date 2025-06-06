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
	logWorker                        *OrganizationLogWorker
}

func NewOrganizationRoleUsecase(
	organizationRoleRepository repository.OrganizationRoleRepository,
	organizationPermissionRepository repository.OrganizationPermissionRepository,
	organizationRepository repository.OrganizationRepository,
	organizationAuthUsecase OrganizationAuthUsecase,
	logWorker *OrganizationLogWorker,
) OrganizationRoleUsecase {
	return &organizationRoleUsecase{
		organizationRoleRepository:       organizationRoleRepository,
		organizationPermissionRepository: organizationPermissionRepository,
		organizationRepository:           organizationRepository,
		organizationAuthUsecase:          organizationAuthUsecase,
		logWorker:                        logWorker,
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

	role, err := u.organizationRoleRepository.Create(ctx, model)
	if err != nil {
		return nil, err
	}

	// Log role creation
	u.logWorker.Push(OrganizationLogEvent{
		OrganizationId: model.OrganizationId,
		ActorId:       currentUserId,
		LogType:       "ROLE_CREATE",
		LogData:       "Created new role in organization",
	})

	return role, nil
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
	hasRole, err := u.organizationAuthUsecase.HasRole(ctx, currentUserId, model.OrganizationId, appconfig.ADMIN_ROLE_KEY)
	if err != nil {
		return nil, err
	}
	if !hasRole {
		return nil, custom_error.Forbidden("You are not authorized to access this resource")
	}

	role, err := u.organizationRoleRepository.Update(ctx, model)
	if err != nil {
		return nil, err
	}

	// Log role update
	u.logWorker.Push(OrganizationLogEvent{
		OrganizationId: model.OrganizationId,
		ActorId:       currentUserId,
		LogType:       "ROLE_UPDATE",
		LogData:       "Updated role details",
	})

	return role, nil
}

func (u *organizationRoleUsecase) DeleteRole(ctx context.Context, currentUserId uint32, roleId uint32) error {
	role, err := u.organizationRoleRepository.FindById(ctx, roleId)
	if err != nil {
		return err
	}
	if role == nil {
		return custom_error.RecordNotFound("Role not found")
	}

	hasRole, err := u.organizationAuthUsecase.HasRole(ctx, currentUserId, role.OrganizationId, appconfig.ADMIN_ROLE_KEY)
	if err != nil {
		return err
	}
	if !hasRole {
		return custom_error.Forbidden("You are not authorized to access this resource")
	}

	err = u.organizationRoleRepository.Delete(ctx, roleId)
	if err != nil {
		return err
	}

	// Log role deletion
	u.logWorker.Push(OrganizationLogEvent{
		OrganizationId: role.OrganizationId,
		ActorId:       currentUserId,
		LogType:       "ROLE_DELETE",
		LogData:       "Deleted role from organization",
	})

	return nil
}
