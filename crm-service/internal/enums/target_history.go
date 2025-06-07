package enums

type ETargetHistory int8

const (
	TargetHistoryCustomer ETargetHistory = 10
	TargetHistoryContact  ETargetHistory = 20
	TargetHistoryPipeline ETargetHistory = 30
	TargetHistoryStage    ETargetHistory = 40
	TargetHistoryRule     ETargetHistory = 50
)

var TargetHistoryMap = map[ETargetHistory]string{
	TargetHistoryCustomer: "Khách hàng",
	TargetHistoryContact:  "Liên hệ",
	TargetHistoryPipeline: "Quy trình",
	TargetHistoryStage:    "Giai đoạn",
	TargetHistoryRule:     "Quy tắc",
}
