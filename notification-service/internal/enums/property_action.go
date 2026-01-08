package enums

type PropertyAction uint32

const (
	PropertyCreated  PropertyAction = 1
	PropertyUpdated  PropertyAction = 2
	PropertyArchived PropertyAction = 3
	PropertyRestored PropertyAction = 4

	PropertyPublished     PropertyAction = 10
	PropertyUnpublished   PropertyAction = 11
	PropertyApproved      PropertyAction = 12
	PropertyRejected      PropertyAction = 13
	PropertyStatusChanged PropertyAction = 14

	PropertyPriceChanged PropertyAction = 20
	PropertyRentChanged  PropertyAction = 21
	PropertyFeeChanged   PropertyAction = 22
)

var PropertyActionText = map[PropertyAction]string{
	PropertyCreated:  "Tạo bất động sản",
	PropertyUpdated:  "Cập nhật thông tin bất động sản",
	PropertyArchived: "Lưu trữ bất động sản",
	PropertyRestored: "Khôi phục bất động sản",

	PropertyPublished:     "Đăng bất động sản",
	PropertyUnpublished:   "Gỡ bất động sản khỏi sàn",
	PropertyApproved:      "Duyệt bất động sản",
	PropertyRejected:      "Từ chối duyệt bất động sản",
	PropertyStatusChanged: "Thay đổi trạng thái bất động sản",

	PropertyPriceChanged: "Thay đổi giá bất động sản",
	PropertyRentChanged:  "Thay đổi giá thuê",
	PropertyFeeChanged:   "Thay đổi phí",
}

func (a PropertyAction) String() string {
	if text, ok := PropertyActionText[a]; ok {
		return text
	}
	return "Unknown"
}
