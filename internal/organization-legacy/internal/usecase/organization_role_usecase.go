package usecase

import (
	_dto "common/domain/dto"
	_utils "common/utils"
	"context"
	"errors"

	"organization/env"
	"organization/internal/custom_error"
	"organization/internal/domain/entity"
	"organization/internal/domain/repository"
	"organization/internal/enums"
	"organization/pkg/utils"

	"gorm.io/gorm"
)

type OrganizationRoleUsecase interface {
	CreateNewRole(ctx context.Context, model *entity.OrganizationRole) (*entity.OrganizationRole, error)
	UpdateRole(ctx context.Context, model *entity.OrganizationRole) (*entity.OrganizationRole, error)
	UpdateRoleWithMembers(ctx context.Context, model *entity.OrganizationRole, memberIds []uint32, permissionIds []uint64) (*entity.OrganizationRole, error)
	DeleteRole(ctx context.Context, roleId uint32) error
	GetRoleByOrganizationId(ctx context.Context) ([]*entity.OrganizationRole, error)
	GetRoleCurrent(ctx context.Context, dto _dto.Pagable) ([]*entity.OrganizationRole, uint32, error)
	GetRoles(ctx context.Context, page, size int) ([]*entity.OrganizationRole, uint32, error)
	GetRole(ctx context.Context, roleId uint32) (*entity.OrganizationRole, error)
	AddPermissionToRole(ctx context.Context, roleId uint32, permissionId uint32) error
	CreateNewRoleWithMembers(ctx context.Context, model *entity.OrganizationRole, memberIds []uint32, permissionIds []uint64) (*entity.OrganizationRole, error)
	GetRoleDealMembers(ctx context.Context) ([]*entity.OrganizationRole, error)
}

type organizationRoleUsecase struct {
	organizationRoleRepository       repository.OrganizationRoleRepository
	organizationPermissionRepository repository.OrganizationPermissionRepository
	organizationRepository           repository.OrganizationRepository
	organizationAuthUsecase          OrganizationAuthUsecase
	logWorker                        *OrganizationLogWorker
	organizationMemberRepository     repository.OrganizationMemberRepository
	organizationMemberUsecase        OrganizationMemberUsecase
}

func NewOrganizationRoleUsecase(
	organizationRoleRepository repository.OrganizationRoleRepository,
	organizationPermissionRepository repository.OrganizationPermissionRepository,
	organizationRepository repository.OrganizationRepository,
	organizationAuthUsecase OrganizationAuthUsecase,
	logWorker *OrganizationLogWorker,
	organizationMemberRepository repository.OrganizationMemberRepository,
	organizationMemberUsecase OrganizationMemberUsecase,
) OrganizationRoleUsecase {
	return &organizationRoleUsecase{
		organizationRoleRepository:       organizationRoleRepository,
		organizationPermissionRepository: organizationPermissionRepository,
		organizationRepository:           organizationRepository,
		organizationAuthUsecase:          organizationAuthUsecase,
		logWorker:                        logWorker,
		organizationMemberRepository:     organizationMemberRepository,
		organizationMemberUsecase:        organizationMemberUsecase,
	}
}

func (u *organizationRoleUsecase) CreateNewRole(ctx context.Context, model *entity.OrganizationRole) (*entity.OrganizationRole, error) {
	currentUserId := utils.GetUserID(ctx, env.USER_CONTEXT)
	organizationId := _utils.GetOrganizationIdFromContext(ctx)
	hasRole, err := u.organizationAuthUsecase.HasPermission(ctx, currentUserId, uint32(organizationId), env.CREATE_ORGANIZATION_ROLE_PERMISSION)
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
	existRole, err := u.organizationRoleRepository.FindByOrganizationIdAndKey(ctx, uint32(organizationId), model.Key)
	if err != nil {
		return nil, err
	}
	if existRole != nil {
		return nil, custom_error.RoleKeyAlreadyExists()
	}
	model.OrganizationId = uint32(organizationId)

	role, err := u.organizationRoleRepository.Create(ctx, model)
	if err != nil {
		return nil, err
	}

	// Log role creation
	u.logWorker.Push(OrganizationLogEvent{
		OrganizationId: uint32(organizationId),
		ActorId:        currentUserId,
		LogType:        "ROLE_CREATE",
		LogData:        "Created new role in organization",
	})

	return role, nil
}

func (u *organizationRoleUsecase) GetRoleCurrent(ctx context.Context, dto _dto.Pagable) ([]*entity.OrganizationRole, uint32, error) {
	// currentUserId := utils.GetUserID(ctx, env.USER_CONTEXT)
	organizationId := _utils.GetOrganizationIdFromContext(ctx)
	// hasRole, err := u.organizationAuthUsecase.IsMemberOrganization(ctx, currentUserId, uint32(organizationId))
	// if err != nil {
	// 	return nil, 0, err
	// }
	// if !hasRole {
	// 	return nil, 0, custom_error.Forbidden("You are not authorized to access this resource")
	// }
	roles, total, err := u.organizationRoleRepository.FindByOrganizationIdWithPagination(ctx, uint32(organizationId), int(dto.GetPage()), int(dto.GetSize()))
	if err != nil {
		return nil, 0, err
	}
	return roles, total, nil
}

func (u *organizationRoleUsecase) GetRoleByOrganizationId(ctx context.Context) ([]*entity.OrganizationRole, error) {
	currentUserId := utils.GetUserID(ctx, env.USER_CONTEXT)
	organizationId := _utils.GetOrganizationIdFromContext(ctx)
	hasRole, err := u.organizationAuthUsecase.HasPermission(ctx, currentUserId, uint32(organizationId), env.GET_ORGANIZATION_ROLES_PERMISSION)
	if err != nil {
		return nil, err
	}
	if !hasRole {
		return nil, custom_error.Forbidden("You are not authorized to access this resource")
	}
	roles, err := u.organizationRoleRepository.FindByOrganizationId(ctx, uint32(organizationId))
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (u *organizationRoleUsecase) UpdateRole(ctx context.Context, model *entity.OrganizationRole) (*entity.OrganizationRole, error) {
	currentUserId := utils.GetUserID(ctx, env.USER_CONTEXT)
	organizationId := _utils.GetOrganizationIdFromContext(ctx)
	hasRole, err := u.organizationAuthUsecase.HasPermission(ctx, currentUserId, uint32(organizationId), env.UPDATE_ORGANIZATION_ROLE_PERMISSION)
	if err != nil {
		return nil, err
	}
	if !hasRole {
		return nil, custom_error.Forbidden("You are not authorized to access this resource")
	}

	model.OrganizationId = uint32(organizationId)
	role, err := u.organizationRoleRepository.Update(ctx, model)
	if err != nil {
		return nil, err
	}

	// Log role update
	u.logWorker.Push(OrganizationLogEvent{
		OrganizationId: uint32(organizationId),
		ActorId:        currentUserId,
		LogType:        "ROLE_UPDATE",
		LogData:        "Updated role details",
	})

	return role, nil
}

func (u *organizationRoleUsecase) UpdateRoleWithMembers(ctx context.Context, model *entity.OrganizationRole, memberIds []uint32, permissionIds []uint64) (*entity.OrganizationRole, error) {
	currentUserId := utils.GetUserID(ctx, env.USER_CONTEXT)
	organizationId := _utils.GetOrganizationIdFromContext(ctx)
	hasRole, err := u.organizationAuthUsecase.HasPermission(ctx, currentUserId, uint32(organizationId), env.UPDATE_ORGANIZATION_ROLE_PERMISSION)
	if err != nil {
		return nil, err
	}
	if !hasRole {
		return nil, custom_error.Forbidden("You are not authorized to access this resource")
	}

	// Check if role exists
	existingRole, err := u.organizationRoleRepository.FindById(ctx, model.Id)
	if err != nil {
		return nil, err
	}
	if existingRole == nil {
		return nil, custom_error.RecordNotFound("Role not found")
	}

	// Check if role belongs to current organization
	if existingRole.OrganizationId != uint32(organizationId) {
		return nil, custom_error.Forbidden("You are not authorized to access this resource")
	}

	model.OrganizationId = uint32(organizationId)
	role, err := u.organizationRoleRepository.Update(ctx, model)
	if err != nil {
		return nil, err
	}

	// Update members for this role
	for _, memberId := range memberIds {
		member, err := u.organizationMemberRepository.FindByUserIdAndOrganizationId(ctx, memberId, uint32(organizationId))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Không tồn tại → tạo mới
				member = &entity.OrganizationMember{
					OrganizationID: uint32(organizationId),
					UserID:         uint64(memberId),
					RoleId:         uint64(role.Id),
				}
				_, err = u.organizationMemberRepository.CreateOrganizationMember(ctx, member)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		} else {
			// Tồn tại → update role nếu cần
			if member.RoleId != uint64(role.Id) {
				member.RoleId = uint64(role.Id)
				_, err = u.organizationMemberRepository.UpdateOrganizationMember(ctx, member)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	// Update permissions for this role
	// Set the permission IDs to the role model so the Update method can handle the association
	role.PermissionIds = permissionIds
	role, err = u.organizationRoleRepository.Update(ctx, role)
	if err != nil {
		return nil, err
	}

	// Log role update
	u.logWorker.Push(OrganizationLogEvent{
		OrganizationId: uint32(organizationId),
		ActorId:        currentUserId,
		LogType:        "ROLE_UPDATE",
		LogData:        "Updated role with members and permissions in organization",
	})

	return role, nil
}

func (u *organizationRoleUsecase) DeleteRole(ctx context.Context, roleId uint32) error {
	currentUserId := utils.GetUserID(ctx, env.USER_CONTEXT)
	organizationId := _utils.GetOrganizationIdFromContext(ctx)
	role, err := u.organizationRoleRepository.FindById(ctx, roleId)
	if err != nil {
		return err
	}
	if role == nil {
		return custom_error.RecordNotFound("Role not found")
	}

	hasRole, err := u.organizationAuthUsecase.HasPermission(ctx, currentUserId, uint32(organizationId), env.DELETE_ORGANIZATION_ROLE_PERMISSION)
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
		OrganizationId: uint32(organizationId),
		ActorId:        currentUserId,
		LogType:        "ROLE_DELETE",
		LogData:        "Deleted role from organization",
	})

	return nil
}

func (u *organizationRoleUsecase) GetRoles(ctx context.Context, page, size int) ([]*entity.OrganizationRole, uint32, error) {
	currentUserId := utils.GetUserID(ctx, env.USER_CONTEXT)
	organizationId := _utils.GetOrganizationIdFromContext(ctx)
	hasRole, err := u.organizationAuthUsecase.HasPermission(ctx, currentUserId, uint32(organizationId), env.GET_ORGANIZATION_ROLES_PERMISSION)
	if err != nil {
		return nil, 0, err
	}
	if !hasRole {
		return nil, 0, custom_error.Forbidden("You are not authorized to access this resource")
	}

	roles, total, err := u.organizationRoleRepository.FindByOrganizationIdWithPagination(ctx, uint32(organizationId), page, size)
	if err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

func (u *organizationRoleUsecase) GetRole(ctx context.Context, roleId uint32) (*entity.OrganizationRole, error) {
	currentUserId := utils.GetUserID(ctx, env.USER_CONTEXT)
	organizationId := _utils.GetOrganizationIdFromContext(ctx)
	hasRole, err := u.organizationAuthUsecase.HasPermission(ctx, currentUserId, uint32(organizationId), env.GET_ORGANIZATION_ROLE_PERMISSION)
	if err != nil {
		return nil, err
	}
	if !hasRole {
		return nil, custom_error.Forbidden("You are not authorized to access this resource")
	}

	role, err := u.organizationRoleRepository.FindById(ctx, roleId)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, custom_error.RecordNotFound("Role not found")
	}

	// Check if role belongs to current organization
	if role.OrganizationId != uint32(organizationId) {
		return nil, custom_error.Forbidden("You are not authorized to access this resource")
	}

	return role, nil
}

func (u *organizationRoleUsecase) AddPermissionToRole(ctx context.Context, roleId uint32, permissionId uint32) error {
	currentUserId := utils.GetUserID(ctx, env.USER_CONTEXT)
	organizationId := _utils.GetOrganizationIdFromContext(ctx)
	hasRole, err := u.organizationAuthUsecase.HasPermission(ctx, currentUserId, uint32(organizationId), env.ADD_PERMISSION_TO_ORGANIZATION_ROLE_PERMISSION)
	if err != nil {
		return err
	}
	if !hasRole {
		return custom_error.Forbidden("You are not authorized to access this resource")
	}
	role, err := u.organizationRoleRepository.FindById(ctx, roleId)
	if err != nil {
		return err
	}
	if role == nil {
		return custom_error.RecordNotFound("Role not found")
	}
	permission, err := u.organizationPermissionRepository.FindById(ctx, permissionId)
	if err != nil {
		return err
	}
	if permission == nil {
		return custom_error.RecordNotFound("Permission not found")
	}
	role.Permissions = append(role.Permissions, permission)
	err = u.organizationRoleRepository.AddPermissionToRole(ctx, roleId, permissionId)
	if err != nil {
		return err
	}
	return nil
}

func (u *organizationRoleUsecase) CreateNewRoleWithMembers(ctx context.Context, model *entity.OrganizationRole, memberIds []uint32, permissionIds []uint64) (*entity.OrganizationRole, error) {
	currentUserId := utils.GetUserID(ctx, env.USER_CONTEXT)
	organizationId := _utils.GetOrganizationIdFromContext(ctx)
	hasRole, err := u.organizationAuthUsecase.HasPermission(ctx, currentUserId, uint32(organizationId), env.CREATE_ORGANIZATION_ROLE_WITH_MEMBERS_PERMISSION)
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
	existRole, err := u.organizationRoleRepository.FindByOrganizationIdAndKey(ctx, uint32(organizationId), model.Key)
	if err != nil {
		return nil, err
	}
	if existRole != nil {
		return nil, custom_error.RoleKeyAlreadyExists()
	}
	model.OrganizationId = uint32(organizationId)

	role, err := u.organizationRoleRepository.Create(ctx, model)
	if err != nil {
		return nil, err
	}

	// Assign members to the role
	for _, memberId := range memberIds {

		member, err := u.organizationMemberRepository.FindByUserIdAndOrganizationId(ctx, memberId, uint32(organizationId))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Không tồn tại → tạo mới
				member = &entity.OrganizationMember{
					OrganizationID: uint32(organizationId),
					UserID:         uint64(memberId),
					RoleId:         uint64(role.Id),
				}
				_, err = u.organizationMemberRepository.CreateOrganizationMember(ctx, member)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		}

		// Tồn tại → update role nếu cần
		if member.RoleId != uint64(role.Id) {
			member.RoleId = uint64(role.Id)
			_, err = u.organizationMemberRepository.UpdateOrganizationMember(ctx, member)
			if err != nil {
				return nil, err
			}
		}
	}

	// Assign permissions to the role
	for _, permissionId := range permissionIds {
		permission, err := u.organizationPermissionRepository.FindById(ctx, uint32(permissionId))
		if err != nil {
			return nil, err
		}
		if permission == nil {
			return nil, custom_error.RecordNotFound("Permission not found")
		}
		err = u.organizationRoleRepository.AddPermissionToRole(ctx, role.Id, uint32(permissionId))
		if err != nil {
			return nil, err
		}
	}

	// Log role creation
	u.logWorker.Push(OrganizationLogEvent{
		OrganizationId: uint32(organizationId),
		ActorId:        currentUserId,
		LogType:        "ROLE_CREATE",
		LogData:        "Created new role with members and permissions in organization",
	})

	return role, nil
}

func (u *organizationRoleUsecase) GetRoleDealMembers(ctx context.Context) ([]*entity.OrganizationRole, error) {
	roles, err := u.organizationRoleRepository.FindByDomain(ctx, enums.DomainTypeDeal)
	if err != nil {
		return nil, err
	}
	return roles, nil
}
