package enums

type DealStatus uint32

const (
	DealStatusNew       DealStatus = 10
	DealStatusDealing   DealStatus = 20
	DealStatusDeposit   DealStatus = 30
	DealStatusCompleted DealStatus = 40
	DealStatusCanceled  DealStatus = 50
)

var DealStatusMap = map[DealStatus]string{
	DealStatusNew:       "Mới tạo",
	DealStatusDealing:   "Đang thương lượng",
	DealStatusDeposit:   "Đặt cọc",
	DealStatusCompleted: "Hoàn tất",
	DealStatusCanceled:  "Hủy",
}
