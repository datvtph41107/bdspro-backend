package enums

type EProductStatus uint32

const (
	EProductNotSold EProductStatus = 10
	EProductSelling EProductStatus = 20
	EProductSold    EProductStatus = 30
	EProductNotRent EProductStatus = 40
	EProductRenting EProductStatus = 50
	EProductRented  EProductStatus = 60
)

var EProductStatusNames = map[EProductStatus]string{
	EProductNotSold: "Chưa bán",
	EProductSelling: "Đang bán",
	EProductSold:    "Đã bán",
	EProductNotRent: "Chưa cho thuê",
	EProductRenting: "Đang cho thuê",
	EProductRented:  "Đã cho thuê",
}
