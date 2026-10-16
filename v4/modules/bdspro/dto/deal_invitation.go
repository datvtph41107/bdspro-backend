package dto

// DealInvitationListRequest is GET /v4/deals/{id}/invitations.
type DealInvitationListRequest struct {
	DealID    uint64 `path:"id|dealId"`
	Page      uint32 `form:"page"`
	Size      uint32 `form:"size"`
	Timestamp int64  `form:"timestamp"`
}

// SendDealInvitationRequest is POST /v4/deals/{id}/invitation.
type SendDealInvitationRequest struct {
	DealID   uint64  `path:"id|dealId"`
	MemberID uint64  `json:"memberId"`
	Message  string  `json:"message"`
	RoleID   uint64  `json:"roleId"`
	RoleKey  uint32  `json:"roleKey"`
	ColorID  *uint32 `json:"colorId"`
}

// AcceptDealInvitationRequest is POST /v4/deals/{id}/invitation/accept.
type AcceptDealInvitationRequest struct {
	DealID       uint64 `path:"id|dealId"`
	InvitationID uint64 `json:"invitationId"`
}

// RejectDealInvitationRequest is POST /v4/deals/{id}/invitation/reject.
type RejectDealInvitationRequest struct {
	DealID       uint64 `path:"id|dealId"`
	InvitationID uint64 `json:"invitationId"`
	Reason       string `json:"reason"`
}

// RemoveFromDealRequest is POST /v4/deals/invitation/remove.
type RemoveFromDealRequest struct {
	InvitationID uint64 `json:"invitationId"`
	Reason       string `json:"reason"`
}

// UpdateDealMemberRoleRequest is PUT /v4/deals/invitation/role.
type UpdateDealMemberRoleRequest struct {
	DealID   uint64 `json:"dealId"`
	MemberID uint64 `json:"memberId"`
	RoleKey  uint32 `json:"roleKey"`
}
