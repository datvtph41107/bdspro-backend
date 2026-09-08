package shared_enum

type EOwnerType = int

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
