package enums

type ESEORefType uint32

const (
	ESEORefTypeAdmUnit  ESEORefType = 100
	ESEORefTypeParcel   ESEORefType = 200
	ESEORefTypeRegion   ESEORefType = 300
	ESEORefTypeProject  ESEORefType = 400
	ESEORefTypeDocument ESEORefType = 500
	ESEORefTypeMap      ESEORefType = 600
	ESEORefTypeReport   ESEORefType = 700
)

func (e ESEORefType) Validate() bool {
	switch e {
	case ESEORefTypeAdmUnit, ESEORefTypeParcel, ESEORefTypeRegion, ESEORefTypeProject, ESEORefTypeDocument, ESEORefTypeMap, ESEORefTypeReport:
		return true
	default:
		return false
	}
}
