package enums

type EMiningScope uint32

const (
	EMiningScopeRoom     EMiningScope = 10
	EMiningScopeFloor    EMiningScope = 20
	EMiningScopeArea     EMiningScope = 30
	EMiningScopePosition EMiningScope = 40
	EMiningScopeOther    EMiningScope = 50
)

var MiningScopeNames = map[EMiningScope]string{
	EMiningScopeRoom:     "Phòng",
	EMiningScopeFloor:    "Tầng",
	EMiningScopeArea:     "Khu vực",
	EMiningScopePosition: "Vị trí",
	EMiningScopeOther:    "Khác",
}
