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
	ESEORefTypeLegacy   ESEORefType = 800
)

func (e ESEORefType) Validate() bool {
	switch e {
	case ESEORefTypeAdmUnit, ESEORefTypeParcel,
		ESEORefTypeRegion, ESEORefTypeProject,
		ESEORefTypeDocument, ESEORefTypeMap, ESEORefTypeReport,
		ESEORefTypeLegacy:
		return true
	default:
		return false
	}
}

// Name returns the display label used by public/admin SEO clients.
func (e ESEORefType) Name() string {
	switch e {
	case ESEORefTypeAdmUnit:
		return "Đơn vị hành chính"
	case ESEORefTypeParcel:
		return "Thửa đất"
	case ESEORefTypeRegion:
		return "Vùng"
	case ESEORefTypeProject:
		return "Đồ án"
	case ESEORefTypeDocument:
		return "Tài liệu"
	case ESEORefTypeMap:
		return "Bản đồ"
	case ESEORefTypeReport:
		return "Báo cáo"
	case ESEORefTypeLegacy:
		return "Nội dung thủ công"
	default:
		return ""
	}
}

// AllSEORefTypes returns the 8 canonical SEO ref types in display order.
func AllSEORefTypes() []ESEORefType {
	return []ESEORefType{
		ESEORefTypeAdmUnit,
		ESEORefTypeParcel,
		ESEORefTypeRegion,
		ESEORefTypeProject,
		ESEORefTypeDocument,
		ESEORefTypeMap,
		ESEORefTypeReport,
		ESEORefTypeLegacy,
	}
}
