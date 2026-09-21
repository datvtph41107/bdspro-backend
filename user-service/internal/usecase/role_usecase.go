package usecase

import (
	"common/case/crud"
	_errors "common/errors"
	_utils "common/utils"
	"context"
	"strings"
	"user/internal"

	"user/internal/domain/access"
	"user/internal/dto"
	"user/internal/interface/providers"
	"user/internal/interface/repo"
)

const systemRootRoleKey = "QHPRO_SYSTEM_ROOT"

type RoleUsecase struct {
	crud.BaseUsecase[access.Role, repo.RoleRepository]
	OrganizationProvider providers.OrganizationProvider
	RoleGroupRegistry    *RoleGroupRegistry
}

func NewRoleUsecase(
	roleRepo repo.RoleRepository,
	organizationProvider providers.OrganizationProvider,
	roleGroupRegistry *RoleGroupRegistry,
) *RoleUsecase {
	return &RoleUsecase{
		BaseUsecase: crud.BaseUsecase[access.Role, repo.RoleRepository]{
			Repo: roleRepo,
		},
		OrganizationProvider: organizationProvider,
		RoleGroupRegistry:    roleGroupRegistry,
	}
}

// Create tạo mới role
func (uc *RoleUsecase) Create(ctx context.Context, role *access.Role) (*access.Role, error) {
	if role.RoleName == "" {
		return nil, _errors.ReturnError(service.RoleNameRequired)
	}
	if role.PermissionIDs == nil {
		return nil, _errors.ReturnError(service.PermissionIDsRequired)
	}
	if isSystemRootRole(role) {
		return nil, _errors.ReturnError(service.RootAdminRoleBootstrapOnly)
	}

	err := uc.Repo.Create(ctx, role)
	if err != nil {
		return nil, err
	}
	return role, nil
}

// GetByID lấy role theo ID
func (uc *RoleUsecase) Detail(ctx context.Context, id uint64) (*access.Role, error) {
	role, err := uc.Repo.GetDetail(ctx, id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, _errors.ReturnError(service.RoleNotFound)
	}
	return role, nil
}

// Update cập nhật role
func (uc *RoleUsecase) Update(ctx context.Context, id uint64, role *access.Role) (*access.Role, error) {
	if role.RoleName == "" {
		return nil, _errors.ReturnError(service.RoleNameRequiredVI)
	}
	// if role.Key == "" {
	// 	return nil, _errors.ReturnError(service.RoleKeyRequired)
	// }

	// Kiểm tra role có tồn tại không
	existingRole, err := uc.Repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existingRole == nil {
		return nil, _errors.ReturnError(service.RoleNotFound)
	}
	if isSystemRootRole(existingRole) || isSystemRootRole(role) {
		return nil, _errors.ReturnError(service.RootAdminRoleMutationDenied)
	}

	role.ID = id
	err = uc.Repo.Update(ctx, id, role)
	if err != nil {
		return nil, err
	}
	return role, nil
}

// Delete xóa role
func (uc *RoleUsecase) Delete(ctx context.Context, id uint64) error {
	// Kiểm tra role có tồn tại không
	existingRole, err := uc.Repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existingRole == nil {
		return _errors.ReturnError(service.RoleNotFound)
	}
	if isSystemRootRole(existingRole) {
		return _errors.ReturnError(service.RootAdminRoleDeleteDenied)
	}

	return uc.Repo.Delete(ctx, id)
}

// GetByOrganizationID lấy danh sách role theo organization ID
func (uc *RoleUsecase) GetByOrganizationID(ctx context.Context, organizationID uint64) ([]*access.Role, error) {
	return uc.Repo.FindByOrganizationId(ctx, organizationID)
}

// GetByOrganizationIDWithPagination lấy danh sách role theo organization ID với phân trang
func (uc *RoleUsecase) GetByOrganizationIDWithPagination(ctx context.Context, organizationID uint64, page, size int) ([]access.Role, uint64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	roles, total, err := uc.Repo.FindByOrganizationIdWithPagination(ctx, organizationID, page, size)
	if err != nil {
		return nil, 0, err
	}
	result := make([]access.Role, 0, len(roles))
	for _, role := range roles {
		if role != nil {
			result = append(result, *role)
		}
	}
	return result, total, nil
}

// ReplacePermissions keeps the compatibility RPC but routes it through the
// same CRUD association owner used by Create/Update. There is no second
// permission mutation repository path.
func (uc *RoleUsecase) ReplacePermissions(ctx context.Context, roleID uint64, permissionIDs []uint64) error {
	existingRole, err := uc.Repo.GetByID(ctx, roleID)
	if err != nil {
		return err
	}
	if existingRole == nil {
		return _errors.ReturnError(service.RoleNotFound)
	}
	if isSystemRootRole(existingRole) {
		return _errors.ReturnError(service.RootAdminRolePermissionMutationDenied)
	}
	existingRole.PermissionIDs = append([]uint64(nil), permissionIDs...)
	return uc.Repo.Update(ctx, roleID, existingRole)
}

// func (uc *RoleUsecase) GetListRole(ctx context.Context, pagable _dto.IPagable) ([]access.Role, int64, error) {
// 	roles, total, err := uc.Repo.GetList(ctx, pagable)
// 	if err != nil {
// 		return nil, 0, err
// 	}
// 	return roles, total, nil
// }

func (uc *RoleUsecase) GetPermissionKeys(ctx context.Context) ([]string, error) {
	// roleIds, err := uc.Repo.GetRoleIdsByProfileId(ctx, profileId)
	profileId := _utils.GetProfileIdWithContext(ctx)
	organizationId := _utils.GetOrganizationIdFromContext(ctx)
	if organizationId == 0 {
		keys, err := uc.Repo.GetPermissionKeysByUserId(ctx, profileId)
		if err != nil {
			return nil, err
		}
		return keys, nil
	}
	organizationMember, err := uc.OrganizationProvider.GetOrganizationMember(ctx, organizationId, profileId)
	if err != nil {
		return nil, err
	}
	if organizationMember == nil {
		return nil, _errors.ReturnError(service.OrganizationMemberNotFound)
	}
	keys, err := uc.Repo.GetPermissionKeysByRoleKey(ctx, organizationMember.RoleKey)
	if err != nil {
		return nil, err
	}

	return keys, nil
}

func (s *RoleUsecase) GetListByGroupKey(c context.Context, groupKey uint32) ([]access.Role, int64, error) {
	entities, total, err := s.Repo.GetListByGroupKey(c, groupKey)
	if err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

// AssignRoleToUser gán role cho user dựa trên roleId
func (uc *RoleUsecase) AssignRoleToUser(ctx context.Context, req *dto.AssignRoleToUserRequest) (*dto.AssignRoleToUserResponse, error) {
	// Validate input
	if req.UserID == 0 {
		return nil, _errors.ReturnError(service.UserIDRequired)
	}
	if req.RoleID == 0 {
		return nil, _errors.ReturnError(service.RoleIDRequired)
	}

	// Lấy role theo roleId
	role, err := uc.Repo.GetByID(ctx, req.RoleID)
	if err != nil {
		return nil, err
	}

	if role == nil {
		return nil, _errors.ReturnError(service.RoleByIDNotFound)
	}
	if isSystemRootRole(role) || !role.AllowAssign {
		return nil, _errors.ReturnError(service.RoleAPIAssignmentDenied)
	}

	// Gán role cho user
	err = uc.Repo.AssignRoleToUser(ctx, req.UserID, role)
	if err != nil {
		return nil, err
	}

	// Convert domain role to DTO
	assignedRole := uc.convertRoleToDTO(role)

	return &dto.AssignRoleToUserResponse{
		AssignedRole: assignedRole,
		Message:      "Gán role cho user thành công",
	}, nil
}

func isSystemRootRole(role *access.Role) bool {
	return role != nil && strings.EqualFold(strings.TrimSpace(role.Key), systemRootRoleKey)
}

// convertRoleToDTO convert access.Role to dto.RoleDTO
func (uc *RoleUsecase) convertRoleToDTO(role *access.Role) dto.RoleDTO {
	roleDTO := dto.RoleDTO{
		ID:              role.ID,
		RoleName:        role.RoleName,
		RoleDescription: role.RoleDescription,
		Key:             role.Key,
		PermissionIDs:   role.PermissionIDs,
		OrganizationID:  role.OrganizationID,
		IsDefault:       role.IsDefault,
		DomainType:      role.DomainType,
		AllowAssign:     role.AllowAssign,
		RoleKey:         role.RoleKey,
		CreatedAt:       role.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:       role.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	// Convert permissions
	if role.Permissions != nil {
		permissions := make([]dto.PermissionDTO, len(role.Permissions))
		for i, perm := range role.Permissions {
			permissions[i] = dto.PermissionDTO{
				ID:             perm.ID,
				Name:           perm.Name,
				Key:            perm.Key,
				Description:    perm.Description,
				PermissionType: perm.PermissionType,
				Module:         perm.Module,
				ParentID:       perm.ParentID,
				Priority:       perm.Priority,
				CreatedAt:      perm.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
				UpdatedAt:      perm.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
			}
		}
		roleDTO.Permissions = permissions
	}

	// Convert color
	if role.Color != nil {
		roleDTO.Color = &dto.ColorDTO{
			ID:              role.Color.ID,
			Name:            role.Color.Name,
			ContentColor:    role.Color.ContentColor,
			BackgroundColor: role.Color.BackgroundColor,
			ColorKey:        role.Color.ColorKey,
			HexCode:         role.Color.HexCode,
			Description:     &role.Color.Description,
			IsActive:        role.Color.IsActive,
			CreatedAt:       role.Color.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:       role.Color.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
			CreatedBy:       0, // Color domain không có CreatedBy/UpdatedBy
			UpdatedBy:       0,
		}
	}

	// Convert role group
	if role.RoleGroupID != nil {
		roleDTO.RoleGroupID = role.RoleGroupID
	}
	// Role domain không có RoleGroupName field

	return roleDTO
}

func (uc *RoleUsecase) GetListByModuleCode(ctx context.Context, code string) ([]access.Role, error) {
	roleGroup, err := uc.RoleGroupRegistry.GetByCode(code)
	if err != nil {
		return nil, err
	}
	if roleGroup == nil {
		return nil, _errors.ReturnError(service.ModuleInvalid)
	}

	entities, err := uc.Repo.GetRoleByGroupId(ctx, roleGroup.ID)
	if err != nil {
		return nil, err
	}
	return entities, nil
}
