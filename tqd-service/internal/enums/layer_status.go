package enums

// LayerType - Loại layer (1-5)
type LayerType uint32

const (
	LayerTypeMasterPlan     LayerType = 10 // Quy hoạch tổng thể
	LayerTypeZoning         LayerType = 20 // Quy hoạch phân khu
	LayerTypeLandUse        LayerType = 30 // Quy hoạch sử dụng đất
	LayerTypeInfrastructure LayerType = 40 // Quy hoạch hạ tầng
	LayerTypeEnvironmental  LayerType = 50 // Quy hoạch môi trường
)

func (t LayerType) IsValid() bool {
	return t == LayerTypeMasterPlan ||
		t == LayerTypeZoning ||
		t == LayerTypeLandUse ||
		t == LayerTypeInfrastructure ||
		t == LayerTypeEnvironmental
}

func (t LayerType) String() string {
	switch t {
	case LayerTypeMasterPlan:
		return "master_plan"
	case LayerTypeZoning:
		return "zoning"
	case LayerTypeLandUse:
		return "land_use"
	case LayerTypeInfrastructure:
		return "infrastructure"
	case LayerTypeEnvironmental:
		return "environmental"
	}
	return "unknown"
}

// LayerStatus - Trạng thái layer (1,10,20,30)
type LayerStatus uint32

const (
	LayerStatusDraft    LayerStatus = 1  // Bản nháp (chỉ admin thấy)
	LayerStatusActive   LayerStatus = 10 // Đang hiệu lực (client thấy)
	LayerStatusExpired  LayerStatus = 20 // Hết hiệu lực
	LayerStatusArchived LayerStatus = 30 // Đã lưu trữ
)

func (s LayerStatus) IsValid() bool {
	return s == 1 || s == 10 || s == 20 || s == 30
}

func (s LayerStatus) IsClientVisible() bool {
	return s == LayerStatusActive
}

func (s LayerStatus) IsActive() bool {
	return s == LayerStatusActive
}

func (s LayerStatus) IsDraft() bool {
	return s == LayerStatusDraft
}

func (s LayerStatus) String() string {
	switch s {
	case LayerStatusDraft:
		return "draft"
	case LayerStatusActive:
		return "active"
	case LayerStatusExpired:
		return "expired"
	case LayerStatusArchived:
		return "archived"
	}
	return "unknown"
}

func (s LayerStatus) DisplayName() string {
	switch s {
	case LayerStatusDraft:
		return "Dự thảo"
	case LayerStatusActive:
		return "Đang hiệu lực"
	case LayerStatusExpired:
		return "Hết hiệu lực"
	case LayerStatusArchived:
		return "Lưu trữ"
	default:
		return "Không xác định"
	}
}
