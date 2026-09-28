package enums

// EPostStatus - Trạng thái cơ bản của tin đăng
type EPostStatus int

const (
	EPostActive   EPostStatus = 10 // Đang hoạt động
	EPostPending  EPostStatus = 20 // Chờ duyệt
	EPostApproved EPostStatus = 30 // Đã duyệt
	EPostRejected EPostStatus = 40 // Từ chối
	EPostHidden   EPostStatus = 50 // Ẩn
)

var EPostStatusNames = map[EPostStatus]string{
	EPostActive:   "Đang hoạt động",
	EPostPending:  "Chờ duyệt",
	EPostApproved: "Đã duyệt",
	EPostRejected: "Từ chối",
	EPostHidden:   "Ẩn",
}

// EPostVisibleStatus - Trạng thái giao dịch của tin đăng
type EPostVisibleStatus int

const (
	EPostTransactionSelling EPostVisibleStatus = 100 // Đang bán
	EPostTransactionSold    EPostVisibleStatus = 110 // Đã bán
	EPostTransactionRenting EPostVisibleStatus = 200 // Đang cho thuê
	EPostTransactionRented  EPostVisibleStatus = 210 // Đã cho thuê
)

var EPostTransactionStatusNames = map[EPostVisibleStatus]string{
	EPostTransactionSelling: "Đang bán",
	EPostTransactionSold:    "Đã bán",
	EPostTransactionRenting: "Đang cho thuê",
	EPostTransactionRented:  "Đã cho thuê",
}
