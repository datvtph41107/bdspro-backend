package base_enum

type EOwnerOf int

const (
	EOwnerOfMember      EOwnerOf = 10
	EOwnerOfGroup       EOwnerOf = 20
	EOwnerOfOrgnization EOwnerOf = 30
)

var OwnerOfNames = map[EOwnerOf]string{
	EOwnerOfMember:      "Thành viên",
	EOwnerOfGroup:       "Nhóm",
	EOwnerOfOrgnization: "Tổ chức",
}

func (e EOwnerOf) IsValid() bool {
	switch e {
	case EOwnerOfMember, EOwnerOfGroup, EOwnerOfOrgnization:
		return true
	default:
		return false
	}
}
