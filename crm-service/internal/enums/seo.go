package enums

// IsValidSEORefType kiểm tra refType nằm trong nhóm entity SEO được hỗ trợ.
// RefType=0 được xem là page SEO độc lập/manual, không gắn entity gốc.
func IsValidSEORefType(v uint32) bool {
	if v == 0 {
		return true
	}
	switch ESEORefType(v) {
	case ESEORefTypeAdmUnit,
		ESEORefTypeRegion,
		ESEORefTypeParcel,
		ESEORefTypeProject,
		ESEORefTypeDocument,
		ESEORefTypeMap,
		ESEORefTypeReport:
		return true
	default:
		return false
	}
}

func SEORefTypeKey(v uint32) string {
	switch ESEORefType(v) {
	case 0:
		return "manual"
	case ESEORefTypeAdmUnit:
		return "adm_unit"
	case ESEORefTypeRegion:
		return "region"
	case ESEORefTypeParcel:
		return "parcel"
	case ESEORefTypeProject:
		return "project"
	case ESEORefTypeDocument:
		return "document"
	case ESEORefTypeMap:
		return "map"
	case ESEORefTypeReport:
		return "report"
	default:
		return "unknown"
	}
}
