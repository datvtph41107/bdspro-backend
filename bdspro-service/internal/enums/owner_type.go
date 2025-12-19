package enums

type EOwnerOf = uint32

const (
	EOwnerOfMember       EOwnerOf = 10
	EOwnerOfGroup        EOwnerOf = 20
	EOwnerOfOrganization EOwnerOf = 30
	EOwnerOfUser         EOwnerOf = 40
)

var OwnerTypeNames = map[EOwnerOf]string{
	EOwnerOfMember:       "Thành viên",
	EOwnerOfGroup:        "Nhóm",
	EOwnerOfOrganization: "Tổ chức",
	EOwnerOfUser:         "Người dùng",
}
var OwnerTypeMap = map[EOwnerOf]string{
	EOwnerOfMember:       "Thành viên",
	EOwnerOfGroup:        "Nhóm",
	EOwnerOfOrganization: "Tổ chức",
	EOwnerOfUser:         "Người dùng",
}
