package enums

type EMiningMode uint32

const (
	EMiningModeAll  EMiningMode = 10
	EMiningModePart EMiningMode = 20
)

var MiningModeNames = map[EMiningMode]string{
	EMiningModeAll:  "Toàn bộ BĐS",
	EMiningModePart: "Một phần BĐS",
}
