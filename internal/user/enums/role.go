package enums

type ERole int

const (
	ERoleUser       ERole = 10
	ERoleAdmin      ERole = 20
	ERoleSuperAdmin ERole = 30
)

var RoleMap = map[ERole]string{
	ERoleUser:       "Người dùng",
	ERoleAdmin:      "Quản trị viên",
	ERoleSuperAdmin: "Quản trị hệ thống",
}
