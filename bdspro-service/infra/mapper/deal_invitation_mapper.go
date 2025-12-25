package mapper

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	_models "common/models"
	_utils "common/utils"
	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"
	"time"
)

type DealInvitationTransformer interface {
	// Request to Entity transformations
	SendInvitationRequestToEntity(req *bdspropb.SendDealInvitationRequest) *domain.DealMember
	AcceptInvitationRequestToEntity(req *bdspropb.AcceptDealInvitationRequest) *domain.DealMember
	RejectInvitationRequestToEntity(req *bdspropb.RejectDealInvitationRequest) *domain.DealMember
	ResendInvitationRequestToEntity(req *bdspropb.ResendDealInvitationRequest) *domain.DealMember
	WithdrawFromDealRequestToEntity(req *bdspropb.WithdrawFromDealRequest) *domain.DealMember
	RemoveFromDealRequestToEntity(req *bdspropb.RemoveFromDealRequest) *domain.DealMember
	ConfirmWithdrawalRequestToEntity(req *bdspropb.ConfirmWithdrawalRequest) *domain.DealMember

	// Entity to Response transformations
	EntityToSendInvitationResponse(invitation *domain.DealMember) *bdspropb.SendDealInvitationResponse
	EntityToAcceptInvitationResponse(invitation *domain.DealMember) *bdspropb.AcceptDealInvitationResponse
	EntityToRejectInvitationResponse(invitation *domain.DealMember) *bdspropb.RejectDealInvitationResponse
	EntityToResendInvitationResponse(invitation *domain.DealMember) *bdspropb.ResendDealInvitationResponse
	EntityToWithdrawFromDealResponse(invitation *domain.DealMember) *bdspropb.WithdrawFromDealResponse
	EntityToRemoveFromDealResponse(invitation *domain.DealMember) *bdspropb.RemoveFromDealResponse
	EntityToConfirmWithdrawalResponse(invitation *domain.DealMember) *bdspropb.ConfirmWithdrawalResponse
	EntityToGetInvitationResponse(invitation *domain.DealMember) *bdspropb.GetDealInvitationResponse
	EntitiesToGetAcceptedMembersResponse(members []*dto.DealInvitationWithProfile) *bdspropb.GetAcceptedMembersResponse
	EntitiesToSearchMembersResponse(members []*dto.DealInvitationWithProfile, total uint32) *bdspropb.SearchDealMembersResponse
	EntitiesToGetPendingInvitationsResponse(invitations []*dto.DealInvitationWithProfile, total uint32) *bdspropb.GetPendingInvitationsResponse
	EntitiesToGetInvitationsByDealResponse(invitations []*dto.DealInvitationWithProfile, total uint32) *bdspropb.GetDealInvitationsByDealResponse
}

type dealInvitationTransformer struct {
	// bdsproRoleTransformer OrganizationRoleTransformer
}

func NewDealInvitationTransformer() DealInvitationTransformer {
	return &dealInvitationTransformer{
		// bdsproRoleTransformer: bdsproRoleTransformer,
	}
}

// Request to Entity transformations
func (t *dealInvitationTransformer) SendInvitationRequestToEntity(req *bdspropb.SendDealInvitationRequest) *domain.DealMember {
	// Convert string role to DealMemberType enum
	// var memberType enums.DealMemberType
	// switch req.Role {
	// case "customer":
	// 	memberType = enums.DealMemberTypeCustomer
	// case "partner":
	// 	memberType = enums.DealMemberTypePartner
	// case "member":
	// 	memberType = enums.DealMemberTypeMember
	// default:
	// 	memberType = enums.DealMemberTypeMember // default
	// }

	invitation := &domain.DealMember{
		DealID:    req.DealId,
		InviterID: req.InviterId,
		MemberID:  req.MemberId,
		Message:   req.Message,
		RoleID:    req.RoleId,
		// MemberType: memberType,
		Status:    domain.DealMemberStatusInvited,
		InvitedAt: time.Now(),
		RoleKey:   enums.RoleKey(req.RoleKey),
		// CreatedAt: time.Now(),
		// UpdatedAt: time.Now(),
	}

	if req.ColorId != nil {
		invitation.ColorId = req.ColorId
	}

	return invitation
}

func (t *dealInvitationTransformer) AcceptInvitationRequestToEntity(req *bdspropb.AcceptDealInvitationRequest) *domain.DealMember {
	return &domain.DealMember{
		BaseEntity: _models.BaseEntity{
			ID: req.InvitationId,
		},
	}
}

func (t *dealInvitationTransformer) RejectInvitationRequestToEntity(req *bdspropb.RejectDealInvitationRequest) *domain.DealMember {
	return &domain.DealMember{
		BaseEntity: _models.BaseEntity{
			ID: req.InvitationId,
		},
		Message: req.Reason,
	}
}

func (t *dealInvitationTransformer) ResendInvitationRequestToEntity(req *bdspropb.ResendDealInvitationRequest) *domain.DealMember {
	return &domain.DealMember{
		BaseEntity: _models.BaseEntity{
			ID: req.InvitationId,
		},
	}
}

func (t *dealInvitationTransformer) WithdrawFromDealRequestToEntity(req *bdspropb.WithdrawFromDealRequest) *domain.DealMember {
	return &domain.DealMember{
		BaseEntity: _models.BaseEntity{
			ID: req.InvitationId,
		},
		Message: req.Reason,
	}
}

func (t *dealInvitationTransformer) RemoveFromDealRequestToEntity(req *bdspropb.RemoveFromDealRequest) *domain.DealMember {
	return &domain.DealMember{
		BaseEntity: _models.BaseEntity{
			ID: req.InvitationId,
		},
		Message: req.Reason,
	}
}

func (t *dealInvitationTransformer) ConfirmWithdrawalRequestToEntity(req *bdspropb.ConfirmWithdrawalRequest) *domain.DealMember {
	return &domain.DealMember{
		BaseEntity: _models.BaseEntity{
			ID: req.InvitationId,
		},
	}
}

// Entity to Response transformations
func (t *dealInvitationTransformer) EntityToSendInvitationResponse(invitation *domain.DealMember) *bdspropb.SendDealInvitationResponse {
	return &bdspropb.SendDealInvitationResponse{
		InvitationId: invitation.ID,
		Success:      true,
		Message:      "Invitation sent successfully",
		RoleId:       invitation.RoleID,
	}
}

func (t *dealInvitationTransformer) EntityToAcceptInvitationResponse(invitation *domain.DealMember) *bdspropb.AcceptDealInvitationResponse {
	return &bdspropb.AcceptDealInvitationResponse{
		InvitationId: invitation.ID,
		Success:      true,
	}
}

func (t *dealInvitationTransformer) EntityToRejectInvitationResponse(invitation *domain.DealMember) *bdspropb.RejectDealInvitationResponse {
	return &bdspropb.RejectDealInvitationResponse{
		InvitationId: invitation.ID,
		Success:      true,
	}
}

func (t *dealInvitationTransformer) EntityToResendInvitationResponse(invitation *domain.DealMember) *bdspropb.ResendDealInvitationResponse {
	return &bdspropb.ResendDealInvitationResponse{
		InvitationId: invitation.ID,
		Success:      true,
		Message:      "Invitation resent successfully",
	}
}

func (t *dealInvitationTransformer) EntityToWithdrawFromDealResponse(invitation *domain.DealMember) *bdspropb.WithdrawFromDealResponse {
	return &bdspropb.WithdrawFromDealResponse{
		InvitationId: invitation.ID,
		Success:      true,
	}
}

func (t *dealInvitationTransformer) EntityToRemoveFromDealResponse(invitation *domain.DealMember) *bdspropb.RemoveFromDealResponse {
	return &bdspropb.RemoveFromDealResponse{
		InvitationId: invitation.ID,
		Success:      true,
	}
}

func (t *dealInvitationTransformer) EntityToConfirmWithdrawalResponse(invitation *domain.DealMember) *bdspropb.ConfirmWithdrawalResponse {
	return &bdspropb.ConfirmWithdrawalResponse{
		InvitationId: invitation.ID,
		Success:      true,
	}
}

func (t *dealInvitationTransformer) EntityToGetInvitationResponse(invitation *domain.DealMember) *bdspropb.GetDealInvitationResponse {
	if invitation == nil {
		return nil
	}

	// Convert DealMemberType to string
	// var role string
	// switch invitation.MemberType {
	// case enums.DealMemberTypeCustomer:
	// 	role = "customer"
	// case enums.DealMemberTypePartner:
	// 	role = "partner"
	// case enums.DealMemberTypeMember:
	// 	role = "member"
	// default:
	// 	role = "member"
	// }

	invitationResponse := &sharepb.DealMember{
		Id:        invitation.ID,
		DealId:    invitation.DealID,
		InviterId: invitation.InviterID,
		MemberId:  invitation.MemberID,
		Status:    uint32(invitation.Status),
		Message:   invitation.Message,
		RoleKey:   uint32(invitation.RoleKey),
		// Role:      role,
		Reason: "", // DealInvitation entity doesn't have Reason field
		// CreatedAt: _utils.FormatTimeToString(&invitation.CreatedAt),
		// UpdatedAt: _utils.FormatTimeToString(&invitation.UpdatedAt),
	}

	if invitation.Color != nil {
		invitationResponse.Color = &sharepb.Color{
			Id:              invitation.Color.ID,
			Name:            invitation.Color.Name,
			ContentColor:    invitation.Color.ContentColor,
			BackgroundColor: invitation.Color.BackgroundColor,
			Description:     invitation.Color.Description,
			IsActive:        invitation.Color.IsActive,
		}
	}

	return &bdspropb.GetDealInvitationResponse{
		Member: &sharepb.DealMember{
			Id:        invitation.ID,
			DealId:    invitation.DealID,
			InviterId: invitation.InviterID,
			MemberId:  invitation.MemberID,
			Status:    uint32(invitation.Status),
			Message:   invitation.Message,
			RoleKey:   uint32(invitation.RoleKey),
			// Role:      role,
			Reason: "", // DealInvitation entity doesn't have Reason field
			// CreatedAt: _utils.FormatTimeToString(&invitation.CreatedAt),
			// UpdatedAt: _utils.FormatTimeToString(&invitation.UpdatedAt),
		},
	}
}

func (t *dealInvitationTransformer) EntitiesToGetAcceptedMembersResponse(members []*dto.DealInvitationWithProfile) *bdspropb.GetAcceptedMembersResponse {
	result := make([]*sharepb.DealMember, len(members))
	for i, member := range members {
		// Convert DealMemberType to string
		// var role string
		// switch member.DealMember.MemberType {
		// case enums.DealMemberTypeCustomer:
		// 	role = "customer"
		// case enums.DealMemberTypePartner:
		// 	role = "partner"
		// case enums.DealMemberTypeMember:
		// 	role = "member"
		// default:
		// 	role = "member"
		// }

		result[i] = &sharepb.DealMember{
			MemberId: member.DealMember.MemberID,
			Profile:  member.Invitee,
			// Name:     member.Invitee.FullName,
			// Email:    "", // ProfileItem doesn't have email field
			RoleId:      &member.DealMember.RoleID,
			RoleKey:     uint32(member.DealMember.RoleKey),
			RespondedAt: _utils.FormatTimeToString(member.DealMember.RespondedAt),
		}
	}

	return &bdspropb.GetAcceptedMembersResponse{
		Members: result,
	}
}

func (t *dealInvitationTransformer) EntitiesToSearchMembersResponse(members []*dto.DealInvitationWithProfile, total uint32) *bdspropb.SearchDealMembersResponse {
	result := make([]*sharepb.DealMember, len(members))
	for i, member := range members {
		// Convert DealMemberType to string
		// var role string
		// switch member.DealMember.MemberType {
		// case enums.DealMemberTypeCustomer:
		// 	role = "customer"
		// case enums.DealMemberTypePartner:
		// 	role = "partner"
		// case enums.DealMemberTypeMember:
		// 	role = "member"
		// default:
		// 	role = "member"
		// }

		result[i] = &sharepb.DealMember{
			MemberId: member.DealMember.MemberID,
			Profile:  member.Invitee,
			// Name:     member.Invitee.FullName,
			// Email:    "", // ProfileItem doesn't have email field
			RoleId:      &member.DealMember.RoleID,
			RoleKey:     uint32(member.DealMember.RoleKey),
			RespondedAt: _utils.FormatTimeToString(member.DealMember.RespondedAt),
		}
	}

	return &bdspropb.SearchDealMembersResponse{
		Members: result,
		Total:   int64(total),
	}
}

func (t *dealInvitationTransformer) EntitiesToGetPendingInvitationsResponse(invitations []*dto.DealInvitationWithProfile, total uint32) *bdspropb.GetPendingInvitationsResponse {
	result := make([]*sharepb.DealMember, len(invitations))
	for i, invitation := range invitations {
		// Convert DealMemberType to string
		// var role string
		// switch invitation.DealMember.MemberType {
		// case enums.DealMemberTypeCustomer:
		// 	role = "customer"
		// case enums.DealMemberTypePartner:
		// 	role = "partner"
		// case enums.DealMemberTypeMember:
		// 	role = "member"
		// default:
		// 	role = "member"
		// }

		invitationResponse := &sharepb.DealMember{
			Id:        invitation.DealMember.ID,
			DealId:    invitation.DealMember.DealID,
			InviterId: invitation.DealMember.InviterID,
			MemberId:  invitation.DealMember.MemberID,
			Status:    uint32(invitation.DealMember.Status),
			Message:   invitation.DealMember.Message,
			RoleKey:   uint32(invitation.DealMember.RoleKey),
			// Role:      role,
			Reason: "", // DealInvitation entity doesn't have Reason field
			// CreatedAt: _utils.FormatTimeToString(&invitation.Create),
			// UpdatedAt: _utils.FormatTimeToString(&invitation.DealMember.UpdatedAt),
		}

		if invitation.DealMember.Color != nil {
			invitationResponse.Color = &sharepb.Color{
				Id:              invitation.DealMember.Color.ID,
				Name:            invitation.DealMember.Color.Name,
				ContentColor:    invitation.DealMember.Color.ContentColor,
				BackgroundColor: invitation.DealMember.Color.BackgroundColor,
				Description:     invitation.DealMember.Color.Description,
				IsActive:        invitation.DealMember.Color.IsActive,
			}
		}

		result[i] = invitationResponse
	}
	return &bdspropb.GetPendingInvitationsResponse{
		Data:  result,
		Total: int64(total),
	}
}

func (t *dealInvitationTransformer) EntitiesToGetInvitationsByDealResponse(invitations []*dto.DealInvitationWithProfile, total uint32) *bdspropb.GetDealInvitationsByDealResponse {
	result := make([]*sharepb.DealMember, len(invitations))
	for i, invitation := range invitations {
		// Convert DealMemberType to string
		// var role string
		// switch invitation.DealMember.MemberType {
		// case enums.DealMemberTypeCustomer:
		// 	role = "customer"
		// case enums.DealMemberTypePartner:
		// 	role = "partner"
		// case enums.DealMemberTypeMember:
		// 	role = "member"
		// default:
		// 	role = "member"
		// }

		invitationResponse := &sharepb.DealMember{
			Id:        invitation.DealMember.ID,
			DealId:    invitation.DealMember.DealID,
			InviterId: invitation.DealMember.InviterID,
			MemberId:  invitation.DealMember.MemberID,
			Status:    uint32(invitation.DealMember.Status),
			Message:   invitation.DealMember.Message,
			RoleKey:   uint32(invitation.DealMember.RoleKey),
			Reason:    "", // DealInvitation entity doesn't have Reason field
			IsOwner:   invitation.DealMember.IsOwner,
			// CreatedAt: _utils.FormatTimeToString(&invitation.DealMember.CreatedAt),
			// UpdatedAt: _utils.FormatTimeToString(&invitation.DealMember.UpdatedAt),
		}

		// if invitation.DealMember.Role != nil {
		// 	invitationResponse.Role = t.bdsproRoleTransformer.EntityToRoleResponse(invitation.DealMember.Role)
		// }

		if invitation.DealMember.Color != nil {
			invitationResponse.Color = &sharepb.Color{
				Id:              invitation.DealMember.Color.ID,
				Name:            invitation.DealMember.Color.Name,
				ContentColor:    invitation.DealMember.Color.ContentColor,
				BackgroundColor: invitation.DealMember.Color.BackgroundColor,
				Description:     invitation.DealMember.Color.Description,
				IsActive:        invitation.DealMember.Color.IsActive,
			}
		}

		result[i] = invitationResponse
	}
	return &bdspropb.GetDealInvitationsByDealResponse{
		Data:  result,
		Total: int64(total),
	}
}
