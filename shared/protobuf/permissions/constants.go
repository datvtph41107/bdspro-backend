package permissions

type PermissionKey string

const (
	ProfileRead   PermissionKey = "profile_read"
	ProfileWrite  PermissionKey = "profile_write"
	ProfileDelete PermissionKey = "profile_delete"
	ProfileCreate PermissionKey = "profile_create"
	ProfileUpdate PermissionKey = "profile_update"
)

var PermissionKeys = []PermissionKey{
	ProfileRead,
	ProfileWrite,
	ProfileDelete,
	ProfileCreate,
	ProfileUpdate,
}
// var PermissionKeys = map[PermissionKey]string{
// 	ProfileRead:   "profile_read",
// 	ProfileWrite:  "profile_write",
// 	ProfileDelete: "profile_delete",
// 	ProfileCreate: "profile_create",
// 	ProfileUpdate: "profile_update",
// }
