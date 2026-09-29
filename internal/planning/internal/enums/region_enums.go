package enums

// LabelStatus - Trạng thái của label
type LabelStatus uint32

const (
	LabelStatusActive   LabelStatus = 10 // Đang sử dụng
	LabelStatusInactive LabelStatus = 20 // Tạm dừng
	LabelStatusDeleted  LabelStatus = 30 // Đã xóa (soft delete)
)

func (l LabelStatus) IsValid() bool {
	return l == LabelStatusActive || l == LabelStatusInactive || l == LabelStatusDeleted
}

func (l LabelStatus) String() string {
	switch l {
	case LabelStatusActive:
		return "active"
	case LabelStatusInactive:
		return "inactive"
	case LabelStatusDeleted:
		return "deleted"
	default:
		return "unknown"
	}
}

// RegionStatus - Trạng thái của region
type RegionStatus uint32

const (
	RegionStatusActive   RegionStatus = 10 // Đang hiệu lực
	RegionStatusInactive RegionStatus = 20 // Hết hiệu lực
	RegionStatusPending  RegionStatus = 30 // Chờ duyệt
	RegionStatusRejected RegionStatus = 40 // Từ chối
)

func (r RegionStatus) IsValid() bool {
	return r == RegionStatusActive || r == RegionStatusInactive ||
		r == RegionStatusPending || r == RegionStatusRejected
}

func (r RegionStatus) String() string {
	switch r {
	case RegionStatusActive:
		return "active"
	case RegionStatusInactive:
		return "inactive"
	case RegionStatusPending:
		return "pending"
	case RegionStatusRejected:
		return "rejected"
	default:
		return "unknown"
	}
}

// ProcessingStatus - Trạng thái xử lý import của region
type ProcessingStatus uint32

const (
	ProcessingStatusPending  ProcessingStatus = 10 // Chờ xử lý (mới import)
	ProcessingStatusMapped   ProcessingStatus = 20 // Đã map label (chờ duyệt)
	ProcessingStatusVerified ProcessingStatus = 30 // Đã xác nhận đúng
	ProcessingStatusRejected ProcessingStatus = 40 // Bị từ chối, cần sửa
)

func (p ProcessingStatus) IsValid() bool {
	return p == ProcessingStatusPending || p == ProcessingStatusMapped ||
		p == ProcessingStatusVerified || p == ProcessingStatusRejected
}

func (p ProcessingStatus) String() string {
	switch p {
	case ProcessingStatusPending:
		return "pending"
	case ProcessingStatusMapped:
		return "mapped"
	case ProcessingStatusVerified:
		return "verified"
	case ProcessingStatusRejected:
		return "rejected"
	default:
		return "unknown"
	}
}
