package enums

type EHistory int

const (
	HistoryCreateStep  EHistory = 10
	HistoryDelete      EHistory = 20
	HistoryAssignRole  EHistory = 30
	HistorySettingAuto EHistory = 40
	HistoryReminder    EHistory = 50
)

var HistoryMap = map[EHistory]string{
	HistoryCreateStep:  "Tạo bước",
	HistoryDelete:      "Xóa",
	HistoryAssignRole:  "Phân quyền",
	HistorySettingAuto: "Cài đặt tự động",
	HistoryReminder:    "Nhắc nhở",
}
