package entity

type OrganizationRole struct {
	Id   uint32
	Name string
	Key  string
	PermissionIds []uint32
	OrganizationId uint32
	Permissions []*OrganizationPermission
}
