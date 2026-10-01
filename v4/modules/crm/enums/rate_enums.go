package enums

type TargetType uint32

const (
	// TargetType - Loại đối tượng được đánh giá
	TargetTypeUser         TargetType = 10 // Người dùng
	TargetTypeGroup        TargetType = 20 // Nhóm
	TargetTypeOrganization TargetType = 30 // Tổ chức
	TargetTypeProduct      TargetType = 40 // Sản phẩm
	TargetTypeService      TargetType = 50 // Dịch vụ
)

var TargetTypeMap = map[TargetType]string{
	TargetTypeUser:         "Người dùng",
	TargetTypeGroup:        "Nhóm",
	TargetTypeOrganization: "Tổ chức",
	TargetTypeProduct:      "Sản phẩm",
	TargetTypeService:      "Dịch vụ",
}

const (
	// RateStatus - Trạng thái đánh giá
	RateStatusPending  uint32 = 10 // Chờ duyệt
	RateStatusApproved uint32 = 20 // Đã duyệt
	RateStatusRejected uint32 = 30 // Từ chối
	RateStatusHidden   uint32 = 40 // Ẩn
)