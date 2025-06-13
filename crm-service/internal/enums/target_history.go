package enums

type ETargetHistory int8

const (
	TargetHistoryLead     ETargetHistory = 10
	TargetHistoryContact  ETargetHistory = 20
	TargetHistoryPipeline ETargetHistory = 30
	TargetHistoryStage    ETargetHistory = 40
	TargetHistoryRule     ETargetHistory = 50
)

var TargetHistoryMap = map[ETargetHistory]string{
	TargetHistoryLead:     "Khách hàng tiềm năng",
	TargetHistoryContact:  "Liên hệ",
	TargetHistoryPipeline: "Quy trình",
	TargetHistoryStage:    "Giai đoạn",
	TargetHistoryRule:     "Quy tắc",
}
