package enums

type NotificationType int32

const (
	NotiDealCreate NotificationType = 0

	NotiMemberJoinDeal   NotificationType = 7001
	NotiCustomerJoinDeal NotificationType = 7002
	NotiPartnerJoinDeal  NotificationType = 7003

	NotiDealUpdateStatus NotificationType = 7004

	NotiInvestmentCreated NotificationType = 7005
	NotiInvestmentApproved NotificationType = 7006
	NotiInvestmentRejected NotificationType = 7007
)

var NotiTypeMap = map[NotificationType]string{
	NotiDealCreate:         "Tạo thương vụ",
	NotiMemberJoinDeal:     "Thêm thành viên vào thương vụ",
	NotiCustomerJoinDeal:   "Thêm khách hàng vào thương vụ",
	NotiPartnerJoinDeal:    "Thêm đối tác vào thương vụ",
	NotiDealUpdateStatus:   "Cập nhật trạng thái thương vụ",
	NotiInvestmentCreated:  "Gửi thông báo góp vốn",
	NotiInvestmentApproved: "Khai báo góp vốn đã được phê duyệt",
	NotiInvestmentRejected: "Khai báo góp vốn đã bị từ chối",
}
