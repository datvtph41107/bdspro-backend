package enums

type ProductStatus uint

const (
	ProductNotSold ProductStatus = 10 // Chưa bán
	ProductSelling ProductStatus = 20 // Đang bán(cọc)
	ProductSold    ProductStatus = 30 // Đã bán
	ProductNotRent ProductStatus = 40 // Chưa cho thuê
	ProductRenting ProductStatus = 50 // Đang cho thuê(cọc)
	ProductRented  ProductStatus = 60 // Đã cho thuê
)

func (s ProductStatus) String() string {
	switch s {
	case ProductNotSold:
		return "Chưa bán"
	case ProductSelling:
		return "Đang bán"
	case ProductSold:
		return "Đã bán"
	case ProductNotRent:
		return "Chưa cho thuê"
	case ProductRenting:
		return "Đang cho thuê"
	case ProductRented:
		return "Đã cho thuê"
	default:
		return ""
	}
}
