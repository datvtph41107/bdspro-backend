package enums

type EProductSaleStatus uint32
type EProductRentStatus uint32

const (
	EProductNotSold EProductSaleStatus = 10
	EProductSelling EProductSaleStatus = 20
	EProductSold    EProductSaleStatus = 30
	EProductNotRent EProductRentStatus = 10
	EProductRenting EProductRentStatus = 20
	EProductRented  EProductRentStatus = 30
)

var EProductSaleStatusNames = map[EProductSaleStatus]string{
	EProductNotSold: "Chưa bán",
	EProductSelling: "Đang bán",
	EProductSold:    "Đã bán",
}

var EProductRentStatusNames = map[EProductRentStatus]string{
	EProductNotRent: "Chưa cho thuê",
	EProductRenting: "Đang cho thuê",
	EProductRented:  "Đã cho thuê",
}
