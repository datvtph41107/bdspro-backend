package repo

import (
	"common/case/crud"
	"context"

	"user/internal/domain/access"

	"user/internal/domain/auth"
)

type RoleRepository interface {
	crud.ICrudRepo[access.Role]
	ExistProfileIdWithKeys(ctx context.Context, profileId uint64, keys []string) (bool, error)
	FindByDomain(ctx context.Context, domain string) ([]*access.Role, error)
	FindByOrganizationId(ctx context.Context, organizationId uint64) ([]*access.Role, error)
	FindByOrganizationIdAndKey(ctx context.Context, organizationId uint64, key string) (*access.Role, error)
	FindByOrganizationIdWithPagination(ctx context.Context, organizationId uint64, page, size int) ([]*access.Role, uint64, error)
	FindByRoleKeys(ctx context.Context, roleKeys []uint32) ([]access.Role, error)
	GetRoleIdsByProfileId(ctx context.Context, profileId uint64) ([]uint64, error)
	GetPermissionIdsByRoleIds(ctx context.Context, roleIds []uint64) (map[uint64][]uint64, error)
	GetPermissionKeysByIds(ctx context.Context, permIds []uint64) ([]string, error)
	GetPermissionKeysByUserId(ctx context.Context, userId uint64) ([]string, error)
	GetPermissionKeysByRoleId(ctx context.Context, roleId uint64) ([]string, error)
	GetPermissionKeysByRoleKey(ctx context.Context, roleKey uint32) ([]string, error)
	GetAdminRoleIdByGroupKey(ctx context.Context, groupKey uint64) (*access.Role, error)
	GetListByGroupKey(ctx context.Context, groupKey uint32) ([]access.Role, int64, error)
	GetRolesByRoleKeysAndOrganization(ctx context.Context, roleKeys []uint32, organizationId uint64) ([]*access.Role, error)
	AssignRoleToUser(ctx context.Context, userId uint64, role *access.Role) error
	GetRoleByRoleKey(ctx context.Context, roleKey uint32) (*access.Role, error)
	GetAuthMethodsByIds(ctx context.Context, userIds []uint64) ([]*auth.AuthMethod, error)
	GetAuthAdminsWithRolesByIds(ctx context.Context, ids []uint64) ([]*access.AuthAdminWithRole, error)
	GetRoleByGroupId(ctx context.Context, groupId uint64) ([]access.Role, error)
}
