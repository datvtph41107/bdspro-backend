package enums

// TxActionStyleMap định nghĩa style cho từng loại action
type TxActionStyle struct {
	Color       string // Màu chữ
	BgColor     string // Màu nền
	BorderColor string // Màu viền
}

var TxActionStyleMap = map[TxAction]TxActionStyle{
	TxActionDeposite: {
		Color:       "#FE9A00",
		BgColor:     "#FFFBEB",
		BorderColor: "#FEE685",
	},
	TxActionSign: {
		Color:       "#00BBA7",
		BgColor:     "#F0FDFA",
		BorderColor: "#96F7E4",
	},
	TxActionPay: {
		Color:       "#2B7FFF",
		BgColor:     "#EFF6FF",
		BorderColor: "#BEDBFF",
	},
	TxActionHandOver: {
		Color:       "#00A63E",
		BgColor:     "#F0FDF4",
		BorderColor: "#B9F8CF",
	},
	TxActionCancel: {
		Color:       "#FF4757",
		BgColor:     "#FFF5F5",
		BorderColor: "#FFB3B3",
	},
	TxActionNegotiate: {
		Color:       "#8E51FF",
		BgColor:     "#F5F3FF",
		BorderColor: "#DDD6FF",
	},
}

// GetTxActionStyle trả về style cho một action
func GetTxActionStyle(action TxAction) TxActionStyle {
	if style, exists := TxActionStyleMap[action]; exists {
		return style
	}
	// Default style nếu không tìm thấy
	return TxActionStyle{
		Color:       "#6B7280",
		BgColor:     "#F9FAFB",
		BorderColor: "#E5E7EB",
	}
}
