package enums

type EViewRole int

const (
	ViewRoleGuest  EViewRole = 10 // Khách
	ViewRoleUser   EViewRole = 20 // Người dùng
	ViewRoleFriend EViewRole = 30 // Bạn bè
	ViewRoleOwner  EViewRole = 40 // Chủ hồ sơ
)

var ViewRoleMap = map[EViewRole]string{
	ViewRoleGuest:  "Khách",
	ViewRoleUser:   "Người dùng",
	ViewRoleFriend: "Bạn bè",
	ViewRoleOwner:  "Chủ hồ sơ",
}

func (e EViewRole) IsValid() bool {
	switch e {
	case ViewRoleGuest, ViewRoleUser, ViewRoleFriend, ViewRoleOwner:
		return true
	default:
		return false
	}
}
