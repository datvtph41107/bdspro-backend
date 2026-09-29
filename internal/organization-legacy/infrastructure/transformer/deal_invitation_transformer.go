package transformer

import (
	_models "common/models"
	_utils "common/utils"
	"organization/internal/domain/entity"
	"organization/internal/dto"
	"organization/internal/enums"
	authpb "pb/types/auth"
	organizationpb "pb/types/organization"
	"time"
)

type DealInvitationTransformer interface {
	// Request to Entity transformations
	SendInvitationRequestToEntity(req *organizationpb.SendDealInvitationRequest) *entity.DealMember
	AcceptInvitationRequestToEntity(req *organizationpb.AcceptDealInvitationRequest) *entity.DealMember
	RejectInvitationRequestToEntity(req *organizationpb.RejectDealInvitationRequest) *entity.DealMember
	ResendInvitationRequestToEntity(req *organizationpb.ResendDealInvitationRequest) *entity.DealMember
	WithdrawFromDealRequestToEntity(req *organizationpb.WithdrawFromDealRequest) *entity.DealMember
	RemoveFromDealRequestToEntity(req *organizationpb.RemoveFromDealRequest) *entity.DealMember
	ConfirmWithdrawalRequestToEntity(req *organizationpb.ConfirmWithdrawalRequest) *entity.DealMember

	// Entity to Response transformations
	EntityToSendInvitationResponse(invitation *entity.DealMember) *organizationpb.SendDealInvitationResponse
	EntityToAcceptInvitationResponse(invitation *entity.DealMember) *organizationpb.AcceptDealInvitationResponse
	EntityToRejectInvitationResponse(invitation *entity.DealMember) *organizationpb.RejectDealInvitationResponse
	EntityToResendInvitationResponse(invitation *entity.DealMember) *organizationpb.ResendDealInvitationResponse
	EntityToWithdrawFromDealResponse(invitation *entity.DealMember) *organizationpb.WithdrawFromDealResponse
	EntityToRemoveFromDealResponse(invitation *entity.DealMember) *organizationpb.RemoveFromDealResponse
	EntityToConfirmWithdrawalResponse(invitation *entity.DealMember) *organizationpb.ConfirmWithdrawalResponse
	EntityToGetInvitationResponse(invitation *entity.DealMember) *organizationpb.GetDealInvitationResponse
	EntitiesToGetAcceptedMembersResponse(members []*dto.DealInvitationWithProfile) *organizationpb.GetAcceptedMembersResponse
	EntitiesToSearchMembersResponse(members []*dto.DealInvitationWithProfile, total uint32) *organizationpb.SearchDealMembersResponse
	EntitiesToGetPendingInvitationsResponse(invitations []*dto.DealInvitationWithProfile, total uint32) *organizationpb.GetPendingInvitationsResponse
	EntitiesToGetInvitationsByDealResponse(invitations []*dto.DealInvitationWithProfile, total uint32) *organizationpb.GetDealInvitationsByDealResponse
}

type dealInvitationTransformer struct {
	organizationRoleTransformer OrganizationRoleTransformer
}

func NewDealInvitationTransformer(organizationRoleTransformer OrganizationRoleTransformer) DealInvitationTransformer {
	return &dealInvitationTransformer{
		organizationRoleTransformer: organizationRoleTransformer,
	}
}

// Request to Entity transformations
func (t *dealInvitationTransformer) SendInvitationRequestToEntity(req *organizationpb.SendDealInvitationRequest) *entity.DealMember {
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

	invitation := &entity.DealMember{
		DealID:    req.DealId,
		InviterID: req.InviterId,
		MemberID:  req.MemberId,
		Message:   req.Message,
		RoleID:    req.RoleId,
		// MemberType: memberType,
		Status:    entity.DealMemberStatusInvited,
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

func (t *dealInvitationTransformer) AcceptInvitationRequestToEntity(req *organizationpb.AcceptDealInvitationRequest) *entity.DealMember {
	return &entity.DealMember{
		BaseEntity: _models.BaseEntity{
			ID: req.InvitationId,
		},
	}
}

func (t *dealInvitationTransformer) RejectInvitationRequestToEntity(req *organizationpb.RejectDealInvitationRequest) *entity.DealMember {
	return &entity.DealMember{
		BaseEntity: _models.BaseEntity{
			ID: req.InvitationId,
		},
		Message: req.Reason,
	}
}

func (t *dealInvitationTransformer) ResendInvitationRequestToEntity(req *organizationpb.ResendDealInvitationRequest) *entity.DealMember {
	return &entity.DealMember{
		BaseEntity: _models.BaseEntity{
			ID: req.InvitationId,
		},
	}
}

func (t *dealInvitationTransformer) WithdrawFromDealRequestToEntity(req *organizationpb.WithdrawFromDealRequest) *entity.DealMember {
	return &entity.DealMember{
		BaseEntity: _models.BaseEntity{
			ID: req.InvitationId,
		},
		Message: req.Reason,
	}
}

func (t *dealInvitationTransformer) RemoveFromDealRequestToEntity(req *organizationpb.RemoveFromDealRequest) *entity.DealMember {
	return &entity.DealMember{
		BaseEntity: _models.BaseEntity{
			ID: req.InvitationId,
		},
		Message: req.Reason,
	}
}

func (t *dealInvitationTransformer) ConfirmWithdrawalRequestToEntity(req *organizationpb.ConfirmWithdrawalRequest) *entity.DealMember {
	return &entity.DealMember{
		BaseEntity: _models.BaseEntity{
			ID: req.InvitationId,
		},
	}
}

// Entity to Response transformations
func (t *dealInvitationTransformer) EntityToSendInvitationResponse(invitation *entity.DealMember) *organizationpb.SendDealInvitationResponse {
	return &organizationpb.SendDealInvitationResponse{
		InvitationId: invitation.ID,
		Success:      true,
		Message:      "Invitation sent successfully",
		RoleId:       invitation.RoleID,
	}
}

func (t *dealInvitationTransformer) EntityToAcceptInvitationResponse(invitation *entity.DealMember) *organizationpb.AcceptDealInvitationResponse {
	return &organizationpb.AcceptDealInvitationResponse{
		InvitationId: invitation.ID,
		Success:      true,
	}
}

func (t *dealInvitationTransformer) EntityToRejectInvitationResponse(invitation *entity.DealMember) *organizationpb.RejectDealInvitationResponse {
	return &organizationpb.RejectDealInvitationResponse{
		InvitationId: invitation.ID,
		Success:      true,
	}
}

func (t *dealInvitationTransformer) EntityToResendInvitationResponse(invitation *entity.DealMember) *organizationpb.ResendDealInvitationResponse {
	return &organizationpb.ResendDealInvitationResponse{
		InvitationId: invitation.ID,
		Success:      true,
		Message:      "Invitation resent successfully",
	}
}

func (t *dealInvitationTransformer) EntityToWithdrawFromDealResponse(invitation *entity.DealMember) *organizationpb.WithdrawFromDealResponse {
	return &organizationpb.WithdrawFromDealResponse{
		InvitationId: invitation.ID,
		Success:      true,
	}
}

func (t *dealInvitationTransformer) EntityToRemoveFromDealResponse(invitation *entity.DealMember) *organizationpb.RemoveFromDealResponse {
	return &organizationpb.RemoveFromDealResponse{
		InvitationId: invitation.ID,
		Success:      true,
	}
}

func (t *dealInvitationTransformer) EntityToConfirmWithdrawalResponse(invitation *entity.DealMember) *organizationpb.ConfirmWithdrawalResponse {
	return &organizationpb.ConfirmWithdrawalResponse{
		InvitationId: invitation.ID,
		Success:      true,
	}
}

func (t *dealInvitationTransformer) EntityToGetInvitationResponse(invitation *entity.DealMember) *organizationpb.GetDealInvitationResponse {
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

	invitationResponse := &organizationpb.DealMember{
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
		invitationResponse.Color = &authpb.Color{
			Id:              invitation.Color.ID,
			Name:            invitation.Color.Name,
			ContentColor:    invitation.Color.ContentColor,
			BackgroundColor: invitation.Color.BackgroundColor,
			Description:     invitation.Color.Description,
			IsActive:        invitation.Color.IsActive,
		}
	}

	return &organizationpb.GetDealInvitationResponse{
		Member: &organizationpb.DealMember{
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

func (t *dealInvitationTransformer) EntitiesToGetAcceptedMembersResponse(members []*dto.DealInvitationWithProfile) *organizationpb.GetAcceptedMembersResponse {
	result := make([]*organizationpb.DealMember, len(members))
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

		result[i] = &organizationpb.DealMember{
			MemberId: member.DealMember.MemberID,
			Profile:  member.Invitee,
			// Name:     member.Invitee.FullName,
			// Email:    "", // ProfileItem doesn't have email field
			RoleId:      &member.DealMember.RoleID,
			RoleKey:     uint32(member.DealMember.RoleKey),
			RespondedAt: _utils.FormatTimeToString(member.DealMember.RespondedAt),
		}
	}

	return &organizationpb.GetAcceptedMembersResponse{
		Members: result,
	}
}

func (t *dealInvitationTransformer) EntitiesToSearchMembersResponse(members []*dto.DealInvitationWithProfile, total uint32) *organizationpb.SearchDealMembersResponse {
	result := make([]*organizationpb.DealMember, len(members))
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

		result[i] = &organizationpb.DealMember{
			MemberId: member.DealMember.MemberID,
			Profile:  member.Invitee,
			// Name:     member.Invitee.FullName,
			// Email:    "", // ProfileItem doesn't have email field
			RoleId:      &member.DealMember.RoleID,
			RoleKey:     uint32(member.DealMember.RoleKey),
			RespondedAt: _utils.FormatTimeToString(member.DealMember.RespondedAt),
		}
	}

	return &organizationpb.SearchDealMembersResponse{
		Members: result,
		Total:   int64(total),
	}
}

func (t *dealInvitationTransformer) EntitiesToGetPendingInvitationsResponse(invitations []*dto.DealInvitationWithProfile, total uint32) *organizationpb.GetPendingInvitationsResponse {
	result := make([]*organizationpb.DealMember, len(invitations))
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

		invitationResponse := &organizationpb.DealMember{
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
			invitationResponse.Color = &authpb.Color{
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
	return &organizationpb.GetPendingInvitationsResponse{
		Data:  result,
		Total: int64(total),
	}
}

func (t *dealInvitationTransformer) EntitiesToGetInvitationsByDealResponse(invitations []*dto.DealInvitationWithProfile, total uint32) *organizationpb.GetDealInvitationsByDealResponse {
	result := make([]*organizationpb.DealMember, len(invitations))
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

		invitationResponse := &organizationpb.DealMember{
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
		// 	invitationResponse.Role = t.organizationRoleTransformer.EntityToRoleResponse(invitation.DealMember.Role)
		// }

		if invitation.DealMember.Color != nil {
			invitationResponse.Color = &authpb.Color{
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
	return &organizationpb.GetDealInvitationsByDealResponse{
		Data:  result,
		Total: int64(total),
	}
}
