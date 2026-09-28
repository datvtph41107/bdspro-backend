package enums

type ESourceType int

const (
	SourceTypeOwner  ESourceType = 10
	SourceTypeBroker ESourceType = 20
	SourceTypeFloor  ESourceType = 30
	SourceTypeOther  ESourceType = 40

	SourceTypeCustomer ESourceType = 50
)

var SourceTypeMap = map[ESourceType]string{
	SourceTypeOwner:  "Chủ nhà trực tiếp",
	SourceTypeBroker: "Môi giới đối tác",
	SourceTypeFloor:  "Sàn nội bộ",
	SourceTypeOther:  "Khác",

	SourceTypeCustomer: "Khách giới thiệu",
}
