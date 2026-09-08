package dto

import (
	"organization/internal/domain/entity"
	sharepb "pb/types/shared"
)

type DealInvitationWithProfile struct {
	*entity.DealMember
	Inviter *sharepb.ProfileItem
	Invitee *sharepb.ProfileItem
}

type CreateDealInvitationRequest struct {
	DealID       uint64   `json:"deal_id" validate:"required"`
	InviteeID    uint64   `json:"invitee_id" validate:"required"`
	MemberType   uint32   `json:"member_type" validate:"required"`
	Message      string   `json:"message"`
	AmountCommit *float64 `json:"amount_commit"`
}

type AcceptInvitationRequest struct {
	InvitationID uint64 `json:"invitation_id" validate:"required"`
}

type RejectInvitationRequest struct {
	InvitationID uint64 `json:"invitation_id" validate:"required"`
	Reason       string `json:"reason"`
}

type ResendInvitationRequest struct {
	InvitationID uint64 `json:"invitation_id" validate:"required"`
}

type WithdrawFromDealRequest struct {
	InvitationID uint64 `json:"invitation_id" validate:"required"`
	Reason       string `json:"reason"`
}

type RemoveFromDealRequest struct {
	InvitationID uint64 `json:"invitation_id" validate:"required"`
	Reason       string `json:"reason"`
}

type ConfirmWithdrawalRequest struct {
	InvitationID uint64 `json:"invitation_id" validate:"required"`
}

type SearchMembersRequest struct {
	DealID  uint64 `json:"deal_id" validate:"required"`
	Keyword string `json:"keyword"`
	Page    int    `json:"page"`
	Size    int    `json:"size"`

	DoneInvestment *bool `json:"done_investment"`
}
