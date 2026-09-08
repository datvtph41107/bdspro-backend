package usecase

import (
	_errors "common/errors"
	_utils "common/utils"
	"context"
	"fmt"
	"sort"
	"user/internal/domain/access"
	"user/internal/interface/providers"
	"user/internal/interface/repo"
)

type PermissionUsecase struct {
	permissionRepo       repo.PermissionRepository
	roleRepo             repo.RoleRepository
	organizationProvider providers.OrganizationProvider
	bdsproProvider       providers.BdsproProvider
}

func NewPermissionUsecase(
	permissionRepo repo.PermissionRepository,
	roleRepo repo.RoleRepository,
	organizationProvider providers.OrganizationProvider,
	bdsproProvider providers.BdsproProvider,
) *PermissionUsecase {
	return &PermissionUsecase{
		permissionRepo:       permissionRepo,
		roleRepo:             roleRepo,
		organizationProvider: organizationProvider,
		bdsproProvider:       bdsproProvider,
	}
}

// GetByID lấy permission theo ID
func (uc *PermissionUsecase) GetByID(ctx context.Context, id uint64) (*access.Permission, error) {
	permission, err := uc.permissionRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	if permission == nil {
		return nil, _errors.ReturnError(404, "Permission không tồn tại")
	}
	return permission, nil
}

// GetByModuleWithPagination phục vụ API permission production cũ nhưng vẫn
// dùng repository/User IAM làm nguồn dữ liệu duy nhất.
func (uc *PermissionUsecase) GetByModuleWithPagination(ctx context.Context, module uint64, page, size int) ([]*access.Permission, uint64, error) {
	if module == 0 {
		return nil, 0, _errors.ReturnError(400, "Module không được để trống")
	}
	return uc.permissionRepo.FindByModuleWithPagination(ctx, fmt.Sprintf("%d", module), page, size)
}

// GetByOrganizationID trả về union permission keys của các role trong tổ chức.
// Đây là read-only compatibility logic; role/permission persistence vẫn do User sở hữu.
func (uc *PermissionUsecase) GetPermissionKeysByOrganizationID(ctx context.Context, organizationID uint64) ([]string, error) {
	if organizationID == 0 {
		return nil, _errors.ReturnError(400, "organizationId là bắt buộc")
	}
	roles, err := uc.roleRepo.FindByOrganizationId(ctx, organizationID)
	if err != nil {
		return nil, err
	}
	set := make(map[string]struct{})
	for _, role := range roles {
		if role == nil {
			continue
		}
		keys, keyErr := uc.roleRepo.GetPermissionKeysByRoleId(ctx, role.ID)
		if keyErr != nil {
			return nil, keyErr
		}
		for _, key := range keys {
			if key != "" {
				set[key] = struct{}{}
			}
		}
	}
	keys := make([]string, 0, len(set))
	for key := range set {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys, nil
}

// Delete xóa permission
// GetAll lấy tất cả permission
func (uc *PermissionUsecase) GetAll(ctx context.Context) ([]*access.Permission, error) {
	permissions, err := uc.permissionRepo.GetPermissions(ctx)
	if err != nil {
		return nil, err
	}

	permMap := make(map[uint64]*access.Permission)
	for i := range permissions {
		permMap[permissions[i].ID] = permissions[i]
	}

	var roots []*access.Permission

	// Xây dựng quan hệ cha-con
	for _, p := range permissions {
		permission := uc.BuildHierarchy(p, permMap)
		if permission.ParentID == nil {
			roots = append(roots, permission)
		}
	}

	return roots, nil
}

func (uc *PermissionUsecase) BuildHierarchy(permission *access.Permission, permMap map[uint64]*access.Permission) *access.Permission {
	if permission == nil {
		return nil
	}

	if permission.ParentID == nil {
		return permission
	}

	if parent, ok := permMap[*permission.ParentID]; ok {
		parent.Children = append(parent.Children, permission)
	}

	return permission
}
func (uc *PermissionUsecase) GetByGroupId(ctx context.Context, groupId uint64) ([]string, error) {
	profileId := _utils.GetProfileIdWithContext(ctx)
	groupMember, err := uc.organizationProvider.GetGroupRoleUser(ctx, groupId, profileId)
	if err != nil {
		return nil, err
	}

	if groupMember == nil {
		return []string{}, nil
	}

	if groupMember.RoleId == nil {
		return []string{}, nil
	}

	permissions, err := uc.roleRepo.GetPermissionKeysByRoleId(ctx, *groupMember.RoleId)
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

func (uc *PermissionUsecase) GetByDealId(ctx context.Context, dealId uint64) ([]string, *access.Role, error) {
	profileId := _utils.GetProfileIdWithContext(ctx)
	dealMember, err := uc.bdsproProvider.GetDealMember(ctx, dealId, profileId)
	if err != nil {
		return nil, nil, err
	}

	if dealMember == nil {
		return []string{}, nil, nil
	}

	if dealMember.RoleId == nil {
		return []string{}, nil, nil
	}

	permissions, err := uc.roleRepo.GetPermissionKeysByRoleKey(ctx, dealMember.RoleKey)
	if err != nil {
		return nil, nil, err
	}

	role, err := uc.roleRepo.GetRoleByRoleKey(ctx, dealMember.RoleKey)
	if err != nil {
		return nil, nil, err
	}

	return permissions, role, nil
}

func (uc *PermissionUsecase) GetByBranchId(ctx context.Context, branchId uint64) ([]string, error) {
	profileId := _utils.GetProfileIdWithContext(ctx)
	branchMember, err := uc.organizationProvider.GetBranchMember(ctx, branchId, profileId)
	if err != nil {
		return nil, err
	}

	if branchMember == nil {
		return []string{}, nil
	}

	if branchMember.RoleId == nil {
		return []string{}, nil
	}

	permissions, err := uc.roleRepo.GetPermissionKeysByRoleId(ctx, *branchMember.RoleId)
	if err != nil {
		return nil, err
	}

	return permissions, nil
}
