package _enum

import "slices"

type ENotificationType int32

const (
	NotificationSystem ENotificationType = 100

	NotificationCreateStep     ENotificationType = 10
	NotificationDelete         ENotificationType = 20
	NotificationAssignRole     ENotificationType = 30
	NotificationSettingAuto    ENotificationType = 40
	NotificationPipelineCreate ENotificationType = 120
	NotificationPipelineUpdate ENotificationType = 130
	NotificationPipelineDelete ENotificationType = 140
	NotificationStageCreate    ENotificationType = 150
	NotificationStageUpdate    ENotificationType = 160
	NotificationStageDelete    ENotificationType = 170
	NotificationRuleCreate     ENotificationType = 180
	NotificationRuleUpdate     ENotificationType = 190
	NotificationRuleDelete     ENotificationType = 200
	NotificationRuleEvent      ENotificationType = 6210

	// -- action for contact --
	NotificationContactCreate  ENotificationType = 6010
	NotificationContactUpdate  ENotificationType = 6020
	NotificationContactDelete  ENotificationType = 6030
	NotificationContactAssign  ENotificationType = 6040
	NotificationReminder       ENotificationType = 6050
	NotificationSwitchStage    ENotificationType = 6060
	NotificationUpdateNote     ENotificationType = 6070
	NotificationContactInstall ENotificationType = 6080
	NotificationContactFriend  ENotificationType = 6090
	NotificationContactStage   ENotificationType = 6100
	NotificationContactMessage ENotificationType = 6110
	NotificationFriendRequest  ENotificationType = 6120
	NotificationFriendResponse ENotificationType = 6130
	NotificationContactFollow  ENotificationType = 6140

	// -- action for asset --
	NotificationAssetUpdate    ENotificationType = 7010
	NotificationAssetDelete    ENotificationType = 7020
	NotificationAssetCreate    ENotificationType = 7030
	NotificationAssetMerge     ENotificationType = 7040
	NotificationAssetSplit     ENotificationType = 7050
	NotificationAssetArchive   ENotificationType = 7060
	NotificationAssetUnarchive ENotificationType = 7070
	NotificationAssetShare     ENotificationType = 7080

	// -- action for product --
	NotificationProductCreate    ENotificationType = 8010
	NotificationProductUpdate    ENotificationType = 8020
	NotificationProductDelete    ENotificationType = 8030
	NotificationProductArchive   ENotificationType = 8060
	NotificationProductUnarchive ENotificationType = 8070
	NotificationProductShare     ENotificationType = 8080
	// -- action for product child --
	NotificationProductChildCreate ENotificationType = 8090
	NotificationProductMerge       ENotificationType = 8100
	NotificationProductDevide      ENotificationType = 8110
	NotificationProductMergeAll    ENotificationType = 8120
	// -- action for investment --
	NotificationInvestmentCreated  ENotificationType = 8130
	NotificationInvestmentApproved ENotificationType = 8140
	NotificationInvestmentRejected ENotificationType = 8150
	// -- action for deal --
	NotificationDealInvitation         ENotificationType = 8160
	NotificationDealInvitationAccepted ENotificationType = 8170
	NotificationDealInvitationRejected ENotificationType = 8180
	NotificationDealMemberWithdrawn    ENotificationType = 8190
	NotificationDealMemberRemoved      ENotificationType = 8200
	NotificationDealCreate             ENotificationType = 8210
	// -- action for deal member --
	NotificationMemberJoinDeal   ENotificationType = 8220
	NotificationCustomerJoinDeal ENotificationType = 8230
	NotificationPartnerJoinDeal  ENotificationType = 8240
	// -- action for deal update status --
	NotificationDealUpdateStatus ENotificationType = 8250

	// -- action for advertising/marketing --
	NotificationCampaignCreate   ENotificationType = 9010
	NotificationCampaignUpdate   ENotificationType = 9020
	NotificationCampaignDelete   ENotificationType = 9030
	NotificationCampaignPayment  ENotificationType = 9040
	NotificationCampaignActivate ENotificationType = 9050
	NotificationCampaignPause    ENotificationType = 9060
	NotificationCampaignResume   ENotificationType = 9070
	NotificationCampaignExtend   ENotificationType = 9080
	NotificationCampaignTopup    ENotificationType = 9090
	NotificationPackageCreate    ENotificationType = 9110
	NotificationPackageUpdate    ENotificationType = 9120
	NotificationPackageDelete    ENotificationType = 9130

	NotificationCommentNewsFeed ENotificationType = 9210 //"COMMENT_NEWS_FEED"
	NotificationCommentChildren ENotificationType = 9220 //"COMMENT_REPLY"
	NotificationLikeNewsFeed    ENotificationType = 9230 //"LIKE_NEWS_FEED"
	NotificationLikeComment     ENotificationType = 9240 //"LIKE_COMMENT"
	NotificationShareNewsFeed   ENotificationType = 9250 //"SHARE_NEWS_FEED"
	NotificationShareProduct    ENotificationType = 9260 //"SHARE_PRODUCT"

	// -- action for admin --
	NotificationAdminCreate ENotificationType = 10010
	NotificationAdminUpdate ENotificationType = 10020
	NotificationAdminDelete ENotificationType = 10030

	// -- action for user management --
	NotificationUserCreate ENotificationType = 11010
	NotificationUserUpdate ENotificationType = 11020
	NotificationUserDelete ENotificationType = 11030

	// -- action for feedback/rate --
	NotificationRateCreate   ENotificationType = 12010
	NotificationRateUpdate   ENotificationType = 12020
	NotificationRateDelete   ENotificationType = 12030
	NotificationReportCreate ENotificationType = 12040
)

var NotificationMap = map[ENotificationType]string{
	NotificationCreateStep:     "Tạo bước",
	NotificationDelete:         "Xóa",
	NotificationAssignRole:     "Phân quyền",
	NotificationSettingAuto:    "Cài đặt tự động",
	NotificationReminder:       "Nhắc nhở",
	NotificationSwitchStage:    "Chuyển giai đoạn",
	NotificationUpdateNote:     "Cập nhật ghi chú",
	NotificationPipelineCreate: "Tạo quy trình",
	NotificationPipelineUpdate: "Cập nhật quy trình",
	NotificationPipelineDelete: "Xóa quy trình",
	NotificationStageCreate:    "Tạo giai đoạn",
	NotificationStageUpdate:    "Cập nhật giai đoạn",
	NotificationStageDelete:    "Xóa giai đoạn",
	NotificationRuleCreate:     "Tạo quy tắc",
	NotificationRuleUpdate:     "Cập nhật quy tắc",
	NotificationRuleDelete:     "Xóa quy tắc",

	// -- action for contact --
	NotificationContactCreate:  "Tạo liên hệ",
	NotificationContactDelete:  "Xóa liên hệ",
	NotificationContactAssign:  "Phân công liên hệ",
	NotificationContactUpdate:  "Cập nhật liên hệ",
	NotificationContactInstall: "Mời cài app",
	NotificationContactFriend:  "Gửi lời kết bạn",
	NotificationContactStage:   "Cập nhật trạng thái CRM",
	NotificationContactMessage: "Gửi tin nhắn",
	NotificationContactFollow:  "Theo dõi",

	// -- action for asset --
	NotificationAssetUpdate:    "Cập nhật tài sản",
	NotificationAssetDelete:    "Xóa tài sản",
	NotificationAssetCreate:    "Tạo tài sản",
	NotificationAssetMerge:     "Gộp tài sản",
	NotificationAssetSplit:     "Tách tài sản",
	NotificationAssetArchive:   "Lưu tài sản",
	NotificationAssetUnarchive: "Mở tài sản",
	NotificationAssetShare:     "Chia sẻ tài sản",

	// -- action for product --
	NotificationProductCreate:    "Tạo sản phẩm",
	NotificationProductUpdate:    "Cập nhật sản phẩm",
	NotificationProductDelete:    "Xóa sản phẩm",
	NotificationProductMerge:     "Gộp sản phẩm",
	NotificationProductDevide:    "Tách sản phẩm",
	NotificationProductArchive:   "Lưu sản phẩm",
	NotificationProductUnarchive: "Mở sản phẩm",
	NotificationProductShare:     "Chia sẻ sản phẩm",

	// -- action for advertising/marketing --
	NotificationCampaignCreate:   "Tạo chiến dịch quảng cáo",
	NotificationCampaignUpdate:   "Cập nhật chiến dịch quảng cáo",
	NotificationCampaignDelete:   "Xóa chiến dịch quảng cáo",
	NotificationCampaignPayment:  "Thanh toán chiến dịch",
	NotificationCampaignActivate: "Kích hoạt chiến dịch",
	NotificationCampaignPause:    "Tạm dừng chiến dịch",
	NotificationCampaignResume:   "Tiếp tục chiến dịch",
	NotificationCampaignExtend:   "Gia hạn chiến dịch",
	NotificationCampaignTopup:    "Nạp tiền chiến dịch",
	NotificationPackageCreate:    "Tạo gói quảng cáo",
	NotificationPackageUpdate:    "Cập nhật gói quảng cáo",
	NotificationPackageDelete:    "Xóa gói quảng cáo",

	// -- action for admin --
	NotificationAdminCreate: "Tạo admin",
	NotificationAdminUpdate: "Cập nhật admin",
	NotificationAdminDelete: "Xóa admin",

	// -- action for user management --
	NotificationUserCreate: "Tạo người dùng",
	NotificationUserUpdate: "Cập nhật người dùng",
	NotificationUserDelete: "Xóa người dùng",

	// -- action for feedback/rate --
	NotificationRateCreate:   "Tạo đánh giá",
	NotificationRateUpdate:   "Cập nhật đánh giá",
	NotificationRateDelete:   "Xóa đánh giá",
	NotificationReportCreate: "Tạo báo cáo",
}

var NotificationContactAction = []ENotificationType{
	NotificationCreateStep,
	NotificationDelete,
	NotificationAssignRole,
	NotificationSettingAuto,
	NotificationReminder,
	NotificationSwitchStage,
	NotificationUpdateNote,
	NotificationPipelineCreate,
	NotificationPipelineUpdate,
	NotificationPipelineDelete,
	NotificationStageCreate,
	NotificationStageUpdate,
	NotificationStageDelete,
	NotificationRuleCreate,
	NotificationRuleUpdate,
	NotificationRuleDelete,
	NotificationContactCreate,
	NotificationContactDelete,
	NotificationContactAssign,
	NotificationContactUpdate,
	NotificationContactInstall,
	NotificationContactFriend,
	NotificationContactStage,
	NotificationContactMessage,
	NotificationContactFollow,
}

var NotificationAssetActionType = []ENotificationType{
	NotificationAssetUpdate,
	NotificationAssetDelete,
	NotificationAssetCreate,
	NotificationAssetMerge,
	NotificationAssetSplit,
	NotificationAssetArchive,
	NotificationAssetUnarchive,
	NotificationAssetShare,
}

var NotificationProductActionType = []ENotificationType{
	NotificationProductCreate,
	NotificationProductUpdate,
	NotificationProductDelete,
	NotificationProductArchive,
	NotificationProductUnarchive,
	NotificationProductShare,
	NotificationProductChildCreate,
	NotificationProductMerge,
	NotificationProductMergeAll,
	NotificationProductDevide,
}

var NotificationProductChild = []ENotificationType{
	NotificationProductMerge,
	NotificationProductMergeAll,
	NotificationProductDevide,
	NotificationProductChildCreate,
}

var NotificationCampaignActionType = []ENotificationType{
	NotificationCampaignCreate,
	NotificationCampaignUpdate,
	NotificationCampaignDelete,
	NotificationCampaignPayment,
	NotificationCampaignActivate,
	NotificationCampaignPause,
	NotificationCampaignResume,
	NotificationCampaignExtend,
	NotificationCampaignTopup,
}

var NotificationPackageActionType = []ENotificationType{
	NotificationPackageCreate,
	NotificationPackageUpdate,
	NotificationPackageDelete,
}

// NotificationMergeTypes - Danh sách các notification types cần merge (xóa thông báo cũ trước khi tạo mới)
var NotificationMergeTypes = []ENotificationType{
	NotificationFriendRequest,
	NotificationFriendResponse,
	NotificationContactFriend,
	NotificationContactFollow,
	NotificationLikeNewsFeed,
	NotificationLikeComment,
	NotificationCommentNewsFeed,
	NotificationCommentChildren,
	NotificationShareNewsFeed,
	NotificationShareProduct,
}

// IsMergeableNotification kiểm tra xem notification type có cần merge không
func (s ENotificationType) IsMergeableNotification() bool {
	return slices.Contains(NotificationMergeTypes, s)
}
