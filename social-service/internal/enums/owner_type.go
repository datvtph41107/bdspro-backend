package enums

type OwnerType int32

const (
	OwnerTypeUser  OwnerType = 10
	OwnerTypeGroup OwnerType = 20
)

var OwnerTypeNames = map[OwnerType]string{
	OwnerTypeUser:  "user",
	OwnerTypeGroup: "group",
}
