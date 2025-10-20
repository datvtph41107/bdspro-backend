package enums

type EEventQueueStatus uint32

const (
	EEventQueueStatusPending    EEventQueueStatus = 10 // Đang chờ xử lý
	EEventQueueStatusProcessing EEventQueueStatus = 20 // Đang xử lý
	EEventQueueStatusCompleted  EEventQueueStatus = 30 // Hoàn thành
	EEventQueueStatusFailed     EEventQueueStatus = 40 // Thất bại
	EEventQueueStatusCancelled  EEventQueueStatus = 50 // Đã hủy
)

var EventQueueStatusMap = map[EEventQueueStatus]string{
	EEventQueueStatusPending:    "Đang chờ xử lý",
	EEventQueueStatusProcessing: "Đang xử lý",
	EEventQueueStatusCompleted:  "Hoàn thành",
	EEventQueueStatusFailed:     "Thất bại",
	EEventQueueStatusCancelled:  "Đã hủy",
}

func (e EEventQueueStatus) String() string {
	if val, ok := EventQueueStatusMap[e]; ok {
		return val
	}
	return "Unknown"
}

func (e EEventQueueStatus) IsValid() bool {
	_, ok := EventQueueStatusMap[e]
	return ok
}

