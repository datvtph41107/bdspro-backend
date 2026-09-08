package enums

type EPositionType uint32

const (
	EPositionTypeRelative EPositionType = 10
	EPositionTypeAbsolute EPositionType = 20
)

var PositionTypeNames = map[EPositionType]string{
	EPositionTypeRelative: "Tương đối",
	EPositionTypeAbsolute: "Tuyệt đối",
}
