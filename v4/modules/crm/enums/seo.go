//go:build ignore

package enums

// IsValidSEORefType kiểm tra refType nằm trong nhóm entity SEO được hỗ trợ.
// RefType=0 được xem là alias cũ của nội dung thủ công (800).
func IsValidSEORefType(v uint32) bool {
	if IsManualSEORefType(v) {
		return true
	}
	switch ESEORefType(v) {
	case ESEORefTypeAdmUnit,
		ESEORefTypeRegion,
		ESEORefTypeParcel,
		ESEORefTypeProject,
		ESEORefTypeDocument,
		ESEORefTypeMap,
		ESEORefTypeReport,
		ESEORefTypeLegacy:
		return true
	default:
		return false
	}
}

// IsManualSEORefType là trang nội dung thủ công, không gắn entity nguồn.
func IsManualSEORefType(v uint32) bool {
	return v == 0 || ESEORefType(v) == ESEORefTypeLegacy
}

// CanonicalSEORefType chuẩn hóa alias cũ 0 thành 800.
func CanonicalSEORefType(v uint32) uint32 {
	if v == 0 {
		return uint32(ESEORefTypeLegacy)
	}
	return v
}

// ExpandSEORefTypes chuẩn hóa danh sách lọc, gồm cả alias 0 của nội dung thủ công.
func ExpandSEORefTypes(values []uint32) []uint32 {
	if len(values) == 0 {
		return nil
	}
	out := make([]uint32, 0, len(values)+1)
	seen := make(map[uint32]struct{}, len(values)+1)
	includeManual := false
	for _, value := range values {
		if IsManualSEORefType(value) {
			includeManual = true
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	if includeManual {
		for _, value := range []uint32{0, uint32(ESEORefTypeLegacy)} {
			if _, ok := seen[value]; ok {
				continue
			}
			seen[value] = struct{}{}
			out = append(out, value)
		}
	}
	return out
}

func SEORefTypeKey(v uint32) string {
	switch ESEORefType(CanonicalSEORefType(v)) {
	case ESEORefTypeLegacy:
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
