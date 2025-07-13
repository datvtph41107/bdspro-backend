package enums

type EPostStatus int

const (
	// EPostPublic          EPostStatus = 10
	// EPostFinished        EPostStatus = 20
	// EPostSold            EPostStatus = 40
	// EPostRent            EPostStatus = 50
	// EPostStop            EPostStatus = 60 // tùy nghiệp vụ về sau
	// EPostStatusCompleted EPostStatus = 70
	EPostSelling EPostStatus = 100
	EPostSold    EPostStatus = 110
	EPostRenting EPostStatus = 200
	EPostRented  EPostStatus = 210
)

var EPostStatusNames = map[EPostStatus]string{
	// EPostPublic:          "Công khai",
	// EPostFinished:        "Hết hạn",
	// EPostSold:            "Đã bán",
	// EPostRent:            "Đã cho thuê",
	// EPostStop:            "Dừng",
	// EPostStatusCompleted: "Đã hoàn thành",
	EPostSelling: "Đang bán",
	EPostSold:    "Đã bán",
	EPostRenting: "Đang cho thuê",
	EPostRented:  "Đã cho thuê",
}
