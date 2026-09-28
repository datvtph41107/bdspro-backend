package enums

// ContributeType — loại đóng góp của user
// Hiển thị trên UI "Chọn loại báo cáo"
type ContributeType uint32

const (
	ContributeTypeCorrectInfo   ContributeType = 10 // Cung cấp thông tin chính xác
	ContributeTypeWrongInfo     ContributeType = 20 // Báo sai thông tin
	ContributeTypeWrongMedia    ContributeType = 30 // Ảnh/video không phù hợp
	ContributeTypeSpam          ContributeType = 40 // Nội dung spam
	ContributeTypeDuplicate     ContributeType = 50 // Trùng lặp
	ContributeTypeWrongLocation ContributeType = 60 // Sai vị trí
	ContributeTypeOther         ContributeType = 99 // Khác
)

func (t ContributeType) String() string {
	switch t {
	case ContributeTypeCorrectInfo:
		return "correct_info"
	case ContributeTypeWrongInfo:
		return "wrong_info"
	case ContributeTypeWrongMedia:
		return "wrong_media"
	case ContributeTypeSpam:
		return "spam"
	case ContributeTypeDuplicate:
		return "duplicate"
	case ContributeTypeWrongLocation:
		return "wrong_location"
	default:
		return "other"
	}
}

func (t ContributeType) Label() string {
	switch t {
	case ContributeTypeCorrectInfo:
		return "Cung cấp thông tin chính xác"
	case ContributeTypeWrongInfo:
		return "Thông tin sai"
	case ContributeTypeWrongMedia:
		return "Ảnh/video không phù hợp"
	case ContributeTypeSpam:
		return "Nội dung spam"
	case ContributeTypeDuplicate:
		return "Trùng lặp bất động sản"
	case ContributeTypeWrongLocation:
		return "Sai vị trí"
	default:
		return "Khác"
	}
}

// IsDataContribute — loại này có kèm dữ liệu đề xuất không
// (phân biệt "báo cáo thuần" vs "đóng góp có dữ liệu")
func (t ContributeType) IsDataContribute() bool {
	return t == ContributeTypeCorrectInfo
}

// ContributeStatus — trạng thái xử lý đóng góp
type ContributeStatus int32

const (
	ContributeStatusPending  ContributeStatus = 10 // Đang chờ xét duyệt
	ContributeStatusApproved ContributeStatus = 20 // Đã duyệt, đã merge vào bản gốc
	ContributeStatusRejected ContributeStatus = 30 // Từ chối
	ContributeStatusMerging  ContributeStatus = 15 // Đang trong quá trình merge
)

func (s ContributeStatus) String() string {
	switch s {
	case ContributeStatusPending:
		return "pending"
	case ContributeStatusApproved:
		return "approved"
	case ContributeStatusRejected:
		return "rejected"
	case ContributeStatusMerging:
		return "merging"
	default:
		return "unknown"
	}
}
