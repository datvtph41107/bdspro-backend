package enums

type ESeoClassify uint32

const (
	// SEO-A: detail page phục vụ long-tail và conversion.
	ESeoClassifyA ESeoClassify = 10
	// SEO-B: guide/category/hub page phục vụ điều hướng và phân bổ authority.
	ESeoClassifyB ESeoClassify = 20
	// SEO-C: search/faceted page cần canonical/dedup rõ.
	ESeoClassifyC ESeoClassify = 30
	// SEO-D: share/preview/viral page, không phải page SEO chính.
	ESeoClassifyD ESeoClassify = 40
	// SEO-N: không phục vụ SEO như login/admin/404/private.
	ESeoClassifyN ESeoClassify = 50
)

func IsValidSeoClassify(v uint32) bool {
	switch ESeoClassify(v) {
	case ESeoClassifyA, ESeoClassifyB, ESeoClassifyC, ESeoClassifyD, ESeoClassifyN:
		return true
	default:
		return false
	}
}

func SeoClassifyFromUint32(v uint32) ESeoClassify {
	if IsValidSeoClassify(v) {
		return ESeoClassify(v)
	}
	return ESeoClassifyN
}

func (e ESeoClassify) Uint32() uint32 {
	return uint32(e)
}

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