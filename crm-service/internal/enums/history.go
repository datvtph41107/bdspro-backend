package enums

type EHistory int

const (
	HistoryCreateStep     EHistory = 10
	HistoryDelete         EHistory = 20
	HistoryAssignRole     EHistory = 30
	HistorySettingAuto    EHistory = 40
	HistoryReminder       EHistory = 50
	HistoryCreateCustomer EHistory = 60
	HistoryUpdateCustomer EHistory = 70
	HistoryDeleteCustomer EHistory = 80
	HistoryAssignCustomer EHistory = 90
	HistorySwitchStage    EHistory = 100
	HistoryUpdateNote     EHistory = 110
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
)

var HistoryMap = map[EHistory]string{
	HistoryCreateStep:     "Tạo bước",
	HistoryDelete:         "Xóa",
	HistoryAssignRole:     "Phân quyền",
	HistorySettingAuto:    "Cài đặt tự động",
	HistoryReminder:       "Nhắc nhở",
	HistoryCreateCustomer: "Tạo khách hàng",
	HistoryUpdateCustomer: "Cập nhật khách hàng",
	HistoryDeleteCustomer: "Xóa khách hàng",
	HistoryAssignCustomer: "Phân công khách hàng",
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
}
