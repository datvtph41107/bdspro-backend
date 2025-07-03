package enums

type EOwnerType = uint32

const (
	EOwnerTypeMember      EOwnerType = 10
	EOwnerTypeGroup       EOwnerType = 20
	EOwnerTypeOrgnization EOwnerType = 30
)

var OwnerTypeNames = map[EOwnerType]string{
	EOwnerTypeMember:      "Thành viên",
	EOwnerTypeGroup:       "Nhóm",
	EOwnerTypeOrgnization: "Tổ chức",
}
