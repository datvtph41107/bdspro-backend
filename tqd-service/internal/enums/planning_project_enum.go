package enums

type PlanningType uint32
type PlanningLevel uint32
type PlanningDocumentType uint32
type PlanningProcessStatus uint32

const (
	PlanningTypeUsed    uint32 = 10
	PlanningTypeGeneral uint32 = 20
	PlanningTypeUrban   uint32 = 30
	PlanningTypeDetail  uint32 = 40

	PlanningLevelNational   uint32 = 10
	PlanningLevelRegional   uint32 = 15
	PlanningLevelProvincial uint32 = 20
	PlanningLevelDistrict   uint32 = 30
	PlanningLevelWard       uint32 = 40
	PlanningLevelProject    uint32 = 45

	PlanningDocumentTypeLegal       uint32 = 10
	PlanningDocumentTypeExplanation uint32 = 20
	PlanningDocumentTypeMap         uint32 = 30
	PlanningDocumentTypeCAD         uint32 = 40
	PlanningDocumentTypeGIS         uint32 = 50
	PlanningDocumentTypeAppendix    uint32 = 60
	PlanningDocumentTypeOther       uint32 = 70

	// PlanningProcessStatus — trạng thái xử lý của đồ án/tài liệu sau khi upload folder.
	PlanningProcessStatusPending    uint32 = 10 // vừa upload, chưa quét
	PlanningProcessStatusProcessing uint32 = 20 // job đang gọi AI phân loại
	PlanningProcessStatusClassified uint32 = 30 // AI đã phân loại xong, chờ admin duyệt
	PlanningProcessStatusApproved   uint32 = 40 // admin đã phê duyệt
	PlanningProcessStatusFailed     uint32 = 50 // lỗi khi upload hoặc phân loại
)

var PlanningDocumentTypeLabelMap = map[uint32]string{
	PlanningDocumentTypeLegal:       "Pháp lý",
	PlanningDocumentTypeExplanation: "Thuyết minh",
	PlanningDocumentTypeMap:         "Bản đồ",
	PlanningDocumentTypeCAD:         "CAD",
	PlanningDocumentTypeGIS:         "GIS gốc",
	PlanningDocumentTypeAppendix:    "Phụ lục",
	PlanningDocumentTypeOther:       "Tệp khác",
}

var PlanningTypeMap = map[uint32]string{
	PlanningTypeUsed:    "QH sử dụng đất",
	PlanningTypeGeneral: "QH chung",
	PlanningTypeUrban:   "QH phân khu",
	PlanningTypeDetail:  "QH chi tiết",
}

var PlanningLevelMap = map[uint32]string{
	PlanningLevelNational:   "Quốc gia",
	PlanningLevelRegional:   "Khu vực",
	PlanningLevelProvincial: "Tỉnh",
	PlanningLevelDistrict:   "Huyện",
	PlanningLevelWard:       "Xã",
	PlanningLevelProject:    "Dự án",
}

var PlanningProcessStatusMap = map[uint32]string{
	PlanningProcessStatusPending:    "Chờ xử lý",
	PlanningProcessStatusProcessing: "Đang xử lý",
	PlanningProcessStatusClassified: "Đã phân loại",
	PlanningProcessStatusApproved:   "Đã phê duyệt",
	PlanningProcessStatusFailed:     "Lỗi",
}

func GetPlanningProcessStatusLabel(status uint32) string {
	return PlanningProcessStatusMap[status]
}

func GetPlanningDocumentTypeLabel(planningDocumentType uint32) string {
	return PlanningDocumentTypeLabelMap[planningDocumentType]
}

func GetPlanningDocumentTypeValue(label string) uint32 {
	for key, value := range PlanningDocumentTypeLabelMap {
		if value == label {
			return key
		}
	}
	return 0
}

func GetPlanningLevelLabel(planningLevel uint32) string {
	return PlanningLevelMap[planningLevel]
}

func GetPlanningLevelValue(label string) uint32 {
	for key, value := range PlanningLevelMap {
		if value == label {
			return key
		}
	}
	return 0
}

func GetPlanningTypeLabel(planningType uint32) string {
	return PlanningTypeMap[planningType]
}

func GetPlanningTypeValue(label string) uint32 {
	for key, value := range PlanningTypeMap {
		if value == label {
			return key
		}
	}
	return 0
}
