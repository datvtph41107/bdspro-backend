package enums

type ProductStatus int

const (
	ProductNotSold ProductStatus = 1 // Chưa bán
	ProductSelling ProductStatus = 2 // Đang bán(cọc)
	ProductSold    ProductStatus = 3 // Đã bán
	ProductNotRent ProductStatus = 4 // Chưa cho thuê
	ProductRenting ProductStatus = 5 // Đang cho thuê(cọc)
	ProductRented  ProductStatus = 6 // Đã cho thuê
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
