package enums

import "fmt"

// ESystemConfigGroup enum cho nhóm system config
type ESystemConfigGroup int32

const (
	SystemConfigDesign      ESystemConfigGroup = 10 // Cấu hình chung
	SystemConfigFormat      ESystemConfigGroup = 20 // Định dạng (date, price, currency)
	SystemConfigGeography   ESystemConfigGroup = 30 // Ngôn ngữ, múi giờ
	SystemConfigParameter   ESystemConfigGroup = 40 // Upload file
	SystemConfigPerformance ESystemConfigGroup = 50 // Email settings
)

var SystemComfigMap = map[string]ESystemConfigGroup{
	"design":      SystemConfigDesign,
	"format":      SystemConfigFormat,
	"geography":   SystemConfigGeography,
	"parameter":   SystemConfigParameter,
	"performance": SystemConfigPerformance,
}

// SystemConfigGroupEnumToKey map từ enum sang key string
var SystemConfigGroupEnumToKey = map[ESystemConfigGroup]string{
	SystemConfigDesign:      "design",
	SystemConfigFormat:      "format",
	SystemConfigGeography:   "geography",
	SystemConfigParameter:   "parameter",
	SystemConfigPerformance: "performance",
}

// SystemConfigGroupMap map từ enum sang tên tiếng Việt
var SystemConfigGroupMap = map[ESystemConfigGroup]string{
	SystemConfigDesign:      "Thương hiệu & Giao diện",
	SystemConfigFormat:      "Ngôn ngữ & Định dạng",
	SystemConfigGeography:   "Phân vùng địa lý",
	SystemConfigParameter:   "Tham số hiển thị & phân trang",
	SystemConfigPerformance: "Tham số backend & hiệu năng",
}

// GetSystemConfigGroupName lấy tên nhóm theo enum
func GetSystemConfigGroupName(group ESystemConfigGroup) string {
	if name, ok := SystemConfigGroupMap[group]; ok {
		return name
	}
	return "Không xác định"
}

// IsValidSystemConfigGroup kiểm tra group có hợp lệ không
func IsValidSystemConfigGroup(group ESystemConfigGroup) bool {
	_, ok := SystemConfigGroupMap[group]
	return ok
}

func GetSystemConfigGroup(group string) (ESystemConfigGroup, error) {
	groupEnum, ok := SystemComfigMap[group]
	if !ok {
		return 0, fmt.Errorf("invalid group value: %v", group)
	}
	return groupEnum, nil
}
