package _enum

type PropertyAction uint32

const (
	ACTION_DISTRIBUTE_CREATE PropertyAction = 1
	ACTION_DISTRIBUTE_UPDATE PropertyAction = 2
)

var PropertyActionText = map[PropertyAction]string{
	ACTION_DISTRIBUTE_CREATE: "Tạo kênh phân phối",
	ACTION_DISTRIBUTE_UPDATE: "Cập nhật kênh phân phối",
}

func (a PropertyAction) String() string {
	if text, ok := PropertyActionText[a]; ok {
		return text
	}
	return "Unknown"
}
