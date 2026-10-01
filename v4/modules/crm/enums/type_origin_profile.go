package enums

type EOriginProfileType uint32

const (
	OriginProfileUser         EOriginProfileType = 10
	OriginProfileGroup        EOriginProfileType = 20
	OriginProfileOrganization EOriginProfileType = 30
)

var OriginProfileTypeNames = map[EOriginProfileType]string{
	OriginProfileUser:         "Cá nhân",
	OriginProfileGroup:        "Nhóm",
	OriginProfileOrganization: "Tổ chức",
}

func (e EOriginProfileType) IsValid() bool {
	switch e {
	case OriginProfileUser, OriginProfileGroup, OriginProfileOrganization:
		return true
	default:
		return false
	}
}