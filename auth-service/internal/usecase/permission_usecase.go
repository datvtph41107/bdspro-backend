package usecase

import (
	"auth/internal/domain"
	"auth/internal/interface/repo"
	_errors "common/errors"
	"context"
)

type PermissionUsecase struct {
	permissionRepo repo.PermissionRepository
}

func NewPermissionUsecase(permissionRepo repo.PermissionRepository) *PermissionUsecase {
	return &PermissionUsecase{
		permissionRepo: permissionRepo,
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
	for i := range permissions {
		p := permissions[i]
		if p.ParentID == nil {
			roots = append(roots, p)
		} else if parent, ok := permMap[*p.ParentID]; ok {
			parent.Children = append(parent.Children, p)
		}
	}

	return roots, nil
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
