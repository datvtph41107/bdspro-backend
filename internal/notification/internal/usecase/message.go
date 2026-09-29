package usecase

import (
	"context"
	"notification/internal/enums"
)

// @deprecated
func (s *NotificationUsecase) GetMessage(c context.Context, typeMessage enums.TypeEnum, attachData []string) []string {
	if typeMessage == enums.ShareNewsFeed {
		return []string{"", attachData[0], " đã chia sẻ bài viết của bạn: \"" + attachData[1] + "\""}
	}
	if typeMessage == enums.CommentNewsFeed {
		return []string{"", attachData[0], " đã bình luận bài viết của bạn: \"" + attachData[1] + "\""}
	}
	if typeMessage == enums.CommentChildren {
		return []string{"", attachData[0], " đã trả lời bình luận của bạn: \"" + attachData[1] + "\""}
	}
	if typeMessage == enums.LikeNewsFeed {
		return []string{"", attachData[0], " thích bài viết của bạn: \"" + attachData[1] + "\""}
	}
	if typeMessage == enums.LikeComment {
		return []string{"", attachData[0], " thích bình luận của bạn: \"" + attachData[1] + "\""}
	}
	if typeMessage == enums.Follow {
		return []string{"", attachData[0], " đã theo dõi bạn"}
	}
	if typeMessage == enums.FriendRequest && len(attachData) > 0 {
		return []string{"", attachData[0], " đã gửi lời mời kết bạn cho bạn"}
	}
	if typeMessage == enums.FriendRequest {
		return []string{"", attachData[0], " đã đồng ý kết bạn"}
	}
	if typeMessage == enums.JoinGroup {
		return []string{"", attachData[0], " đã thêm bạn vào nhóm: \"" + attachData[1] + "\""}
	}
	if typeMessage == enums.JoinOrganization {
		return []string{"", "Quản trị viên đã thêm bạn vào tổ chức: \"" + attachData[1] + "\""}
	}
	if typeMessage == enums.Deposit {
		return []string{"", attachData[0], " đã nạp tiền vào tài khoản của bạn"}
	}
	if typeMessage == enums.Withdraw {
		return []string{"", attachData[0], " đã rút tiền từ tài khoản của bạn"}
	}
	if typeMessage == enums.Payment {
		return []string{"", attachData[0], " đã thanh toán cho bạn"}
	}

	// Tổ chức
	if typeMessage == enums.MemberJoinDeal {
		return []string{"Bạn đã trở thành thành viên của thương vụ: \"", attachData[0], "\""}
	}
	if typeMessage == enums.CustomerJoinDeal {
		return []string{"Bạn đã trở thành khách hàng của thương vụ: \"", attachData[0], "\""}
	}
	if typeMessage == enums.PartnerJoinDeal {
		return []string{"Bạn đã trở thành đối tác của thương vụ: \"", attachData[0], "\""}
	}
	if typeMessage == enums.DealUpdateStatus {
		return []string{"Thương vụ \"", attachData[0], "\" đã được cập nhật trạng thái thành: \"", attachData[1], "\""}
	}
	if typeMessage == enums.InvestmentApproved {
		// return []string{"Bạn đã được phê duyệt khai báo góp vốn cho thương vụ: \"", attachData[0], "\" với số tiền: ", attachData[1]}
		return []string{"Bạn đã đồng ý tham gia thương vụ"}
	}
	if typeMessage == enums.InvestmentRejected {
		// return []string{"Bạn đã bị từ chối khai báo góp vốn cho thương vụ: \"", attachData[0], "\" với số tiền: ", attachData[1]}
		return []string{"Bạn đã bị chối khai báo góp vốn"}
	}

	if typeMessage == enums.TransactionDeposit {
		return []string{"Bạn đã nạp tiền vào tài khoản của bạn: ", attachData[0]}
	}
	if typeMessage == enums.TransactionWithdraw {
		return []string{"Bạn đã rút tiền từ tài khoản của bạn: ", attachData[0]}
	}
	if typeMessage == enums.TransactionPayment {
		return []string{"Bạn đã thanh toán: ", attachData[0]}
	}

	if typeMessage == enums.NotiDealInvitation {
		return []string{"Bạn đã nhận được lời mời tham gia thương vụ: \"", attachData[0], "\""}
	}
	if typeMessage == enums.NotiDealInvitationAccepted {
		return []string{"Bạn đã chấp nhận lời mời tham gia thương vụ: \"", attachData[0], "\""}
	}
	if typeMessage == enums.NotiDealInvitationRejected {
		return []string{"Bạn đã từ chối lời mời tham gia thương vụ: \"", attachData[0], "\""}
	}
	if typeMessage == enums.NotiDealMemberWithdrawn {
		return []string{"Bạn đã rút khỏi thương vụ: \"", attachData[0], "\""}
	}
	if typeMessage == enums.NotiDealMemberRemoved {
		return []string{"Bạn đã bị gỡ khỏi thương vụ: \"", attachData[0], "\""}
	}

	return []string{}
}
