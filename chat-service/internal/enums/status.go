package enums

type Status uint32

const (
	StatusUnknown   Status = 0
	STATUS_ACTIVE   Status = 1
	STATUS_INACTIVE Status = 2
	STATUS_BANNED   Status = 3
)

func (s Status) String() string {
	switch s {
	case STATUS_ACTIVE:
		return "active"
	case STATUS_INACTIVE:
		return "inactive"
	case STATUS_BANNED:
		return "banned"
	default:
		return "unknown"
	}
}
