// internal/enums/group_code.go
package enums

// GroupCode - Nhóm đất (uint32)
type GroupCode uint32

const (
	GroupCodeUnknown     GroupCode = 0
	GroupCodeResidential GroupCode = 1
	GroupCodeTransport   GroupCode = 2
	GroupCodePublic      GroupCode = 3
	GroupCodeGreen       GroupCode = 4
	GroupCodeWater       GroupCode = 5
	GroupCodeCemetery    GroupCode = 6
	GroupCodeIndustrial  GroupCode = 7
	GroupCodeCommercial  GroupCode = 8
	GroupCodeAgriculture GroupCode = 9
	GroupCodeForestry    GroupCode = 10
	GroupCodeDefense     GroupCode = 11
	GroupCodeHistorical  GroupCode = 12
)

// String - Chuyển sang string (dùng cho API)
func (g GroupCode) String() string {
	switch g {
	case GroupCodeResidential:
		return "residential"
	case GroupCodeTransport:
		return "transport"
	case GroupCodePublic:
		return "public"
	case GroupCodeGreen:
		return "green"
	case GroupCodeWater:
		return "water"
	case GroupCodeCemetery:
		return "cemetery"
	case GroupCodeIndustrial:
		return "industrial"
	case GroupCodeCommercial:
		return "commercial"
	case GroupCodeAgriculture:
		return "agriculture"
	case GroupCodeForestry:
		return "forestry"
	case GroupCodeDefense:
		return "defense"
	case GroupCodeHistorical:
		return "historical"
	default:
		return "unknown"
	}
}

// DisplayName - Tên hiển thị tiếng Việt
func (g GroupCode) DisplayName() string {
	switch g {
	case GroupCodeResidential:
		return "Đất ở"
	case GroupCodeTransport:
		return "Đất giao thông"
	case GroupCodePublic:
		return "Đất công cộng"
	case GroupCodeGreen:
		return "Đất cây xanh"
	case GroupCodeWater:
		return "Đất mặt nước"
	case GroupCodeCemetery:
		return "Đất nghĩa trang"
	case GroupCodeIndustrial:
		return "Đất công nghiệp"
	case GroupCodeCommercial:
		return "Đất thương mại"
	case GroupCodeAgriculture:
		return "Đất nông nghiệp"
	case GroupCodeForestry:
		return "Đất lâm nghiệp"
	case GroupCodeDefense:
		return "Đất quốc phòng"
	case GroupCodeHistorical:
		return "Đất di tích"
	default:
		return "Không xác định"
	}
}

// CanBuild - Xác định có được xây dựng không (theo luật mặc định)
func (g GroupCode) CanBuild() bool {
	switch g {
	case GroupCodeResidential, GroupCodePublic, GroupCodeIndustrial, GroupCodeCommercial:
		return true
	default:
		return false
	}
}
