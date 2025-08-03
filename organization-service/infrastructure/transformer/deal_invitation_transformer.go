package transformer

import (
	_utils "common/utils"
	"organization/internal/domain/entity"
	"organization/internal/dto"
	"organization/internal/enums"
	organizationpb "pb/types/organization"
	"time"
)

type DealInvitationTransformer interface {
	// Request to Entity transformations
	SendInvitationRequestToEntity(req *organizationpb.SendInvitationRequest) *entity.DealInvitation
	AcceptInvitationRequestToEntity(req *organizationpb.AcceptInvitationRequest) *entity.DealInvitation
	RejectInvitationRequestToEntity(req *organizationpb.RejectInvitationRequest) *entity.DealInvitation
	ResendInvitationRequestToEntity(req *organizationpb.ResendInvitationRequest) *entity.DealInvitation
	WithdrawFromDealRequestToEntity(req *organizationpb.WithdrawFromDealRequest) *entity.DealInvitation
	RemoveFromDealRequestToEntity(req *organizationpb.RemoveFromDealRequest) *entity.DealInvitation
	ConfirmWithdrawalRequestToEntity(req *organizationpb.ConfirmWithdrawalRequest) *entity.DealInvitation

	// Entity to Response transformations
	EntityToSendInvitationResponse(invitation *entity.DealInvitation) *organizationpb.SendInvitationResponse
	EntityToAcceptInvitationResponse(invitation *entity.DealInvitation) *organizationpb.AcceptInvitationResponse
	EntityToRejectInvitationResponse(invitation *entity.DealInvitation) *organizationpb.RejectInvitationResponse
	EntityToResendInvitationResponse(invitation *entity.DealInvitation) *organizationpb.ResendInvitationResponse
	EntityToWithdrawFromDealResponse(invitation *entity.DealInvitation) *organizationpb.WithdrawFromDealResponse
	EntityToRemoveFromDealResponse(invitation *entity.DealInvitation) *organizationpb.RemoveFromDealResponse
	EntityToConfirmWithdrawalResponse(invitation *entity.DealInvitation) *organizationpb.ConfirmWithdrawalResponse
	EntityToGetInvitationResponse(invitation *entity.DealInvitation) *organizationpb.GetInvitationResponse
	EntitiesToGetAcceptedMembersResponse(members []*dto.DealInvitationWithProfile) *organizationpb.GetAcceptedMembersResponse
	EntitiesToSearchMembersResponse(members []*dto.DealInvitationWithProfile, total uint32) *organizationpb.SearchMembersResponse
	EntitiesToGetPendingInvitationsResponse(invitations []*dto.DealInvitationWithProfile, total uint32) *organizationpb.GetPendingInvitationsResponse
	EntitiesToGetInvitationsByDealResponse(invitations []*dto.DealInvitationWithProfile, total uint32) *organizationpb.GetInvitationsByDealResponse
}

type dealInvitationTransformer struct{}

func NewDealInvitationTransformer() DealInvitationTransformer {
	return &dealInvitationTransformer{}
}

// Request to Entity transformations
func (t *dealInvitationTransformer) SendInvitationRequestToEntity(req *organizationpb.SendInvitationRequest) *entity.DealInvitation {
	// Convert string role to DealMemberType enum
	var memberType enums.DealMemberType
	switch req.Role {
	case "customer":
		memberType = enums.DealMemberTypeCustomer
	case "partner":
		memberType = enums.DealMemberTypePartner
	case "member":
		memberType = enums.DealMemberTypeMember
	default:
		memberType = enums.DealMemberTypeMember // default
	}

	invitation := &entity.DealInvitation{
		DealID:     req.DealId,
		InviterID:  req.InviterId,
		InviteeID:  req.InviteeId,
		Message:    req.Message,
		MemberType: memberType,
		Status:     entity.DealInvitationStatusInvited,
		InvitedAt:  time.Now(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	
	if req.ColorId != nil {
		invitation.ColorId = req.ColorId
	}
	
	return invitation
}

func (t *dealInvitationTransformer) AcceptInvitationRequestToEntity(req *organizationpb.AcceptInvitationRequest) *entity.DealInvitation {
	return &entity.DealInvitation{
		ID: req.InvitationId,
	}
}

func (t *dealInvitationTransformer) RejectInvitationRequestToEntity(req *organizationpb.RejectInvitationRequest) *entity.DealInvitation {
	return &entity.DealInvitation{
		ID:      req.InvitationId,
		Message: req.Reason,
	}
}

func (t *dealInvitationTransformer) ResendInvitationRequestToEntity(req *organizationpb.ResendInvitationRequest) *entity.DealInvitation {
	return &entity.DealInvitation{
		ID: req.InvitationId,
	}
}

func (t *dealInvitationTransformer) WithdrawFromDealRequestToEntity(req *organizationpb.WithdrawFromDealRequest) *entity.DealInvitation {
	return &entity.DealInvitation{
		ID:      req.InvitationId,
		Message: req.Reason,
	}
}

func (t *dealInvitationTransformer) RemoveFromDealRequestToEntity(req *organizationpb.RemoveFromDealRequest) *entity.DealInvitation {
	return &entity.DealInvitation{
		ID:      req.InvitationId,
		Message: req.Reason,
	}
}

func (t *dealInvitationTransformer) ConfirmWithdrawalRequestToEntity(req *organizationpb.ConfirmWithdrawalRequest) *entity.DealInvitation {
	return &entity.DealInvitation{
		ID: req.InvitationId,
	}
}

// Entity to Response transformations
func (t *dealInvitationTransformer) EntityToSendInvitationResponse(invitation *entity.DealInvitation) *organizationpb.SendInvitationResponse {
	return &organizationpb.SendInvitationResponse{
		InvitationId: invitation.ID,
		Success:      true,
		Message:      "Invitation sent successfully",
	}
}

func (t *dealInvitationTransformer) EntityToAcceptInvitationResponse(invitation *entity.DealInvitation) *organizationpb.AcceptInvitationResponse {
	return &organizationpb.AcceptInvitationResponse{
		InvitationId: invitation.ID,
		Success:      true,
	}
}

func (t *dealInvitationTransformer) EntityToRejectInvitationResponse(invitation *entity.DealInvitation) *organizationpb.RejectInvitationResponse {
	return &organizationpb.RejectInvitationResponse{
		InvitationId: invitation.ID,
		Success:      true,
	}
}

func (t *dealInvitationTransformer) EntityToResendInvitationResponse(invitation *entity.DealInvitation) *organizationpb.ResendInvitationResponse {
	return &organizationpb.ResendInvitationResponse{
		InvitationId: invitation.ID,
		Success:      true,
		Message:      "Invitation resent successfully",
	}
}

func (t *dealInvitationTransformer) EntityToWithdrawFromDealResponse(invitation *entity.DealInvitation) *organizationpb.WithdrawFromDealResponse {
	return &organizationpb.WithdrawFromDealResponse{
		InvitationId: invitation.ID,
		Success:      true,
	}
}

func (t *dealInvitationTransformer) EntityToRemoveFromDealResponse(invitation *entity.DealInvitation) *organizationpb.RemoveFromDealResponse {
	return &organizationpb.RemoveFromDealResponse{
		InvitationId: invitation.ID,
		Success:      true,
	}
}

func (t *dealInvitationTransformer) EntityToConfirmWithdrawalResponse(invitation *entity.DealInvitation) *organizationpb.ConfirmWithdrawalResponse {
	return &organizationpb.ConfirmWithdrawalResponse{
		InvitationId: invitation.ID,
		Success:      true,
	}
}

func (t *dealInvitationTransformer) EntityToGetInvitationResponse(invitation *entity.DealInvitation) *organizationpb.GetInvitationResponse {
	if invitation == nil {
		return nil
	}

	// Convert DealMemberType to string
	var role string
	switch invitation.MemberType {
	case enums.DealMemberTypeCustomer:
		role = "customer"
	case enums.DealMemberTypePartner:
		role = "partner"
	case enums.DealMemberTypeMember:
		role = "member"
	default:
		role = "member"
	}

	invitationResponse := &organizationpb.DealInvitation{
		Id:        invitation.ID,
		DealId:    invitation.DealID,
		InviterId: invitation.InviterID,
		InviteeId: invitation.InviteeID,
		Status:    uint32(invitation.Status),
		Message:   invitation.Message,
		Role:      role,
		Reason:    "", // DealInvitation entity doesn't have Reason field
		CreatedAt: _utils.FormatTimeToString(&invitation.CreatedAt),
		UpdatedAt: _utils.FormatTimeToString(&invitation.UpdatedAt),
	}
	
	if invitation.Color != nil {
		invitationResponse.Color = &organizationpb.Color{
			Id:              uint32(invitation.Color.ID),
			Name:            invitation.Color.Name,
			ContentColor:    invitation.Color.ContentColor,
			BackgroundColor: invitation.Color.BackgroundColor,
			Description:     invitation.Color.Description,
			IsActive:        invitation.Color.IsActive,
		}
	}
	
	return &organizationpb.GetInvitationResponse{
		Invitation: invitationResponse,
	}
}

func (t *dealInvitationTransformer) EntitiesToGetAcceptedMembersResponse(members []*dto.DealInvitationWithProfile) *organizationpb.GetAcceptedMembersResponse {
	result := make([]*organizationpb.DealInvitationMember, len(members))
	for i, member := range members {
		// Convert DealMemberType to string
		var role string
		switch member.DealInvitation.MemberType {
		case enums.DealMemberTypeCustomer:
			role = "customer"
		case enums.DealMemberTypePartner:
			role = "partner"
		case enums.DealMemberTypeMember:
			role = "member"
		default:
			role = "member"
		}

		result[i] = &organizationpb.DealInvitationMember{
			UserId:   member.DealInvitation.InviteeID,
			Name:     member.Invitee.FullName,
			Email:    "", // ProfileItem doesn't have email field
			Role:     role,
			JoinedAt: _utils.FormatTimeToString(&member.DealInvitation.CreatedAt),
		}
	}

	return &organizationpb.GetAcceptedMembersResponse{
		Members: result,
	}
}

func (t *dealInvitationTransformer) EntitiesToSearchMembersResponse(members []*dto.DealInvitationWithProfile, total uint32) *organizationpb.SearchMembersResponse {
	result := make([]*organizationpb.DealInvitationMember, len(members))
	for i, member := range members {
		// Convert DealMemberType to string
		var role string
		switch member.DealInvitation.MemberType {
		case enums.DealMemberTypeCustomer:
			role = "customer"
		case enums.DealMemberTypePartner:
			role = "partner"
		case enums.DealMemberTypeMember:
			role = "member"
		default:
			role = "member"
		}

		result[i] = &organizationpb.DealInvitationMember{
			UserId:   member.DealInvitation.InviteeID,
			Name:     member.Invitee.FullName,
			Email:    "", // ProfileItem doesn't have email field
			Role:     role,
			JoinedAt: _utils.FormatTimeToString(&member.DealInvitation.CreatedAt),
		}
	}

	return &organizationpb.SearchMembersResponse{
		Members: result,
		Total:   int64(total),
	}
}

func (t *dealInvitationTransformer) EntitiesToGetPendingInvitationsResponse(invitations []*dto.DealInvitationWithProfile, total uint32) *organizationpb.GetPendingInvitationsResponse {
	result := make([]*organizationpb.DealInvitation, len(invitations))
	for i, invitation := range invitations {
		// Convert DealMemberType to string
		var role string
		switch invitation.DealInvitation.MemberType {
		case enums.DealMemberTypeCustomer:
			role = "customer"
		case enums.DealMemberTypePartner:
			role = "partner"
		case enums.DealMemberTypeMember:
			role = "member"
		default:
			role = "member"
		}

		invitationResponse := &organizationpb.DealInvitation{
			Id:        invitation.DealInvitation.ID,
			DealId:    invitation.DealInvitation.DealID,
			InviterId: invitation.DealInvitation.InviterID,
			InviteeId: invitation.DealInvitation.InviteeID,
			Status:    uint32(invitation.DealInvitation.Status),
			Message:   invitation.DealInvitation.Message,
			Role:      role,
			Reason:    "", // DealInvitation entity doesn't have Reason field
			CreatedAt: _utils.FormatTimeToString(&invitation.DealInvitation.CreatedAt),
			UpdatedAt: _utils.FormatTimeToString(&invitation.DealInvitation.UpdatedAt),
		}
		
		if invitation.DealInvitation.Color != nil {
			invitationResponse.Color = &organizationpb.Color{
				Id:              uint32(invitation.DealInvitation.Color.ID),
				Name:            invitation.DealInvitation.Color.Name,
				ContentColor:    invitation.DealInvitation.Color.ContentColor,
				BackgroundColor: invitation.DealInvitation.Color.BackgroundColor,
				Description:     invitation.DealInvitation.Color.Description,
				IsActive:        invitation.DealInvitation.Color.IsActive,
			}
		}
		
		result[i] = invitationResponse
	}
	return &organizationpb.GetPendingInvitationsResponse{
		Data:  result,
		Total: int64(total),
	}
}

func (t *dealInvitationTransformer) EntitiesToGetInvitationsByDealResponse(invitations []*dto.DealInvitationWithProfile, total uint32) *organizationpb.GetInvitationsByDealResponse {
	result := make([]*organizationpb.DealInvitation, len(invitations))
	for i, invitation := range invitations {
		// Convert DealMemberType to string
		var role string
		switch invitation.DealInvitation.MemberType {
		case enums.DealMemberTypeCustomer:
			role = "customer"
		case enums.DealMemberTypePartner:
			role = "partner"
		case enums.DealMemberTypeMember:
			role = "member"
		default:
			role = "member"
		}

		invitationResponse := &organizationpb.DealInvitation{
			Id:        invitation.DealInvitation.ID,
			DealId:    invitation.DealInvitation.DealID,
			InviterId: invitation.DealInvitation.InviterID,
			InviteeId: invitation.DealInvitation.InviteeID,
			Status:    uint32(invitation.DealInvitation.Status),
			Message:   invitation.DealInvitation.Message,
			Role:      role,
			Reason:    "", // DealInvitation entity doesn't have Reason field
			CreatedAt: _utils.FormatTimeToString(&invitation.DealInvitation.CreatedAt),
			UpdatedAt: _utils.FormatTimeToString(&invitation.DealInvitation.UpdatedAt),
		}
		
		if invitation.DealInvitation.Color != nil {
			invitationResponse.Color = &organizationpb.Color{
				Id:              uint32(invitation.DealInvitation.Color.ID),
				Name:            invitation.DealInvitation.Color.Name,
				ContentColor:    invitation.DealInvitation.Color.ContentColor,
				BackgroundColor: invitation.DealInvitation.Color.BackgroundColor,
				Description:     invitation.DealInvitation.Color.Description,
				IsActive:        invitation.DealInvitation.Color.IsActive,
			}
		}
		
		result[i] = invitationResponse
	}
	return &organizationpb.GetInvitationsByDealResponse{
		Data:  result,
		Total: int64(total),
	}
}
