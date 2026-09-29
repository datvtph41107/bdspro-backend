package enums

// TypeEnum định nghĩa các loại thông báo
type TypeEnum int16

const (
	// Crm
	FriendRequest    TypeEnum = 10  //"FRIEND_REQUEST"
	FriendResponse   TypeEnum = 20  //"FRIEND_RESPONSE"
	Follow           TypeEnum = 30  //"FOLLOW"
	JoinGroup        TypeEnum = 80  //"JOIN_GROUP"
	JoinOrganization TypeEnum = 90  //"JOIN_ORGANIZATION"
	Deposit          TypeEnum = 110 //"DEPOSIT"
	Withdraw         TypeEnum = 111 //"WITHDRAW"
	Payment          TypeEnum = 112 //"PAYMENT"

	// Social
	CommentNewsFeed TypeEnum = 40  //"COMMENT_NEWS_FEED"
	CommentChildren TypeEnum = 41  //"COMMENT_REPLY"
	LikeNewsFeed    TypeEnum = 50  //"LIKE_NEWS_FEED"
	LikeComment     TypeEnum = 60  //"LIKE_COMMENT"
	ShareNewsFeed   TypeEnum = 70  //"SHARE_NEWS_FEED"
	ShareProduct    TypeEnum = 100 //"SHARE_PRODUCT"

	// Organization
	MemberJoinDeal     TypeEnum = 7001 //"MEMBER_JOIN_DEAL"
	CustomerJoinDeal   TypeEnum = 7002 //"CUSTOMER_JOIN_DEAL"
	PartnerJoinDeal    TypeEnum = 7003 //"PARTNER_JOIN_DEAL"
	DealUpdateStatus   TypeEnum = 7004 //"DEAL_UPDATE_STATUS"
	InvestmentApproved TypeEnum = 7005 //"INVESTMENT_APPROVED"
	InvestmentRejected TypeEnum = 7006 //"INVESTMENT_REJECTED"

	NotiDealInvitation         TypeEnum = 7008
	NotiDealInvitationAccepted TypeEnum = 7009
	NotiDealInvitationRejected TypeEnum = 7010
	NotiDealMemberWithdrawn    TypeEnum = 7011
	NotiDealMemberRemoved      TypeEnum = 7012

	// Transaction
	TransactionDeposit  TypeEnum = 5001 //"TRANSACTION_DEPOSIT"
	TransactionWithdraw TypeEnum = 5002 //"TRANSACTION_WITHDRAW"
	TransactionPayment  TypeEnum = 5003 //"TRANSACTION_PAYMENT"

	// Warning/Alert Types
	WarningViolation TypeEnum = 8001 //"WARNING_VIOLATION" - Cảnh báo vi phạm
	WarningReminder  TypeEnum = 8002 //"WARNING_REMINDER" - Nhắc nhở chung
	WarningGuidance  TypeEnum = 8003 //"WARNING_GUIDANCE" - Hướng dẫn sử dụng
)

var TypeEnumMap = map[TypeEnum]string{
	FriendRequest:    "Gửi yêu cầu kết bạn",
	FriendResponse:   "Kết bạn",
	Follow:           "Theo dõi",
	CommentNewsFeed:  "Bình luận bài viết",
	CommentChildren:  "Bình luận bài viết",
	LikeNewsFeed:     "Thích bài viết",
	LikeComment:      "Thích bình luận",
	ShareNewsFeed:    "Chia sẻ bài viết",
	JoinGroup:        "Tham gia nhóm",
	JoinOrganization: "Tham gia tổ chức",
	Deposit:          "Nạp tiền",
	Withdraw:         "Rút tiền",
	Payment:          "Thanh toán",

	// Tổ chức
	MemberJoinDeal:     "Thêm thành viên vào thương vụ",
	CustomerJoinDeal:   "Thêm khách hàng vào thương vụ",
	PartnerJoinDeal:    "Thêm đối tác vào thương vụ",
	DealUpdateStatus:   "Cập nhật trạng thái thương vụ",
	InvestmentApproved: "Khai báo góp vốn đã được phê duyệt",
	InvestmentRejected: "Khai báo góp vốn đã bị từ chối",

	// Cảnh báo
	WarningViolation: "Cảnh báo vi phạm",
	WarningReminder:  "Nhắc nhở chung",
	WarningGuidance:  "Hướng dẫn sử dụng",
}

// IsValid kiểm tra xem giá trị có hợp lệ không
func (t TypeEnum) IsValid() bool {
	switch t {
	case FriendRequest, FriendResponse, Follow, CommentNewsFeed, CommentChildren, LikeNewsFeed, LikeComment, ShareNewsFeed, JoinGroup, JoinOrganization, Deposit, Withdraw, Payment, MemberJoinDeal, CustomerJoinDeal, PartnerJoinDeal, DealUpdateStatus, InvestmentApproved, InvestmentRejected, NotiDealInvitation, NotiDealInvitationAccepted, NotiDealInvitationRejected, NotiDealMemberWithdrawn, NotiDealMemberRemoved, TransactionDeposit, TransactionWithdraw, TransactionPayment, WarningViolation, WarningReminder, WarningGuidance:
		return true
	default:
		return false
	}
}
