// platform/errors/kind.go
package apperr

// Kind phân loại nguồn gốc lỗi — quyết định log level và response behavior.
//
// Tại sao uint8?
// - Nhỏ gọn, so sánh O(1)
// - Đủ cho 255 loại lỗi (thực tế chỉ cần 4)
// - Không dùng string để tránh typo: KindBussiness != KindBusiness
type Kind uint8

const (
	// KindValidation: lỗi do client gửi data sai — không log ở server (expected behavior)
	KindValidation Kind = iota + 1
	// KindBusiness: lỗi domain logic — log Info nếu cần audit
	KindBusiness
	// KindInternal: bug server — log Error + stack trace
	KindInternal
	// KindExternal: downstream fail — log Warn + retry info
	KindExternal
)

// String implement fmt.Stringer cho log output.
// Dùng trong LogFields: {"kind": "internal"} — dễ filter trên ELK/Loki.
func (k Kind) String() string {
	switch k {
	case KindValidation:
		return "validation"
	case KindBusiness:
		return "business"
	case KindInternal:
		return "internal"
	case KindExternal:
		return "external"
	default:
		return "unknown"
	}
}

// IsClientFault trả true khi lỗi do client — không cần alert ops.
// Dùng để quyết định có increment error counter hay không.
func (k Kind) IsClientFault() bool {
	return k == KindValidation || k == KindBusiness
}

// IsSystemFault trả true khi lỗi là server bug.
// Dùng để trigger alert/PagerDuty.
func (k Kind) IsSystemFault() bool {
	return k == KindInternal
}
