package entity

import "organization/internal/enums"

type OrganizationRole struct {
	Id             uint32
	Name           string
	Key            string
	PermissionIds  []uint64
	OrganizationId uint32
	Permissions    []*OrganizationPermission
	IsDefault      bool
	DomainType     enums.DomainType
	ColorId        *uint32
	Color          *Color
	RoleKey        uint32
}
