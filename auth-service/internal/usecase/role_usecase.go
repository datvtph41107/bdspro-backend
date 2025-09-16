package usecase

import (
	"common/case/crud"
	_errors "common/errors"
	_utils "common/utils"
	"context"

	"auth/internal/domain"
	"auth/internal/dto"
	"auth/internal/interface/provider"
	"auth/internal/interface/repo"
)

type RoleUsecase struct {
	crud.BaseUsecase[domain.Role, repo.RoleRepository]
	OrganizationProvider provider.OrganizationProvider
}

func NewRoleUsecase(roleRepo repo.RoleRepository, organizationProvider provider.OrganizationProvider) *RoleUsecase {
	return &RoleUsecase{
		BaseUsecase: crud.BaseUsecase[domain.Role, repo.RoleRepository]{
			Repo: roleRepo,
		},
		OrganizationProvider: organizationProvider,
	}
}

// Create tạo mới role
func (uc *RoleUsecase) Create(ctx context.Context, role *domain.Role) (*domain.Role, error) {
	if role.RoleName == "" {
		return nil, _errors.ReturnError(400, "name không được để trống")
	}
	if role.PermissionIDs == nil {
		return nil, _errors.ReturnError(400, "permissionIds không được để trống")
	}

	err := uc.Repo.Create(ctx, role)
	if err != nil {
		return nil, err
	}
	return role, nil
}

// GetByID lấy role theo ID
func (uc *RoleUsecase) Detail(ctx context.Context, id uint64) (*domain.Role, error) {
	role, err := uc.Repo.GetDetail(ctx, id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, _errors.ReturnError(404, "Role không tồn tại")
	}
	return role, nil
}

// Update cập nhật role
func (uc *RoleUsecase) Update(ctx context.Context, id uint64, role *domain.Role) (*domain.Role, error) {
	if role.RoleName == "" {
		return nil, _errors.ReturnError(400, "Tên role không được để trống")
	}
	// if role.Key == "" {
	// 	return nil, _errors.ReturnError(400, "Key role không được để trống")
	// }

	// Kiểm tra role có tồn tại không
	existingRole, err := uc.Repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existingRole == nil {
		return nil, _errors.ReturnError(404, "Role không tồn tại")
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
		return _errors.ReturnError(404, "Role không tồn tại")
	}

	return uc.Repo.Delete(ctx, id)
}

// GetByOrganizationID lấy danh sách role theo organization ID
func (uc *RoleUsecase) GetByOrganizationID(ctx context.Context, organizationID uint64) ([]*domain.Role, error) {
	return uc.Repo.FindByOrganizationId(ctx, organizationID)
}

// GetByOrganizationIDWithPagination lấy danh sách role theo organization ID với phân trang
func (uc *RoleUsecase) GetByOrganizationIDWithPagination(ctx context.Context, organizationID uint64, page, size int) ([]domain.Role, uint64, error) {
	// return uc.roleRepo.FindByOrganizationIdWithPagination(ctx, organizationID, page, size)
	return nil, 0, nil
}

// AddPermissionToRole thêm permission vào role (legacy - giữ lại để tương thích)
func (uc *RoleUsecase) AddPermissionToRole(ctx context.Context, roleID uint64, permissionID uint64) error {
	return uc.Repo.AddPermissionToRole(ctx, roleID, permissionID)
}

// AddPermissionsToRole thêm nhiều permission vào role
func (uc *RoleUsecase) AddPermissionsToRole(ctx context.Context, roleID uint64, permissionIDs []uint64) error {
	// Kiểm tra role có tồn tại không
	existingRole, err := uc.Repo.GetByID(ctx, roleID)
	if err != nil {
		return err
	}
	if existingRole == nil {
		return _errors.ReturnError(404, "Role không tồn tại")
	}

	// Thêm từng permission vào role
	for _, permissionID := range permissionIDs {
		err := uc.Repo.AddPermissionToRole(ctx, roleID, permissionID)
		if err != nil {
			return err
		}
	}

	return nil
}

// func (uc *RoleUsecase) GetListRole(ctx context.Context, pagable _dto.IPagable) ([]domain.Role, int64, error) {
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
		return nil, _errors.ReturnError(404, "Organization member not found")
	}
	keys, err := uc.Repo.GetPermissionKeysByRoleKey(ctx, organizationMember.RoleKey)
	if err != nil {
		return nil, err
	}

	return keys, nil
}

func (s *RoleUsecase) GetListByGroupKey(c context.Context, groupKey uint32) ([]domain.Role, int64, error) {
	entities, total, err := s.Repo.GetListByGroupKey(c, groupKey)
	if err != nil {
		return nil, 0, _errors.ReturnError(500, err.Error())
	}

	return entities, total, nil
}

// AssignRoleToUser gán role cho user dựa trên roleId
func (uc *RoleUsecase) AssignRoleToUser(ctx context.Context, req *dto.AssignRoleToUserRequest) (*dto.AssignRoleToUserResponse, error) {
	// Validate input
	if req.UserID == 0 {
		return nil, _errors.ReturnError(400, "UserID không được để trống")
	}
	if req.RoleID == 0 {
		return nil, _errors.ReturnError(400, "RoleID không được để trống")
	}

	// Lấy role theo roleId
	role, err := uc.Repo.GetByID(ctx, req.RoleID)
	if err != nil {
		return nil, _errors.ReturnError(500, "Lỗi khi lấy role: "+err.Error())
	}

	if role == nil {
		return nil, _errors.ReturnError(404, "Không tìm thấy role với roleId đã cung cấp")
	}

	// Gán role cho user
	err = uc.Repo.AssignRoleToUser(ctx, req.UserID, role)
	if err != nil {
		return nil, _errors.ReturnError(500, "Lỗi khi gán role cho user: "+err.Error())
	}

	// Convert domain role to DTO
	assignedRole := uc.convertRoleToDTO(role)

	return &dto.AssignRoleToUserResponse{
		AssignedRole: assignedRole,
		Message:      "Gán role cho user thành công",
	}, nil
}

// convertRoleToDTO convert domain.Role to dto.RoleDTO
func (uc *RoleUsecase) convertRoleToDTO(role *domain.Role) dto.RoleDTO {
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
