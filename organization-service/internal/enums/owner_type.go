package enums

type OwnerType uint

const (
	OwnerTypeUser         OwnerType = 10
	OwnerTypeGroup        OwnerType = 20
	OwnerTypeOrganization OwnerType = 30
)

var OwnerTypeMap = map[OwnerType]string{
	OwnerTypeUser:         "Người dùng",
	OwnerTypeGroup:        "Nhóm",
	OwnerTypeOrganization: "Công ty",
}
