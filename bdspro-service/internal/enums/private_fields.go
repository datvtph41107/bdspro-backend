package enums

var PrivateFields = []struct {
	Value string `json:"value"`
	Label string `json:"label"`
}{
	{Value: "ImportPrice", Label: "Giá nhập"},
	{Value: "OperatingCost", Label: "Chi phí vận hành"},
	{Value: "InternalNote", Label: "Ghi chú nội bộ"},
	{Value: "TargetProfit", Label: "Mục tiêu lợi nhuận"},
	{Value: "PrivateDocs", Label: "Xem file"},
}
