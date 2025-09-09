package enums

// ActionStyleMap định nghĩa style cho từng loại action
type ActionStyle struct {
	Color       string // Màu chữ
	BgColor     string // Màu nền
	BorderColor string // Màu viền
}

var ActionStyleMap = map[EAction]ActionStyle{
	EActionDeposite: {
		Color:       "#FE9A00",
		BgColor:     "#FFFBEB",
		BorderColor: "#FEE685",
	},
	EActionSign: {
		Color:       "#00BBA7",
		BgColor:     "#F0FDFA",
		BorderColor: "#96F7E4",
	},
	EActionPay: {
		Color:       "#2B7FFF",
		BgColor:     "#EFF6FF",
		BorderColor: "#BEDBFF",
	},
	EActionHandOver: {
		Color:       "#00A63E",
		BgColor:     "#F0FDF4",
		BorderColor: "#B9F8CF",
	},
	EActionCancel: {
		Color:       "#FF4757",
		BgColor:     "#FFF5F5",
		BorderColor: "#FFB3B3",
	},
	EActionNegotiate: {
		Color:       "#8E51FF",
		BgColor:     "#F5F3FF",
		BorderColor: "#DDD6FF",
	},
}

// GetActionStyle trả về style cho một action
func GetActionStyle(action EAction) ActionStyle {
	if style, exists := ActionStyleMap[action]; exists {
		return style
	}
	// Default style nếu không tìm thấy
	return ActionStyle{
		Color:       "#6B7280",
		BgColor:     "#F9FAFB",
		BorderColor: "#E5E7EB",
	}
}
