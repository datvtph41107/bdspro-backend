package enums

type ActorType uint32

const (
	ActorUser   ActorType = 1
	ActorAdmin  ActorType = 2
	ActorSystem ActorType = 3
)

var ActorTypeText = map[ActorType]string{
	ActorUser:   "Người dùng",
	ActorAdmin:  "Quản trị viên",
	ActorSystem: "Quản trị hệ thống",
}

func (a ActorType) String() string {
	if text, ok := ActorTypeText[a]; ok {
		return text
	}
	return "Unknown"
}
