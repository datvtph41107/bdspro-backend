package enums

// DealMemberRoleKey enum cho role trong deal member
type DealMemberRoleKey uint32

const (
	DealMemberRoleKeyOwner    DealMemberRoleKey = 410 // Chủ thương vụ
	DealMemberRoleKeyAdmin    DealMemberRoleKey = 420 // Quản trị viên
	DealMemberRoleKeyMember   DealMemberRoleKey = 430 // Thành viên
	DealMemberRoleKeyPartner  DealMemberRoleKey = 440 // Đối tác
	DealMemberRoleKeyCustomer DealMemberRoleKey = 450 // Khách hàng
)

var DealMemberRoleKeyMap = map[DealMemberRoleKey]string{
	DealMemberRoleKeyOwner:    "Chủ sở hữu",
	DealMemberRoleKeyAdmin:    "Quản lý",
	DealMemberRoleKeyMember:   "Thành viên",
	DealMemberRoleKeyCustomer: "Đối tác",
	DealMemberRoleKeyPartner:  "Khách hàng",
}
