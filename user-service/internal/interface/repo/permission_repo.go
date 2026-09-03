package repo

import (
	"context"
	"user/internal/domain/access"
)

// PermissionRepository exposes the permission reads owned by User IAM.
// Organization-scoped permission management belongs to Organization Service.
type PermissionRepository interface {
	GetPermissions(ctx context.Context) ([]*access.Permission, error)
	FindById(ctx context.Context, id uint64) (*access.Permission, error)
	FindByModuleWithPagination(ctx context.Context, module string, page, size int) ([]*access.Permission, uint64, error)
}
