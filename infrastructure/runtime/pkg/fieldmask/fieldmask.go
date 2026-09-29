// Package fieldmask giải quyết bài toán:
// "Field này có được phép update trong request này không?"
// # Bản chất
// FieldMask là authorization layer — client tường minh khai báo field nào
// được update. Tách biệt hoàn toàn với:
//   - Presence detection: proto3 optional → *T (wire layer)
//   - Format validity:    pkg/validate (handler layer)
//   - Diff check:         pkg/patch (usecase layer)
//
// # Zero value
//
// Zero value (FieldMask{}) = AllowAll.
// Lý do: internal calls không gửi mask → phải allow all.
// Khi paths != nil → strict whitelist mode.
//
// # Path convention
//
// snake_case, dot-separated, align với proto field name:
//
//	"info.title"                   → UpdateInfoPayload.title
//	"landInfo.area.area_total"    → UpdateLandInfoPayload.area.area_total
//	"lineage_meta.national_id"     → UpdateLineageMetaPayload.national_id
//
// # Thread safety
//
// FieldMask là value type bất biến sau khi khởi tạo — safe for concurrent use.
package fieldmask

import (
	"sort"
	"strings"
)

// FieldMask là value type bất biến.
// Zero value = AllowAll (paths == nil).
// Khởi tạo 1 lần tại Handler, truyền xuống qua DTO, tiêu thụ tại patch.Builder.
type FieldMask struct {
	// nil  → allowAll (zero value)
	// !nil → strict whitelist, key = exact path
	paths map[string]struct{}
}

// FromPaths tạo FieldMask từ proto paths.
//
//	nil / rỗng → AllowAll (zero value, backward compat với client cũ)
//	có paths   → strict whitelist
func FromPaths(paths []string) FieldMask {
	if len(paths) == 0 {
		return FieldMask{} // zero value = allowAll
	}
	m := make(map[string]struct{}, len(paths))
	for _, p := range paths {
		if p != "" {
			m[p] = struct{}{}
		}
	}
	if len(m) == 0 {
		return FieldMask{}
	}
	return FieldMask{paths: m}
}

// Only tạo FieldMask với danh sách paths tường minh.
// Dùng trong code (không phải từ proto) — tên gọi rõ ý định hơn FromPaths.
//
//	mask := fieldmask.Only("info.title", "info.note")
func Only(paths ...string) FieldMask {
	return FromPaths(paths)
}

// Allows kiểm tra path có được phép update không.
//
// Hỗ trợ 2 kiểu match:
//  1. Exact:  paths=["info.title"]     → Allows("info.title") = true
//  2. Prefix: paths=["landInfo.area"] → Allows("landInfo.area.area_total") = true
//
// Prefix match cho phép allow cả sub-group với 1 path ngắn hơn.
func (f FieldMask) Allows(path string) bool {
	if f.paths == nil {
		return true // zero value = allowAll
	}
	if _, ok := f.paths[path]; ok {
		return true
	}
	// prefix match: "landInfo.area" covers "landInfo.area.area_total"
	for p := range f.paths {
		if strings.HasPrefix(path, p+".") {
			return true
		}
	}
	return false
}

// AllowsAny trả về true nếu bất kỳ path nào được phép.
// Dùng để kiểm tra nhanh "sub-domain này có cần xử lý không?".
func (f FieldMask) AllowsAny(paths ...string) bool {
	for _, p := range paths {
		if f.Allows(p) {
			return true
		}
	}
	return false
}

// IsAllowAll trả về true nếu mask cho phép tất cả.
func (f FieldMask) IsAllowAll() bool { return f.paths == nil }

// Intersect tạo mask con — chỉ giữ paths nằm trong allowed list.
//
// Dùng ở gateway để strip quyền:
//
//	userMask  = ["info.title", "lineage_meta.national_id"]
//	allowed   = ["info.title", "info.note"]  // user thường không có lineage_meta
//	→ result  = ["info.title"]
func (f FieldMask) Intersect(allowed []string) FieldMask {
	if f.paths == nil {
		// allowAll intersect allowed = allowed
		return FromPaths(allowed)
	}
	result := make([]string, 0, len(allowed))
	for _, a := range allowed {
		if f.Allows(a) {
			result = append(result, a)
		}
	}
	// FromPaths([]) intentionally means AllowAll for backwards compatibility.
	// An intersection with no common path is instead a strict deny-all mask.
	if len(result) == 0 {
		return FieldMask{paths: map[string]struct{}{}}
	}
	return FromPaths(result)
}

// Paths trả về sorted paths — deterministic, dùng cho logging/debugging/testing.
//
// Trả về ["*"] nếu allowAll.
func (f FieldMask) Paths() []string {
	if f.paths == nil {
		return []string{"*"}
	}
	result := make([]string, 0, len(f.paths))
	for p := range f.paths {
		result = append(result, p)
	}
	sort.Strings(result) // deterministic order
	return result
}
