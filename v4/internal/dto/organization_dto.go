package dto

type GroupMember struct {
	ID      uint64  `json:"id"`
	GroupId uint64  `json:"groupId"`
	UserId  uint64  `json:"userId"`
	Role    string  `json:"role"`
	RoleId  *uint64 `json:"roleId"`
}

type DealMember struct {
	DealId uint64  `json:"dealId"`
	UserId uint64  `json:"userId"`
	RoleId *uint64 `json:"roleId"`

	RoleKey uint32 `json:"roleKey"`
}

type BranchMember struct {
	BranchId uint64  `json:"branchId"`
	UserId   uint64  `json:"userId"`
	RoleId   *uint64 `json:"roleId"`
}

type OrganizationMember struct {
	OrganizationId uint64  `json:"organizationId"`
	UserId         uint64  `json:"userId"`
	RoleId         *uint64 `json:"roleId"`
	RoleKey        uint32  `json:"roleKey"`
	Status         uint32  `json:"status"`
}
