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

	// advertising/marketing
	TargetHistoryCampaign ETargetHistory = 80
	TargetHistoryPackage  ETargetHistory = 90

	// admin
	TargetHistoryAdmin ETargetHistory = 100
)

var TargetHistoryMap = map[ETargetHistory]string{
	TargetHistoryLead:     "Khách hàng tiềm năng",
	TargetHistoryContact:  "Liên hệ",
	TargetHistoryPipeline: "Quy trình",
	TargetHistoryStage:    "Giai đoạn",
	TargetHistoryRule:     "Quy tắc",
	TargetHistoryProduct:  "Sản phẩm",
	TargetHistoryAsset:    "Tài sản",
	TargetHistoryCampaign: "Chiến dịch quảng cáo",
	TargetHistoryPackage:  "Gói quảng cáo",
	TargetHistoryAdmin:    "Admin",
}
