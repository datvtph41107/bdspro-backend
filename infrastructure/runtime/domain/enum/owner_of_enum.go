package _enum

type EOwnerOf int

const (
	EOwnerOfMember      EOwnerOf = 10
	EOwnerOfGroup       EOwnerOf = 20
	EOwnerOfOrgnization EOwnerOf = 30
	EOwnerOfAdmin       EOwnerOf = 40
)

var OwnerOfNames = map[EOwnerOf]string{
	EOwnerOfMember:      "Thành viên",
	EOwnerOfGroup:       "Nhóm",
	EOwnerOfOrgnization: "Tổ chức",
	EOwnerOfAdmin:       "Admin",
}

func (e EOwnerOf) IsValid() bool {
	switch e {
	case EOwnerOfMember, EOwnerOfGroup, EOwnerOfOrgnization, EOwnerOfAdmin:
		return true
	default:
		return false
	}
}
