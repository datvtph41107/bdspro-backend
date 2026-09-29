package enums

type ProductPriceStatus uint32

const (
	PriceStatusUnknown    ProductPriceStatus = 0
	PriceStatusFixed      ProductPriceStatus = 1 // Cố định
	PriceStatusNegotiable ProductPriceStatus = 2 // Thương lượng
	PriceStatusContact    ProductPriceStatus = 3 // Liên hệ
	PriceStatusHidden     ProductPriceStatus = 4 // Ẩn giá
)

var ProductPriceStatusToString = map[ProductPriceStatus]string{
	PriceStatusFixed:      "cố định",
	PriceStatusNegotiable: "Thương lượng",
	PriceStatusContact:    "Liên hệ",
	PriceStatusHidden:     "Ẩn giá",
}
