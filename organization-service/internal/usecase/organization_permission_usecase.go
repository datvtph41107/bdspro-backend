package usecase

import (
	_utils "common/utils"
	"context"

	"github.com/tuandq2112/organization-service/appconfig"
	"github.com/tuandq2112/organization-service/internal/custom_error"
	"github.com/tuandq2112/organization-service/internal/domain/entity"
	"github.com/tuandq2112/organization-service/internal/domain/repository"
	"github.com/tuandq2112/organization-service/pkg/utils"
)

type OrganizationPermissionUsecase interface {
	CreateNewPermission(ctx context.Context, permission *entity.OrganizationPermission) (*entity.OrganizationPermission, error)
	UpdatePermission(ctx context.Context, permission *entity.OrganizationPermission) (*entity.OrganizationPermission, error)
	GetPermissionsByOrganizationId(ctx context.Context, page, size int) ([]*entity.OrganizationPermission, uint32, error)
	DeletePermission(ctx context.Context, permissionId uint32) error
}

type organizationPermissionUsecase struct {
	organizationPermissionRepository repository.OrganizationPermissionRepository
	organizationRoleRepository repository.OrganizationRoleRepository
	organizationRepository repository.OrganizationRepository
	organizationAuthUsecase OrganizationAuthUsecase
}

func NewOrganizationPermissionUsecase(
	organizationPermissionRepository repository.OrganizationPermissionRepository, 
	organizationRoleRepository repository.OrganizationRoleRepository,
	organizationRepository repository.OrganizationRepository,
	organizationAuthUsecase OrganizationAuthUsecase,
) OrganizationPermissionUsecase {
	return &organizationPermissionUsecase{
		organizationPermissionRepository: organizationPermissionRepository,
		organizationRoleRepository: organizationRoleRepository,
		organizationRepository: organizationRepository,
		organizationAuthUsecase: organizationAuthUsecase,
	}
}

func (u *organizationPermissionUsecase) CreateNewPermission(ctx context.Context, permission *entity.OrganizationPermission) (*entity.OrganizationPermission, error) {
	currentUserId := utils.GetUserID(ctx, appconfig.USER_CONTEXT)
	organizationId := _utils.GetOrganizationIdFromContext(ctx)
	hasRole, err := u.organizationAuthUsecase.HasRole(ctx, currentUserId, uint32(organizationId), appconfig.ADMIN_ROLE_KEY)
	if err != nil {
		return nil, err
	}
	if !hasRole {
		return nil, custom_error.Forbidden("You are not authorized to access this resource")
	}
	organization, err := u.organizationRepository.FindById(ctx, uint32(organizationId))
	if err != nil {
		return nil, err
	}
	if organization == nil {
		return nil, custom_error.RecordNotFound("Organization not found")
	}
	existPermissionByKey, err := u.organizationPermissionRepository.FindByOrganizationIdAndKey(ctx, uint32(organizationId), permission.Key)
	if err != nil {
		return nil, err
	}
	if existPermissionByKey != nil {
		return nil, custom_error.PermissionKeyAlreadyExists()
	}
	permission.OrganizationId = uint32(organizationId)
	return u.organizationPermissionRepository.Create(ctx, permission)
}

func (u *organizationPermissionUsecase) UpdatePermission(ctx context.Context, permission *entity.OrganizationPermission) (*entity.OrganizationPermission, error) {
	currentUserId := utils.GetUserID(ctx, appconfig.USER_CONTEXT)
	organizationId := _utils.GetOrganizationIdFromContext(ctx)
	existPermission, err := u.organizationPermissionRepository.FindById(ctx, permission.Id)
	if err != nil {
		return nil, err
	}
	if existPermission == nil {
		return nil, custom_error.RecordNotFound("Permission not found")
	}
	organization, err := u.organizationRepository.FindById(ctx, uint32(organizationId))
	if err != nil {
		return nil, err
	}
	if organization == nil {
		return nil, custom_error.RecordNotFound("Organization not found")
	}
	hasRole, err := u.organizationAuthUsecase.HasRole(ctx, currentUserId, uint32(organizationId), appconfig.ADMIN_ROLE_KEY)
	if err != nil {
		return nil, err
	}
	if !hasRole {
		return nil, custom_error.Forbidden("You are not authorized to access this resource")
	}

	permission.OrganizationId = uint32(organizationId)
	return u.organizationPermissionRepository.Update(ctx, permission)
}

func (u *organizationPermissionUsecase) GetPermissionsByOrganizationId(ctx context.Context, page, size int) ([]*entity.OrganizationPermission, uint32, error) {
	currentUserId := utils.GetUserID(ctx, appconfig.USER_CONTEXT)
	organizationId := _utils.GetOrganizationIdFromContext(ctx)
	hasRole, err := u.organizationAuthUsecase.HasRole(ctx, currentUserId, uint32(organizationId), appconfig.ADMIN_ROLE_KEY)
	if err != nil {
		return nil, 0, err
	}
	if !hasRole {
		return nil, 0, custom_error.Forbidden("You are not authorized to access this resource")
	}
	return u.organizationPermissionRepository.FindByOrganizationIdWithPagination(ctx, uint32(organizationId), page, size)
}

func (u *organizationPermissionUsecase) DeletePermission(ctx context.Context, permissionId uint32) error {
	currentUserId := utils.GetUserID(ctx, appconfig.USER_CONTEXT)
	organizationId := _utils.GetOrganizationIdFromContext(ctx)
	_, err := u.organizationPermissionRepository.FindById(ctx, permissionId)
	if err != nil {
		return err
	}
	organization, err := u.organizationRepository.FindById(ctx, uint32(organizationId))
	if err != nil {
		return err
	}
	if organization == nil {
		return custom_error.RecordNotFound("Organization not found")
	}
	hasRole, err := u.organizationAuthUsecase.HasRole(ctx, currentUserId, uint32(organizationId), appconfig.ADMIN_ROLE_KEY)
	if err != nil {
		return err
	}
	if !hasRole {
		return custom_error.Forbidden("You are not authorized to access this resource")
	}
	return u.organizationPermissionRepository.Delete(ctx, permissionId)
}


