package enums

type DealStatus uint32

const (
	DealStatusProcessing DealStatus = 10
	DealStatusCompleted  DealStatus = 20
	DealStatusCanceled   DealStatus = 30
	// DealStatusDealing   DealStatus = 20
	// DealStatusDeposit   DealStatus = 30
)

var DealStatusMap = map[DealStatus]string{
	DealStatusProcessing: "Đang xử lý",
	DealStatusCompleted:  "Hoàn tất",
	DealStatusCanceled:   "Hủy",
	// DealStatusDealing:   "Đang thương lượng",
	// DealStatusDeposit:   "Đặt cọc",
}

func (s DealStatus) IsValid() bool {
	switch s {
	case DealStatusProcessing,
		DealStatusCompleted,
		DealStatusCanceled:
		return true
	}
	return false
}

