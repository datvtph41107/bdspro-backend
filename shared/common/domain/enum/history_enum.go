package _enum

type EHistory int32

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

	// -- action for asset --
	HistoryAssetUpdate    EHistory = 7010
	HistoryAssetDelete    EHistory = 7020
	HistoryAssetCreate    EHistory = 7030
	HistoryAssetMerge     EHistory = 7040
	HistoryAssetSplit     EHistory = 7050
	HistoryAssetArchive   EHistory = 7060
	HistoryAssetUnarchive EHistory = 7070
	HistoryAssetShare     EHistory = 7080

	// -- action for product --
	HistoryProductCreate    EHistory = 8010
	HistoryProductUpdate    EHistory = 8020
	HistoryProductDelete    EHistory = 8030
	HistoryProductArchive   EHistory = 8060
	HistoryProductUnarchive EHistory = 8070
	HistoryProductShare     EHistory = 8080
	// -- action for product child --
	HistoryProductChildCreate EHistory = 8090
	HistoryProductMerge       EHistory = 8100
	HistoryProductDevide      EHistory = 8110
	HistoryProductMergeAll    EHistory = 8120

	// -- action for advertising/marketing --
	HistoryCampaignCreate   EHistory = 9010
	HistoryCampaignUpdate   EHistory = 9020
	HistoryCampaignDelete   EHistory = 9030
	HistoryCampaignPayment  EHistory = 9040
	HistoryCampaignActivate EHistory = 9050
	HistoryCampaignPause    EHistory = 9060
	HistoryCampaignResume   EHistory = 9070
	HistoryCampaignExtend   EHistory = 9080
	HistoryCampaignTopup    EHistory = 9090
	HistoryPackageCreate    EHistory = 9110
	HistoryPackageUpdate    EHistory = 9120
	HistoryPackageDelete    EHistory = 9130

	// -- action for admin --
	HistoryAdminCreate EHistory = 10010
	HistoryAdminUpdate EHistory = 10020
	HistoryAdminDelete EHistory = 10030

	// -- action for user management --
	HistoryUserCreate EHistory = 11010
	HistoryUserUpdate EHistory = 11020
	HistoryUserDelete EHistory = 11030
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

	// -- action for asset --
	HistoryAssetUpdate:    "Cập nhật tài sản",
	HistoryAssetDelete:    "Xóa tài sản",
	HistoryAssetCreate:    "Tạo tài sản",
	HistoryAssetMerge:     "Gộp tài sản",
	HistoryAssetSplit:     "Tách tài sản",
	HistoryAssetArchive:   "Lưu tài sản",
	HistoryAssetUnarchive: "Mở tài sản",
	HistoryAssetShare:     "Chia sẻ tài sản",

	// -- action for product --
	HistoryProductCreate:    "Tạo sản phẩm",
	HistoryProductUpdate:    "Cập nhật sản phẩm",
	HistoryProductDelete:    "Xóa sản phẩm",
	HistoryProductMerge:     "Gộp sản phẩm",
	HistoryProductDevide:    "Tách sản phẩm",
	HistoryProductArchive:   "Lưu sản phẩm",
	HistoryProductUnarchive: "Mở sản phẩm",
	HistoryProductShare:     "Chia sẻ sản phẩm",

	// -- action for advertising/marketing --
	HistoryCampaignCreate:   "Tạo chiến dịch quảng cáo",
	HistoryCampaignUpdate:   "Cập nhật chiến dịch quảng cáo",
	HistoryCampaignDelete:   "Xóa chiến dịch quảng cáo",
	HistoryCampaignPayment:  "Thanh toán chiến dịch",
	HistoryCampaignActivate: "Kích hoạt chiến dịch",
	HistoryCampaignPause:    "Tạm dừng chiến dịch",
	HistoryCampaignResume:   "Tiếp tục chiến dịch",
	HistoryCampaignExtend:   "Gia hạn chiến dịch",
	HistoryCampaignTopup:    "Nạp tiền chiến dịch",
	HistoryPackageCreate:    "Tạo gói quảng cáo",
	HistoryPackageUpdate:    "Cập nhật gói quảng cáo",
	HistoryPackageDelete:    "Xóa gói quảng cáo",

	// -- action for admin --
	HistoryAdminCreate: "Tạo admin",
	HistoryAdminUpdate: "Cập nhật admin",
	HistoryAdminDelete: "Xóa admin",

	// -- action for user management --
	HistoryUserCreate: "Tạo người dùng",
	HistoryUserUpdate: "Cập nhật người dùng",
	HistoryUserDelete: "Xóa người dùng",
}

var HistoryContactAction = []EHistory{
	HistoryCreateStep,
	HistoryDelete,
	HistoryAssignRole,
	HistorySettingAuto,
	HistoryReminder,
	HistorySwitchStage,
	HistoryUpdateNote,
	HistoryPipelineCreate,
	HistoryPipelineUpdate,
	HistoryPipelineDelete,
	HistoryStageCreate,
	HistoryStageUpdate,
	HistoryStageDelete,
	HistoryRuleCreate,
	HistoryRuleUpdate,
	HistoryRuleDelete,
	HistoryContactCreate,
	HistoryContactDelete,
	HistoryContactAssign,
	HistoryContactUpdate,
	HistoryContactInstall,
	HistoryContactFriend,
	HistoryContactStage,
	HistoryContactMessage,
}

var HistoryAssetActionType = []EHistory{
	HistoryAssetUpdate,
	HistoryAssetDelete,
	HistoryAssetCreate,
	HistoryAssetMerge,
	HistoryAssetSplit,
	HistoryAssetArchive,
	HistoryAssetUnarchive,
	HistoryAssetShare,
}

var HistoryProductActionType = []EHistory{
	HistoryProductCreate,
	HistoryProductUpdate,
	HistoryProductDelete,
	HistoryProductArchive,
	HistoryProductUnarchive,
	HistoryProductShare,
	HistoryProductChildCreate,
	HistoryProductMerge,
	HistoryProductMergeAll,
	HistoryProductDevide,
}

var HistoryProductChild = []EHistory{
	HistoryProductMerge,
	HistoryProductMergeAll,
	HistoryProductDevide,
	HistoryProductChildCreate,
}

var HistoryCampaignActionType = []EHistory{
	HistoryCampaignCreate,
	HistoryCampaignUpdate,
	HistoryCampaignDelete,
	HistoryCampaignPayment,
	HistoryCampaignActivate,
	HistoryCampaignPause,
	HistoryCampaignResume,
	HistoryCampaignExtend,
	HistoryCampaignTopup,
}

var HistoryPackageActionType = []EHistory{
	HistoryPackageCreate,
	HistoryPackageUpdate,
	HistoryPackageDelete,
}
