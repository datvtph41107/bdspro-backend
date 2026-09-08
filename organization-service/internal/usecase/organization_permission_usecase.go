package usecase

import (
	_utils "common/utils"
	"context"

	"organization/env"
	"organization/internal/custom_error"
	"organization/internal/domain/entity"
	"organization/internal/domain/repository"
	"organization/internal/enums"
	"organization/pkg/utils"
)

type OrganizationPermissionUsecase interface {
	CreateNewPermission(ctx context.Context, permission *entity.OrganizationPermission) (*entity.OrganizationPermission, error)
	UpdatePermission(ctx context.Context, permission *entity.OrganizationPermission) (*entity.OrganizationPermission, error)
	GetPermissionOrganization(ctx context.Context) ([]*entity.OrganizationPermission, error)
	DeletePermission(ctx context.Context, permissionId uint32) error
}

type organizationPermissionUsecase struct {
	organizationPermissionRepository repository.OrganizationPermissionRepository
	organizationRoleRepository       repository.OrganizationRoleRepository
	organizationRepository           repository.OrganizationRepository
	organizationAuthUsecase          OrganizationAuthUsecase
}

func NewOrganizationPermissionUsecase(
	organizationPermissionRepository repository.OrganizationPermissionRepository,
	organizationRoleRepository repository.OrganizationRoleRepository,
	organizationRepository repository.OrganizationRepository,
	organizationAuthUsecase OrganizationAuthUsecase,
) OrganizationPermissionUsecase {
	return &organizationPermissionUsecase{
		organizationPermissionRepository: organizationPermissionRepository,
		organizationRoleRepository:       organizationRoleRepository,
		organizationRepository:           organizationRepository,
		organizationAuthUsecase:          organizationAuthUsecase,
	}
}

func (u *organizationPermissionUsecase) CreateNewPermission(ctx context.Context, permission *entity.OrganizationPermission) (*entity.OrganizationPermission, error) {
	if err := u.organizationAuthUsecase.IsSystemAdmin(ctx); err != nil {
		return nil, err
	}
	currentUserId := utils.GetUserID(ctx, env.USER_CONTEXT)
	organizationId := _utils.GetOrganizationIdFromContext(ctx)
	hasRole, err := u.organizationAuthUsecase.HasPermission(ctx, currentUserId, uint32(organizationId), env.CREATE_ORGANIZATION_PERMISSION)
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
	// permission.OrganizationId = uint32(organizationId)
	permission.PermissionType = enums.PermissionTypeOrganization
	return u.organizationPermissionRepository.Create(ctx, permission)
}

func (u *organizationPermissionUsecase) UpdatePermission(ctx context.Context, permission *entity.OrganizationPermission) (*entity.OrganizationPermission, error) {
	if err := u.organizationAuthUsecase.IsSystemAdmin(ctx); err != nil {
		return nil, err
	}
	currentUserId := utils.GetUserID(ctx, env.USER_CONTEXT)
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
	hasRole, err := u.organizationAuthUsecase.HasPermission(ctx, currentUserId, uint32(organizationId), env.UPDATE_ORGANIZATION_PERMISSION)
	if err != nil {
		return nil, err
	}
	if !hasRole {
		return nil, custom_error.Forbidden("You are not authorized to access this resource")
	}

	// permission.OrganizationId = uint32(organizationId)
	return u.organizationPermissionRepository.Update(ctx, permission)
}

func (u *organizationPermissionUsecase) GetPermissionOrganization(ctx context.Context) ([]*entity.OrganizationPermission, error) {
	// currentUserId := utils.GetUserID(ctx, env.USER_CONTEXT)
	// organizationId := _utils.GetOrganizationIdFromContext(ctx)
	// hasRole, err := u.organizationAuthUsecase.HasRole(ctx, currentUserId, uint32(organizationId), env.ADMIN_ROLE_KEY)
	// if err != nil {
	// 	return nil, 0, err
	// }
	// if !hasRole {
	// 	return nil, 0, custom_error.Forbidden("You are not authorized to access this resource")
	// }
	permissions, err := u.organizationPermissionRepository.GetPermissionOrganization(ctx)
	if err != nil {
		return nil, err
	}
	permMap := make(map[uint32]*entity.OrganizationPermission)
	for i := range permissions {
		permMap[permissions[i].Id] = permissions[i]
	}

	var roots []*entity.OrganizationPermission

	// Xây dựng quan hệ cha-con
	for i := range permissions {
		p := permissions[i]
		if p.ParentId == 0 {
			roots = append(roots, p)
		} else if parent, ok := permMap[p.ParentId]; ok {
			parent.Childrens = append(parent.Childrens, p)
		}
	}
	return roots, nil
}

func (u *organizationPermissionUsecase) DeletePermission(ctx context.Context, permissionId uint32) error {
	if err := u.organizationAuthUsecase.IsSystemAdmin(ctx); err != nil {
		return err
	}
	currentUserId := utils.GetUserID(ctx, env.USER_CONTEXT)
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
	hasRole, err := u.organizationAuthUsecase.HasPermission(ctx, currentUserId, uint32(organizationId), env.DELETE_ORGANIZATION_PERMISSION)
	if err != nil {
		return err
	}
	if !hasRole {
		return custom_error.Forbidden("You are not authorized to access this resource")
	}
	return u.organizationPermissionRepository.Delete(ctx, permissionId)
}
