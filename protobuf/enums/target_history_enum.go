package shared_enum

type ETargetHistory int8

const (
	// crm
	TargetHistoryLead     ETargetHistory = 10
	TargetHistoryContact  ETargetHistory = 20
	TargetHistoryPipeline ETargetHistory = 30
	TargetHistoryStage    ETargetHistory = 40
	TargetHistoryRule     ETargetHistory = 50

	// bdspro
	TargetHistoryProduct ETargetHistory = 60
	TargetHistoryAsset   ETargetHistory = 70
)

var TargetHistoryMap = map[ETargetHistory]string{
	TargetHistoryLead:     "Khách hàng tiềm năng",
	TargetHistoryContact:  "Liên hệ",
	TargetHistoryPipeline: "Quy trình",
	TargetHistoryStage:    "Giai đoạn",
	TargetHistoryRule:     "Quy tắc",
	TargetHistoryProduct:  "Sản phẩm",
	TargetHistoryAsset:    "Tài sản",
}
