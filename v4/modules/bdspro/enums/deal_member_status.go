package enums

// DealMemberStatus matches bdspro domain.DealMemberStatus.
type DealMemberStatus uint32

const (
	DealMemberStatusInvited   DealMemberStatus = 10 // Đã mời
	DealMemberStatusAccepted  DealMemberStatus = 20 // Đã tham gia
	DealMemberStatusRejected  DealMemberStatus = 30 // Từ chối
	DealMemberStatusWithdrawn DealMemberStatus = 40 // Đã rút khỏi thương vụ
)

var DealMemberStatusMap = map[DealMemberStatus]string{
	DealMemberStatusInvited:   "Đã mời",
	DealMemberStatusAccepted:  "Đã tham gia",
	DealMemberStatusRejected:  "Từ chối",
	DealMemberStatusWithdrawn: "Đã rút khỏi thương vụ",
}
