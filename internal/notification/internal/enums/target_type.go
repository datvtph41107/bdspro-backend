package enums

type OwnerType int32

const (
	OwnerTypeUser         OwnerType = 10 // User
	OwnerTypeGroup        OwnerType = 20 // Group
	OwnerTypeOrganization OwnerType = 30 // Organization
)

var OwnerTypeMap = map[OwnerType]string{
	OwnerTypeUser:         "Người dùng",
	OwnerTypeGroup:        "Nhóm",
	OwnerTypeOrganization: "Tổ chức",
}
