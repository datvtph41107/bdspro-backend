package enums

type Role uint32

const (
	ROLE_SUPER_ADMIN Role = 1
	ROLE_ADMIN       Role = 2
	ROLE_MEMBER      Role = 3
)

func (r Role) String() string {
	switch r {
	case ROLE_SUPER_ADMIN:
		return "SUPER_ADMIN"
	case ROLE_ADMIN:
		return "ADMIN"
	case ROLE_MEMBER:
		return "MEMBER"
	default:
		return "unknown"
	}
}

func (r Role) IsValid() bool {
	switch r {
	case ROLE_SUPER_ADMIN, ROLE_MEMBER, ROLE_ADMIN:
		return true
	default:
		return false
	}
}
