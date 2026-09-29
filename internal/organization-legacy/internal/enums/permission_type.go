package enums

type PermissionType uint

const (
	PermissionTypeOrganization PermissionType = 10
	PermissionTypeBranch       PermissionType = 20
)

var PermissionTypeMap = map[PermissionType]string{
	PermissionTypeOrganization: "Tổ chức",
	PermissionTypeBranch:       "Chi nhánh",
}
