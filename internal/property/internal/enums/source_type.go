package enums

type ESourceType int

const (
	SourceTypeOwner  ESourceType = 10
	SourceTypeBroker ESourceType = 20
	SourceTypeFloor  ESourceType = 30
	SourceTypeOther  ESourceType = 40
)

var SourceTypeMap = map[ESourceType]string{
	SourceTypeOwner:  "Chủ đất",
	SourceTypeBroker: "Môi giới",
	SourceTypeFloor:  "Sàn",
	SourceTypeOther:  "Khác",
}
