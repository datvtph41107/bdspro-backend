package enums

// ClickURLDomainType loại domain khi track click URL.
type ClickURLDomainType uint32

const (
	ClickURLDomainTypeUnspecified ClickURLDomainType = 0
	ClickURLDomainTypeAdmUnit     ClickURLDomainType = 100
	ClickURLDomainTypeParcel      ClickURLDomainType = 200
	ClickURLDomainTypeRegion      ClickURLDomainType = 300
	ClickURLDomainTypeProject     ClickURLDomainType = 400
	ClickURLDomainTypeDocument    ClickURLDomainType = 500
	ClickURLDomainTypeMap         ClickURLDomainType = 600
	ClickURLDomainTypeReport      ClickURLDomainType = 700
	ClickURLDomainTypeLegacy      ClickURLDomainType = 800
)

func (e ClickURLDomainType) Validate() bool {
	switch e {
	case ClickURLDomainTypeAdmUnit,
		ClickURLDomainTypeParcel,
		ClickURLDomainTypeRegion,
		ClickURLDomainTypeProject,
		ClickURLDomainTypeDocument,
		ClickURLDomainTypeMap,
		ClickURLDomainTypeReport,
		ClickURLDomainTypeLegacy:
		return true
	default:
		return false
	}
}
