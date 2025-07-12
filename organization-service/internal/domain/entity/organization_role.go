package entity

import "organization/internal/enums"

type OrganizationRole struct {
	Id             uint32
	Name           string
	Key            string
	PermissionIds  []uint32
	OrganizationId uint32
	Permissions    []*OrganizationPermission
	IsDefault      bool
	DomainType     enums.DomainType
}
