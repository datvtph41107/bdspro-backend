package repo

import (
	"common/case/crud"
	"context"

	"auth/internal/domain"
)

type RoleRepository interface {
	crud.ICrudRepo[domain.Role]
	ExistProfileIdWithKeys(ctx context.Context, profileId uint64, keys []string) (bool, error)
	FindByDomain(ctx context.Context, domain string) ([]*domain.Role, error)
	FindByOrganizationId(ctx context.Context, organizationId uint64) ([]*domain.Role, error)
	FindByOrganizationIdAndKey(ctx context.Context, organizationId uint64, key string) (*domain.Role, error)
	FindByOrganizationIdWithPagination(ctx context.Context, organizationId uint64, page, size int) ([]*domain.Role, uint64, error)
	FindByRoleKeys(ctx context.Context, roleKeys []uint32) ([]domain.Role, error)
	AddPermissionToRole(ctx context.Context, roleID uint64, permissionID uint64) error
	GetRoleIdsByProfileId(ctx context.Context, profileId uint64) ([]uint64, error)
	GetPermissionIdsByRoleIds(ctx context.Context, roleIds []uint64) (map[uint64][]uint64, error)
	GetPermissionKeysByIds(ctx context.Context, permIds []uint64) ([]string, error)
	GetPermissionKeysByUserId(ctx context.Context, userId uint64) ([]string, error)
	GetPermissionKeysByRoleId(ctx context.Context, roleId uint64) ([]string, error)
	GetPermissionKeysByRoleKey(ctx context.Context, roleKey uint32) ([]string, error)
	GetAdminRoleIdByGroupKey(ctx context.Context, groupKey uint64) (*domain.Role, error)
	GetListByGroupKey(ctx context.Context, groupKey uint32) ([]domain.Role, int64, error)
	GetRolesByRoleKeysAndOrganization(ctx context.Context, roleKeys []uint32, organizationId uint64) ([]*domain.Role, error)
	AssignRoleToUser(ctx context.Context, userId uint64, role *domain.Role) error
	GetRoleByRoleKey(ctx context.Context, roleKey uint32) (*domain.Role, error)
}
