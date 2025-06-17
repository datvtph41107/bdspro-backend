package enums

type EHistory int16

const (
	HistoryCreateStep     EHistory = 10
	HistoryDelete         EHistory = 20
	HistoryAssignRole     EHistory = 30
	HistorySettingAuto    EHistory = 40
	HistoryPipelineCreate EHistory = 120
	HistoryPipelineUpdate EHistory = 130
	HistoryPipelineDelete EHistory = 140
	HistoryStageCreate    EHistory = 150
	HistoryStageUpdate    EHistory = 160
	HistoryStageDelete    EHistory = 170
	HistoryRuleCreate     EHistory = 180
	HistoryRuleUpdate     EHistory = 190
	HistoryRuleDelete     EHistory = 200
	HistoryRuleEvent      EHistory = 210

	// -- action for contact --
	HistoryContactCreate  EHistory = 6010
	HistoryContactUpdate  EHistory = 6020
	HistoryContactDelete  EHistory = 6030
	HistoryContactAssign  EHistory = 6040
	HistoryReminder       EHistory = 6050
	HistorySwitchStage    EHistory = 6060
	HistoryUpdateNote     EHistory = 6070
	HistoryContactInstall EHistory = 6080
	HistoryContactFriend  EHistory = 6090
	HistoryContactStage   EHistory = 6100
	HistoryContactMessage EHistory = 6110
)

var HistoryMap = map[EHistory]string{
	HistoryCreateStep:     "Tạo bước",
	HistoryDelete:         "Xóa",
	HistoryAssignRole:     "Phân quyền",
	HistorySettingAuto:    "Cài đặt tự động",
	HistoryReminder:       "Nhắc nhở",
	HistorySwitchStage:    "Chuyển giai đoạn",
	HistoryUpdateNote:     "Cập nhật ghi chú",
	HistoryPipelineCreate: "Tạo quy trình",
	HistoryPipelineUpdate: "Cập nhật quy trình",
	HistoryPipelineDelete: "Xóa quy trình",
	HistoryStageCreate:    "Tạo giai đoạn",
	HistoryStageUpdate:    "Cập nhật giai đoạn",
	HistoryStageDelete:    "Xóa giai đoạn",
	HistoryRuleCreate:     "Tạo quy tắc",
	HistoryRuleUpdate:     "Cập nhật quy tắc",
	HistoryRuleDelete:     "Xóa quy tắc",

	// -- action for contact --
	HistoryContactCreate:  "Tạo liên hệ",
	HistoryContactDelete:  "Xóa liên hệ",
	HistoryContactAssign:  "Phân công liên hệ",
	HistoryContactUpdate:  "Cập nhật liên hệ",
	HistoryContactInstall: "Mời cài app",
	HistoryContactFriend:  "Gửi lời kết bạn",
	HistoryContactStage:   "Cập nhật trạng thái CRM",
	HistoryContactMessage: "Gửi tin nhắn",
}
