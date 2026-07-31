package usecase

import (
	_errors "common/errors"
	_utils "common/utils"
	"context"
	"user/internal/domain"
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

// Create tạo mới permission
func (uc *PermissionUsecase) Create(ctx context.Context, permission *domain.Permission) (*domain.Permission, error) {
	if permission.Name == "" {
		return nil, _errors.ReturnError(400, "Tên permission không được để trống")
	}
	if permission.Key == "" {
		return nil, _errors.ReturnError(400, "Key permission không được để trống")
	}
	if permission.Module == "" {
		return nil, _errors.ReturnError(400, "Module không được để trống")
	}

	// Kiểm tra key đã tồn tại chưa
	existingPermission, err := uc.permissionRepo.GetByKey(ctx, permission.Key)
	if err == nil && existingPermission != nil {
		return nil, _errors.ReturnError(400, "Key permission đã tồn tại")
	}

	return uc.permissionRepo.Create(ctx, permission)
}

// GetByID lấy permission theo ID
func (uc *PermissionUsecase) GetByID(ctx context.Context, id uint64) (*domain.Permission, error) {
	permission, err := uc.permissionRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	if permission == nil {
		return nil, _errors.ReturnError(404, "Permission không tồn tại")
	}
	return permission, nil
}

// Update cập nhật permission
func (uc *PermissionUsecase) Update(ctx context.Context, id uint64, permission *domain.Permission) (*domain.Permission, error) {
	if permission.Name == "" {
		return nil, _errors.ReturnError(400, "Tên permission không được để trống")
	}
	if permission.Key == "" {
		return nil, _errors.ReturnError(400, "Key permission không được để trống")
	}
	if permission.Module == "" {
		return nil, _errors.ReturnError(400, "Module không được để trống")
	}

	// Kiểm tra permission có tồn tại không
	existingPermission, err := uc.permissionRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	if existingPermission == nil {
		return nil, _errors.ReturnError(404, "Permission không tồn tại")
	}

	permission.ID = id
	return uc.permissionRepo.Update(ctx, permission)
}

// Delete xóa permission
func (uc *PermissionUsecase) Delete(ctx context.Context, id uint64) error {
	// Kiểm tra permission có tồn tại không
	existingPermission, err := uc.permissionRepo.FindById(ctx, id)
	if err != nil {
		return err
	}
	if existingPermission == nil {
		return _errors.ReturnError(404, "Permission không tồn tại")
	}

	return uc.permissionRepo.Delete(ctx, id)
}

// GetAll lấy tất cả permission
func (uc *PermissionUsecase) GetAll(ctx context.Context) ([]*domain.Permission, error) {
	permissions, err := uc.permissionRepo.GetPermissions(ctx)
	if err != nil {
		return nil, err
	}

	permMap := make(map[uint64]*domain.Permission)
	for i := range permissions {
		permMap[permissions[i].ID] = permissions[i]
	}

	var roots []*domain.Permission

	// Xây dựng quan hệ cha-con
	for _, p := range permissions {
		permission := uc.BuildHierarchy(p, permMap)
		if permission.ParentID == nil {
			roots = append(roots, permission)
		}
	}

	return roots, nil
}

func (uc *PermissionUsecase) BuildHierarchy(permission *domain.Permission, permMap map[uint64]*domain.Permission) *domain.Permission {
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

// GetByKey lấy permission theo key
func (uc *PermissionUsecase) GetByKey(ctx context.Context, key string) (*domain.Permission, error) {
	return uc.permissionRepo.GetByKey(ctx, key)
}

// GetByOrganizationID lấy danh sách permission theo organization ID
func (uc *PermissionUsecase) GetByOrganizationID(ctx context.Context, organizationID uint64) ([]*domain.Permission, error) {
	return uc.permissionRepo.FindByOrganizationId(ctx, organizationID)
}

// GetByOrganizationIDWithPagination lấy danh sách permission theo organization ID với phân trang
func (uc *PermissionUsecase) GetByOrganizationIDWithPagination(ctx context.Context, organizationID uint64, page, size int) ([]*domain.Permission, uint64, error) {
	return uc.permissionRepo.FindByOrganizationIdWithPagination(ctx, organizationID, page, size)
}

// GetByModule lấy danh sách permission theo module
func (uc *PermissionUsecase) GetByModule(ctx context.Context, module string) ([]*domain.Permission, error) {
	if module == "" {
		return nil, _errors.ReturnError(400, "Module không được để trống")
	}
	return uc.permissionRepo.FindByModule(ctx, module)
}

// GetByModuleWithPagination lấy danh sách permission theo module với phân trang
func (uc *PermissionUsecase) GetByModuleWithPagination(ctx context.Context, module uint32, page, size int) ([]*domain.Permission, uint64, error) {
	if module <= 0 {
		return nil, 0, _errors.ReturnError(400, "Module không được để trống")
	}
	return uc.permissionRepo.FindByModuleWithPagination(ctx, module, page, size)
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

func (uc *PermissionUsecase) GetByDealId(ctx context.Context, dealId uint64) ([]string, *domain.Role, error) {
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
