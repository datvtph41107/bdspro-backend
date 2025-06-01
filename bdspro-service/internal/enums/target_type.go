package enums

type ETargetType = int

const (
	ETargetTypeMember       ETargetType = 10
	ETargetTypeGroup        ETargetType = 20
	ETargetTypeOrganization ETargetType = 30
)

var TargetTypeNames = map[ETargetType]string{
	ETargetTypeMember:       "Thành viên",
	ETargetTypeGroup:        "Nhóm",
	ETargetTypeOrganization: "Tổ chức",
}
