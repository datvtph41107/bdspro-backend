package enums

// WarningTypeEnum định nghĩa các loại cảnh báo
type WarningTypeEnum string

const (
	// Cảnh báo vi phạm
	WarningTypeViolation WarningTypeEnum = "violation"

	// Nhắc nhở chung
	WarningTypeReminder WarningTypeEnum = "reminder"

	// Hướng dẫn sử dụng
	WarningTypeGuide WarningTypeEnum = "guide"

	// Cảnh báo bảo mật
	WarningTypeSecurity WarningTypeEnum = "security"

	// Thông báo hệ thống
	WarningTypeSystem WarningTypeEnum = "system"
)

// SeverityEnum định nghĩa mức độ nghiêm trọng
type SeverityEnum string

const (
	SeverityLow      SeverityEnum = "low"
	SeverityMedium   SeverityEnum = "medium"
	SeverityHigh     SeverityEnum = "high"
	SeverityCritical SeverityEnum = "critical"
)

// TargetTypeEnum định nghĩa loại đối tượng nhận cảnh báo
type TargetTypeEnum string

const (
	TargetTypeUser         TargetTypeEnum = "user"
	TargetTypeOrganization TargetTypeEnum = "organization"
	TargetTypeGroup        TargetTypeEnum = "group"
)

// WarningStatusEnum định nghĩa trạng thái cảnh báo
type WarningStatusEnum string

const (
	WarningStatusSent         WarningStatusEnum = "sent"
	WarningStatusRead         WarningStatusEnum = "read"
	WarningStatusAcknowledged WarningStatusEnum = "acknowledged"
	WarningStatusExpired      WarningStatusEnum = "expired"
)

// WarningActionEnum định nghĩa các hành động với cảnh báo
type WarningActionEnum string

const (
	WarningActionSent         WarningActionEnum = "sent"
	WarningActionRead         WarningActionEnum = "read"
	WarningActionAcknowledged WarningActionEnum = "acknowledged"
	WarningActionDeleted      WarningActionEnum = "deleted"
)

// IsValid kiểm tra WarningType có hợp lệ không
func (w WarningTypeEnum) IsValid() bool {
	switch w {
	case WarningTypeViolation, WarningTypeReminder, WarningTypeGuide, WarningTypeSecurity, WarningTypeSystem:
		return true
	default:
		return false
	}
}

// IsValid kiểm tra Severity có hợp lệ không
func (s SeverityEnum) IsValid() bool {
	switch s {
	case SeverityLow, SeverityMedium, SeverityHigh, SeverityCritical:
		return true
	default:
		return false
	}
}

// IsValid kiểm tra TargetType có hợp lệ không
func (t TargetTypeEnum) IsValid() bool {
	switch t {
	case TargetTypeUser, TargetTypeOrganization, TargetTypeGroup:
		return true
	default:
		return false
	}
}

// IsValid kiểm tra WarningStatus có hợp lệ không
func (w WarningStatusEnum) IsValid() bool {
	switch w {
	case WarningStatusSent, WarningStatusRead, WarningStatusAcknowledged, WarningStatusExpired:
		return true
	default:
		return false
	}
}

// GetWarningTypeMap trả về map mô tả các loại cảnh báo
func GetWarningTypeMap() map[WarningTypeEnum]string {
	return map[WarningTypeEnum]string{
		WarningTypeViolation: "Cảnh báo vi phạm",
		WarningTypeReminder:  "Nhắc nhở chung",
		WarningTypeGuide:     "Hướng dẫn sử dụng",
		WarningTypeSecurity:  "Cảnh báo bảo mật",
		WarningTypeSystem:    "Thông báo hệ thống",
	}
}

// GetSeverityMap trả về map mô tả các mức độ nghiêm trọng
func GetSeverityMap() map[SeverityEnum]string {
	return map[SeverityEnum]string{
		SeverityLow:      "Thấp",
		SeverityMedium:   "Trung bình",
		SeverityHigh:     "Cao",
		SeverityCritical: "Nghiêm trọng",
	}
}
