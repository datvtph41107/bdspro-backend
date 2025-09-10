package dto

// UpdateDealMemberRoleRequest request để cập nhật vai trò thành viên thương vụ
type UpdateDealMemberRoleRequest struct {
	DealID   uint64 `json:"deal_id" validate:"required"`
	MemberID uint64 `json:"member_id" validate:"required"`
	RoleKey  uint32 `json:"role_key" validate:"required"`
}

// UpdateDealMemberRoleResponse response khi cập nhật vai trò thành viên thương vụ
type UpdateDealMemberRoleResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	DealID   uint64 `json:"deal_id"`
	MemberID uint64 `json:"member_id"`
	RoleKey  uint32 `json:"role_key"`
}
