package enums

type InviteType uint32

const (
	INVITE_TYPE_NONE    InviteType = 0
	INVITE_TYPE_REFERAL InviteType = 1
	INVITE_TYPE_DIRECT  InviteType = 2
)

func (r InviteType) String() string {
	switch r {
	case INVITE_TYPE_NONE:
		return "NONE"
	case INVITE_TYPE_REFERAL:
		return "REFERAL"
	case INVITE_TYPE_DIRECT:
		return "DIRECT"
	default:
		return "unknown"
	}
}
