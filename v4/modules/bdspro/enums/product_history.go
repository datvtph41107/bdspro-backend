package enums

type EProductHistory uint

const (
	EProductHistory_NewChild EProductHistory = 10
	EProductHistory_Merge    EProductHistory = 20
	EProductHistory_MergeAll EProductHistory = 30
	EProductHistory_Devide   EProductHistory = 40
)

var EProductHistoryNames = map[EProductHistory]string{
	EProductHistory_NewChild: "Tạo con",
	EProductHistory_Merge:    "Gộp",
	EProductHistory_MergeAll: "Gộp tất cả",
	EProductHistory_Devide:   "Chia",
}
