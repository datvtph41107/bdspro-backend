package enums

// RoleKey enum cho role trong deal member
type RoleKey uint32

const (

	// thương vụ: key nhóm 400
	RoleKeyDealOwner    RoleKey = 410 // Chủ thương vụ
	RoleKeyDealAdmin    RoleKey = 420 // Quản trị viên
	RoleKeyDealMember   RoleKey = 430 // Thành viên
	RoleKeyDealPartner  RoleKey = 440 // Đối tác
	RoleKeyDealCustomer RoleKey = 450 // Khách hàng

	// tổ chức: key nhóm 500
	RoleKeyOrganizationOwner  RoleKey = 510 // Chủ sở hữu
	RoleKeyOrganizationAdmin  RoleKey = 520 // Quản trị viên
	RoleKeyOrganizationMember RoleKey = 530 // Thành viên
)

var RoleKeyMap = map[RoleKey]string{
	RoleKeyDealOwner:    "Chủ sở hữu",
	RoleKeyDealAdmin:    "Quản lý",
	RoleKeyDealMember:   "Thành viên",
	RoleKeyDealCustomer: "Đối tác",
	RoleKeyDealPartner:  "Khách hàng",

	RoleKeyOrganizationOwner:  "Chủ sở hữu",
	RoleKeyOrganizationAdmin:  "Quản lý",
	RoleKeyOrganizationMember: "Thành viên",
}
