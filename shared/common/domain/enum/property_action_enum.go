package _enum

type PropertyAction uint32

const (
	ACTION_DISTRIBUTE_CREATE  PropertyAction = 1
	ACTION_DISTRIBUTE_UPDATE  PropertyAction = 2
	ACTION_ADD_PRODUCT_TODEAL PropertyAction = 3

	ACTION_CREATE_NOTE_PRODUCT PropertyAction = 4
	ACTION_APPLY_UPDATE_PRICE  PropertyAction = 5
	ACTION_PROPERTY_CREATE     PropertyAction = 6
)

var PropertyActionText = map[PropertyAction]string{
	ACTION_DISTRIBUTE_CREATE:   "Tạo kênh phân phối",
	ACTION_DISTRIBUTE_UPDATE:   "Cập nhật kênh phân phối",
	ACTION_ADD_PRODUCT_TODEAL:  "Thêm sản phẩm vào thương vụ",
	ACTION_CREATE_NOTE_PRODUCT: "Ghi chú vào sản phẩm",
	ACTION_APPLY_UPDATE_PRICE:  "Cập nhật giá sản phẩm/kênh giá",
	ACTION_PROPERTY_CREATE:     "Sản phẩm được tạo ở đây",
}

func (a PropertyAction) String() string {
	if text, ok := PropertyActionText[a]; ok {
		return text
	}
	return "Unknown"
}
